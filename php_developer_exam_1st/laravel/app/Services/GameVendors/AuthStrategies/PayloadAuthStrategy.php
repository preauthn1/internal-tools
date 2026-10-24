<?php

namespace App\Services\GameVendors\AuthStrategies;

class PayloadAuthStrategy implements AuthStrategyInterface
{
    public function preparePayload(array $payload, array $authCredentials): array
    {
        return array_merge($authCredentials, $payload);
    }

    public function prepareHeaders(array $authCredentials): array
    {
        return [
            'Content-Type' => 'application/json',
            'Accept' => 'application/json',
        ];
    }
}
