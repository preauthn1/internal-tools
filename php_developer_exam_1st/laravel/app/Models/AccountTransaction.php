<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class AccountTransaction extends Model
{
    const TYPE_DEPOSIT = 0;
    const TYPE_WITHDRAW = 1;
    const STATUS_PENDING = 0;
    const STATUS_SUCCEED = 1;
    const STATUS_FAILED = 2;

    use HasFactory;

    protected $table = 'account_transaction';

    protected $fillable = [
        'transaction_no',
        'userid',
        'amount',
        'type',
        'status',
        'vendor',
    ];

    protected function casts(): array
    {
        return [
            'userid' => 'integer',
            'amount' => 'decimal:2',
            'type' => 'integer',
            'status' => 'integer',
            'created_at' => 'datetime',
            'updated_at' => 'datetime',
        ];
    }

    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }

    public function isDeposit(): bool
    {
        return $this->type === self::TYPE_DEPOSIT;
    }

    public function isWithdraw(): bool
    {
        return $this->type === self::TYPE_WITHDRAW;
    }

    public function isPending(): bool
    {
        return $this->status === self::STATUS_PENDING;
    }

    public function isSucceed(): bool
    {
        return $this->status === self::STATUS_SUCCEED;
    }

    public function isFailed(): bool
    {
        return $this->status === self::STATUS_FAILED;
    }
}
