<?php

namespace App\Services\GameVendors\AuthStrategies;

interface AuthStrategyInterface
{
    /**
     * Prepare request payload with authentication
     */
    public function preparePayload(array $payload, array $authCredentials): array;

    /**
     * Prepare request headers with authentication
     */
    public function prepareHeaders(array $authCredentials): array;
}
