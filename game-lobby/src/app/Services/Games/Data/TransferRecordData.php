<?php

namespace App\Services\Games\Data;

use App\Enums\GameTransferDirection;

final readonly class TransferRecordData
{
    /** @param array<string, mixed> $raw */
    public function __construct(
        public string $providerTransactionNo,
        public string $merchantTransactionNo,
        public int $transferType,
        public string $playerName,
        public MoneyAmount $amount,
        public array $raw = [],
    ) {}

    /** @param array<string, mixed> $raw */
    public static function fromPgs(array $raw): self
    {
        return new self(
            providerTransactionNo: (string) ($raw['pgsTransactionNo'] ?? $raw['transactionNo'] ?? ''),
            merchantTransactionNo: (string) ($raw['merchantTransactionNo'] ?? $raw['merchantOrderNo'] ?? ''),
            transferType: (int) ($raw['transferType'] ?? 0),
            playerName: (string) ($raw['playerName'] ?? ''),
            amount: MoneyAmount::from($raw['amount'] ?? '0'),
            raw: $raw,
        );
    }

    public function matches(
        string $reference,
        string $playerName,
        GameTransferDirection $direction,
        MoneyAmount $amount,
    ): bool {
        return hash_equals($reference, $this->merchantTransactionNo)
            && hash_equals($playerName, $this->playerName)
            && $this->transferType === $direction->providerType()
            && $this->amount->equals($amount);
    }
}
