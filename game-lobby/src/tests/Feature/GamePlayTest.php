<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Http;
use Tests\TestCase;

class GamePlayTest extends TestCase
{
    use RefreshDatabase;

    private function fakePgs(bool $transferFails = false): void
    {
        Http::fake([
            '*/player/exist*' => Http::response([
                'code' => 0, 'message' => 'OK', 'data' => ['isExist' => false, 'createTime' => null],
            ]),
            '*/player' => Http::response(['code' => 0, 'message' => 'OK', 'data' => []]),
            '*/wallet/transfer-in*' => $transferFails
                ? Http::response(['code' => 500001, 'message' => '系統錯誤', 'data' => null])
                : Http::response([
                    'code' => 0, 'message' => 'OK',
                    'data' => ['transactionNo' => 'tx-1', 'amount' => 10000.0, 'beforeAmount' => 0.0, 'afterAmount' => 10000.0],
                ]),
            '*/game/link*' => Http::response([
                'code' => 0, 'message' => 'OK', 'data' => ['link' => 'https://game.test/launch?token=[REDACTED]'],
            ]),
        ]);
    }

    public function test_entering_a_game_binds_player_transfers_all_in_and_redirects_to_launch_url()
    {
        $this->fakePgs();
        $user = User::factory()->create(['account' => 'leo']);

        $response = $this->actingAs($user)->post('/play/target/1004');

        // 直接 302 導向遊戲商啟動連結（新分頁開啟）
        $response->assertRedirect('https://game.test/launch?token=[REDACTED]');

        // 綁定建立，player_name = {prefix}_{英數亂數}_{user_id}，不含平台 account
        $binding = $user->gameAccounts()->where('provider', 'target')->firstOrFail();
        $this->assertMatchesRegularExpression('/^pgs_[A-Za-z0-9]+_'.$user->id.'$/', $binding->player_name);
        $this->assertStringNotContainsString('leo', $binding->player_name);

        // 全額轉入 → 平台餘額歸零
        $this->assertSame('0.00', $user->wallet()->first()->balance);

        // 轉帳流水完成
        $this->assertDatabaseHas('game_transfers', [
            'user_id' => $user->id,
            'provider' => 'target',
            'direction' => 'in',
            'amount' => '10000.00',
            'status' => 'completed',
        ]);
    }

    public function test_re_entering_with_zero_balance_does_not_transfer_again()
    {
        $this->fakePgs();
        $user = User::factory()->create();
        $user->wallet()->update(['balance' => 0]);

        $response = $this->actingAs($user)->post('/play/target/1004');

        $response->assertRedirect('https://game.test/launch?token=[REDACTED]');
        $this->assertDatabaseCount('game_transfers', 0);
    }

    public function test_entering_while_binding_is_locked_does_not_create_an_external_player(): void
    {
        $this->fakePgs();
        $user = User::factory()->create();
        $lock = Cache::lock("game-binding:target:{$user->id}", 30);
        $this->assertTrue($lock->get());

        try {
            $this->actingAs($user)
                ->post('/play/target/1004')
                ->assertRedirect(route('home'))
                ->assertSessionHasErrors('play');
        } finally {
            $lock->release();
        }

        Http::assertNothingSent();
        $this->assertDatabaseCount('game_provider_accounts', 0);
    }

    public function test_provider_failure_keeps_balance_and_flashes_error()
    {
        $this->fakePgs(transferFails: true);
        $user = User::factory()->create();

        $response = $this->actingAs($user)->post('/play/target/1004');

        $response->assertRedirect(route('home'));
        $response->assertSessionHasErrors('play');

        // 轉入失敗 → 交易回滾，餘額不變、無完成的流水
        $this->assertSame('10000.00', $user->wallet()->first()->balance);
        $this->assertDatabaseMissing('game_transfers', ['status' => 'completed']);
    }

    public function test_guest_is_redirected_to_login()
    {
        $this->post('/play/target/1004')->assertRedirect(route('login'));
    }

    public function test_get_play_route_is_not_available()
    {
        $user = User::factory()->create();

        $this->actingAs($user)->get('/play/target/1004')->assertMethodNotAllowed();
    }
}
