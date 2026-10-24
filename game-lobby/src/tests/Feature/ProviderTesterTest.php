<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Http;
use Inertia\Testing\AssertableInertia;
use Tests\TestCase;

class ProviderTesterTest extends TestCase
{
    use RefreshDatabase;

    public function test_index_renders_for_authenticated_user(): void
    {
        $user = User::factory()->create();

        $this->actingAs($user)
            ->get('/dev/provider-tester')
            ->assertOk()
            ->assertInertia(fn (AssertableInertia $page) => $page
                ->component('ProviderTester')
                ->has('providers')
                ->has('actions')
            );
    }

    public function test_call_returns_provider_result(): void
    {
        Http::fake([
            '*/games*' => Http::response([
                'code' => 0, 'message' => 'OK',
                'data' => ['games' => [['gameID' => '1004', 'gameCode' => 'abc', 'gameName' => '測試遊戲']]],
            ]),
        ]);
        $user = User::factory()->create();

        $this->actingAs($user)
            ->postJson('/dev/provider-tester/call', ['provider' => 'target', 'action' => 'games'])
            ->assertOk()
            ->assertJson([
                'ok' => true,
                'action' => 'games',
                'apiCode' => 0,
                'data' => [['gameId' => '1004', 'code' => 'abc', 'name' => '測試遊戲']],
            ]);
    }

    public function test_call_surfaces_provider_error_code(): void
    {
        Http::fake([
            '*/game/link*' => Http::response(['code' => 403508, 'message' => '狀態不合法', 'data' => null]),
        ]);
        $user = User::factory()->create();

        $this->actingAs($user)
            ->postJson('/dev/provider-tester/call', [
                'provider' => 'target', 'action' => 'gameLink',
                'playerName' => 'pgs_abc_1', 'gameId' => '1004',
            ])
            ->assertOk()
            ->assertJson(['ok' => false, 'apiCode' => 403508, 'providerMessage' => '狀態不合法']);
    }

    public function test_call_surfaces_raw_body_on_gateway_error(): void
    {
        Http::fake([
            '*/games*' => Http::response('error code: 502', 502),
        ]);
        $user = User::factory()->create();

        $this->actingAs($user)
            ->postJson('/dev/provider-tester/call', ['provider' => 'target', 'action' => 'games'])
            ->assertOk()
            ->assertJson([
                'ok' => false,
                'apiCode' => 502,
                'message' => 'target 服務暫時不可用（HTTP 502）。',
                'providerMessage' => 'error code: 502',
            ]);
    }

    public function test_call_rejects_unknown_action(): void
    {
        $user = User::factory()->create();

        $this->actingAs($user)
            ->postJson('/dev/provider-tester/call', ['provider' => 'target', 'action' => 'nope'])
            ->assertStatus(422);
    }

    public function test_tester_is_hidden_outside_local_and_testing(): void
    {
        $this->app['env'] = 'production';
        $user = User::factory()->create();

        $this->actingAs($user)->get('/dev/provider-tester')->assertNotFound();
    }
}
