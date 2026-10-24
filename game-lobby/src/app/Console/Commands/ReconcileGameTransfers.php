<?php

namespace App\Console\Commands;

use App\Enums\GameTransferStatus;
use App\Models\GameTransfer;
use App\Services\Games\GameTransferReconcileService;
use Illuminate\Console\Attributes\Description;
use Illuminate\Console\Attributes\Signature;
use Illuminate\Console\Command;

#[Signature('game:reconcile-transfers {--reference= : 僅重查指定 merchantOrderNo}')]
#[Description('查詢並收斂待處理或狀態未知的遊戲轉帳')]
class ReconcileGameTransfers extends Command
{
    public function handle(GameTransferReconcileService $service): int
    {
        $manualReference = $this->option('reference');
        $statuses = [GameTransferStatus::Pending->value, GameTransferStatus::Unknown->value];
        if ($manualReference) {
            $statuses[] = GameTransferStatus::ManualReview->value;
        }
        $query = GameTransfer::query()->whereIn('status', $statuses);

        if ($reference = $manualReference) {
            $query->where('reference', (string) $reference);
        } else {
            $query->where(fn ($due) => $due->whereNull('next_retry_at')->orWhere('next_retry_at', '<=', now()));
        }

        $processed = 0;
        $query->orderBy('id')->chunkById(100, function ($transfers) use ($service, $manualReference, &$processed): void {
            foreach ($transfers as $transfer) {
                $status = $service->reconcile($transfer, $manualReference !== null);
                $this->line("{$transfer->reference}: {$status->value}");
                $processed++;
            }
        });

        $this->info("processed={$processed}");

        return self::SUCCESS;
    }
}
