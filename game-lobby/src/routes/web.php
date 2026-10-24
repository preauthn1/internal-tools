<?php

use App\Http\Controllers\GamePlayController;
use App\Http\Controllers\GameRecordController;
use App\Http\Controllers\LobbyController;
use App\Http\Controllers\ProviderTesterController;
use App\Http\Controllers\WalletTransactionController;
use Illuminate\Support\Facades\Route;

// 首頁 = 遊戲大廳（公開瀏覽；餘額/會員名稱由共享 props 於前端依登入狀態顯示）
Route::get('/', [LobbyController::class, 'index'])->name('home');

Route::middleware(['auth', 'verified'])->group(function () {
    // 遊戲紀錄（會員檢視自己的下注紀錄）
    Route::get('records', [GameRecordController::class, 'index'])->name('records');

    // 錢包紀錄（會員檢視自己的錢包帳本）
    Route::get('wallet', [WalletTransactionController::class, 'index'])->name('wallet');
});

// 進入遊戲：登入即可（不強制 email 驗證，符合 prototype 需求）。
// 以 GET 於新分頁直接開啟：伺服端確保綁定 + 全額轉入後，直接 302 導向遊戲商啟動連結，
// 新分頁全程顯示瀏覽器原生載入中（無空白頁 / 無前端時序問題）。
Route::middleware('auth')->group(function () {
    Route::post('/play/{provider}/{gameId}', [GamePlayController::class, 'launch'])
        ->whereNumber('gameId')
        ->name('play');

    // 結算（#5 轉出）：離開遊戲回到大廳時自動觸發，把遊戲商餘額轉回平台
    Route::post('/settle', [GamePlayController::class, 'settle'])->name('settle');

    // 遊戲商 API 測試頁（開發工具，僅非正式環境；控制器內再次守門）
    Route::get('/dev/provider-tester', [ProviderTesterController::class, 'index'])->name('provider-tester');
    Route::post('/dev/provider-tester/call', [ProviderTesterController::class, 'call'])->name('provider-tester.call');
});

require __DIR__.'/settings.php';
