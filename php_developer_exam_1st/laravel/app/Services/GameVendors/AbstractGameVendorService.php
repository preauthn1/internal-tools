<?php

namespace App\Services\GameVendors;

use App\Services\GameVendors\GameVendorServiceInterface;
use App\Exceptions\UnsupportedVendorFeatureException;
use App\Exceptions\GameVendorException;
use App\Services\GameVendors\AuthStrategies\AuthStrategyInterface;
use App\Services\GameVendors\ResponseInterface;
use Illuminate\Http\Client\PendingRequest;
use Illuminate\Http\Client\Response as HttpResponse;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;

abstract class AbstractGameVendorService implements GameVendorServiceInterface
{
    protected array $config;
    protected string $vendorName;
    protected ?PendingRequest $httpClient = null;
    protected ?AuthStrategyInterface $authStrategy = null;

    abstract protected function authenticate(): array;
    abstract protected function getAuthStrategy(): AuthStrategyInterface;

    abstract public function getGameList(): array;

    abstract public function isUserExist(string $username): bool;

    abstract public function createUser(string $username): bool;

    abstract public function createUserIfNotExist(string $username): bool;

    abstract public function getBalance(string $username): float;

    abstract public function deposit(string $username, float $amount, string $transactionId): array;

    abstract public function withdraw(string $username, float $amount, string $transactionId): array;

    abstract public function getTransactionById(array $filters = []): array;

    abstract public function getGameRecords(array $filters = []): array;

    abstract public function getGameUrl(string $username, string $gameId, array $options = []): string;

    protected function getVendorResponse(): ResponseInterface
    {
        $vendorDir = str_replace(' ', '', ucwords(str_replace(['-', '_'], ' ', $this->vendorName)));
        $vendorClass = "App\\Services\\GameVendors\\{$vendorDir}\\{$vendorDir}Response";

        if (class_exists($vendorClass)) {
            return new $vendorClass();
        }

        throw new \InvalidArgumentException(
            "Response handler not found for vendor '{$this->vendorName}'. Expected {$vendorClass}."
        );
    }

    public function __construct(string $vendorName, array $config)
    {
        $this->vendorName = $vendorName;
        $this->config = $config;
    }

    protected function getHttpClient(): PendingRequest
    {
        if ($this->httpClient === null) {
            $this->httpClient = Http::timeout($this->config['api_timeout'] ?? 30)
                ->baseUrl($this->config['api_url'] ?? '');
        }

        return $this->httpClient;
    }

    protected function getConfig(string $key)
    {
        return $this->config[$key];
    }

    public function getVendorName(): string
    {
        return $this->vendorName;
    }

    protected function logRequest(string $method, string $url, array $data = []): void
    {
        Log::debug("Game Vendor API Request [{$this->vendorName}]", [
            'method' => $method,
            'url' => $url,
            'data' => $this->sanitizeLogData($data),
        ]);
    }

    protected function logResponse(string $method, string $url, $response): void
    {
        Log::debug("Game Vendor API Response [{$this->vendorName}]", [
            'method' => $method,
            'url' => $url,
            'status' => $response->status() ?? 'N/A',
            'body' => $this->sanitizeLogData($response->json() ?? []),
        ]);
    }

    protected function sanitizeLogData(array $data): array
    {
        $sensitiveKeys = ['secret', 'password', 'token', 'api_key', 'merchant_secret'];

        foreach ($data as $key => $value) {
            if (in_array(strtolower($key), $sensitiveKeys)) {
                $data[$key] = '***REDACTED***';
            } elseif (is_array($value)) {
                $data[$key] = $this->sanitizeLogData($value);
            }
        }

        return $data;
    }

    protected function handleApiError($response, string $operation): void
    {
        if ($response->failed()) {
            Log::error("Game Vendor API Error [{$this->vendorName}] - {$operation}", [
                'status' => $response->status(),
                'body' => $response->body(),
            ]);

            throw new \Exception(
                "Failed to {$operation} for vendor {$this->vendorName}: " . $response->body(),
                $response->status()
            );
        }
    }

    /**
     * Get the API endpoint URL for a given operation.
     * @param string $operation
     * @return string|null
     * @throws UnsupportedVendorFeatureException
     */
    protected function getEndpoint(string $operation): ?string
    {
        return null;
    }

    protected function makeGetRequest(string $operation, array $params = []): HttpResponse
    {
        return $this->makeRequest('get', $operation, $params);
    }

    protected function makePostRequest(string $operation, array $payload = []): HttpResponse
    {
        return $this->makeRequest('post', $operation, $payload);
    }

    protected function makeRequest(string $httpMethod, string $operation, array $payload = []): HttpResponse
    {
        if (!in_array(strtolower($httpMethod), ['post', 'get'])) {
            throw new \InvalidArgumentException("Unknown HTTP method given: [{$httpMethod}]");
        }

        $url = $this->getEndpoint($operation);

        if ($url === null) {
            throw new UnsupportedVendorFeatureException($this->vendorName, $operation);
        }

        $authCredentials = $this->authenticate();
        $authStrategy = $this->getAuthStrategy();

        // Use strategy to prepare payload and headers
        $requestPayload = $authStrategy->preparePayload($payload, $authCredentials);
        $headers = $authStrategy->prepareHeaders($authCredentials);

        $this->logRequest($httpMethod, $url, $requestPayload);
        $rawResponse = $this->getHttpClient()
            ->withHeaders($headers)
            ->{strtolower($httpMethod)}($url, $requestPayload);

        $this->logResponse($httpMethod, $url, $rawResponse);
        $this->handleApiError($rawResponse, $operation);

        return $rawResponse;
    }

    public function getMerchantSummaryReport(array $filters = []): array
    {
        $response = $this->makeGetRequest('getMerchantSummaryReport', $filters);
        return $this->processMerchantSummaryReportResponse($response->json() ?? []);
    }

    public function getEventAwards(array $filters = []): array
    {
        $response = $this->makeGetRequest('getEventAwards', $filters);
        return $this->processEventAwardsResponse($response->json() ?? []);
    }

    protected function processMerchantSummaryReportResponse(array $response): array
    {
        return $response['data'] ?? $response['report'] ?? $response;
    }

    protected function processEventAwardsResponse(array $response): array
    {
        return $response['data'] ?? $response['awards'] ?? $response;
    }

}
