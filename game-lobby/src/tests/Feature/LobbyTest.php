<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Http;
use Inertia\Testing\AssertableInertia as Assert;
use Tests\TestCase;

class LobbyTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        Cache::flush();
    }

    private function fakeGames(): void
    {
        // 反映真實 target /games 回應：data 為扁平陣列，圖片欄位為 imageUrl
        Http::fake([
            '*/games*' => Http::response([
                'code' => 0,
                'message' => 'OK',
                'data' => [
                    ['gameID' => 1004, 'gameCode' => '1004', 'gameName' => '金吉報喜', 'imageUrl' => 'https://img.test/1004.png', 'isPrePaidMode' => false],
                    ['gameID' => 1013, 'gameCode' => '1013', 'gameName' => '88Fortune', 'imageUrl' => 'https://img.test/1013.png', 'isPrePaidMode' => false],
                ],
            ]),
        ]);
    }

    public function test_lobby_is_public_and_lists_games_for_guests()
    {
        $this->fakeGames();

        $this->get(route('home'))
            ->assertOk()
            ->assertInertia(fn (Assert $page) => $page
                ->component('Lobby')
                ->where('gamesUnavailable', false)
                ->has('games', 2)
                ->where('games.0.name', '金吉報喜')
                ->where('games.0.provider', 'target')
                ->where('auth.user', null)
                ->where('auth.balance', null),
            );
    }

    public function test_lobby_shows_balance_for_authenticated_user()
    {
        $this->fakeGames();
        $user = User::factory()->create();

        $this->actingAs($user)
            ->get(route('home'))
            ->assertOk()
            ->assertInertia(fn (Assert $page) => $page
                ->component('Lobby')
                ->where('auth.user.name', $user->name)
                ->where('auth.balance', '10000.00'),
            );
    }

    public function test_lobby_degrades_gracefully_when_provider_fails()
    {
        Http::fake([
            '*/games*' => Http::response([
                'code' => 500001,
                'message' => '系統錯誤',
                'data' => null,
            ]),
        ]);

        $this->get(route('home'))
            ->assertOk()
            ->assertInertia(fn (Assert $page) => $page
                ->component('Lobby')
                ->where('gamesUnavailable', true)
                ->has('games', 0),
            );
    }
}
