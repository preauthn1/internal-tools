<?php

namespace App\Services\Games\Exceptions;

use RuntimeException;
use Throwable;

/**
 * 遊戲商 API 錯誤（回應 code!=0、連線失敗、回應格式異常等）。
 */
class GameProviderException extends RuntimeException
{
    public function __construct(
        public readonly string $provider,
        public readonly int $apiCode,
        string $message,
        ?Throwable $previous = null,
        // 遊戲商回傳的原始訊息（JSON envelope 的 message，或非 JSON 時的原始 body），供上層完整呈現。
        public readonly ?string $providerMessage = null,
    ) {
        parent::__construct($message, 0, $previous);
    }
}
