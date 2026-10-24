<?php

namespace App\Http\Controllers;

use App\Services\Games\Exceptions\GameProviderException;
use App\Services\Games\GameProviderManager;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Log;
use Inertia\Inertia;
use Inertia\Response;

class LobbyController extends Controller
{
    /**
     * 遊戲大廳（公開首頁）。遊戲列表逐遊戲商取得並快取；
     * 任一來源失敗僅略過該來源並標記 gamesUnavailable，不影響整頁。
     */
    public function index(GameProviderManager $manager): Response
    {
        $games = [];
        $gamesUnavailable = false;

        foreach ($manager->enabled() as $provider) {
            try {
                $ttl = (int) $provider->setting('games_cache_ttl', 300);

                $providerGames = Cache::remember(
                    "games.{$provider->key()}",
                    $ttl,
                    fn (): array => $provider->games(),
                );

                $games = [...$games, ...$providerGames];
            } catch (GameProviderException $e) {
                $gamesUnavailable = true;
                Log::warning("取得遊戲列表失敗 [{$e->provider}] code={$e->apiCode}: {$e->getMessage()}");
            }
        }

        return Inertia::render('Lobby', [
            'games' => $games,
            'gamesUnavailable' => $gamesUnavailable,
            'csrfToken' => csrf_token(),
        ]);
    }
}
