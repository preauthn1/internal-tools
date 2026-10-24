<?php

namespace App\Services\Games\Contracts;

use App\Services\Games\Data\GameData;
use App\Services\Games\Data\GameRecordData;
use App\Services\Games\Data\TransferRecordData;

/**
 * 遊戲商接入契約。每個遊戲商（target、未來其他家）各自實作此介面，
 * 上層（大廳、進入遊戲流程）只依賴此契約，達成多遊戲商可插拔。
 */
interface GameProvider
{
    /**
     * 遊戲商代號（如 target），對應 game_providers.key。
     */
    public function key(): string;

    /**
     * 讀取此遊戲商的設定值（來自 game_providers.settings），
     * 供上層取用如 games_cache_ttl、records_* 等參數。
     */
    public function setting(string $key, mixed $default = null): mixed;

    /**
     * 取得遊戲列表。
     *
     * @return array<int, GameData>
     */
    public function games(): array;

    // ── 以下為進入遊戲流程 / 轉帳 / 紀錄（#3~#5 實作，介面先固定契約）──

    public function playerExists(string $playerName): bool;

    /**
     * 於遊戲商建立玩家，回傳遊戲商玩家參數。
     *
     * @return array<string, mixed>
     */
    public function createPlayer(string $playerName): array;

    public function walletAmount(string $playerName): string;

    /**
     * @return array<string, mixed>
     */
    public function transferIn(string $playerName, string $amount, string $reference): array;

    /**
     * @return array<string, mixed>
     */
    public function transferOut(string $playerName, string $amount, string $reference): array;

    public function transferRecord(string $merchantTransactionNo): ?TransferRecordData;

    public function logoutPlayer(string $playerName): void;

    /**
     * 取得指定遊戲的啟動連結。
     *
     * @param  string  $gameId  遊戲商的數字 gameID（非 gameCode）
     */
    public function gameLink(string $playerName, string $gameId): string;

    /**
     * 取得遊戲（下注）紀錄的單頁結果。
     *
     * @param  array{startTime?: string, endTime?: string, playerName?: string, pageIndex?: int, pageSize?: int}  $query
     * @return array{items: array<int, GameRecordData>, pageIndex: int, pageSize: int, totalCount: int, totalPage: int}
     */
    public function gameRecords(array $query = []): array;
}
