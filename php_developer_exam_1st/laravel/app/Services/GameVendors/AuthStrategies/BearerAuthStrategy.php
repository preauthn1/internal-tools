<?php

namespace App\Services\GameVendors\AuthStrategies;

class BearerAuthStrategy implements AuthStrategyInterface
{
    public function preparePayload(array $payload, array $authCredentials): array
    {
        return $payload;
    }

    public function prepareHeaders(array $authCredentials): array
    {
        $id = $authCredentials['id'] ?? $authCredentials['merchant_id'] ?? '';
        $secret = [REDACTED]'secret'] ?? $authCredentials['merchant_secret'] ?? '';

        return [
            'Authorization' => 'Bearer ' . base64_encode($id . ':' . $secret),
            'Content-Type' => 'application/json',
            'Accept' => 'application/json',
        ];
    }
}
