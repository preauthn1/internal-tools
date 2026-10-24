<?php

namespace App\Services\Games;

use App\Enums\GameTransferAuditStatus;
use App\Enums\GameTransferDirection;
use App\Enums\GameTransferStatus;
use App\Models\GameRecord;
use App\Models\GameTransfer;
use App\Models\GameTransferAudit;
use App\Services\Games\Data\MoneyAmount;
use Carbon\CarbonImmutable;
use Illuminate\Support\Facades\Cache;
use InvalidArgumentException;

class GameTransferAuditService
{
    public function __construct(private GameProviderManager $manager) {}

    public function createForTransferOut(GameTransfer $transferOut): GameTransferAudit
    {
        if ($transferOut->direction !== GameTransferDirection::Out || $transferOut->status !== GameTransferStatus::Completed) {
            throw new InvalidArgumentException('只有已完成的遊戲轉出流水可以建立財務稽核。');
        }

        $transferIn = GameTransfer::query()
            ->where('user_id', $transferOut->user_id)
            ->where('provider', $transferOut->provider)
            ->where('direction', GameTransferDirection::In->value)
            ->where('status', GameTransferStatus::Completed->value)
            ->where('created_at', '<=', $transferOut->created_at)
            ->latest('created_at')
            ->first();

        return GameTransferAudit::query()->firstOrCreate(
            ['transfer_out_id' => $transferOut->id],
            [
                'user_id' => $transferOut->user_id,
                'provider' => $transferOut->provider,
                'player_name' => (string) ($transferOut->meta['player_name'] ?? $transferIn?->meta['player_name'] ?? ''),
                'transfer_in_id' => $transferIn?->id,
                'period_started_at' => $transferIn?->created_at,
                'period_ended_at' => $transferOut->created_at,
                'actual_transfer_out' => $transferOut->amount,
                'status' => GameTransferAuditStatus::PendingRecords,
                'next_retry_at' => now(),
            ],
        );
    }

    public function audit(GameTransferAudit $audit, bool $force = false): GameTransferAuditStatus
    {
        $audit->refresh();
        if ($audit->status === GameTransferAuditStatus::Matched) {
            return $audit->status;
        }
        if ($audit->status === GameTransferAuditStatus::ManualReview && ! $force) {
            return $audit->status;
        }

        $audit->forceFill(['attempts' => $audit->attempts + 1])->save();
        if ($audit->transfer_in_id === null || $audit->period_started_at === null) {
            return $this->deferOrReview($audit, 'missing_transfer_in');
        }
        if (! $this->recordsAreReady($audit)) {
            return $this->deferOrReview($audit, 'records_not_ready');
        }

        /** @var GameTransfer $transferIn */
        $transferIn = $audit->transferIn()->firstOrFail();
        $openingBalance = $this->openingBalance($transferIn);
        if ($openingBalance === null) {
            return $this->deferOrReview($audit, 'missing_opening_balance');
        }

        $records = GameRecord::query()
            ->where('user_id', $audit->user_id)
            ->where('provider', $audit->provider)
            ->where('player_name', $audit->player_name)
            ->whereNotNull('bet_at')
            ->whereBetween('bet_at', [$audit->period_started_at, $audit->period_ended_at]);

        $recordTotals = (clone $records)->toBase()
            ->selectRaw('COALESCE(SUM(bet_amount), 0) AS total_bet')
            ->selectRaw('COALESCE(SUM(payout_amount), 0) AS total_payout')
            ->selectRaw('COUNT(*) AS record_count')
            ->first();
        $currencies = (clone $records)->whereNotNull('currency')
            ->where('currency', '!=', '')
            ->distinct()
            ->orderBy('currency')
            ->pluck('currency')
            ->all();

        $totalTransferIn = $this->sumTransfers($audit, GameTransferDirection::In);
        $priorTransferOut = $this->sumTransfers($audit, GameTransferDirection::Out, $audit->transfer_out_id);
        $totalBet = MoneyAmount::from((string) $recordTotals->total_bet)->value();
        $totalPayout = MoneyAmount::from((string) $recordTotals->total_payout)->value();
        /** @var numeric-string $expected */
        $expected = bcadd($openingBalance, $totalTransferIn, 2);
        $expected = bcsub($expected, $priorTransferOut, 2);
        $expected = bcsub($expected, $totalBet, 2);
        $expected = bcadd($expected, $totalPayout, 2);
        $expected = bcadd($expected, $audit->total_adjustment, 2);
        /** @var numeric-string $difference */
        $difference = bcsub($audit->actual_transfer_out, $expected, 2);

        $audit->forceFill([
            'opening_balance' => $openingBalance,
            'total_transfer_in' => $totalTransferIn,
            'prior_transfer_out' => $priorTransferOut,
            'total_bet' => $totalBet,
            'total_payout' => $totalPayout,
            'expected_transfer_out' => $expected,
            'difference' => $difference,
            'record_count' => (int) $recordTotals->record_count,
            'audited_at' => now(),
            'meta' => ['currencies' => $currencies],
        ])->save();

        if (bccomp($difference, '0', 2) === 0) {
            if (! $this->reviewWindowElapsed($audit)) {
                $audit->forceFill([
                    'status' => GameTransferAuditStatus::PendingRecords,
                    'reason' => 'stabilizing_records',
                    'next_retry_at' => now()->addMinutes((int) $this->setting($audit->provider, 'audit_retry_minutes', 5)),
                ])->save();

                return GameTransferAuditStatus::PendingRecords;
            }

            $audit->forceFill([
                'status' => GameTransferAuditStatus::Matched,
                'reason' => null,
                'next_retry_at' => null,
            ])->save();

            return GameTransferAuditStatus::Matched;
        }

        return $this->deferOrReview($audit, 'amount_mismatch');
    }

    private function recordsAreReady(GameTransferAudit $audit): bool
    {
        $cursor = Cache::get("game-records:cursor:{$audit->provider}");
        if (! is_string($cursor)) {
            return false;
        }

        $delay = (int) $this->setting($audit->provider, 'audit_delay_minutes', 30);

        return CarbonImmutable::parse($cursor)->utc()
            ->greaterThanOrEqualTo(CarbonImmutable::instance($audit->period_ended_at)->utc()->addMinutes($delay));
    }

    /** @return numeric-string|null */
    private function openingBalance(GameTransfer $transferIn): ?string
    {
        $response = $transferIn->meta['provider_response'] ?? null;
        if (! is_array($response)) {
            return null;
        }

        $before = $response['beforeAmount'] ?? null;
        if (is_int($before) || is_float($before) || is_string($before)) {
            return MoneyAmount::from($before)->value();
        }

        $after = $response['afterAmount'] ?? null;
        if (! is_int($after) && ! is_float($after) && ! is_string($after)) {
            return null;
        }

        $afterAmount = MoneyAmount::from($after)->value();
        if (bccomp($afterAmount, $transferIn->amount, 2) < 0) {
            return null;
        }

        /** @var numeric-string $opening */
        $opening = bcsub($afterAmount, $transferIn->amount, 2);

        return $opening;
    }

    /** @return numeric-string */
    private function sumTransfers(
        GameTransferAudit $audit,
        GameTransferDirection $direction,
        ?int $excludedTransferId = null,
    ): string {
        $query = GameTransfer::query()
            ->where('user_id', $audit->user_id)
            ->where('provider', $audit->provider)
            ->where('direction', $direction->value)
            ->where('status', GameTransferStatus::Completed->value)
            ->whereBetween('created_at', [$audit->period_started_at, $audit->period_ended_at]);
        if ($excludedTransferId !== null) {
            $query->whereKeyNot($excludedTransferId);
        }

        return MoneyAmount::from((string) $query->sum('amount'))->value();
    }

    private function deferOrReview(GameTransferAudit $audit, string $reason): GameTransferAuditStatus
    {
        $maxAttempts = (int) $this->setting($audit->provider, 'audit_max_attempts', 24);
        $shouldReview = $this->reviewWindowElapsed($audit)
            || $audit->attempts >= $maxAttempts;
        $status = $shouldReview
            ? GameTransferAuditStatus::ManualReview
            : GameTransferAuditStatus::PendingRecords;

        $audit->forceFill([
            'status' => $status,
            'reason' => $reason,
            'next_retry_at' => $shouldReview
                ? null
                : now()->addMinutes((int) $this->setting($audit->provider, 'audit_retry_minutes', 5)),
        ])->save();

        return $status;
    }

    private function reviewWindowElapsed(GameTransferAudit $audit): bool
    {
        $reviewAfter = (int) $this->setting($audit->provider, 'audit_review_after_minutes', 120);

        return $audit->period_ended_at->lte(now()->subMinutes($reviewAfter));
    }

    private function setting(string $provider, string $key, mixed $default): mixed
    {
        return $this->manager->provider($provider)->setting($key, $default);
    }
}
