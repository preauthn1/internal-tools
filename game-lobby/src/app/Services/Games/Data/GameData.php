<?php

namespace App\Services\Games\Data;

use JsonSerializable;

/**
 * 正規化後的遊戲項目，供前端大廳顯示。
 * 各遊戲商回應欄位不同，於 from*() 做防禦性對應。
 */
class GameData implements JsonSerializable
{
    public function __construct(
        public string $provider,
        public string $gameId,
        public string $code,
        public string $name,
        public bool $isPrePaidMode = false,
        public bool $isLaunchAvailable = true,
        public ?string $image = null,
    ) {}

    /**
     * 由 target /games 單筆回應建立。
     * 注意：gameID（數字）與 gameCode（字串）常不同，進入遊戲須用 gameID。
     *
     * @param  array<string, mixed>  $raw
     */
    public static function fromPgs(string $provider, array $raw): self
    {
        return new self(
            provider: $provider,
            gameId: (string) ($raw['gameID'] ?? $raw['gameId'] ?? $raw['id'] ?? ''),
            code: (string) ($raw['gameCode'] ?? $raw['gameID'] ?? $raw['code'] ?? ''),
            name: (string) ($raw['gameName'] ?? $raw['name'] ?? ''),
            isPrePaidMode: (bool) ($raw['isPrePaidMode'] ?? false),
            isLaunchAvailable: (bool) ($raw['isLaunchAvailable'] ?? true),
            image: isset($raw['imageUrl']) || isset($raw['image']) || isset($raw['logo'])
                ? (string) ($raw['imageUrl'] ?? $raw['image'] ?? $raw['logo'])
                : null,
        );
    }

    /**
     * @return array<string, mixed>
     */
    public function jsonSerialize(): array
    {
        return [
            'provider' => $this->provider,
            'gameId' => $this->gameId,
            'code' => $this->code,
            'name' => $this->name,
            'isPrePaidMode' => $this->isPrePaidMode,
            'isLaunchAvailable' => $this->isLaunchAvailable,
            'image' => $this->image,
        ];
    }
}
