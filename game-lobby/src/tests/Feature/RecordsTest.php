<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Http;
use Inertia\Testing\AssertableInertia as Assert;
use Tests\TestCase;

class RecordsTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        // /records 會取遊戲列表把編號解析成名稱；fake 掉避免測試打真實遊戲商 API。
        Http::fake([
            '*/games*' => Http::response(['code' => 0, 'message' => 'OK', 'data' => ['games' => [
                ['gameID' => '290187', 'gameCode' => 'lion', 'gameName' => '舞獅爭霸'],
            ]]]),
        ]);
    }

    private function makeRecord(User $user, string $recordId, string $gameName): void
    {
        $user->gameRecords()->create([
            'provider' => 'target',
            'player_name' => 'gff'.$user->id,
            'record_id' => $recordId,
            'game_id' => '1004',
            'game_code' => 'jjbx',
            'game_name' => $gameName,
            'bet_amount' => '10.50',
            'payout_amount' => '20.00',
            'valid_amount' => '10.50',
            'currency' => 'BETA',
            'bet_at' => '2026-07-17 08:05:00',
        ]);
    }

    public function test_guest_is_redirected_to_login()
    {
        $this->get(route('records'))->assertRedirect(route('login'));
    }

    public function test_member_sees_only_their_own_records()
    {
        $user = User::factory()->create();
        $other = User::factory()->create();
        $this->makeRecord($user, 'bet-mine', '金吉報喜');
        $this->makeRecord($other, 'bet-theirs', '別人的遊戲');

        $this->actingAs($user)
            ->get(route('records'))
            ->assertOk()
            ->assertInertia(fn (Assert $page) => $page
                ->component('Records')
                ->has('records.data', 1)
                ->where('records.data.0.record_id', 'bet-mine')
                ->where('records.data.0.game_name', '金吉報喜'),
            );
    }

    public function test_game_number_is_resolved_to_name_via_games_list()
    {
        $user = User::factory()->create();
        // 遊戲商紀錄未帶名稱（game_name 為 null），僅有數字 game_id。
        $user->gameRecords()->create([
            'provider' => 'target',
            'player_name' => 'target'.$user->id,
            'record_id' => 'bet-no-name',
            'game_id' => '290187',
            'game_code' => null,
            'game_name' => null,
            'bet_amount' => '0.18',
            'payout_amount' => '0.00',
            'valid_amount' => '0.18',
            'currency' => 'BETA',
            'bet_at' => '2026-07-21 06:24:00',
        ]);

        $this->actingAs($user)
            ->get(route('records'))
            ->assertOk()
            ->assertInertia(fn (Assert $page) => $page
                ->component('Records')
                ->where('records.data.0.game_name', null)
                // 前端據此對照表把 290187 顯示成「舞獅爭霸」。
                ->where('gameNames.target.290187', '舞獅爭霸'),
            );
    }
}
