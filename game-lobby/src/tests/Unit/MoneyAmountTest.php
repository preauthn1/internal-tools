<?php

namespace Tests\Unit;

use App\Services\Games\Data\MoneyAmount;
use InvalidArgumentException;
use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;

class MoneyAmountTest extends TestCase
{
    #[DataProvider('validAmounts')]
    public function test_accepts_valid_amounts(string $input, string $expected): void
    {
        $this->assertSame($expected, MoneyAmount::from($input)->value());
    }

    public static function validAmounts(): array
    {
        return [['0', '0.00'], ['0.01', '0.01'], ['0.09', '0.09'], ['0.10', '0.10'], ['10.5', '10.50']];
    }

    #[DataProvider('invalidAmounts')]
    public function test_rejects_invalid_amounts(string $input): void
    {
        $this->expectException(InvalidArgumentException::class);
        MoneyAmount::from($input);
    }

    public static function invalidAmounts(): array
    {
        return [['-1'], ['0.001'], ['abc'], ['1e2']];
    }
}
