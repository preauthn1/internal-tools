<?php

namespace App\Models;

use App\Enums\WalletTransactionType;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Support\Carbon;

/**
 * 錢包帳本：記錄每一筆平台錢包餘額異動（初始金 / 轉入 / 轉出 / 退回），含異動後結餘。
 *
 * @property int $id
 * @property int $user_id
 * @property WalletTransactionType $type
 * @property int|null $game_provider_id
 * @property numeric-string $amount
 * @property numeric-string $balance_after
 * @property string|null $reference
 * @property array<string, mixed>|null $meta
 * @property Carbon|null $created_at
 */
class WalletTransaction extends Model
{
    protected $fillable = [
        'user_id',
        'type',
        'game_provider_id',
        'amount',
        'balance_after',
        'reference',
        'meta',
    ];

    protected function casts(): array
    {
        return [
            'type' => WalletTransactionType::class,
            'amount' => 'decimal:2',
            'balance_after' => 'decimal:2',
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

    /**
     * 這筆帳本對應的遊戲商（轉帳類才有；初始金為 null）。
     *
     * @return BelongsTo<GameProvider, $this>
     */
    public function gameProvider(): BelongsTo
    {
        return $this->belongsTo(GameProvider::class);
    }
}
