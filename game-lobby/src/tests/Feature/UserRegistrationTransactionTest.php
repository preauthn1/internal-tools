<?php

namespace Tests\Feature;

use App\Actions\Fortify\CreateNewUser;
use App\Models\WalletTransaction;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Event;
use RuntimeException;
use Tests\TestCase;

class UserRegistrationTransactionTest extends TestCase
{
    use RefreshDatabase;

    public function test_registration_rolls_back_user_when_initial_ledger_creation_fails(): void
    {
        $event = 'eloquent.creating: '.WalletTransaction::class;
        Event::listen($event, static function (): never {
            throw new RuntimeException('simulated ledger failure');
        });

        try {
            app(CreateNewUser::class)->create([
                'account' => 'rollback-user',
                'email' => 'user@example.invalid',
                'password' => 'Password123!',
                'password_confirmation' => 'Password123!',
            ]);
            $this->fail('Expected ledger creation to fail.');
        } catch (RuntimeException $exception) {
            $this->assertSame('simulated ledger failure', $exception->getMessage());
        } finally {
            Event::forget($event);
        }

        $this->assertDatabaseMissing('users', ['email' => 'user@example.invalid']);
        $this->assertDatabaseCount('wallets', 0);
        $this->assertDatabaseCount('wallet_transactions', 0);
    }
}
