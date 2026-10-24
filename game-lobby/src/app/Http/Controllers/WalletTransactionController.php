<?php

namespace App\Http\Controllers;

use App\Models\WalletTransaction;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class WalletTransactionController extends Controller
{
    /**
     * 錢包紀錄頁（會員檢視自己的錢包帳本）。
     */
    public function index(Request $request): Response
    {
        $transactions = $request->user()
            ->walletTransactions()
            ->with('gameProvider')
            ->latest('id')
            ->paginate(20)
            ->through(fn (WalletTransaction $tx): array => [
                'id' => $tx->id,
                'type' => $tx->type->value,
                'type_label' => $tx->type->label(),
                'provider' => $tx->gameProvider?->key,
                'provider_name' => $tx->gameProvider?->name,
                'amount' => $tx->amount,
                'balance_after' => $tx->balance_after,
                'created_at' => $tx->created_at,
            ])
            ->withQueryString();

        return Inertia::render('WalletHistory', [
            'transactions' => $transactions,
            'balance' => (string) $request->user()->wallet->balance,
        ]);
    }
}
