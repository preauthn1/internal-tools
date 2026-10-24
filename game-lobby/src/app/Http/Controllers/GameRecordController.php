<?php

namespace App\Http\Controllers;

use App\Services\Games\Exceptions\GameProviderException;
use App\Services\Games\GameProviderManager;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Inertia\Inertia;
use Inertia\Response;

class GameRecordController extends Controller
{
    /**
     * 會員的遊戲紀錄頁（僅顯示自己的紀錄）。
     */
    public function index(Request $request, GameProviderManager $manager): Response
    {
        $records = $request->user()
            ->gameRecords()
            ->latest('bet_at')
            ->latest('id')
            ->paginate(20)
            ->withQueryString();

        return Inertia::render('Records', [
            'records' => $records,
            // 遊戲商紀錄多半不帶遊戲名稱，附上 gameId/code → 名稱對照，前端把編號顯示成真正的遊戲名。
            'gameNames' => $this->gameNames($manager, $records->getCollection()->pluck('provider')->unique()->all()),
        ]);
    }

    /**
     * @param  array<int, string>  $providerKeys
     * @return array<string, array<string, string>>
     */
    private function gameNames(GameProviderManager $manager, array $providerKeys): array
    {
        $map = [];
        foreach ($providerKeys as $key) {
            $names = $this->providerGameNames($manager, $key);
            if ($names !== []) {
                $map[$key] = $names;
            }
        }

        return $map;
    }

    /**
     * 取得單一遊戲商的 gameId/code → 名稱對照。
     *
     * 快取「純字串陣列」而非 GameData 物件——物件序列化進 Redis 後反序列化容易踩到
     * incomplete object，純陣列則無此風險，也與大廳的 games 物件快取解耦。
     *
     * @return array<string, string>
     */
    private function providerGameNames(GameProviderManager $manager, string $key): array
    {
        try {
            $provider = $manager->provider($key);
            $ttl = (int) $provider->setting('games_cache_ttl', 300);

            return Cache::remember("games.names.{$key}", $ttl, function () use ($provider): array {
                $names = [];
                foreach ($provider->games() as $game) {
                    if ($game->name === '') {
                        continue;
                    }
                    if ($game->gameId !== '') {
                        $names[$game->gameId] = $game->name;
                    }
                    if ($game->code !== '') {
                        $names[$game->code] = $game->name;
                    }
                }

                return $names;
            });
        } catch (GameProviderException) {
            // 取不到遊戲列表就略過，前端會退回顯示遊戲編號。
            return [];
        }
    }
}
