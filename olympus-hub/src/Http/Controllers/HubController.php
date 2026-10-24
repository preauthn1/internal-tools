<?php

namespace OlympusHub\Http\Controllers;

use Illuminate\Routing\Controller;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Artisan;
use OlympusHub\Exceptions\ConsulConnectionException;
use OlympusHub\Exceptions\InvalidTetherDataException;
use OlympusHub\Exceptions\RealmTetherNotFoundException;
use OlympusHub\Services\ConsulRealmTetherService;

class HubController extends Controller
{
    /** @var ConsulRealmTetherService */
    private $tetherService;

    public function __construct(ConsulRealmTetherService $tetherService)
    {
        $this->tetherService = $tetherService;
    }

    /**
     * 聖所首頁 (Sanctum Home)
     */
    public function index()
    {
        $basePath = base_path();

        // 檢測宿主專案以決定顯示哪個門戶
        $showCashFlow = str_contains($basePath, 'apis');
        $showGame = str_contains($basePath, 'RD-Studio-01');

        return view('olympus::index', compact('showCashFlow', 'showGame'));
    }

    /**
     * 命運之門 (The Gates of Destiny) - 金流連線選擇頁面
     */
    public function gates()
    {
        $realms = $this->getCashFlowRealms();
        $title = "金流命運之門 (Cash Flow Realms)";
        return view('olympus::gates', compact('realms', 'title'));
    }

    /**
     * 遊戲核心聖域 (Game Core Realm) - 遊戲連線選擇頁面
     */
    public function gameGates()
    {
        $realms = $this->getGameRealms();
        $title = "遊戲核心聖域 (Game Core Realms)";
        return view('olympus::gates', compact('realms', 'title'));
    }

    /**
     * 神聖祭典 (The Sacred Ritual) - 指令輸入頁面
     */
    public function ritual($realmId)
    {
        $allRealms = array_merge($this->getCashFlowRealms(), $this->getGameRealms());
        $realm = collect($allRealms)->firstWhere('id', $realmId);
        
        if (!$realm) {
            return redirect()->route('olympus.hub');
        }

        // 根據位面類型決定指令列表
        $commands = $this->getCommandsForRealm($realmId);

        return view('olympus::ritual', compact('realm', 'commands'));
    }

    /**
     * 降下神諭 (Invoke Oracle) - 執行指令
     */
    public function invoke(Request $request, $realmId)
    {
        // 聖域安全協議：僅限於本地環境啟動神諭
        if (!app()->environment('local')) {
            return response()->json(['status' => 'error', 'message' => '❌ 聖域拒絕降臨：神諭僅能在本地位面 (APP_ENV=local) 啟動。']);
        }

        $command = $request->input('command');
        $arguments = $request->input('arguments', '');

        // 安全檢查：只允許特定的指令
        $allowedCommands = [];
        foreach (config('olympus.commands') as $type => $groups) {
            foreach ($groups as $group => $commands) {
                $allowedCommands = array_merge($allowedCommands, array_keys($commands));
            }
        }

        if (!in_array($command, $allowedCommands)) {
            return response()->json(['status' => 'error', 'message' => '❌ 非法指令！聖域拒絕執行未經授權的神諭。']);
        }

        try {
            // 遊戲位面跳過神聖鏈結活化
            $gameRealmIds = collect($this->getGameRealms())->pluck('id')->toArray();
            if (!in_array($realmId, $gameRealmIds)) {
                // 活化神聖鏈結 (Activate Consul Realm Tether)
                $this->tetherService->activate($realmId);
            }

            // 執行神諭
            $exitCode = Artisan::call("{$command} {$arguments}");
            $output = Artisan::output();

            return response()->json([
                'status' => $exitCode === 0 ? 'success' : 'error',
                'chronicle' => $output,
                'message' => $exitCode === 0 
                    ? "⚡ 神諭已在 [{$realmId}] 成功降臨！" 
                    : "❌ [{$realmId}] 位面發生了災厄 (Calamity)",
            ]);
        } catch (RealmTetherNotFoundException $e) {
            return response()->json(['status' => 'error', 'message' => "❌ 神聖鏈結不存在: {$e->getMessage()}"]);
        } catch (ConsulConnectionException $e) {
            return response()->json(['status' => 'error', 'message' => "❌ 無法連線至 Consul: {$e->getMessage()}"]);
        } catch (InvalidTetherDataException $e) {
            return response()->json(['status' => 'error', 'message' => "❌ 鏈結資料異常: {$e->getMessage()}"]);
        } catch (\Exception $e) {
            return response()->json([
                'status' => 'error',
                'message' => "💥 聖域發生崩塌: " . $e->getMessage(),
            ]);
        }
    }

    /**
     * 獲取金流位面
     */
    protected function getCashFlowRealms()
    {
        return [
            ['id' => 'localhost', 'name' => '本地神域 (Localhost)', 'description' => '凡間的棲息之地，一切的起點。', 'status' => 'Stable', 'power' => 'Normal'],
            ['id' => 'dev', 'name' => '眾神試煉場 (DEV)', 'description' => '凡人進階神域的必經之路。', 'status' => 'Stable', 'power' => 'Extreme'],
            ['id' => 'tpower', 'name' => '雷霆之力 (T-Power)', 'description' => '金流如同雷霆般奔湧。', 'status' => 'Stable', 'power' => 'Extreme'],
            ['id' => 'vworld', 'name' => '幻象世界 (V-World)', 'description' => '虛實交錯的數位領域。', 'status' => 'Active', 'power' => 'High'],
            ['id' => 'alfawin77', 'name' => '阿爾法勝境 (Alfawin77)', 'description' => '最初的起源之地。', 'status' => 'Stable', 'power' => 'Ancient'],
            ['id' => 'yes2win', 'name' => '必勝神域 (Yes2Win)', 'description' => '命運女神眷顧之地。', 'status' => 'Stable', 'power' => 'Blessed'],
            ['id' => 'vplus33', 'name' => '極致增益 (V-Plus33)', 'description' => '不斷自我進化的領域。', 'status' => 'Active', 'power' => 'Transcendent'],
            ['id' => 'me77', 'name' => '自我聖所 (ME77)', 'description' => '金流因果在此靜靜律動。', 'status' => 'Resting', 'power' => 'Stable'],
        ];
    }

    /**
     * 獲獲遊戲位面
     */
    protected function getGameRealms()
    {
        return [
            ['id' => 'fon666', 'name' => '豐收 666 (Fon666)', 'description' => '無盡的財富在此匯聚。', 'status' => 'Active', 'power' => 'Radiant'],
            ['id' => 'vplus33-game', 'name' => '增益核心 (VPlus33 Core)', 'description' => '遊戲邏輯的核心運算位面。', 'status' => 'Stable', 'power' => 'Maximum'],
        ];
    }

    /**
     * 根據 Realm 獲取對應指令
     */
    protected function getCommandsForRealm($realmId)
    {
        $gameRealms = collect($this->getGameRealms())->pluck('id')->toArray();

        if (in_array($realmId, $gameRealms)) {
            return config('olympus.commands.game');
        }

        return config('olympus.commands.finance');
    }
}
