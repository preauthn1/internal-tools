<?php

namespace App\Services\Games;

use App\Enums\GameTransferDirection;
use App\Enums\GameTransferStatus;
use App\Models\GameTransfer;
use App\Services\Games\Data\MoneyAmount;
use App\Services\Games\Data\TransferRecordData;
use App\Services\Games\Exceptions\GameProviderException;
use App\Systems\Games\GameTransferSystem;
use Illuminate\Support\Facades\Cache;

class GameTransferReconcileService
{
    private const MAX_ATTEMPTS = 10;

    private const INVALID_STATE_MAX_ATTEMPTS = 3;

    public function __construct(
        private GameProviderManager $manager,
        private GameTransferSystem $transfers,
    ) {}

    public function reconcile(GameTransfer $transfer, bool $force = false): GameTransferStatus
    {
        $lock = Cache::lock("game-transfer:reconcile:{$transfer->id}", 60);
        if (! $lock->get()) {
            return $transfer->fresh()->status;
        }

        try {
            return $this->reconcileLocked($transfer, $force);
        } finally {
            $lock->release();
        }
    }

    private function reconcileLocked(GameTransfer $transfer, bool $force): GameTransferStatus
    {
        $transfer->refresh();
        if (! $transfer->status->isUnresolved() || ($transfer->status === GameTransferStatus::ManualReview && ! $force)) {
            return $transfer->status;
        }
        if ($transfer->attempts >= self::MAX_ATTEMPTS && ! $force) {
            $this->transfers->markManualReview($transfer, 'max_attempts');

            return GameTransferStatus::ManualReview;
        }

        $transfer->forceFill(['attempts' => $transfer->attempts + 1])->save();
        $provider = $this->manager->provider($transfer->provider);
        $playerName = (string) ($transfer->meta['player_name'] ?? '');

        try {
            $record = $provider->transferRecord($transfer->reference);
        } catch (GameProviderException $exception) {
            $this->transfers->scheduleRetry($transfer, GameTransferStatus::Unknown, (string) $exception->apiCode);

            return GameTransferStatus::Unknown;
        }

        if ($record !== null) {
            return $this->resolveRecord($transfer, $record, $playerName);
        }

        if ($this->transfers->isInvalidStateCode($provider, (int) $transfer->error_code)) {
            if ($transfer->attempts > self::INVALID_STATE_MAX_ATTEMPTS) {
                $this->transfers->markManualReview($transfer, 'invalid_state');

                return GameTransferStatus::ManualReview;
            }
            try {
                $provider->logoutPlayer($playerName);
                $this->transfers->execute($transfer, $provider, $playerName);
            } catch (GameProviderException) {
                $this->transfers->scheduleRetry($transfer, GameTransferStatus::Pending, 'logout_pending');
            }

            return $transfer->fresh()->status;
        }

        if ($transfer->created_at !== null && $transfer->created_at->diffInSeconds(now()) < 120) {
            $this->transfers->scheduleRetry($transfer, GameTransferStatus::Unknown, 'not_found');

            return GameTransferStatus::Unknown;
        }

        if ($transfer->direction === GameTransferDirection::Out) {
            try {
                $provider->logoutPlayer($playerName);
            } catch (GameProviderException) {
                $this->transfers->scheduleRetry($transfer, GameTransferStatus::Pending, 'logout_pending');

                return GameTransferStatus::Pending;
            }
        }

        $this->transfers->execute($transfer, $provider, $playerName);

        return $transfer->fresh()->status;
    }

    private function resolveRecord(
        GameTransfer $transfer,
        TransferRecordData $record,
        string $playerName,
    ): GameTransferStatus {
        $amount = MoneyAmount::from((string) $transfer->amount);
        if (! $record->matches($transfer->reference, $playerName, $transfer->direction, $amount)) {
            $this->transfers->markManualReview($transfer, 'record_mismatch');

            return GameTransferStatus::ManualReview;
        }

        $this->transfers->finalize($transfer, $record->raw, $record->providerTransactionNo);

        return GameTransferStatus::Completed;
    }
}
