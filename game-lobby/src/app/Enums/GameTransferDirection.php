<?php

namespace App\Enums;

enum GameTransferDirection: string
{
    case In = 'in';
    case Out = 'out';

    public function providerType(): int
    {
        return match ($this) {
            self::In => 1,
            self::Out => 2,
        };
    }
}
