<?php

namespace App\Services\Games\Exceptions;

use Throwable;

class GameProviderConnectionException extends GameProviderException
{
    public function __construct(string $provider, ?Throwable $previous = null)
    {
        parent::__construct($provider, 0, '遊戲商暫時無法連線。', $previous);
    }
}
