<?php

namespace App\Services\GameVendors\AuthStrategies;

use Illuminate\Support\Facades\Crypt;

class DesAuthStrategy implements AuthStrategyInterface
{
    protected string $encryptedKeyName;

    public function __construct(string $encryptedKeyName = 'encrypted_data')
    {
        $this->encryptedKeyName = $encryptedKeyName;
    }

    public function preparePayload(array $payload, array $authCredentials): array
    {
        $secret = [REDACTED]'secret'] ?? $authCredentials['merchant_secret'] ?? '';

        // Convert payload to JSON string
        $jsonPayload = json_encode($payload);

        // Encrypt using DES (using Laravel's encryption which uses AES-256-CBC by default)
        // For true DES, you would need to use openssl_encrypt with DES-CBC
        // This is a simplified version - adjust based on actual DES requirements
        $encrypted = $this->encryptDes($jsonPayload, $secret);

        return [
            $this->encryptedKeyName => $encrypted,
            'id' => $authCredentials['id'] ?? $authCredentials['merchant_id'] ?? '',
        ];
    }

    public function prepareHeaders(array $authCredentials): array
    {
        return [
            'Content-Type' => 'application/json',
            'Accept' => 'application/json',
        ];
    }

    /**
     * Encrypt data using DES
     */
    protected function encryptDes(string $data, string $key): string
    {
        // Pad key to 8 bytes for DES
        $key = substr(str_pad($key, 8, "\0"), 0, 8);

        // Use OpenSSL for DES encryption
        $iv = openssl_random_pseudo_bytes(8);
        $encrypted = openssl_encrypt($data, 'DES-CBC', $key, OPENSSL_RAW_DATA, $iv);

        // Return base64 encoded IV + encrypted data
        return base64_encode($iv . $encrypted);
    }
}
