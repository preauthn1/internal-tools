<?php

namespace Database\Seeders;

use App\Models\GameProvider;
use App\Services\Games\target\PgsProvider;
use Illuminate\Database\Seeder;

class GameProviderSeeder extends Seeder
{
    public function run(): void
    {
        GameProvider::query()->updateOrCreate(
            ['key' => 'target'],
            [
                'name' => 'target',
                'driver' => PgsProvider::class,
                'enabled' => true,
                'settings' => config('games.provider_defaults.target', []),
            ],
        );
    }
}
