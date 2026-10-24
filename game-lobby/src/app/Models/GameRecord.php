<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Support\Carbon;

/**
 * @property int $id
 * @property int $user_id
 * @property string $provider
 * @property string $player_name
 * @property string $record_id
 * @property string|null $game_id
 * @property string|null $game_code
 * @property string|null $game_name
 * @property string $bet_amount
 * @property string $payout_amount
 * @property string $valid_amount
 * @property string|null $currency
 * @property Carbon|null $bet_at
 * @property array<string, mixed>|null $raw
 */
class GameRecord extends Model
{
    protected $fillable = [
        'user_id',
        'provider',
        'player_name',
        'record_id',
        'game_id',
        'game_code',
        'game_name',
        'bet_amount',
        'payout_amount',
        'valid_amount',
        'currency',
        'bet_at',
        'raw',
    ];

    protected function casts(): array
    {
        return [
            'bet_amount' => 'decimal:2',
            'payout_amount' => 'decimal:2',
            'valid_amount' => 'decimal:2',
            'bet_at' => 'datetime',
            'raw' => 'array',
        ];
    }

    /**
     * @return BelongsTo<User, $this>
     */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }
}
