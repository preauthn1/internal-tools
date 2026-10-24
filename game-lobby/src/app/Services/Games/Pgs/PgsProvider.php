<?php

namespace App\Services\Games\target;

use App\Services\Games\Contracts\GameProvider;
use App\Services\Games\Data\GameData;
use App\Services\Games\Data\GameRecordData;
use App\Services\Games\Data\MoneyAmount;
use App\Services\Games\Data\TransferRecordData;
use App\Services\Games\Exceptions\GameProviderConnectionException;
use App\Services\Games\Exceptions\GameProviderException;
use Illuminate\Http\Client\ConnectionException;
use Illuminate\Support\Facades\Http;
use InvalidArgumentException;

class PgsProvider implements GameProvider
{
    /** @param array<string, mixed> $config */
    public function __construct(private array $config) {}

    public function key(): string
    {
        return 'target';
    }

    public function setting(string $key, mixed $default = null): mixed
    {
        return $this->config[$key] ?? $default;
    }

    public function games(): array
    {
        $data = $this->request('GET', '/games') ?? [];
        $list = $data['games'] ?? $data['list'] ?? ($this->isList($data) ? $data : []);

        return array_map(
            fn (array $game): GameData => GameData::fromPgs($this->key(), $game),
            is_array($list) ? $list : [],
        );
    }

    public function playerExists(string $playerName): bool
    {
        $this->assertPlayerName($playerName);
        $data = $this->request('GET', '/player/exist', ['playerName' => $playerName]);

        return (bool) ($data['isExist'] ?? false);
    }

    public function createPlayer(string $playerName): array
    {
        $this->assertPlayerName($playerName);

        return $this->request('POST', '/player', ['playerName' => $playerName]) ?? [];
    }

    public function walletAmount(string $playerName): string
    {
        $this->assertPlayerName($playerName);
        $data = $this->request('GET', '/wallet/amount', ['playerName' => $playerName]);

        return MoneyAmount::from($data['amount'] ?? '0')->value();
    }

    public function transferIn(string $playerName, string $amount, string $reference): array
    {
        return $this->transfer('/wallet/transfer-in', $playerName, $amount, $reference);
    }

    public function transferOut(string $playerName, string $amount, string $reference): array
    {
        return $this->transfer('/wallet/transfer-out', $playerName, $amount, $reference);
    }

    public function transferRecord(string $merchantTransactionNo): ?TransferRecordData
    {
        $this->assertReference($merchantTransactionNo);
        $data = $this->request('GET', '/transfer', ['merchantTransactionNo' => $merchantTransactionNo]);

        return $data === null || $data === [] ? null : TransferRecordData::fromPgs($data);
    }

    public function logoutPlayer(string $playerName): void
    {
        $this->assertPlayerName($playerName);
        $this->request('POST', '/player/logout', ['playerName' => $playerName]);
    }

    public function gameLink(string $playerName, string $gameId): string
    {
        $this->assertPlayerName($playerName);
        $data = $this->request('GET', '/game/link', ['playerName' => $playerName, 'gameID' => $gameId]);
        $link = (string) ($data['link'] ?? '');
        $parts = parse_url($link);
        $host = (string) ($parts['host'] ?? '');
        $scheme = strtolower((string) ($parts['scheme'] ?? ''));

        // 連結由已認證的 target API 回傳，host 為動態的第三方遊戲網域，不做白名單；
        // 僅要求為合法 https 連結。
        if ($host === '' || $scheme !== 'https') {
            throw new GameProviderException($this->key(), 0, 'target 啟動網址不合法。');
        }

        return $link;
    }

    public function gameRecords(array $query = []): array
    {
        $data = $this->request('GET', '/game-records', $query) ?? [];
        $rawItems = is_array($data['items'] ?? null) ? $data['items'] : [];
        $items = array_map(fn (array $item): GameRecordData => GameRecordData::fromPgs($this->key(), $item), $rawItems);

        return [
            'items' => $items,
            'pageIndex' => (int) ($data['pageIndex'] ?? 1),
            'pageSize' => (int) ($data['pageSize'] ?? 0),
            'totalCount' => (int) ($data['totalCount'] ?? 0),
            'totalPage' => (int) ($data['totalPage'] ?? 0),
        ];
    }

    /** @return array<string, mixed> */
    private function transfer(string $path, string $playerName, string $amount, string $reference): array
    {
        $this->assertPlayerName($playerName);
        $this->assertReference($reference);
        $money = MoneyAmount::from($amount);

        if ($money->isZero() || $money->isBelowProviderMinimum()) {
            throw new InvalidArgumentException('target 轉帳金額最低為 0.10。');
        }

        return $this->request('POST', $path, [
            'playerName' => $playerName,
            'amount' => $money->value(),
            'merchantOrderNo' => $reference,
        ]) ?? [];
    }

    private function assertReference(string $reference): void
    {
        if (! preg_match('/^[A-Za-z0-9_-]{1,50}$/', $reference)) {
            throw new InvalidArgumentException('target reference 格式不合法。');
        }
    }

    private function assertPlayerName(string $playerName): void
    {
        if (! preg_match('/^[A-Za-z0-9_-]{1,50}$/', $playerName)) {
            throw new InvalidArgumentException('target playerName 格式不合法。');
        }
    }

    private function isList(mixed $value): bool
    {
        return is_array($value) && array_is_list($value);
    }

    /**
     * @param  array<string, mixed>  $params
     * @return array<string, mixed>|null
     */
    private function request(string $method, string $path, array $params = []): ?array
    {
        $credentials = [
            'merchantID' => $this->config['merchant_id'],
            'merchantSecretKey' => $this->config['merchant_secret'],
        ];
        $url = rtrim((string) $this->config['base_url'], '/').$path;

        try {
            $request = Http::acceptJson()->timeout((int) ($this->config['timeout'] ?? 10))
                ->withHeaders(['locale' => $this->config['locale'] ?? 'zh-cn']);
            $response = strtoupper($method) === 'GET'
                ? $request->get($url, [...$credentials, ...$params])
                : $request->post($url, [...$credentials, ...$params]);
        } catch (ConnectionException $exception) {
            throw new GameProviderConnectionException($this->key(), $exception);
        }

        $json = $response->json();
        if (! is_array($json) || ! array_key_exists('code', $json)) {
            // 回應非預期 JSON envelope：多半是遊戲商閘道 5xx（服務暫時不可用）或 4xx，
            // 而非我們解析錯誤；依 HTTP 狀態給出可辨識的訊息，apiCode 帶入 HTTP 狀態碼。
            $status = $response->status();
            $message = match (true) {
                $status >= 500 => "target 服務暫時不可用（HTTP {$status}）。",
                $status >= 400 => "target 拒絕請求（HTTP {$status}）。",
                default => 'target 回應格式錯誤。',
            };
            $body = trim($response->body());

            throw new GameProviderException($this->key(), $status, $message, null, $body !== '' ? $body : null);
        }
        if ((int) $json['code'] !== 0) {
            $providerMessage = isset($json['message']) ? (string) $json['message'] : null;

            throw new GameProviderException($this->key(), (int) $json['code'], 'target 拒絕此項操作。', null, $providerMessage);
        }

        $data = $json['data'] ?? null;

        return is_array($data) ? $data : null;
    }
}
