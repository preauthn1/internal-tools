<?php

namespace App\Console\Commands;

use App\Enums\GameTransferDirection;
use App\Enums\GameTransferStatus;
use App\Enums\WalletTransactionType;
use App\Models\GameProvider as GameProviderModel;
use App\Models\GameTransfer;
use App\Models\User;
use App\Services\Games\Data\MoneyAmount;
use Illuminate\Console\Attributes\Description;
use Illuminate\Console\Attributes\Signature;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

/**
 * 為錢包帳本上線前既有的帳號回填 wallet_transactions。
 * 依「初始金 + 已完成的 game_transfers（時間序）」重建每筆異動與結餘。
 * 冪等：只處理「有錢包且尚無任何帳本」的帳號，重跑不會重複。
 */
#[Signature('wallet:backfill-ledger')]
#[Description('為既有帳號回填錢包帳本（wallet_transactions）')]
class BackfillWalletLedger extends Command
{
    public function handle(): int
    {
        $initial = MoneyAmount::from((string) config('games.initial_balance', 10000))->value();
        $backfilled = 0;

        // provider key → game_providers.id 對照
        $providerIds = GameProviderModel::query()->pluck('id', 'key')->all();

        User::query()
            ->has('wallet')
            ->doesntHave('walletTransactions')
            ->with('wallet')
            ->chunkById(100, function ($users) use ($initial, $providerIds, &$backfilled): void {
                foreach ($users as $user) {
                    $this->backfillUser($user, $initial, $providerIds);
                    $backfilled++;
                }
            });

        $this->info("已回填 {$backfilled} 個帳號的錢包帳本。");

        return self::SUCCESS;
    }

    /**
     * @param  numeric-string  $initial
     * @param  array<string, int>  $providerIds
     */
    private function backfillUser(User $user, string $initial, array $providerIds): void
    {
        DB::transaction(function () use ($user, $initial, $providerIds): void {
            $rows = [];
            /** @var numeric-string $running */
            $running = $initial;

            // 初始金
            $rows[] = [
                'type' => WalletTransactionType::Initial->value,
                'game_provider_id' => null,
                'amount' => $initial,
                'balance_after' => $running,
                'reference' => null,
                'meta' => json_encode(['backfilled' => true], JSON_THROW_ON_ERROR),
                'created_at' => $user->created_at,
                'updated_at' => $user->created_at,
            ];

            // 已完成的轉帳（時間序）；失敗+退回互相抵銷，不影響結餘，故僅取 completed
            $transfers = GameTransfer::query()
                ->where('user_id', $user->id)
                ->where('status', GameTransferStatus::Completed->value)
                ->orderBy('resolved_at')
                ->orderBy('id')
                ->get();

            foreach ($transfers as $transfer) {
                $isIn = $transfer->direction === GameTransferDirection::In;
                $running = $isIn
                    ? bcsub($running, $transfer->amount, 2)
                    : bcadd($running, $transfer->amount, 2);

                $rows[] = [
                    'type' => ($isIn ? WalletTransactionType::TransferIn : WalletTransactionType::TransferOut)->value,
                    'game_provider_id' => $providerIds[$transfer->provider] ?? null,
                    'amount' => $isIn ? '-'.$transfer->amount : (string) $transfer->amount,
                    'balance_after' => $running,
                    'reference' => $transfer->reference,
                    'meta' => json_encode(['backfilled' => true], JSON_THROW_ON_ERROR),
                    'created_at' => $transfer->resolved_at ?? $transfer->created_at,
                    'updated_at' => $transfer->resolved_at ?? $transfer->created_at,
                ];
            }

            $user->walletTransactions()->insert(
                array_map(fn (array $row): array => [...$row, 'user_id' => $user->id], $rows),
            );

            if (bccomp($running, $user->wallet->balance, 2) !== 0) {
                $this->warn("[user {$user->id}] 重建結餘 {$running} 與目前餘額 {$user->wallet->balance} 不一致。");
            }
        });
    }
}
