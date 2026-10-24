<?php

namespace OlympusHub\Services;

use OlympusHub\Consul\ConsulKvClient;
use OlympusHub\Support\DatabaseConfigOverrider;
use OlympusHub\Support\EnvParser;
use Illuminate\Support\Facades\Log;

class ConsulRealmTetherService
{
    /** @var ConsulKvClient */
    private $client;

    /** @var EnvParser */
    private $parser;

    /** @var DatabaseConfigOverrider */
    private $overrider;

    public function __construct(
        ConsulKvClient $client,
        EnvParser $parser,
        DatabaseConfigOverrider $overrider
    ) {
        $this->client   = $client;
        $this->parser   = $parser;
        $this->overrider = $overrider;
    }

    public function activate(string $realmId): void
    {
        if ($realmId === 'localhost') {
            Log::info("⚡ [Olympus Hub] 使用本地神域配置 (Localhost): 直接採用現有 .env 設定。");
            return;
        }

        $key  = 'DB/' . strtoupper($realmId);
        $raw  = $this->client->fetch($key);

        Log::info("⚡ [Olympus Hub] 神聖鏈結已活化 (Realm Tether Activated): [{$realmId}]", [
            'realmId' => $realmId,
            'raw'     => $raw,
        ]);

        $vars = $this->parser->parse($raw);
        $this->overrider->apply($vars);
    }

    public function import(string $realmId, string $content): void
    {
        $key = 'DB/' . strtoupper($realmId);
        $this->client->put($key, $content);
    }
}
