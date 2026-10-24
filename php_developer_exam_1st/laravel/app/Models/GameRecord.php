<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class GameRecord extends Model
{
    /**
     * The table associated with the model.
     *
     * @var string
     */
    protected $table = 'game_records';

    /**
     * Indicates if the model should be timestamped.
     *
     * @var bool
     */
    public $timestamps = false;

    /**
     * The attributes that are mass assignable.
     *
     * @var array<int, string>
     */
    protected $fillable = [
        'vendor',
        'username',
        'user_id',
        'game_id',
        'bet_id',
        'bet_time',
        'bet_result_amount',
        'bet_amount',
        'payout_amount',
        'valid_amount',
        'tax_amount',
        'extra_info',
    ];

    /**
     * The attributes that should be cast.
     *
     * @var array<string, string>
     */
    protected $casts = [
        'bet_result_amount' => 'float',
        'bet_amount' => 'float',
        'payout_amount' => 'float',
        'valid_amount' => 'float',
        'tax_amount' => 'float',
        'user_id' => 'integer',
    ];

    /**
     * Get the user that owns the game record using Laravel's built-in relationship
     *
     * @return BelongsTo
     */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class, 'user_id');
    }
}
