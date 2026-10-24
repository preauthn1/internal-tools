<?php

namespace Tests\Feature;

use App\Enums\GameTransferDirection;
use App\Enums\GameTransferStatus;
use App\Models\GameTransfer;
use App\Models\GameTransferAudit;
use App\Models\User;
use Carbon\CarbonImmutable;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Cache;
use Tests\TestCase;

class AuditGameWalletSessionsTest extends TestCase
{
    use RefreshDatabase;

    public function test_completed_session_matches_game_record_accounting(): void
    {
        [$user, , $transferOut] = $this->createSession('1000.00', '950.00', CarbonImmutable::now()->subMinutes(130));
        $this->record($user, 'bet-1', '300.00', '250.00', $transferOut->created_at->subMinute());
        $this->recordsReady($transferOut);

        $this->artisan('game:audit-wallet-sessions')->assertSuccessful();

        $this->assertDatabaseHas('game_transfer_audits', [
            'transfer_out_id' => $transferOut->id,
            'opening_balance' => '0.00',
            'total_transfer_in' => '1000.00',
            'total_bet' => '300.00',
            'total_payout' => '250.00',
            'expected_transfer_out' => '950.00',
            'actual_transfer_out' => '950.00',
            'difference' => '0.00',
            'record_count' => 1,
            'status' => 'matched',
        ]);

        $this->artisan('game:audit-wallet-sessions')->assertSuccessful();
        $this->assertDatabaseCount('game_transfer_audits', 1);
    }

    public function test_audit_waits_until_record_sync_cursor_covers_delay_window(): void
    {
        [, , $transferOut] = $this->createSession('100.00', '100.00', CarbonImmutable::now()->subMinutes(3));
        Cache::forever('game-records:cursor:target', $transferOut->created_at->addMinutes(4)->toIso8601String());

        $this->artisan('game:audit-wallet-sessions')->assertSuccessful();

        $this->assertDatabaseHas('game_transfer_audits', [
            'transfer_out_id' => $transferOut->id,
            'status' => 'pending_records',
            'reason' => 'records_not_ready',
            'expected_transfer_out' => null,
        ]);
    }

    public function test_matching_recent_session_stays_pending_until_records_stabilize(): void
    {
        [, , $transferOut] = $this->createSession('100.00', '100.00', CarbonImmutable::now()->subMinutes(20));
        $this->recordsReady($transferOut);

        $this->artisan('game:audit-wallet-sessions')->assertSuccessful();

        $this->assertDatabaseHas('game_transfer_audits', [
            'transfer_out_id' => $transferOut->id,
            'expected_transfer_out' => '100.00',
            'difference' => '0.00',
            'status' => 'pending_records',
            'reason' => 'stabilizing_records',
        ]);
    }

    public function test_late_record_can_recover_manual_review_with_reference(): void
    {
        [$user, , $transferOut] = $this->createSession('1000.00', '950.00', CarbonImmutable::now()->subMinutes(130));
        $this->record($user, 'bet-1', '300.00', '0.00', $transferOut->created_at->subMinutes(2));
        $this->recordsReady($transferOut);

        $this->artisan('game:audit-wallet-sessions')->assertSuccessful();
        $this->assertDatabaseHas('game_transfer_audits', [
            'transfer_out_id' => $transferOut->id,
            'status' => 'manual_review',
            'expected_transfer_out' => '700.00',
            'difference' => '250.00',
            'reason' => 'amount_mismatch',
        ]);

        $this->record($user, 'payout-1', '0.00', '250.00', $transferOut->created_at->subMinute());
        $this->artisan('game:audit-wallet-sessions', ['--reference' => $transferOut->reference])->assertSuccessful();

        $audit = GameTransferAudit::query()->where('transfer_out_id', $transferOut->id)->firstOrFail();
        $this->assertSame('matched', $audit->status->value);
        $this->assertSame('950.00', $audit->expected_transfer_out);
        $this->assertSame('0.00', $audit->difference);
        $this->assertSame(2, $audit->record_count);
    }

    public function test_expected_amount_includes_opening_balance_and_prior_transfer_out(): void
    {
        [$user, $transferIn, $transferOut] = $this->createSession('1000.00', '550.05', CarbonImmutable::now()->subMinutes(130));
        $transferIn->forceFill(['meta' => [
            'player_name' => 'target'.$user->id,
            'provider_response' => ['beforeAmount' => '0.05', 'afterAmount' => '1000.05'],
        ]])->save();
        $priorOut = $user->gameTransfers()->create([
            'provider' => 'target',
            'direction' => GameTransferDirection::Out,
            'amount' => '400.00',
            'reference' => 'audit-prior-out-'.$user->id,
            'status' => GameTransferStatus::Completed,
            'resolved_at' => $transferOut->created_at->subMinutes(10),
            'meta' => ['player_name' => 'target'.$user->id],
        ]);
        $priorOut->forceFill([
            'created_at' => $transferOut->created_at->subMinutes(10),
            'updated_at' => $transferOut->created_at->subMinutes(10),
        ])->save();
        $this->record($user, 'bet-1', '50.00', '0.00', $transferOut->created_at->subMinute());
        $this->recordsReady($transferOut);

        $this->artisan('game:audit-wallet-sessions', ['--reference' => $transferOut->reference])->assertSuccessful();

        $this->assertDatabaseHas('game_transfer_audits', [
            'transfer_out_id' => $transferOut->id,
            'opening_balance' => '0.05',
            'prior_transfer_out' => '400.00',
            'expected_transfer_out' => '550.05',
            'difference' => '0.00',
            'status' => 'matched',
        ]);
    }

    /** @return array{User, GameTransfer, GameTransfer} */
    private function createSession(string $transferInAmount, string $transferOutAmount, CarbonImmutable $endedAt): array
    {
        $startedAt = $endedAt->subMinutes(30);
        $user = User::factory()->create();
        $playerName = 'target'.$user->id;
        $user->gameAccounts()->create(['provider' => 'target', 'player_name' => $playerName]);

        $transferIn = $user->gameTransfers()->create([
            'provider' => 'target',
            'direction' => GameTransferDirection::In,
            'amount' => $transferInAmount,
            'reference' => 'audit-in-'.$user->id,
            'status' => GameTransferStatus::Completed,
            'resolved_at' => $startedAt,
            'meta' => [
                'player_name' => $playerName,
                'provider_response' => ['beforeAmount' => '0.00', 'afterAmount' => $transferInAmount],
            ],
        ]);
        $transferIn->forceFill(['created_at' => $startedAt, 'updated_at' => $startedAt])->save();

        $transferOut = $user->gameTransfers()->create([
            'provider' => 'target',
            'direction' => GameTransferDirection::Out,
            'amount' => $transferOutAmount,
            'reference' => 'audit-out-'.$user->id,
            'status' => GameTransferStatus::Completed,
            'resolved_at' => $endedAt,
            'meta' => ['player_name' => $playerName, 'provider_response' => ['afterAmount' => '0.00']],
        ]);
        $transferOut->forceFill(['created_at' => $endedAt, 'updated_at' => $endedAt])->save();

        return [$user, $transferIn, $transferOut];
    }

    private function record(User $user, string $recordId, string $bet, string $payout, CarbonImmutable $betAt): void
    {
        $user->gameRecords()->create([
            'provider' => 'target',
            'player_name' => 'target'.$user->id,
            'record_id' => $recordId,
            'bet_amount' => $bet,
            'payout_amount' => $payout,
            'valid_amount' => $bet,
            'currency' => 'BETA',
            'bet_at' => $betAt,
        ]);
    }

    private function recordsReady(GameTransfer $transferOut): void
    {
        // 游標須超過 audit_delay_minutes（預設 30）才算紀錄齊全，取 60 分留足餘裕。
        Cache::forever('game-records:cursor:target', $transferOut->created_at->addMinutes(60)->toIso8601String());
    }
}
