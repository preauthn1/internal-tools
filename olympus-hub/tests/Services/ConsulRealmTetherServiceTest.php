<?php

namespace OlympusHub\Tests\Services;

use OlympusHub\Consul\ConsulKvClient;
use OlympusHub\Exceptions\RealmTetherNotFoundException;
use OlympusHub\Services\ConsulRealmTetherService;
use OlympusHub\Support\DatabaseConfigOverrider;
use OlympusHub\Support\EnvParser;
use PHPUnit\Framework\TestCase;
use Illuminate\Support\Facades\Log;
use Mockery;

/**
 * ⚡ 聖域考驗：位面鏈結服務測試 (Divine Trial: Realm Tether Service Test)
 * 驗證位面 (Realm) 的活化儀式是否能正確導引能量 (Config)。
 */
class ConsulRealmTetherServiceTest extends TestCase
{
    /**
     * 結束儀式，清理 Mock 環境。
     */
    protected function tearDown(): void
    {
        Mockery::close();
    }

    /** @var string 模擬的原始能量數據 (Raw Env Data) */
    private string $rawValue = "DB_CONNECTION=mysql\nDB_HOST_READ=0.0.0.0\nDB_HOST_WRITE=0.0.0.0\nDB_PORT=3306\nDB_DATABASE=xsgcasino\nDB_USERNAME_READ=reader\nDB_PASSWORD_READ=rpass\nDB_USERNAME_WRITE=writer\nDB_PASSWORD_WRITE=wpass";

    /** @var array 解析後的能量矩陣 (Parsed Config Array) */
    private array $parsedVars = [
        'DB_CONNECTION'     => 'mysql',
        'DB_HOST_READ'      => '0.0.0.0',
        'DB_HOST_WRITE'     => '0.0.0.0',
        'DB_PORT'           => '3306',
        'DB_DATABASE'       => 'xsgcasino',
        'DB_USERNAME_READ'  => 'reader',
        'DB_PASSWORD_READ'  => 'rpass',
        'DB_USERNAME_WRITE' => 'writer',
        'DB_PASSWORD_WRITE' => 'wpass',
    ];

    /**
     * 測試：完整的活化儀式 (Divine Success: Full Activation)
     * 驗證服務是否能正確：獲取 -> 解析 -> 覆寫位面配置。
     */
    public function test_activate_fetches_parses_and_applies_config(): void
    {
        Log::shouldReceive('info')
            ->once()
            ->withArgs(fn($msg, $context) => str_contains($msg, 'tpower') && $context['raw'] === $this->rawValue);

        $client   = $this->createMock(ConsulKvClient::class);
        $parser   = $this->createMock(EnvParser::class);
        $overrider = $this->createMock(DatabaseConfigOverrider::class);

        $client->expects($this->once())->method('fetch')->with('DB/TPOWER')->willReturn($this->rawValue);
        $parser->expects($this->once())->method('parse')->with($this->rawValue)->willReturn($this->parsedVars);
        $overrider->expects($this->once())->method('apply')->with($this->parsedVars);

        (new ConsulRealmTetherService($client, $parser, $overrider))->activate('tpower');
    }

    /**
     * 測試：位面標識符自動轉化為大寫 (Divine Precision: Uppercase Realm ID)
     * 確保神諭請求時，位面 ID 會被正確轉化為大寫以符合觀測位面的標準。
     */
    public function test_activate_uppercases_realm_id_for_consul_key(): void
    {
        Log::shouldReceive('info')->once();

        $client   = $this->createMock(ConsulKvClient::class);
        $parser   = $this->createMock(EnvParser::class);
        $overrider = $this->createMock(DatabaseConfigOverrider::class);

        $client->expects($this->once())->method('fetch')->with('DB/ALFAWIN77')->willReturn($this->rawValue);
        $parser->method('parse')->willReturn($this->parsedVars);
        $overrider->method('apply');

        (new ConsulRealmTetherService($client, $parser, $overrider))->activate('alfawin77');
    }

    /**
     * 測試：本地位面跳過獲取邏輯 (Divine Sanctuary: Skip Fetching for Localhost)
     * 驗證當指定位面為 localhost 時，聖域會直接保留凡間原本的配置。
     */
    public function test_activate_skips_fetching_for_localhost(): void
    {
        Log::shouldReceive('info')
            ->once()
            ->with("⚡ [Olympus Hub] 使用本地神域配置 (Localhost): 直接採用現有 .env 設定。");

        $client   = $this->createMock(ConsulKvClient::class);
        $parser   = $this->createMock(EnvParser::class);
        $overrider = $this->createMock(DatabaseConfigOverrider::class);

        $client->expects($this->never())->method('fetch');
        $parser->expects($this->never())->method('parse');
        $overrider->expects($this->never())->method('apply');

        (new ConsulRealmTetherService($client, $parser, $overrider))->activate('localhost');
    }

    /**
     * 測試：傳遞位面不存在的異常 (Calamity Propagation: Realm Not Found)
     * 確保底層客戶端的災厄能被正確傳遞至上層。
     */
    public function test_activate_propagates_realm_tether_not_found_exception(): void
    {
        $this->expectException(RealmTetherNotFoundException::class);

        $client   = $this->createMock(ConsulKvClient::class);
        $parser   = $this->createMock(EnvParser::class);
        $overrider = $this->createMock(DatabaseConfigOverrider::class);

        $client->method('fetch')->willThrowException(new RealmTetherNotFoundException('not found'));

        (new ConsulRealmTetherService($client, $parser, $overrider))->activate('tpower');
    }
}
