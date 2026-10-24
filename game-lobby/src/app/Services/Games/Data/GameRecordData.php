<?php

namespace App\Services\Games\Data;

/**
 * 正規化後的單筆遊戲（下注）紀錄。
 * 各遊戲商回應欄位不同，於 from*() 做防禦性對應並保留原始 raw。
 */
class GameRecordData
{
    /**
     * @param  array<string, mixed>  $raw
     */
    public function __construct(
        public string $provider,
        public string $playerName,
        public string $recordId,
        public ?string $gameId = null,
        public ?string $gameCode = null,
        public ?string $gameName = null,
        public string $betAmount = '0',
        public string $payoutAmount = '0',
        public string $validAmount = '0',
        public ?string $currency = null,
        public ?string $betAt = null,
        public array $raw = [],
    ) {}

    /**
     * 由 target /game-records 單筆回應建立。
     * 實際 item 欄位待有真實下注時校正，此處防禦性對應常見命名並保留 raw。
     *
     * @param  array<string, mixed>  $raw
     */
    public static function fromPgs(string $provider, array $raw): self
    {
        return new self(
            provider: $provider,
            playerName: (string) ($raw['playerName'] ?? $raw['playerId'] ?? ''),
            recordId: (string) ($raw['betID'] ?? $raw['betId'] ?? $raw['roundID'] ?? $raw['roundId'] ?? $raw['recordID'] ?? $raw['recordId'] ?? $raw['id'] ?? $raw['transactionNo'] ?? ''),
            gameId: isset($raw['gameID']) || isset($raw['gameId'])
                ? (string) ($raw['gameID'] ?? $raw['gameId'])
                : null,
            gameCode: isset($raw['gameCode']) ? (string) $raw['gameCode'] : null,
            gameName: isset($raw['gameName']) ? (string) $raw['gameName'] : null,
            betAmount: (string) ($raw['betAmount'] ?? '0'),
            payoutAmount: (string) ($raw['payoutAmount'] ?? '0'),
            validAmount: (string) ($raw['validAmount'] ?? '0'),
            currency: isset($raw['currency']) ? (string) $raw['currency'] : null,
            betAt: (string) ($raw['betTime'] ?? $raw['createTime'] ?? $raw['betAt'] ?? $raw['wagerTime'] ?? '') ?: null,
            raw: $raw,
        );
    }

    /**
     * 轉為 game_records 欄位（不含 user_id，由呼叫端依綁定填入）。
     *
     * @return array<string, mixed>
     */
    public function toAttributes(): array
    {
        return [
            'provider' => $this->provider,
            'player_name' => $this->playerName,
            'record_id' => $this->recordId,
            'game_id' => $this->gameId,
            'game_code' => $this->gameCode,
            'game_name' => $this->gameName,
            'bet_amount' => $this->betAmount,
            'payout_amount' => $this->payoutAmount,
            'valid_amount' => $this->validAmount,
            'currency' => $this->currency,
            'bet_at' => $this->betAt,
            'raw' => $this->raw,
        ];
    }
}
