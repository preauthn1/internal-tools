<?php

namespace App\Http\Controllers;

use App\Services\Games\Contracts\GameProvider;
use App\Services\Games\Exceptions\GameProviderConnectionException;
use App\Services\Games\Exceptions\GameProviderException;
use App\Services\Games\GameProviderManager;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Str;
use Inertia\Inertia;
use Inertia\Response;
use InvalidArgumentException;
use Throwable;

/**
 * 遊戲商 API 測試頁（開發工具）。
 *
 * 直接呼叫 provider driver 的各端點並回傳原始結果 / 錯誤碼，方便互動式驗證遊戲商 API。
 * 注意：此頁直呼遊戲商，不經過 PlayGameService，故 transferIn/Out 只動遊戲商錢包，
 * 不會寫入平台的 wallet_transactions / game_transfers 帳本。僅限非正式環境開放。
 */
class ProviderTesterController extends Controller
{
    /** 可測試的 provider 動作與所需欄位。 */
    private const ACTIONS = [
        'games', 'playerExists', 'createPlayer', 'walletAmount',
        'gameLink', 'transferIn', 'transferOut', 'transferRecord',
        'logoutPlayer', 'gameRecords',
    ];

    public function index(Request $request, GameProviderManager $manager): Response
    {
        abort_unless(app()->environment(['local', 'testing']), 404);

        $providers = array_map(
            fn (GameProvider $p): array => ['key' => $p->key()],
            $manager->enabled(),
        );

        // 預設帶入目前登入者於第一個遊戲商的綁定 playerName，省得手打。
        $defaultPlayerName = $request->user()?->gameAccounts()->value('player_name');

        return Inertia::render('ProviderTester', [
            'providers' => $providers,
            'actions' => self::ACTIONS,
            'defaultPlayerName' => $defaultPlayerName,
            'csrfToken' => csrf_token(),
        ]);
    }

    public function call(Request $request, GameProviderManager $manager): JsonResponse
    {
        abort_unless(app()->environment(['local', 'testing']), 404);

        $validated = $request->validate([
            'provider' => ['required', 'string'],
            'action' => ['required', 'string', 'in:'.implode(',', self::ACTIONS)],
            'playerName' => ['nullable', 'string', 'max:50'],
            'gameId' => ['nullable', 'string', 'max:50'],
            'amount' => ['nullable', 'string', 'max:20'],
            'reference' => ['nullable', 'string', 'max:50'],
            'startTime' => ['nullable', 'string', 'max:40'],
            'endTime' => ['nullable', 'string', 'max:40'],
            'pageIndex' => ['nullable', 'integer', 'min:1'],
            'pageSize' => ['nullable', 'integer', 'min:1', 'max:1000'],
        ]);

        try {
            $provider = $manager->provider($validated['provider']);
        } catch (Throwable) {
            return response()->json(['ok' => false, 'apiCode' => null, 'message' => "未知的遊戲商：{$validated['provider']}"], 422);
        }

        $start = hrtime(true);
        try {
            $data = $this->dispatch($provider, $validated);

            return response()->json([
                'ok' => true,
                'action' => $validated['action'],
                'elapsedMs' => $this->elapsed($start),
                'apiCode' => 0,
                'data' => $data,
            ]);
        } catch (GameProviderConnectionException $e) {
            return $this->fail($validated['action'], $start, null, '連線失敗（逾時 / 無法連到遊戲商）：'.$e->getMessage(), $e->providerMessage);
        } catch (GameProviderException $e) {
            return $this->fail($validated['action'], $start, $e->apiCode, $e->getMessage(), $e->providerMessage);
        } catch (InvalidArgumentException $e) {
            return $this->fail($validated['action'], $start, null, '參數錯誤：'.$e->getMessage());
        } catch (Throwable $e) {
            return $this->fail($validated['action'], $start, null, '未預期錯誤：'.$e->getMessage());
        }
    }

    /**
     * @param  array<string, mixed>  $v
     */
    private function dispatch(GameProvider $provider, array $v): mixed
    {
        $playerName = (string) ($v['playerName'] ?? '');

        return match ($v['action']) {
            'games' => array_map(
                fn ($g): array => ['gameId' => $g->gameId, 'code' => $g->code, 'name' => $g->name],
                $provider->games(),
            ),
            'playerExists' => ['isExist' => $provider->playerExists($playerName)],
            'createPlayer' => $provider->createPlayer($playerName),
            'walletAmount' => ['amount' => $provider->walletAmount($playerName)],
            'gameLink' => ['link' => $provider->gameLink($playerName, (string) ($v['gameId'] ?? ''))],
            'transferIn' => $this->doTransfer($provider, 'in', $playerName, $v),
            'transferOut' => $this->doTransfer($provider, 'out', $playerName, $v),
            'transferRecord' => $provider->transferRecord((string) ($v['reference'] ?? ''))?->raw,
            'logoutPlayer' => $this->doLogout($provider, $playerName),
            'gameRecords' => $this->doGameRecords($provider, $v),
            default => throw new InvalidArgumentException("未知的動作：{$v['action']}"),
        };
    }

    /**
     * @param  array<string, mixed>  $v
     * @return array<string, mixed>
     */
    private function doTransfer(GameProvider $provider, string $direction, string $playerName, array $v): array
    {
        $reference = (string) ($v['reference'] ?? '') !== '' ? (string) $v['reference'] : (string) Str::uuid();
        $amount = (string) ($v['amount'] ?? '');
        $response = $direction === 'in'
            ? $provider->transferIn($playerName, $amount, $reference)
            : $provider->transferOut($playerName, $amount, $reference);

        // 回傳實際使用的 reference，方便之後用 transferRecord 對帳。
        return ['reference' => $reference, 'response' => $response];
    }

    /**
     * @return array{loggedOut: bool}
     */
    private function doLogout(GameProvider $provider, string $playerName): array
    {
        $provider->logoutPlayer($playerName);

        return ['loggedOut' => true];
    }

    /**
     * @param  array<string, mixed>  $v
     * @return array<string, mixed>
     */
    private function doGameRecords(GameProvider $provider, array $v): array
    {
        $query = array_filter([
            'startTime' => $v['startTime'] ?? null,
            'endTime' => $v['endTime'] ?? null,
            'playerName' => ($v['playerName'] ?? '') !== '' ? $v['playerName'] : null,
            'pageIndex' => $v['pageIndex'] ?? null,
            'pageSize' => $v['pageSize'] ?? null,
        ], fn ($value): bool => $value !== null);

        $result = $provider->gameRecords($query);

        return [
            'pageIndex' => $result['pageIndex'],
            'pageSize' => $result['pageSize'],
            'totalCount' => $result['totalCount'],
            'totalPage' => $result['totalPage'],
            'items' => array_map(fn ($item): array => $item->toAttributes(), $result['items']),
        ];
    }

    private function elapsed(float $start): int
    {
        return (int) round((hrtime(true) - $start) / 1e6);
    }

    private function fail(string $action, float $start, ?int $apiCode, string $message, ?string $providerMessage = null): JsonResponse
    {
        return response()->json([
            'ok' => false,
            'action' => $action,
            'elapsedMs' => $this->elapsed($start),
            'apiCode' => $apiCode,
            'message' => $message,
            'providerMessage' => $providerMessage,
        ]);
    }
}
