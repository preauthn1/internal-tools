<?php

namespace Tests\Feature;

use App\Models\GameRecord;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Http;
use Tests\TestCase;

class SyncGameRecordsTest extends TestCase
{
    use RefreshDatabase;

    private function bindPlayer(User $user): void
    {
        $user->gameAccounts()->create([
            'provider' => 'target',
            'player_name' => 'gff'.$user->id,
        ]);
    }

    private function fakeRecords(string $playerName): void
    {
        Http::fake([
            '*/game-records*' => Http::response([
                'code' => 0,
                'message' => 'OK',
                'data' => [
                    'items' => [
                        // 反映真實 target 欄位：betID / betTime（ISO8601 帶 T/Z）
                        [
                            'playerName' => $playerName, 'betID' => 'bet-001',
                            'gameID' => 1004, 'gameCode' => 'jjbx', 'gameName' => '金吉報喜',
                            'betAmount' => 10.5, 'payoutAmount' => 20.0, 'validAmount' => 10.5,
                            'currency' => 'BETA', 'betTime' => '2026-07-17T08:05:00Z',
                        ],
                        // 非本平台玩家 → 應被略過
                        [
                            'playerName' => 'stranger999', 'betID' => 'bet-002',
                            'gameID' => 1004, 'betAmount' => 5.0, 'payoutAmount' => 0.0, 'validAmount' => 5.0,
                            'betTime' => '2026-07-17T08:06:00Z',
                        ],
                    ],
                    'pageIndex' => 1, 'pageSize' => 200, 'totalCount' => 2, 'totalPage' => 1,
                ],
            ]),
        ]);
    }

    public function test_sync_stores_matched_records_and_skips_unknown_players()
    {
        $user = User::factory()->create();
        $this->bindPlayer($user);
        $this->fakeRecords('gff'.$user->id);

        $this->artisan('game:sync-records', [
            '--since' => '2026-07-17T08:00:00Z',
            '--until' => '2026-07-17T08:30:00Z',
        ])->assertSuccessful();

        $this->assertDatabaseCount('game_records', 1);
        $record = GameRecord::first();
        $this->assertSame($user->id, $record->user_id);
        $this->assertSame('bet-001', $record->record_id);
        $this->assertSame('1004', $record->game_id);
        $this->assertSame('10.50', $record->bet_amount);
        // ISO8601（帶 T/Z）需正規化為 DB datetime
        $this->assertSame('2026-07-17 08:05:00', $record->bet_at->format('Y-m-d H:i:s'));
        $this->assertNotNull(Cache::get('game-records:cursor:target'));
    }

    public function test_sync_is_idempotent_on_rerun()
    {
        $user = User::factory()->create();
        $this->bindPlayer($user);
        $this->fakeRecords('gff'.$user->id);

        $args = ['--since' => '2026-07-17T08:00:00Z', '--until' => '2026-07-17T08:30:00Z'];

        $this->artisan('game:sync-records', $args)->assertSuccessful();
        $this->artisan('game:sync-records', $args)->assertSuccessful();

        $this->assertDatabaseCount('game_records', 1);
    }
}
