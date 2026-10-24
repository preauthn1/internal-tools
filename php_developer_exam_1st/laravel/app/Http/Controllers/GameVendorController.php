<?php

namespace App\Http\Controllers;

use \Exception;
use App\Exceptions\GameVendorRuntimeException;
use App\Models\AccountTransaction;
use App\Models\GameRecord;
use App\Services\AccountTransactionService;
use App\Services\GameRecordService;
use App\Services\GameVendors\GameVendorFactory;
use App\Services\UserService;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\View\View;

class GameVendorController extends Controller
{
    public function __construct(
        protected AccountTransactionService $accountTransactionService,
        protected UserService $userService,
        protected GameRecordService $gameRecordService
    ) {
    }

    public function index(Request $request, string $vendorId): View
    {
        try {
            $vendor = GameVendorFactory::make($vendorId);
            $vendorConfig = config("game_vendors.{$vendorId}", []);
            $gameList = $vendor->getGameList();

            return view("game-vendors.index", [
                "vendorId" => $vendorId,
                "vendorIcon" => $vendorConfig["icon"],
                "vendorName" => $vendorConfig["name"],
                "gameList" => $gameList,
            ]);

        } catch (\Exception $e) {
            LOG::ERROR("Error in GameVendorController@index: {$e->getMessage()}");
            return $this->showErrorPage($vendorId);
        }
    }

    public function game(Request $request, string $vendorId, string $gameId): View
    {
        try {
            $vendor = GameVendorFactory::make($vendorId);
            $vendorConfig = config("game_vendors.{$vendorId}", []);
            $user = $request->user();
            $username = $user->name;

            try {
                $vendor->createUserIfNotExist($username);
            } catch (\Exception $e) {
                LOG::ERROR("User {$username} createUserIfNotExist failed: {$e->getMessage()}");
                throw new \Exception("建立遊戲帳號失敗，請稍後再試。");
            }

            $user_balance = $user->balance;
            $vendor_user_balance = $vendor->getBalance($username);

            LOG::DEBUG("User {$username} balance: {$user_balance}, vendor {$vendorId} balance: {$vendor_user_balance}");

            if ($user->balance <= 0 and $vendor_user_balance <= 0) {
                $errorMessage = "餘額不足，請儲值";
                return $this->showErrorPage($vendorId, $errorMessage);
            } elseif ($user_balance > 0) {
                DB::beginTransaction();

                if ($user_balance > 0) {
                    try {
                        #deposit all balance to vendor
                        LOG::DEBUG("User {$username} start creating deposit transaction");
                        $transaction = $this->accountTransactionService->createPendingDepositTransaction(
                            userId: $user->id,
                            vendor: $vendorId,
                            amount: $user_balance,
                        );

                        // Subtract balance using UserService
                        $this->userService->subBalance($user, $user_balance);

                        DB::commit();
                        LOG::DEBUG("User {$username} transaction_no {$transaction->transaction_no} created successfully" );

                    } catch (\Exception $e) {
                        DB::rollBack();
                        LOG::ERROR("User {$username} create transaction_no failed: {$e->getMessage()}");
                        throw new \Exception("建立遊戲轉帳交易失敗，請稍後再試。");
                    }

                    # Create user on vendor if not exist
                    $vendor_deposit_data = $vendor->deposit(
                        username: $username,
                        amount: $user_balance,
                        accountTransactionNo: $transaction->transaction_no,
                    );
                    $transaction->vendor_transaction_no = $vendor_deposit_data['transactionNo'] ?? null;

                    #checki if  vendor response data matches
                    if (floatval($vendor_deposit_data['amount']) !== floatval($user_balance)) {
                        LOG::ERROR("User {$username} deposit amount mismatch: local {$user_balance}, vendor {$vendor_deposit_data['amount']}, transaction_no: {$transaction->transaction_no}");
                        #update account transaction status to failed
                        $transaction->status = AccountTransaction::STATUS_FAILED;
                        $transaction->save();

                        throw new \Exception("儲值金額異常，請檢查LOG");
                    } else {
                        $transaction->status = AccountTransaction::STATUS_SUCCEED;
                        $transaction->save();
                    }

                }
            }

            $game_url = $vendor->getGameUrl($username, $gameId);

            return view("game-vendors.show", [
                "vendorId" => $vendorId,
                "gameId" => $gameId,
                "vendorIcon" => $vendorConfig["icon"] ?? null,
                "vendorName" => $vendorConfig["name"] ?? null,
                "user" => $user,
                "gameUrl" => $game_url,
            ]);

        } catch (\Exception $e) {
            Log::ERROR("Error in GameVendorController@game: ", ['message' => $e->getMessage(), 'file' => $e->getFile(), 'line' => $e->getLine()]);
            return $this->showErrorPage($vendorId);
        }
    }

    public function withdraw(Request $request, string $vendorId): View
    {
        try {
            // vendorId is already validated by GlobalParamValidator middleware
            $vendor = GameVendorFactory::make($vendorId);
            $vendorConfig = config("game_vendors.{$vendorId}", []);
            $user = $request->user();
            $username = $user->name;

            // Get vendor balance using Laravel's built-in service method
            $vendorBalance = $vendor->getBalance($username);

            return view("game-vendors.withdraw", [
                "vendorId" => $vendorId,
                "vendorIcon" => $vendorConfig["icon"] ?? null,
                "vendorName" => $vendorConfig["name"] ?? null,
                "user" => $user,
                "vendorBalance" => $vendorBalance,
            ]);
        } catch (\Exception $e) {
            Log::ERROR("Error in GameVendorController@withdraw: {$e->getMessage()}");
            return $this->showErrorPage($vendorId);
        }
    }

    public function processWithdraw(Request $request, string $vendorId): RedirectResponse
    {
        try {
            // vendorId is already validated by GlobalParamValidator middleware
            $vendor = GameVendorFactory::make($vendorId);
            $user = $request->user();
            $username = $user->name;
            $vendorBalance = $vendor->getBalance($username);

            //1. 建立單號
            $transaction = $this->accountTransactionService->createPendingWithdrawTransaction(
                userId: $user->id,
                vendor: $vendorId,
                amount: $vendorBalance,
            );

            //2. 執行提款
            $vendorWithdrawData = $vendor->withdraw(
                username: $username,
                amount: $vendorBalance,
                accountTransactionNo: $transaction->transaction_no,
            );

            //3. 驗證回傳資料並更新使用者餘額及交易狀態
            if (!isset($vendorWithdrawData['amount']) or floatval($vendorWithdrawData['amount']) <= 0) {
                Log::ERROR("User {$username} withdraw amount : {$vendorBalance}, vendor returned invalid amount: {$vendorWithdrawData['amount']}, transaction_no: {$transaction->transaction_no}");
                $transaction->status = AccountTransaction::STATUS_FAILED;
                $transaction->save();

                return redirect()
                    ->route('game-vendors.withdraw', $vendorId)
                    ->withErrors(['amount' => '提款異常，請聯繫客服']);
            }

            DB::beginTransaction();

            try {
                $transaction->vendor_transaction_no = $vendorWithdrawData['transactionNo'];
                $transaction->status = AccountTransaction::STATUS_SUCCEED;
                $transaction->save();

                $this->userService->addBalance($user, $vendorBalance);

                DB::commit();

                Log::INFO("User {$username} withdraw successful: amount {$vendorBalance}, transaction_no: {$transaction->transaction_no}");

                return redirect()
                    ->route('game-vendors.withdraw', $vendorId)
                    ->with('success', "成功提款 {$vendorBalance}。交易編號: {$transaction->transaction_no}");
            } catch (\Exception $e) {
                DB::rollBack();
                Log::ERROR("User {$username} withdraw failed: {$e->getMessage()}");

                return redirect()
                    ->route('game-vendors.withdraw', $vendorId)
                    ->withErrors(['amount' => '提款失敗，聯絡客服。']);
            }
        } catch (\Exception $e) {
            Log::ERROR("Error in GameVendorController@processWithdraw: {$e->getMessage()}");
            return redirect()
                ->route('game-vendors.withdraw', $vendorId)
                ->withErrors(['amount' => '系統異常，請稍後再試。']);
        }
    }

    public function gameRecords(Request $request, string $vendorId): View
    {
        try {
            // vendorId is already validated by GlobalParamValidator middleware
            $vendor = GameVendorFactory::make($vendorId);
            $vendorConfig = config("game_vendors.{$vendorId}", []);
            $user = $request->user();

            // Get filters from request query parameters using Laravel's built-in request methods
            $filters = [
                'vendor' => $vendorId,
                'username' => $user->name,
                'per_page' => $request->get('per_page', 15),
            ];

            // Add date filters if provided
            if ($request->has('start_time')) {
                $filters['start_time'] = $request->get('start_time');
            }

            if ($request->has('end_time')) {
                $filters['end_time'] = $request->get('end_time');
            }

            // Fetch game records from database using GameRecordService
            $gameRecords = $this->gameRecordService->getByUser($user, $filters);

            // Get game list from cache to map game_id to zh-CN names
            $gameList = $vendor->getGameList();
            $gameNameMap = [];
            
            // Recursively collect all games and build game_id => zh-CN name mapping
            $collectGames = function ($data) use (&$collectGames, &$gameNameMap) {
                if (!is_array($data)) {
                    return;
                }

                // Check if this array itself looks like a game
                $hasId = isset($data['id']) || isset($data['game_id']) || isset($data['gameId']) || isset($data['gameID']) || isset($data['game_code']) || isset($data['gameCode']) || isset($data['code']);
                
                if ($hasId) {
                    $gameId = $data['gameID'] ?? $data['gameCode'] ?? $data['id'] ?? $data['game_id'] ?? $data['gameId'] ?? $data['game_code'] ?? $data['code'] ?? null;
                    
                    if ($gameId !== null) {
                        // Extract zh-CN name from gameNameMappings
                        $gameName = null;
                        if (isset($data['gameNameMappings']) && is_array($data['gameNameMappings'])) {
                            foreach ($data['gameNameMappings'] as $mapping) {
                                if (isset($mapping['languageCode']) && $mapping['languageCode'] === 'zh-cn' && isset($mapping['text'])) {
                                    $gameName = $mapping['text'];
                                    break;
                                }
                            }
                        }
                        
                        // Fallback to other name fields if zh-CN not found
                        if ($gameName === null) {
                            $gameName = $data['name'] ?? $data['game_name'] ?? $data['title'] ?? $data['gameName'] ?? $data['game_title'] ?? null;
                        }
                        
                        if ($gameName !== null) {
                            $gameNameMap[(string) $gameId] = $gameName;
                        }
                    }
                }

                // Recurse into array items
                foreach ($data as $item) {
                    if (is_array($item)) {
                        $collectGames($item);
                    }
                }
            };
            
            $collectGames($gameList);

            return view("game-vendors.game-records", [
                "vendorId" => $vendorId,
                "vendorIcon" => $vendorConfig["icon"] ?? null,
                "vendorName" => $vendorConfig["name"] ?? null,
                "user" => $user,
                "gameRecords" => $gameRecords,
                "gameNameMap" => $gameNameMap,
            ]);
        } catch (\Exception $e) {
            Log::ERROR("Error in GameVendorController@gameRecords: {$e->getMessage()}");
            return $this->showErrorPage($vendorId);
        }
    }

    private function showErrorPage(string $vendorId, string $errorMessage = ""): View
    {
        return view("game-vendors.error", [
            "vendorId" => $vendorId,
            "errorMessage" => $errorMessage ?? "系統異常，請稍後再試。",
        ]);
    }
}
