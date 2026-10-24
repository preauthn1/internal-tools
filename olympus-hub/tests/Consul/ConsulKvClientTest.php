<?php

namespace OlympusHub\Tests\Consul;

use Consul\ClientInterface;
use Consul\ConsulResponse;
use Consul\Exception\ClientException;
use Consul\Exception\ServerException;
use Consul\Services\KV;
use OlympusHub\Consul\ConsulKvClient;
use OlympusHub\Exceptions\ConsulConnectionException;
use OlympusHub\Exceptions\RealmTetherNotFoundException;
use PHPUnit\Framework\TestCase;

/**
 * ⚡ 聖域考驗：Consul KV 客戶端測試 (Divine Trial: Consul KV Client Test)
 * 驗證與奧林帕斯觀測位面 (Consul) 的通訊是否穩固。
 */
class ConsulKvClientTest extends TestCase
{
    /**
     * 測試：成功獲取原始配置內容 (Divine Success: Fetch Raw Body)
     * 確保當神諭請求正確時，能完整回傳位面的能量數據 (Raw Data)。
     */
    public function test_fetch_returns_raw_body_on_success(): void
    {
        $expected = "DB_HOST_READ=0.0.0.0\nDB_PORT=3306";

        $client = $this->createMock(ClientInterface::class);
        $client->expects($this->once())
            ->method('get')
            ->with('v1/kv/DB/TPOWER', ['query' => ['raw' => true]])
            ->willReturn(new ConsulResponse([], $expected, 200));

        $kv = new KV($client);
        $this->assertSame($expected, (new ConsulKvClient($kv))->fetch('DB/TPOWER'));
    }

    /**
     * 測試：當位面不存在時拋出異常 (Calamity: Realm Not Found 404)
     * 驗證當尋求不存在的位面鏈結時，聖域會降下 RealmTetherNotFoundException。
     */
    public function test_fetch_throws_realm_tether_not_found_on_404(): void
    {
        $this->expectException(RealmTetherNotFoundException::class);

        $client = $this->createMock(ClientInterface::class);
        $client->method('get')->willThrowException(
            new ClientException('Not Found', 404)
        );

        $kv = new KV($client);
        (new ConsulKvClient($kv))->fetch('DB/UNKNOWN');
    }

    /**
     * 測試：當 Consul 伺服器異常時拋出異常 (Calamity: Server Error 500)
     * 驗證當觀測位面崩塌時，聖域會正確轉化為 ConsulConnectionException。
     */
    public function test_fetch_throws_consul_connection_exception_on_server_error(): void
    {
        $this->expectException(ConsulConnectionException::class);

        $client = $this->createMock(ClientInterface::class);
        $client->method('get')->willThrowException(
            new ServerException('Internal Server Error', 500)
        );

        $kv = new KV($client);
        (new ConsulKvClient($kv))->fetch('DB/TPOWER');
    }

    /**
     * 測試：當客戶端權限不足或其他錯誤時拋出異常 (Calamity: Forbidden/Other 403)
     * 驗證非 404 的客戶端錯誤亦會被視為連線災厄。
     */
    public function test_fetch_throws_consul_connection_exception_on_non_404_client_error(): void
    {
        $this->expectException(ConsulConnectionException::class);

        $client = $this->createMock(ClientInterface::class);
        $client->method('get')->willThrowException(
            new ClientException('Forbidden', 403)
        );

        $kv = new KV($client);
        (new ConsulKvClient($kv))->fetch('DB/TPOWER');
    }
}
