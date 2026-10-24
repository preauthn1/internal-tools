<?php

namespace App\Console\Commands;

use App\Enums\GameTransferAuditStatus;
use App\Enums\GameTransferDirection;
use App\Enums\GameTransferStatus;
use App\Models\GameTransfer;
use App\Models\GameTransferAudit;
use App\Services\Games\GameTransferAuditService;
use Illuminate\Console\Attributes\Description;
use Illuminate\Console\Attributes\Signature;
use Illuminate\Console\Command;

#[Signature('game:audit-wallet-sessions {--reference= : 僅稽核指定轉出 merchantOrderNo}')]
#[Description('以遊戲紀錄驗算已完成的遊戲錢包轉出金額')]
class AuditGameWalletSessions extends Command
{
    public function handle(GameTransferAuditService $service): int
    {
        $reference = $this->option('reference');
        $this->discoverCompletedTransfers($service, is_string($reference) ? $reference : null);

        $statuses = [GameTransferAuditStatus::PendingRecords->value];
        if (is_string($reference)) {
            $statuses[] = GameTransferAuditStatus::ManualReview->value;
        }

        $query = GameTransferAudit::query()
            ->with('transferOut:id,reference')
            ->whereIn('status', $statuses);
        if (is_string($reference)) {
            $query->whereHas('transferOut', fn ($transfer) => $transfer->where('reference', $reference));
        } else {
            $query->where(fn ($due) => $due->whereNull('next_retry_at')->orWhere('next_retry_at', '<=', now()));
        }

        $processed = 0;
        $query->orderBy('id')->chunkById(100, function ($audits) use ($service, $reference, &$processed): void {
            foreach ($audits as $audit) {
                // 先讀 eager-load 的 reference，因為 audit() 內部會 refresh() 重載關聯。
                $transferReference = $audit->transferOut?->reference;
                $status = $service->audit($audit, is_string($reference));
                $this->line("{$transferReference}: {$status->value}");
                $processed++;
            }
        });

        $this->info("processed={$processed}");

        return self::SUCCESS;
    }

    private function discoverCompletedTransfers(GameTransferAuditService $service, ?string $reference): void
    {
        $query = GameTransfer::query()
            ->where('direction', GameTransferDirection::Out->value)
            ->where('status', GameTransferStatus::Completed->value)
            ->whereDoesntHave('audit');
        if ($reference !== null) {
            $query->where('reference', $reference);
        }

        $query->orderBy('id')->chunkById(100, function ($transfers) use ($service): void {
            foreach ($transfers as $transfer) {
                $service->createForTransferOut($transfer);
            }
        });
    }
}
