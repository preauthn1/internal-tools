<?php

namespace OlympusHub;

use Illuminate\Support\Facades\Route;
use Illuminate\Support\ServiceProvider;
use OlympusHub\Console\Commands\ImportRealmCommand;
use OlympusHub\Consul\ConsulKvClient;
use OlympusHub\Services\ConsulRealmTetherService;
use OlympusHub\Support\DatabaseConfigOverrider;
use OlympusHub\Support\EnvParser;
use Consul\Services\KV;
use Consul\Client;

class OlympusHubServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        $this->mergeConfigFrom(__DIR__ . '/../config/olympus.php', 'olympus');

        // 只有在本地開發環境才會啟用聖域門戶
        if ($this->app->environment('local')) {
            $this->registerRoutes();
            $this->registerViews();
            $this->registerServices();
        }
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        // 這裡處理資源的發布
        if ($this->app->runningInConsole()) {
            $this->commands([
                ImportRealmCommand::class,
            ]);

            $this->publishes([
                __DIR__ . '/../config/olympus.php' => config_path('olympus.php'),
            ], 'olympus-config');
        }
    }

    /**
     * 註冊奧林帕斯的路由
     */
    protected function registerRoutes(): void
    {
        Route::group([
            'prefix'     => 'olympus',
            'namespace'  => 'OlympusHub\Http\Controllers',
            'middleware' => ['web'],
        ], function () {
            $this->loadRoutesFrom(__DIR__ . '/../routes/web.php');
        });
    }

    /**
     * 註冊奧林帕斯的視圖
     */
    protected function registerViews(): void
    {
        $this->loadViewsFrom(__DIR__ . '/../resources/views', 'olympus');
    }

    /**
     * 註冊奧林帕斯所需的服務
     */
    protected function registerServices(): void
    {
        $this->app->singleton(ConsulRealmTetherService::class, function ($app) {
            $config = [
                'base_uri' => config('olympus.consul.base_url'),
            ];

            if ($token = [REDACTED]'olympus.consul.token')) {
                $config['headers'] = [
                    'X-Consul-Token' => $token,
                ];
            }

            $client = new Client($config);

            return new ConsulRealmTetherService(
                new ConsulKvClient(new KV($client)),
                new EnvParser(),
                new DatabaseConfigOverrider()
            );
        });
    }
}
