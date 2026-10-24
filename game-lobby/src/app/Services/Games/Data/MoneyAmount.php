<?php

namespace App\Services\Games\Data;

use InvalidArgumentException;

final readonly class MoneyAmount
{
    /** @var numeric-string */
    private string $value;

    /** @param numeric-string $value */
    private function __construct(string $value)
    {
        $this->value = $value;
    }

    public static function from(int|float|string $value): self
    {
        $raw = is_float($value) ? rtrim(rtrim(sprintf('%.2F', $value), '0'), '.') : (string) $value;

        if (! is_numeric($raw) || ! preg_match('/^\d+(?:\.\d{1,2})?$/', $raw)) {
            throw new InvalidArgumentException('金額必須為非負數，且最多兩位小數。');
        }

        /** @var numeric-string $normalized */
        $normalized = bcadd($raw, '0', 2);

        return new self($normalized);
    }

    /** @return numeric-string */
    public function value(): string
    {
        return $this->value;
    }

    public function isZero(): bool
    {
        return bccomp($this->value, '0', 2) === 0;
    }

    public function isBelowProviderMinimum(): bool
    {
        return ! $this->isZero() && bccomp($this->value, '0.10', 2) < 0;
    }

    public function equals(self $other): bool
    {
        return bccomp($this->value, $other->value, 2) === 0;
    }
}
