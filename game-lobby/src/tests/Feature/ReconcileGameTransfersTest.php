<?php

namespace Tests\Feature;

use App\Enums\GameTransferDirection;
use App\Enums\GameTransferStatus;
use App\Models\GameTransfer;
use App\Models\User;
use App\Services\Games\GameTransferReconcileService;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Http;
use Tests\TestCase;

class ReconcileGameTransfersTest extends TestCase
{
    use RefreshDatabase;

    public function test_matching_transfer_record_finalizes_transfer_out_exactly_once(): void
    {
        [$user, $transfer] = $this->unknownTransfer(GameTransferDirection::Out, '12.30');
        $this->fakeRecord($transfer, 2, '12.30');
        $service = app(GameTransferReconcileService::class);

        $this->assertSame(GameTransferStatus::Completed, $service->reconcile($transfer));
        $this->assertSame(GameTransferStatus::Completed, $service->reconcile($transfer));
        $this->assertSame('12.30', $user->wallet()->firstOrFail()->balance);
    }

    public function test_mismatched_record_goes_to_manual_review_without_changing_wallet(): void
    {
        [$user, $transfer] = $this->unknownTransfer(GameTransferDirection::Out, '12.30');
        $this->fakeRecord($transfer, 1, '12.30');

        $status = app(GameTransferReconcileService::class)->reconcile($transfer);

        $this->assertSame(GameTransferStatus::ManualReview, $status);
        $this->assertSame('0.00', $user->wallet()->firstOrFail()->balance);
        $this->assertDatabaseHas('game_transfers', ['id' => $transfer->id, 'error_code' => 'record_mismatch']);
    }

    public function test_tenth_unresolved_attempt_is_sent_to_manual_review(): void
    {
        [, $transfer] = $this->unknownTransfer(GameTransferDirection::In, '5.00', 10);

        $this->assertSame(GameTransferStatus::ManualReview, app(GameTransferReconcileService::class)->reconcile($transfer));
        Http::assertNothingSent();
    }

    public function test_invalid_state_is_sent_to_manual_review_after_limited_retries(): void
    {
        // 「狀態不合法」(403601) 不必然是玩家在線，也可能是帳號/遊戲被停用等無法自動恢復者，
        // 故僅重試有限次數（INVALID_STATE_MAX_ATTEMPTS=3）後即轉人工，不空轉到 MAX_ATTEMPTS。
        [$user, $transfer] = $this->unknownTransfer(GameTransferDirection::Out, '10.00', 3);
        $transfer->forceFill(['error_code' => '403601'])->save();
        Http::fake([
            '*/transfer*' => Http::response(['code' => 0, 'data' => null]), // 查無成功紀錄
            '*/player/logout' => Http::response(['code' => 0, 'data' => []]),
            '*/wallet/transfer-out*' => Http::response(['code' => 403601, 'message' => 'invalid state']),
        ]);

        $status = app(GameTransferReconcileService::class)->reconcile($transfer);

        $this->assertSame(GameTransferStatus::ManualReview, $status);
        $this->assertDatabaseHas('game_transfers', [
            'id' => $transfer->id,
            'status' => 'manual_review',
            'error_code' => 'invalid_state',
        ]);
        // 已達重試上限：不再嘗試登出或轉出
        Http::assertNotSent(fn ($request) => str_contains($request->url(), 'wallet/transfer-out'));
        $this->assertSame('0.00', $user->wallet()->firstOrFail()->balance);
    }

    public function test_duplicate_transfer_response_stays_unknown_and_is_not_compensated(): void
    {
        Http::fake([
            '*/player/exist*' => Http::response(['code' => 0, 'data' => ['isExist' => true]]),
            '*/wallet/transfer-in*' => Http::response(['code' => 500029, 'message' => 'duplicate', 'data' => null]),
        ]);
        $user = User::factory()->create();

        $this->actingAs($user)->post('/play/target/1004')->assertSessionHasErrors('play');

        $this->assertSame('0.00', $user->wallet()->firstOrFail()->balance);
        $this->assertDatabaseHas('game_transfers', ['user_id' => $user->id, 'status' => 'unknown', 'error_code' => '500029']);
    }

    public function test_connection_timeout_keeps_debit_in_unknown_state_for_query_first_reconciliation(): void
    {
        Http::fake([
            '*/player/exist*' => Http::response(['code' => 0, 'data' => ['isExist' => true]]),
            '*/wallet/transfer-in*' => Http::failedConnection('secret must not escape'),
        ]);
        $user = User::factory()->create();

        $this->actingAs($user)->post('/play/target/1004')->assertSessionHasErrors('play');

        $this->assertSame('0.00', $user->wallet()->firstOrFail()->balance);
        $this->assertDatabaseHas('game_transfers', [
            'user_id' => $user->id,
            'status' => 'unknown',
            'attempts' => 1,
            'error_code' => 'connection',
        ]);
    }

    /** @return array{User, GameTransfer} */
    private function unknownTransfer(GameTransferDirection $direction, string $amount, int $attempts = 1): array
    {
        $user = User::factory()->create();
        $user->wallet()->update(['balance' => 0]);
        $playerName = 'gff'.$user->id;
        $user->gameAccounts()->create(['provider' => 'target', 'player_name' => $playerName]);
        $transfer = $user->gameTransfers()->create([
            'provider' => 'target',
            'direction' => $direction,
            'amount' => $amount,
            'reference' => 'ref-'.$user->id.'-'.$direction->value,
            'status' => GameTransferStatus::Unknown,
            'attempts' => $attempts,
            'last_attempt_at' => now()->subSeconds(121),
            'next_retry_at' => now()->subSecond(),
            'meta' => ['player_name' => $playerName],
        ]);
        $transfer->forceFill(['created_at' => now()->subSeconds(121)])->save();

        return [$user, $transfer];
    }

    private function fakeRecord(GameTransfer $transfer, int $type, string $amount): void
    {
        Http::fake([
            '*/transfer*' => Http::response(['code' => 0, 'data' => [
                'pgsTransactionNo' => 'target-123',
                'merchantTransactionNo' => $transfer->reference,
                'transferType' => $type,
                'playerName' => $transfer->meta['player_name'],
                'amount' => $amount,
            ]]),
        ]);
    }
}
