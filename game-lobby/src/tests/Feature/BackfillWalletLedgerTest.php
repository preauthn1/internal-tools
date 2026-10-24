<?php

namespace Tests\Feature;

use App\Models\GameProvider;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class BackfillWalletLedgerTest extends TestCase
{
    use RefreshDatabase;

    public function test_backfill_reconstructs_ledger_for_legacy_user()
    {
        $user = User::factory()->create();
        // 模擬帳本上線前的舊帳號：清掉自動產生的帳本
        $user->walletTransactions()->delete();

        // 歷史轉帳：in 10000（→0）、out 9994.20（→9994.20）
        $user->gameTransfers()->create([
            'provider' => 'target', 'direction' => 'in', 'amount' => '10000.00',
            'reference' => 'r1', 'status' => 'completed', 'resolved_at' => now()->subMinutes(10),
        ]);
        $user->gameTransfers()->create([
            'provider' => 'target', 'direction' => 'out', 'amount' => '9994.20',
            'reference' => 'r2', 'status' => 'completed', 'resolved_at' => now()->subMinutes(5),
        ]);
        $user->wallet()->update(['balance' => '9994.20']);

        $this->artisan('wallet:backfill-ledger')->assertSuccessful();

        $pgsId = GameProvider::where('key', 'target')->value('id');
        $this->assertSame(3, $user->walletTransactions()->count());
        $this->assertDatabaseHas('wallet_transactions', ['user_id' => $user->id, 'type' => 'initial', 'game_provider_id' => null, 'balance_after' => '10000.00']);
        $this->assertDatabaseHas('wallet_transactions', ['user_id' => $user->id, 'type' => 'transfer_in', 'game_provider_id' => $pgsId, 'amount' => '-10000.00', 'balance_after' => '0.00']);
        $this->assertDatabaseHas('wallet_transactions', ['user_id' => $user->id, 'type' => 'transfer_out', 'game_provider_id' => $pgsId, 'amount' => '9994.20', 'balance_after' => '9994.20']);
    }

    public function test_backfill_skips_users_that_already_have_a_ledger()
    {
        $user = User::factory()->create(); // 已有 initial 帳本一筆

        $this->artisan('wallet:backfill-ledger')->assertSuccessful();

        $this->assertSame(1, $user->walletTransactions()->count());
    }
}
