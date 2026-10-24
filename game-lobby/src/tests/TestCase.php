<?php

namespace Tests;

use Database\Seeders\GameProviderSeeder;
use Illuminate\Foundation\Testing\TestCase as BaseTestCase;

abstract class TestCase extends BaseTestCase
{
    // 遊戲商設定改由 game_providers 資料表提供；用 RefreshDatabase 的測試自動 seed target。
    protected bool $seed = true;

    protected string $seeder = GameProviderSeeder::class;
}
