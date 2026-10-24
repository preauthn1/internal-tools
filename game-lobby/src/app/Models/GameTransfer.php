<?php

namespace App\Models;

use App\Enums\GameTransferDirection;
use App\Enums\GameTransferStatus;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasOne;
use Illuminate\Support\Carbon;

/**
 * @property int $id
 * @property int $user_id
 * @property string $provider
 * @property GameTransferDirection $direction
 * @property numeric-string $amount
 * @property string $reference
 * @property GameTransferStatus $status
 * @property int $attempts
 * @property Carbon|null $last_attempt_at
 * @property Carbon|null $next_retry_at
 * @property Carbon|null $resolved_at
 * @property Carbon|null $created_at
 * @property string|null $error_code
 * @property array<string, mixed>|null $meta
 */
class GameTransfer extends Model
{
    protected $fillable = [
        'user_id',
        'provider',
        'direction',
        'amount',
        'reference',
        'status',
        'attempts',
        'last_attempt_at',
        'next_retry_at',
        'resolved_at',
        'error_code',
        'provider_transaction_no',
        'meta',
    ];

    protected function casts(): array
    {
        return [
            'amount' => 'decimal:2',
            'direction' => GameTransferDirection::class,
            'status' => GameTransferStatus::class,
            'last_attempt_at' => 'datetime',
            'next_retry_at' => 'datetime',
            'resolved_at' => 'datetime',
            'meta' => 'array',
        ];
    }

    /**
     * @return BelongsTo<User, $this>
     */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }

    /** @return HasOne<GameTransferAudit, $this> */
    public function audit(): HasOne
    {
        return $this->hasOne(GameTransferAudit::class, 'transfer_out_id');
    }
}
