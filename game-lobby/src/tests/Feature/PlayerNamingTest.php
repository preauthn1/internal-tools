<?php

namespace Tests\Feature;

use App\Models\GameProvider;
use App\Models\User;
use App\Services\Games\target\PgsProvider;
use App\Systems\Games\PlayerBindingSystem;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class PlayerNamingTest extends TestCase
{
    use RefreshDatabase;

    public function test_player_name_is_prefix_random_userid_without_account(): void
    {
        $user = User::factory()->create(['account' => 'leo']); // seeder 已建 target（prefix=target）
        $name = app(PlayerBindingSystem::class)->playerNameFor($user, 'target');

        // {prefix}_{英數亂數}_{user_id}；不得含平台 account
        $this->assertMatchesRegularExpression('/^pgs_[A-Za-z0-9]+_'.$user->id.'$/', $name);
        $this->assertStringNotContainsString('leo', $name);
        $this->assertLessThanOrEqual(50, strlen($name));
    }

    public function test_random_middle_differs_between_generations(): void
    {
        $user = User::factory()->create();
        $service = app(PlayerBindingSystem::class);

        $this->assertNotSame(
            $service->playerNameFor($user, 'target'),
            $service->playerNameFor($user, 'target'),
        );
    }

    public function test_stays_within_provider_limit_with_long_prefix(): void
    {
        $user = User::factory()->create();
        GameProvider::factory()->create([
            'key' => 'longp',
            'driver' => PgsProvider::class,
            'settings' => ['player_prefix' => str_repeat('p', 45)],
        ]);

        $name = app(PlayerBindingSystem::class)->playerNameFor($user, 'longp');

        $this->assertLessThanOrEqual(50, strlen($name));
        $this->assertMatchesRegularExpression('/^[A-Za-z0-9_-]{1,50}$/', $name);
        $this->assertStringEndsWith('_'.$user->id, $name);
    }

    public function test_each_provider_uses_its_own_prefix(): void
    {
        $user = User::factory()->create();
        $service = app(PlayerBindingSystem::class);

        // 第二個遊戲商：自帶不同前綴
        GameProvider::factory()->create([
            'key' => 'acme',
            'driver' => PgsProvider::class,
            'settings' => ['player_prefix' => 'ace'],
        ]);
        // 第三個遊戲商：未設定前綴 → 退回以 provider key 當前綴
        GameProvider::factory()->create([
            'key' => 'zeta',
            'driver' => PgsProvider::class,
            'settings' => [],
        ]);

        $this->assertMatchesRegularExpression('/^pgs_[A-Za-z0-9]+_'.$user->id.'$/', $service->playerNameFor($user, 'target'));
        $this->assertMatchesRegularExpression('/^ace_[A-Za-z0-9]+_'.$user->id.'$/', $service->playerNameFor($user, 'acme'));
        $this->assertMatchesRegularExpression('/^zeta_[A-Za-z0-9]+_'.$user->id.'$/', $service->playerNameFor($user, 'zeta'));
    }
}
