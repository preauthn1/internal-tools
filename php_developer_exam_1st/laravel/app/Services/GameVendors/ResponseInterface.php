<?php

namespace App\Services\GameVendors;

interface ResponseInterface
{
    public function getGameList(): array;

    public function isUserExist(): bool;

    public function createUser(): bool;

    public function getBalance(): float;

    public function deposit(): array;

    public function withdraw(): array;

    public function getTransactionById(): array;

    public function getGameRecords(): array;

    public function getGameUrl(): string;

    public function getMerchantSummaryReport(): array;

    public function getEventAwards(): array;

    public function forceQuitUser(): bool;
}
