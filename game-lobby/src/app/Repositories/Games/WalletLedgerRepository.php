<?php

namespace App\Repositories\Games;

use App\Enums\WalletTransactionType;
use App\Models\GameProvider;
use App\Models\WalletTransaction;

class WalletLedgerRepository
{
    /** @var array<string, int|null> */
    private array $providerIds = [];

    public function record(
        int $userId,
        WalletTransactionType $type,
        ?string $providerKey,
        string $signedAmount,
        string $balanceAfter,
        ?string $reference = null,
    ): void {
        WalletTransaction::create([
            'user_id' => $userId,
            'type' => $type,
            'game_provider_id' => $providerKey !== null ? $this->providerId($providerKey) : null,
            'amount' => $signedAmount,
            'balance_after' => $balanceAfter,
            'reference' => $reference,
        ]);
    }

    private function providerId(string $key): ?int
    {
        return $this->providerIds[$key] ??= GameProvider::query()->where('key', $key)->value('id');
    }
}
