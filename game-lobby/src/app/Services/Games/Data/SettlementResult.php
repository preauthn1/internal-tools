<?php

namespace App\Services\Games\Data;

final readonly class SettlementResult
{
    /** @param numeric-string $amount */
    public function __construct(
        public string $amount,
        public bool $pending = false,
        public bool $residual = false,
    ) {}
}
