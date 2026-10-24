<?php

namespace App\Systems\Games;

use App\Models\GameProviderAccount;
use App\Models\User;
use App\Services\Games\Contracts\GameProvider;
use App\Services\Games\Exceptions\GameTransferPendingException;
use App\Services\Games\GameProviderManager;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Str;

class PlayerBindingSystem
{
    public function __construct(private GameProviderManager $manager) {}

    public function playerNameFor(User $user, string $providerKey): string
    {
        $prefix = (string) $this->manager->provider($providerKey)
            ->setting('player_prefix', $providerKey);
        $suffix = '_'.$user->id;
        $randomLength = max(1, min(8, 50 - strlen($prefix) - 1 - strlen($suffix)));

        return $prefix.'_'.Str::random($randomLength).$suffix;
    }

    public function ensure(User $user, GameProvider $provider): GameProviderAccount
    {
        $binding = $user->gameAccounts()->where('provider', $provider->key())->first();
        if ($binding !== null) {
            return $binding;
        }

        $lock = Cache::lock("game-binding:{$provider->key()}:{$user->id}", 30);
        if (! $lock->get()) {
            throw new GameTransferPendingException('遊戲帳號綁定處理中，請稍後再試。');
        }

        try {
            $binding = $user->gameAccounts()->where('provider', $provider->key())->first();
            if ($binding !== null) {
                return $binding;
            }

            $playerName = $this->playerNameFor($user, $provider->key());
            $meta = $provider->playerExists($playerName) ? null : $provider->createPlayer($playerName);

            return $user->gameAccounts()->firstOrCreate(
                ['provider' => $provider->key()],
                ['player_name' => $playerName, 'meta' => $meta],
            );
        } finally {
            $lock->release();
        }
    }
}
