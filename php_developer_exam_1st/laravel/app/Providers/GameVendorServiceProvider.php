<?php

namespace App\Providers;

use App\Services\GameVendors\GameVendorFactory;
use Illuminate\Support\ServiceProvider;

class GameVendorServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        $this->app->singleton(GameVendorFactory::class, function ($app) {
            return new GameVendorFactory();
        });

        $this->app->bind('game.vendor', function ($app) {
            return new GameVendorFactory();
        });
    }

    public function boot(): void {}
}
