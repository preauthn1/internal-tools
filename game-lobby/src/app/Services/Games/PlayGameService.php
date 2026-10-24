<?php

namespace App\Services\Games;

use App\Models\User;
use App\Services\Games\Data\SettlementResult;
use App\Services\Games\Exceptions\GameProviderConnectionException;
use App\Services\Games\Exceptions\GameProviderException;
use App\Systems\Games\GameTransferSystem;
use App\Systems\Games\PlayerBindingSystem;
use Illuminate\Support\Facades\Log;

class PlayGameService
{
    public function __construct(
        private GameProviderManager $manager,
        private PlayerBindingSystem $bindings,
        private GameTransferSystem $transfers,
    ) {}

    public function enter(User $user, string $providerKey, string $gameId): string
    {
        $provider = $this->manager->provider($providerKey);
        $binding = $this->bindings->ensure($user, $provider);
        $this->transfers->transferAllIn($user, $provider, $binding->player_name);

        return $provider->gameLink($binding->player_name, $gameId);
    }

    public function settle(User $user): SettlementResult
    {
        /** @var numeric-string $total */
        $total = '0';
        $pending = false;
        $residual = false;

        foreach ($user->gameAccounts()->get() as $binding) {
            $provider = $this->manager->provider($binding->provider);
            try {
                $provider->logoutPlayer($binding->player_name);
            } catch (GameProviderConnectionException) {
                $pending = true;

                continue;
            } catch (GameProviderException $exception) {
                Log::warning('遊戲商登出失敗', [
                    'provider' => $provider->key(),
                    'api_code' => $exception->apiCode,
                    'category' => 'logout',
                ]);
            }

            $result = $this->transfers->transferAllOut($user, $provider, $binding->player_name);
            /** @var numeric-string $total */
            $total = bcadd($total, $result->amount, 2);
            $pending = $pending || $result->pending;
            $residual = $residual || $result->residual;
        }

        return new SettlementResult($total, $pending, $residual);
    }
}
