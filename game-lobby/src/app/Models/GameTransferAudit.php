<?php

namespace App\Models;

use App\Enums\GameTransferAuditStatus;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Support\Carbon;

/**
 * @property int $id
 * @property int $user_id
 * @property string $provider
 * @property string $player_name
 * @property int|null $transfer_in_id
 * @property int $transfer_out_id
 * @property Carbon|null $period_started_at
 * @property Carbon $period_ended_at
 * @property numeric-string|null $opening_balance
 * @property numeric-string $total_transfer_in
 * @property numeric-string $prior_transfer_out
 * @property numeric-string $total_bet
 * @property numeric-string $total_payout
 * @property numeric-string $total_adjustment
 * @property numeric-string|null $expected_transfer_out
 * @property numeric-string $actual_transfer_out
 * @property numeric-string|null $difference
 * @property int $record_count
 * @property GameTransferAuditStatus $status
 * @property int $attempts
 * @property Carbon|null $next_retry_at
 * @property Carbon|null $audited_at
 * @property string|null $reason
 * @property array<string, mixed>|null $meta
 */
class GameTransferAudit extends Model
{
    protected $fillable = [
        'user_id',
        'provider',
        'player_name',
        'transfer_in_id',
        'transfer_out_id',
        'period_started_at',
        'period_ended_at',
        'opening_balance',
        'total_transfer_in',
        'prior_transfer_out',
        'total_bet',
        'total_payout',
        'total_adjustment',
        'expected_transfer_out',
        'actual_transfer_out',
        'difference',
        'record_count',
        'status',
        'attempts',
        'next_retry_at',
        'audited_at',
        'reason',
        'meta',
    ];

    protected function casts(): array
    {
        return [
            'period_started_at' => 'datetime',
            'period_ended_at' => 'datetime',
            'opening_balance' => 'decimal:2',
            'total_transfer_in' => 'decimal:2',
            'prior_transfer_out' => 'decimal:2',
            'total_bet' => 'decimal:2',
            'total_payout' => 'decimal:2',
            'total_adjustment' => 'decimal:2',
            'expected_transfer_out' => 'decimal:2',
            'actual_transfer_out' => 'decimal:2',
            'difference' => 'decimal:2',
            'status' => GameTransferAuditStatus::class,
            'next_retry_at' => 'datetime',
            'audited_at' => 'datetime',
            'meta' => 'array',
        ];
    }

    /** @return BelongsTo<User, $this> */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }

    /** @return BelongsTo<GameTransfer, $this> */
    public function transferIn(): BelongsTo
    {
        return $this->belongsTo(GameTransfer::class, 'transfer_in_id');
    }

    /** @return BelongsTo<GameTransfer, $this> */
    public function transferOut(): BelongsTo
    {
        return $this->belongsTo(GameTransfer::class, 'transfer_out_id');
    }
}
