<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class GameRecord extends Model
{
    use HasFactory;

    /**
     * The attributes that are mass assignable.
     *
     * @var array<int, string>
     */
    protected $fillable = [
        'name',
        'gameID',
        'beforeAmount',
        'winLossAmount',
        'roundID',
        'betID',
        'roundStartTime',
        'roundEndTime',
        'transactionNo',
        'betAmount',
        'betTime',
        'payoutAmount',
        'taxAmount',
        'validAmount',
    ];
}
