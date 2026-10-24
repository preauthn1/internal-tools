<?php

namespace App\Services\GameVendors;

interface GameVendorServiceInterface
{
    public function getGameList(): array;

    public function isUserExist(string $username): bool;

    public function createUser(string $username): bool;

    public function createUserIfNotExist(string $username): bool;

    public function getBalance(string $username): float;

    public function deposit(string $username, float $amount, string $transactionId): array;

    public function withdraw(string $username, float $amount, string $transactionId): array;

    public function getTransactionById(array $filters = []): array;

    public function getGameRecords(array $filters = []): array;

    public function getGameUrl(string $username, string $gameId, array $options = []): string;

    public function getVendorName(): string;
}
