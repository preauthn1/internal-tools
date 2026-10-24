<?php

namespace App\Models;

use Database\Factories\GameProviderFactory;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

/**
 * 遊戲商註冊資料（取代 config/games.php 的 providers）。
 * 注意：此為「持久化的遊戲商設定列」，與 App\Services\Games\Contracts\GameProvider
 * （driver 契約/實例）是不同概念——後者由 GameProviderManager 依本表的 driver + settings 解析而來。
 *
 * @property int $id
 * @property string $key
 * @property string $name
 * @property string $driver
 * @property array<string, mixed> $settings
 * @property bool $enabled
 */
class GameProvider extends Model
{
    /** @use HasFactory<GameProviderFactory> */
    use HasFactory;

    protected $fillable = [
        'key',
        'name',
        'driver',
        'settings',
        'enabled',
    ];

    protected function casts(): array
    {
        return [
            'settings' => 'array',
            'enabled' => 'boolean',
        ];
    }
}
