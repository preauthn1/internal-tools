<?php

namespace Tests\Feature;

use App\Services\Games\Data\GameData;
use App\Services\Games\Exceptions\GameProviderException;
use App\Services\Games\target\PgsProvider;
use Illuminate\Support\Facades\Http;
use InvalidArgumentException;
use Tests\TestCase;

class PgsProviderTest extends TestCase
{
    private function provider(): PgsProvider
    {
        return new PgsProvider([
            'merchant_id' => '175888',
            'merchant_secret' => 'test-secret',
            'base_url' => 'https://vendor-api.example.test/v1',
            'timeout' => 10,
            'locale' => 'zh-cn',
        ]);
    }

    public function test_games_parses_envelope_into_game_data()
    {
        Http::fake([
            '*/games*' => Http::response([
                'code' => 0,
                'message' => 'OK',
                'data' => [
                    ['gameID' => 1004, 'gameCode' => '1004', 'gameName' => '金吉報喜', 'imageUrl' => 'https://img.test/1004.png', 'isPrePaidMode' => false],
                ],
            ]),
        ]);

        $games = $this->provider()->games();

        $this->assertCount(1, $games);
        $this->assertInstanceOf(GameData::class, $games[0]);
        $this->assertSame('1004', $games[0]->code);
        $this->assertSame('金吉報喜', $games[0]->name);
        $this->assertSame('target', $games[0]->provider);
    }

    public function test_non_zero_code_throws_provider_exception()
    {
        Http::fake([
            '*/games*' => Http::response([
                'code' => 400001,
                'message' => '參數錯誤',
                'data' => null,
            ]),
        ]);

        $this->expectException(GameProviderException::class);

        $this->provider()->games();
    }

    public function test_request_carries_merchant_credentials()
    {
        Http::fake([
            '*/games*' => Http::response(['code' => 0, 'message' => null, 'data' => ['games' => []]]),
        ]);

        $this->provider()->games();

        Http::assertSent(function ($request) {
            return str_contains($request->url(), 'merchantID=175888')
                && str_contains($request->url(), 'merchantSecretKey=test-secret');
        });
    }

    public function test_player_exists_reads_is_exist_flag()
    {
        Http::fake([
            '*/player/exist*' => Http::response([
                'code' => 0, 'message' => 'OK', 'data' => ['isExist' => true, 'createTime' => null],
            ]),
        ]);

        $this->assertTrue($this->provider()->playerExists('gff1'));
    }

    public function test_create_player_posts_player_name()
    {
        Http::fake([
            '*/player' => Http::response(['code' => 0, 'message' => 'OK', 'data' => []]),
        ]);

        $this->provider()->createPlayer('gff1');

        Http::assertSent(fn ($request) => str_ends_with(strtok($request->url(), '?'), '/player')
            && $request['playerName'] === 'gff1'
            && $request['merchantID'] === '175888');
    }

    public function test_wallet_amount_reads_amount()
    {
        Http::fake([
            '*/wallet/amount*' => Http::response([
                'code' => 0, 'message' => 'OK', 'data' => ['currency' => 'BETA', 'amount' => 1234.5],
            ]),
        ]);

        $this->assertSame('1234.50', $this->provider()->walletAmount('gff1'));
    }

    public function test_transfer_in_sends_amount_and_merchant_order_no()
    {
        Http::fake([
            '*/wallet/transfer-in*' => Http::response([
                'code' => 0, 'message' => 'OK',
                'data' => ['transactionNo' => 'tx1', 'amount' => 100.0, 'beforeAmount' => 0.0, 'afterAmount' => 100.0],
            ]),
        ]);

        $result = $this->provider()->transferIn('gff1', '100.00', 'ref-abc');

        $this->assertSame('tx1', $result['transactionNo']);
        Http::assertSent(fn ($request) => $request['playerName'] === 'gff1'
            && (float) $request['amount'] === 100.0
            && $request['merchantOrderNo'] === 'ref-abc');
    }

    public function test_transfer_out_sends_amount_and_merchant_order_no()
    {
        Http::fake([
            '*/wallet/transfer-out*' => Http::response([
                'code' => 0, 'message' => 'OK',
                'data' => ['transactionNo' => 'tx2', 'amount' => 50.0, 'beforeAmount' => 50.0, 'afterAmount' => 0.0],
            ]),
        ]);

        $result = $this->provider()->transferOut('gff1', '50.00', 'ref-out');

        $this->assertSame('tx2', $result['transactionNo']);
        Http::assertSent(fn ($request) => $request['playerName'] === 'gff1'
            && (float) $request['amount'] === 50.0
            && $request['merchantOrderNo'] === 'ref-out');
    }

    public function test_game_link_uses_numeric_game_id_and_returns_link()
    {
        Http::fake([
            '*/game/link*' => Http::response([
                'code' => 0, 'message' => 'OK', 'data' => ['link' => 'https://game.test/launch?token=x'],
            ]),
        ]);

        $link = $this->provider()->gameLink('gff1', '1004');

        $this->assertSame('https://game.test/launch?token=x', $link);
        Http::assertSent(fn ($request) => str_contains($request->url(), 'gameID=1004')
            && ! str_contains($request->url(), 'gameCode='));
    }

    public function test_game_records_maps_items_and_pagination()
    {
        Http::fake([
            '*/game-records*' => Http::response([
                'code' => 0,
                'message' => 'OK',
                'data' => [
                    'items' => [
                        [
                            'playerName' => 'gff1', 'betId' => 'bet-001',
                            'gameID' => 1004, 'gameCode' => 'jjbx', 'gameName' => '金吉報喜',
                            'betAmount' => 10.5, 'payoutAmount' => 20.0, 'validAmount' => 10.5,
                            'currency' => 'BETA', 'betTime' => '2026-07-17T08:05:00Z',
                        ],
                    ],
                    'pageIndex' => 1, 'pageSize' => 200, 'totalCount' => 1, 'totalPage' => 1,
                ],
            ]),
        ]);

        $result = $this->provider()->gameRecords([
            'startTime' => '2026-07-17T08:00:00Z',
            'endTime' => '2026-07-17T08:30:00Z',
        ]);

        $this->assertSame(1, $result['totalPage']);
        $this->assertCount(1, $result['items']);
        $record = $result['items'][0];
        $this->assertSame('gff1', $record->playerName);
        $this->assertSame('bet-001', $record->recordId);
        $this->assertSame('1004', $record->gameId);
        $this->assertSame('10.5', $record->betAmount);
    }

    public function test_transfer_record_maps_query_response()
    {
        Http::fake([
            '*/transfer*' => Http::response(['code' => 0, 'data' => [
                'pgsTransactionNo' => 'target-1', 'merchantTransactionNo' => 'ref-1',
                'transferType' => 2, 'playerName' => 'gff1', 'amount' => '10.10',
            ]]),
        ]);

        $record = $this->provider()->transferRecord('ref-1');

        $this->assertNotNull($record);
        $this->assertSame(2, $record->transferType);
        $this->assertSame('10.10', $record->amount->value());
    }

    public function test_game_link_accepts_any_https_host()
    {
        // 連結由已認證的 target API 回傳，host 為動態第三方遊戲網域，不做白名單。
        Http::fake(['*/game/link*' => Http::response([
            'code' => 0, 'data' => ['link' => 'https://any-vendor.example/launch?token=x'],
        ])]);

        $this->assertSame(
            'https://any-vendor.example/launch?token=x',
            $this->provider()->gameLink('gff1', '1004'),
        );
    }

    public function test_game_link_rejects_non_https_link()
    {
        Http::fake(['*/game/link*' => Http::response([
            'code' => 0, 'data' => ['link' => 'http://insecure.test/launch'],
        ])]);

        $this->expectException(GameProviderException::class);
        $this->provider()->gameLink('gff1', '1004');
    }

    public function test_transfer_rejects_amount_below_provider_minimum()
    {
        $this->expectException(InvalidArgumentException::class);
        $this->provider()->transferIn('gff1', '0.09', 'ref-min');
    }

    public function test_transfer_rejects_invalid_reference()
    {
        $this->expectException(InvalidArgumentException::class);
        $this->provider()->transferIn('gff1', '0.10', 'invalid reference!');
    }
}
