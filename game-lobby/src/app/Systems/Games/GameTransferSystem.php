<?php

namespace App\Systems\Games;

use App\Enums\GameTransferDirection;
use App\Enums\GameTransferStatus;
use App\Enums\WalletTransactionType;
use App\Models\GameTransfer;
use App\Models\User;
use App\Models\Wallet;
use App\Repositories\Games\WalletLedgerRepository;
use App\Services\Games\Contracts\GameProvider;
use App\Services\Games\Data\MoneyAmount;
use App\Services\Games\Data\SettlementResult;
use App\Services\Games\Exceptions\GameProviderConnectionException;
use App\Services\Games\Exceptions\GameProviderException;
use App\Services\Games\Exceptions\GameTransferPendingException;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Str;
use Throwable;

class GameTransferSystem
{
    private const INVALID_STATE_CODE = 403601;

    private const DUPLICATE_CODE = 500029;

    public function __construct(private WalletLedgerRepository $ledger) {}

    public function transferAllIn(User $user, GameProvider $provider, string $playerName): void
    {
        $transfer = DB::transaction(function () use ($user, $provider, $playerName): ?GameTransfer {
            /** @var Wallet $wallet */
            $wallet = $user->wallet()->lockForUpdate()->firstOrFail();
            $amount = MoneyAmount::from((string) $wallet->balance);
            if ($amount->isZero()) {
                return null;
            }
            if ($amount->isBelowProviderMinimum()) {
                throw new GameTransferPendingException('平台餘額低於遊戲商最低轉帳額 0.10，無法啟動付費遊戲。');
            }

            $existing = $user->gameTransfers()->where('provider', $provider->key())
                ->where('direction', GameTransferDirection::In->value)
                ->whereIn('status', [GameTransferStatus::Pending->value, GameTransferStatus::Unknown->value])
                ->lockForUpdate()->first();
            if ($existing !== null) {
                throw new GameTransferPendingException;
            }

            $transfer = $user->gameTransfers()->create([
                'provider' => $provider->key(),
                'direction' => GameTransferDirection::In,
                'amount' => $amount->value(),
                'reference' => (string) Str::uuid(),
                'status' => GameTransferStatus::Pending,
                'next_retry_at' => now()->addSeconds(30),
                'meta' => ['player_name' => $playerName],
            ]);
            $wallet->balance = bcsub($wallet->balance, $amount->value(), 2);
            $wallet->save();

            $this->ledger->record(
                $user->id,
                WalletTransactionType::TransferIn,
                $provider->key(),
                '-'.$amount->value(),
                (string) $wallet->balance,
                $transfer->reference,
            );

            return $transfer;
        });

        if ($transfer !== null) {
            $this->execute($transfer, $provider, $playerName);
            $transfer->refresh();
            if ($transfer->status !== GameTransferStatus::Completed) {
                throw new GameTransferPendingException;
            }
        }
    }

    public function transferAllOut(User $user, GameProvider $provider, string $playerName): SettlementResult
    {
        $existing = $user->gameTransfers()->where('provider', $provider->key())
            ->where('direction', GameTransferDirection::Out->value)
            ->whereIn('status', [GameTransferStatus::Pending->value, GameTransferStatus::Unknown->value])
            ->first();
        if ($existing !== null) {
            return new SettlementResult('0', true);
        }

        $amount = MoneyAmount::from($provider->walletAmount($playerName));
        if ($amount->isZero()) {
            return new SettlementResult('0');
        }
        if ($amount->isBelowProviderMinimum()) {
            return new SettlementResult('0', false, true);
        }

        $transfer = DB::transaction(function () use ($user, $provider, $playerName, $amount): ?GameTransfer {
            $user->wallet()->lockForUpdate()->firstOrFail();
            $unresolved = $user->gameTransfers()->where('provider', $provider->key())
                ->where('direction', GameTransferDirection::Out->value)
                ->whereIn('status', [GameTransferStatus::Pending->value, GameTransferStatus::Unknown->value])
                ->lockForUpdate()->first();
            if ($unresolved !== null) {
                return null;
            }

            return $user->gameTransfers()->create([
                'provider' => $provider->key(),
                'direction' => GameTransferDirection::Out,
                'amount' => $amount->value(),
                'reference' => (string) Str::uuid(),
                'status' => GameTransferStatus::Pending,
                'next_retry_at' => now()->addSeconds(30),
                'meta' => ['player_name' => $playerName],
            ]);
        });
        if ($transfer === null) {
            return new SettlementResult('0', true);
        }

        $this->execute($transfer, $provider, $playerName);
        $transfer->refresh();

        return new SettlementResult(
            $transfer->status === GameTransferStatus::Completed ? $amount->value() : '0',
            $transfer->status->isUnresolved(),
        );
    }

    public function execute(GameTransfer $transfer, GameProvider $provider, string $playerName): void
    {
        $transfer->forceFill(['attempts' => $transfer->attempts + 1, 'last_attempt_at' => now()])->save();

        try {
            $result = $transfer->direction === GameTransferDirection::In
                ? $provider->transferIn($playerName, (string) $transfer->amount, $transfer->reference)
                : $provider->transferOut($playerName, (string) $transfer->amount, $transfer->reference);
        } catch (GameProviderConnectionException) {
            $this->scheduleRetry($transfer, GameTransferStatus::Unknown, 'connection');
            Log::warning('遊戲轉帳連線狀態未知', ['provider' => $transfer->provider, 'reference' => $transfer->reference, 'category' => 'connection']);

            return;
        } catch (GameProviderException $exception) {
            if ($exception->apiCode === self::DUPLICATE_CODE) {
                $this->scheduleRetry($transfer, GameTransferStatus::Unknown, (string) $exception->apiCode);
            } elseif ($this->isInvalidStateCode($provider, $exception->apiCode)) {
                $this->scheduleRetry($transfer, GameTransferStatus::Pending, (string) $exception->apiCode);
            } else {
                $this->fail($transfer, (string) $exception->apiCode);
            }
            Log::warning('遊戲轉帳遭拒', ['provider' => $transfer->provider, 'reference' => $transfer->reference, 'api_code' => $exception->apiCode, 'category' => 'provider']);

            return;
        }

        try {
            $this->finalize($transfer, $result);
        } catch (Throwable $exception) {
            $this->scheduleRetry($transfer, GameTransferStatus::Unknown, 'local_finalize');
            report($exception);
        }
    }

    /** @param array<string, mixed> $result */
    public function finalize(GameTransfer $transfer, array $result, ?string $providerTransactionNo = null): void
    {
        DB::transaction(function () use ($transfer, $result, $providerTransactionNo): void {
            /** @var GameTransfer $locked */
            $locked = GameTransfer::query()->lockForUpdate()->findOrFail($transfer->id);
            if ($locked->status === GameTransferStatus::Completed) {
                return;
            }
            if ($locked->direction === GameTransferDirection::Out) {
                /** @var Wallet $wallet */
                $wallet = Wallet::query()->where('user_id', $locked->user_id)->lockForUpdate()->firstOrFail();
                $wallet->balance = bcadd($wallet->balance, $locked->amount, 2);
                $wallet->save();

                $this->ledger->record(
                    $locked->user_id,
                    WalletTransactionType::TransferOut,
                    $locked->provider,
                    (string) $locked->amount,
                    (string) $wallet->balance,
                    $locked->reference,
                );
            }
            $locked->forceFill([
                'status' => GameTransferStatus::Completed,
                'resolved_at' => now(),
                'next_retry_at' => null,
                'error_code' => null,
                'provider_transaction_no' => $providerTransactionNo ?? (string) ($result['pgsTransactionNo'] ?? $result['transactionNo'] ?? ''),
                'meta' => [...($locked->meta ?? []), 'provider_response' => $result],
            ])->save();
        });
    }

    public function scheduleRetry(GameTransfer $transfer, GameTransferStatus $status, string $errorCode): void
    {
        $delay = match (true) {
            $transfer->attempts <= 2 => 30,
            $transfer->attempts === 3 => 60,
            default => 120,
        };
        $transfer->forceFill([
            'status' => $status,
            'error_code' => $errorCode,
            'next_retry_at' => now()->addSeconds($delay),
        ])->save();
    }

    public function markManualReview(GameTransfer $transfer, string $errorCode): void
    {
        $transfer->forceFill([
            'status' => GameTransferStatus::ManualReview,
            'error_code' => $errorCode,
            'resolved_at' => now(),
            'next_retry_at' => null,
        ])->save();
    }

    public function isInvalidStateCode(GameProvider $provider, int $apiCode): bool
    {
        /** @var array<int, int> $codes */
        $codes = (array) $provider->setting('invalid_state_codes', [self::INVALID_STATE_CODE]);

        return in_array($apiCode, $codes, true);
    }

    private function fail(GameTransfer $transfer, string $errorCode): void
    {
        DB::transaction(function () use ($transfer, $errorCode): void {
            /** @var GameTransfer $locked */
            $locked = GameTransfer::query()->lockForUpdate()->findOrFail($transfer->id);
            if (! $locked->status->isUnresolved()) {
                return;
            }
            if ($locked->direction === GameTransferDirection::In) {
                /** @var Wallet $wallet */
                $wallet = Wallet::query()->where('user_id', $locked->user_id)->lockForUpdate()->firstOrFail();
                $wallet->balance = bcadd($wallet->balance, $locked->amount, 2);
                $wallet->save();

                $this->ledger->record(
                    $locked->user_id,
                    WalletTransactionType::Refund,
                    $locked->provider,
                    (string) $locked->amount,
                    (string) $wallet->balance,
                    $locked->reference,
                );
            }
            $locked->forceFill([
                'status' => GameTransferStatus::Failed,
                'error_code' => $errorCode,
                'resolved_at' => now(),
                'next_retry_at' => null,
            ])->save();
        });
    }
}
