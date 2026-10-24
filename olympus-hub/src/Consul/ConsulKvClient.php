<?php

namespace OlympusHub\Consul;

use Consul\Exception\ClientException;
use Consul\Exception\ServerException;
use Consul\Services\KV;
use OlympusHub\Exceptions\ConsulConnectionException;
use OlympusHub\Exceptions\RealmTetherNotFoundException;

class ConsulKvClient
{
    /** @var KV */
    private $kv;

    public function __construct(KV $kv)
    {
        $this->kv = $kv;
    }

    public function fetch(string $key): string
    {
        try {
            $response = $this->kv->get($key, ['raw' => true]);
            return $response->getBody();
        } catch (ClientException $e) {
            if ($e->getCode() === 404) {
                throw new RealmTetherNotFoundException("Consul key '{$key}' not found");
            }
            throw new ConsulConnectionException($e->getMessage());
        } catch (ServerException $e) {
            throw new ConsulConnectionException("Consul server error: {$e->getMessage()}");
        }
    }

    public function put(string $key, string $value): bool
    {
        try {
            $response = $this->kv->put($key, $value);
            return (string)$response->getBody() === 'true';
        } catch (\Exception $e) {
            throw new ConsulConnectionException("Failed to put key '{$key}': {$e->getMessage()}");
        }
    }
}
