<?php

namespace App\Models;

// use Illuminate\Contracts\Auth\MustVerifyEmail;
use App\Observers\UserObserver;
use Database\Factories\UserFactory;
use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Hidden;
use Illuminate\Database\Eloquent\Attributes\ObservedBy;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\Relations\HasOne;
use Illuminate\Foundation\Auth\User as Authenticatable;
use Illuminate\Notifications\Notifiable;
use Illuminate\Support\Carbon;

/**
 * @property int $id
 * @property string $name
 * @property string $account
 * @property string $email
 * @property Carbon|null $email_verified_at
 * @property string $password
 * @property string|null $remember_token
 * @property Carbon|null $created_at
 * @property Carbon|null $updated_at
 * @property-read Wallet|null $wallet
 */
#[ObservedBy(UserObserver::class)]
#[Fillable(['name', 'account', 'email', 'password'])]
#[Hidden(['password', 'remember_token'])]
class User extends Authenticatable
{
    /** @use HasFactory<UserFactory> */
    use HasFactory, Notifiable;

    /**
     * Get the attributes that should be cast.
     *
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'email_verified_at' => 'datetime',
            'password' => 'hashed',
        ];
    }

    /**
     * 平台錢包（1:1）。
     *
     * @return HasOne<Wallet, $this>
     */
    public function wallet(): HasOne
    {
        return $this->hasOne(Wallet::class);
    }

    /**
     * 各遊戲商的玩家帳號綁定（1:N）。
     *
     * @return HasMany<GameProviderAccount, $this>
     */
    public function gameAccounts(): HasMany
    {
        return $this->hasMany(GameProviderAccount::class);
    }

    /**
     * 平台與遊戲商之間的轉帳流水（1:N）。
     *
     * @return HasMany<GameTransfer, $this>
     */
    public function gameTransfers(): HasMany
    {
        return $this->hasMany(GameTransfer::class);
    }

    /**
     * 自遊戲商同步下來的遊戲（下注）紀錄（1:N）。
     *
     * @return HasMany<GameRecord, $this>
     */
    public function gameRecords(): HasMany
    {
        return $this->hasMany(GameRecord::class);
    }

    /**
     * 錢包帳本（每筆餘額異動，1:N）。
     *
     * @return HasMany<WalletTransaction, $this>
     */
    public function walletTransactions(): HasMany
    {
        return $this->hasMany(WalletTransaction::class);
    }

    /** @return HasMany<GameTransferAudit, $this> */
    public function gameTransferAudits(): HasMany
    {
        return $this->hasMany(GameTransferAudit::class);
    }
}
