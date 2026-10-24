<?php

declare(strict_types=1);

namespace OlympusHub\Console\Commands;

use Illuminate\Console\Command;
use OlympusHub\Services\ConsulRealmTetherService;
use Illuminate\Support\Facades\File;

class ImportRealmCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'olympus:import {realm? : 要匯入的環境 ID (若不指定則匯入目錄下所有 .env 檔案)}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '將 realms 目錄下的 .env 檔案匯入至 Consul KV';

    /**
     * Execute the console command.
     */
    public function handle(ConsulRealmTetherService $service): void
    {
        $realmsPath = config('olympus.consul.realms_path', base_path('realms'));
        
        if (!is_dir($realmsPath)) {
            $this->error("錯誤：位面目錄 [{$realmsPath}] 不存在。");
            return;
        }

        $realmId = $this->argument('realm');

        if ($realmId) {
            $this->importSingleFile($service, $realmsPath, $realmId);
        } else {
            $this->importAllFiles($service, $realmsPath);
        }
    }

    /**
     * 匯入單一檔案
     */
    protected function importSingleFile(ConsulRealmTetherService $service, string $path, string $realmId): void
    {
        $filePath = "{$path}/{$realmId}.env";

        if (!File::exists($filePath)) {
            $this->error("錯誤：找不到配置檔案 [{$filePath}]");
            return;
        }

        $this->processImport($service, $realmId, $filePath);
    }

    /**
     * 匯入目錄下所有 .env 檔案
     */
    protected function importAllFiles(ConsulRealmTetherService $service, string $path): void
    {
        $files = File::glob("{$path}/*.env");

        if (empty($files)) {
            $this->warn("在 [{$path}] 目錄下找不到任何 .env 檔案。");
            return;
        }

        foreach ($files as $file) {
            $realmId = File::name($file);
            $this->processImport($service, $realmId, $file);
        }
    }

    /**
     * 執行匯入動作
     */
    protected function processImport(ConsulRealmTetherService $service, string $realmId, string $filePath): void
    {
        $content = File::get($filePath);
        $consulKey = 'DB/' . strtoupper($realmId);

        $this->info("正在匯入 [{$realmId}] -> Consul Key: [{$consulKey}]...");

        try {
            $service->import($realmId, $content);
            $this->info("✅ 匯入成功！");
        } catch (\Exception $e) {
            $this->error("❌ 匯入失敗：{$e->getMessage()}");
        }
    }
}
