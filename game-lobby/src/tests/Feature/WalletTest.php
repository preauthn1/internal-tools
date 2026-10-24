<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class WalletTest extends TestCase
{
    use RefreshDatabase;

    public function test_new_user_gets_a_wallet_with_initial_balance()
    {
        $user = User::factory()->create();

        $this->assertNotNull($user->wallet);
        $this->assertSame(
            (string) config('games.initial_balance'),
            (string) (int) $user->wallet->balance,
        );
    }

    public function test_new_user_gets_an_initial_wallet_ledger_entry()
    {
        $user = User::factory()->create();

        $this->assertDatabaseHas('wallet_transactions', [
            'user_id' => $user->id,
            'type' => 'initial',
            'amount' => '10000.00',
            'balance_after' => '10000.00',
        ]);
    }

    public function test_registration_grants_initial_balance_of_10000()
    {
        $this->post(route('register.store'), [
            'email' => 'user@example.invalid',
            'account' => 'player001',
            'password' => 'password',
            'password_confirmation' => 'password',
        ]);

        $user = User::where('account', 'player001')->firstOrFail();

        $this->assertSame('10000.00', $user->wallet->balance);
    }
}
