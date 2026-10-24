<?php

namespace App\Services\GameVendors\AuthStrategies;

class HashAuthStrategy implements AuthStrategyInterface
{
    protected string $hashAlgorithm;
    protected string $hashKeyName;

    public function __construct(string $hashAlgorithm = 'sha256', string $hashKeyName = 'hash')
    {
        $this->hashAlgorithm = $hashAlgorithm;
        $this->hashKeyName = $hashKeyName;
    }

    public function preparePayload(array $payload, array $authCredentials): array
    {
        $secret = [REDACTED]'secret'] ?? $authCredentials['merchant_secret'] ?? '';

        // Sort payload by keys for consistent hashing
        ksort($payload);

        // Create query string from payload
        $queryString = http_build_query($payload);

        // Calculate hash
        $hash = hash_hmac($this->hashAlgorithm, $queryString, $secret);

        // Add hash to payload
        $payload[$this->hashKeyName] = $hash;

        // Add auth credentials (id) to payload if needed
        if (isset($authCredentials['id'])) {
            $payload['id'] = $authCredentials['id'];
        } elseif (isset($authCredentials['merchant_id'])) {
            $payload['merchant_id'] = $authCredentials['merchant_id'];
        }

        return $payload;
    }

    public function prepareHeaders(array $authCredentials): array
    {
        return [
            'Content-Type' => 'application/json',
            'Accept' => 'application/json',
        ];
    }
}
