<?php

namespace App\Services\Games;

use App\Models\GameProvider as GameProviderModel;
use App\Services\Games\Contracts\GameProvider;
use InvalidArgumentException;

/**
 * 依 game_providers 資料表解析並快取各遊戲商 driver 實例。
 * 新增遊戲商 = 在 game_providers 新增一列（driver implements GameProvider + settings），無需改動上層。
 */
class GameProviderManager
{
    /** @var array<string, GameProvider> */
    private array $instances = [];

    public function provider(?string $key = null): GameProvider
    {
        $key ??= (string) config('games.default');

        if (isset($this->instances[$key])) {
            return $this->instances[$key];
        }

        $model = GameProviderModel::query()->where('key', $key)->first();

        if ($model === null) {
            throw new InvalidArgumentException("未定義的遊戲商：{$key}");
        }

        return $this->instances[$key] = $this->make($model);
    }

    /**
     * 所有已啟用的遊戲商。
     *
     * @return array<int, GameProvider>
     */
    public function enabled(): array
    {
        return GameProviderModel::query()
            ->where('enabled', true)
            ->orderBy('key')
            ->get()
            ->map(fn (GameProviderModel $model): GameProvider => $this->resolve($model))
            ->all();
    }

    public function default(): GameProvider
    {
        return $this->provider();
    }

    /**
     * 已解析者取快取，否則建立並快取。
     */
    private function resolve(GameProviderModel $model): GameProvider
    {
        return $this->instances[$model->key] ??= $this->make($model);
    }

    private function make(GameProviderModel $model): GameProvider
    {
        $driver = $model->driver;

        if (! is_a($driver, GameProvider::class, true)) {
            throw new InvalidArgumentException("遊戲商 {$model->key} 的 driver 必須實作 GameProvider");
        }

        /** @var class-string<GameProvider> $driver */
        return new $driver($model->settings);
    }
}
