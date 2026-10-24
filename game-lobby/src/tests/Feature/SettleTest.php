<?php

namespace Tests\Feature;

use App\Models\User;
use App\Services\Games\GameTransferReconcileService;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Http;
use Tests\TestCase;

class SettleTest extends TestCase
{
    use RefreshDatabase;

    private function bind(User $user): void
    {
        $user->gameAccounts()->create([
            'provider' => 'target',
            'player_name' => 'gff'.$user->id,
        ]);
    }

    private function fakeWallet(float $amount, bool $transferFails = false): void
    {
        Http::fake([
            '*/player/logout' => Http::response(['code' => 0, 'message' => 'OK', 'data' => []]),
            '*/wallet/amount*' => Http::response([
                'code' => 0, 'message' => 'OK',
                'data' => ['currency' => 'BETA', 'amount' => $amount],
            ]),
            '*/wallet/transfer-out*' => $transferFails
                ? Http::response(['code' => 500001, 'message' => '系統錯誤', 'data' => null])
                : Http::response([
                    'code' => 0, 'message' => 'OK',
                    'data' => ['transactionNo' => 'tx-out', 'amount' => $amount, 'beforeAmount' => $amount, 'afterAmount' => 0.0],
                ]),
        ]);
    }

    public function test_settle_pulls_provider_balance_back_to_platform()
    {
        $user = User::factory()->create();
        $this->bind($user);
        $user->wallet()->update(['balance' => 0]); // 進遊戲後平台餘額為 0
        $this->fakeWallet(2500.0);

        $this->actingAs($user)->post(route('settle'))->assertRedirect();

        // 遊戲商餘額轉回平台
        $this->assertSame('2500.00', $user->wallet()->first()->balance);
        $this->assertDatabaseHas('game_transfers', [
            'user_id' => $user->id,
            'provider' => 'target',
            'direction' => 'out',
            'amount' => '2500.00',
            'status' => 'completed',
        ]);
    }

    public function test_settle_does_nothing_when_provider_balance_is_zero()
    {
        $user = User::factory()->create();
        $this->bind($user);
        $user->wallet()->update(['balance' => 0]);
        $this->fakeWallet(0.0);

        $this->actingAs($user)->post(route('settle'))->assertRedirect();

        $this->assertSame('0.00', $user->wallet()->first()->balance);
        $this->assertDatabaseCount('game_transfers', 0);
    }

    public function test_settle_keeps_balance_when_provider_fails()
    {
        $user = User::factory()->create();
        $this->bind($user);
        $user->wallet()->update(['balance' => 0]);
        $this->fakeWallet(2500.0, transferFails: true);

        $this->actingAs($user)->post(route('settle'))->assertRedirect();

        // 轉出失敗 → 交易回滾，餘額不變、無完成的 out 流水
        $this->assertSame('0.00', $user->wallet()->first()->balance);
        $this->assertDatabaseMissing('game_transfers', ['status' => 'completed']);
    }

    public function test_guest_is_redirected_to_login()
    {
        $this->post(route('settle'))->assertRedirect(route('login'));
    }

    public function test_player_still_online_is_kept_pending_and_later_reconciled()
    {
        $user = User::factory()->create();
        $this->bind($user);
        $user->wallet()->update(['balance' => 0]);
        Http::fake([
            '*/player/logout' => Http::response(['code' => 0, 'data' => []]),
            '*/wallet/amount*' => Http::response(['code' => 0, 'data' => ['amount' => '25.00']]),
            '*/wallet/transfer-out*' => Http::sequence()
                ->push(['code' => 403601, 'message' => 'player active'])
                ->push(['code' => 0, 'data' => ['transactionNo' => 'tx-retry']]),
            '*/transfer*' => Http::response(['code' => 0, 'data' => null]),
        ]);

        $this->actingAs($user)->post(route('settle'))->assertSessionHasErrors('settlement');
        $transfer = $user->gameTransfers()->firstOrFail();
        $this->assertSame('pending', $transfer->status->value);
        $this->assertSame('403601', $transfer->error_code);
        $this->assertSame('0.00', $user->wallet()->firstOrFail()->balance);

        app(GameTransferReconcileService::class)->reconcile($transfer);

        $this->assertSame('completed', $transfer->fresh()->status->value);
        $this->assertSame('25.00', $user->wallet()->firstOrFail()->balance);
    }
}
