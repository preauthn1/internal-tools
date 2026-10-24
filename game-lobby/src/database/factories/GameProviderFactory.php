<?php

namespace Database\Factories;

use App\Models\GameProvider;
use App\Services\Games\target\PgsProvider;
use Illuminate\Database\Eloquent\Factories\Factory;

/**
 * @extends Factory<GameProvider>
 */
class GameProviderFactory extends Factory
{
    /**
     * Define the model's default state.
     *
     * @return array<string, mixed>
     */
    public function definition(): array
    {
        return [
            'key' => 'target',
            'name' => 'target',
            'driver' => PgsProvider::class,
            'enabled' => true,
            'settings' => [
                'merchant_id' => '175888',
                'merchant_secret' => 'test-secret',
                'base_url' => 'https://vendor-api.example.test/v1',
                'timeout' => 10,
                'locale' => 'zh-cn',
                'player_prefix' => 'target',
                'games_cache_ttl' => 300,
                'invalid_state_codes' => [403601],
                'records_window_minutes' => 60,
                'records_lookback_minutes' => 120,
                'records_page_size' => 200,
            ],
        ];
    }

    public function disabled(): static
    {
        return $this->state(fn (): array => ['enabled' => false]);
    }
}
