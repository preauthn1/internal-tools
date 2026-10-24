<?php

namespace App\Observers;

use App\Enums\WalletTransactionType;
use App\Models\User;

class UserObserver
{
    /**
     * Handle the User "created" event.
     *
     * 每位使用者建立時自動開立平台錢包，給予設定的初始金額，並記入錢包帳本。
     */
    public function created(User $user): void
    {
        $initial = (string) config('games.initial_balance', 10000);

        $user->wallet()->create(['balance' => $initial]);

        $user->walletTransactions()->create([
            'type' => WalletTransactionType::Initial,
            'amount' => $initial,
            'balance_after' => $initial,
        ]);
    }
}
