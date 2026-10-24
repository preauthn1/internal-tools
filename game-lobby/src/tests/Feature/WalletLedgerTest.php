<?php

namespace Tests\Feature;

use App\Models\GameProvider;
use App\Models\User;
use App\Services\Games\PlayGameService;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Http;
use Inertia\Testing\AssertableInertia as Assert;
use Tests\TestCase;

class WalletLedgerTest extends TestCase
{
    use RefreshDatabase;

    private function fakePgs(): void
    {
        Http::fake([
            '*/player/exist*' => Http::response(['code' => 0, 'data' => ['isExist' => false]]),
            '*/player/logout' => Http::response(['code' => 0, 'data' => []]),
            '*/player' => Http::response(['code' => 0, 'data' => []]),
            '*/wallet/transfer-in*' => Http::response(['code' => 0, 'data' => ['transactionNo' => 'in-1', 'beforeAmount' => 0.0, 'afterAmount' => 10000.0]]),
            '*/wallet/amount*' => Http::response(['code' => 0, 'data' => ['amount' => 10000.0]]),
            '*/wallet/transfer-out*' => Http::response(['code' => 0, 'data' => ['transactionNo' => 'out-1', 'beforeAmount' => 10000.0, 'afterAmount' => 0.0]]),
            '*/game/link*' => Http::response(['code' => 0, 'data' => ['link' => 'https://game.test/launch?token=[REDACTED]']]),
        ]);
    }

    public function test_transfer_in_and_out_are_recorded_in_the_ledger()
    {
        $this->fakePgs();
        $user = User::factory()->create();
        $service = app(PlayGameService::class);
        $pgsId = GameProvider::where('key', 'target')->value('id');

        // 進入遊戲：全額轉入 → 帳本記 transfer_in（-10000，結餘 0）
        $service->enter($user, 'target', '1004');

        $this->assertDatabaseHas('wallet_transactions', [
            'user_id' => $user->id,
            'type' => 'transfer_in',
            'game_provider_id' => $pgsId,
            'amount' => '-10000.00',
            'balance_after' => '0.00',
        ]);

        // 離開遊戲結算：全額轉回 → 帳本記 transfer_out（+10000，結餘 10000）
        $service->settle($user);

        $this->assertDatabaseHas('wallet_transactions', [
            'user_id' => $user->id,
            'type' => 'transfer_out',
            'game_provider_id' => $pgsId,
            'amount' => '10000.00',
            'balance_after' => '10000.00',
        ]);

        // 帳本共三筆：initial + transfer_in + transfer_out
        $this->assertSame(3, $user->walletTransactions()->count());
    }

    public function test_wallet_history_page_shows_only_own_transactions()
    {
        $user = User::factory()->create();
        $other = User::factory()->create();

        $this->actingAs($user)
            ->get(route('wallet'))
            ->assertOk()
            ->assertInertia(fn (Assert $page) => $page
                ->component('WalletHistory')
                ->where('balance', '10000.00')
                ->has('transactions.data', 1)
                ->where('transactions.data.0.type', 'initial'),
            );
    }

    public function test_guest_is_redirected_to_login()
    {
        $this->get(route('wallet'))->assertRedirect(route('login'));
    }
}
