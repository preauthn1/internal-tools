<?php

namespace App\Services\Games\Exceptions;

use RuntimeException;

class GameTransferPendingException extends RuntimeException
{
    public function __construct(string $message = '款項正在處理中，請稍後重試。')
    {
        parent::__construct($message);
    }
}
