<?php

use App\Http\Controllers\DashboardController;
use App\Http\Controllers\GameVendorController;
use App\Http\Controllers\ProfileController;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return view('welcome');
});

Route::middleware('auth')->group(function () {
    Route::get('/dashboard', [DashboardController::class, 'index'])->name('dashboard');

        Route::prefix('game-vendors')->name('game-vendors.')->middleware('validate.vendor')->group(function () {
            Route::get('/{vendorId}', [GameVendorController::class, 'index'])->name('index');
            Route::get('/{vendorId}/withdraw', [GameVendorController::class, 'withdraw'])->name('withdraw');
            Route::post('/{vendorId}/withdraw', [GameVendorController::class, 'processWithdraw'])->name('withdraw.process');
            Route::get('/{vendorId}/game-records', [GameVendorController::class, 'gameRecords'])->name('game-records');
            Route::get('/{vendorId}/{gameId}', [GameVendorController::class, 'game'])->name('game');
        });

    Route::get('/profile', [ProfileController::class, 'edit'])->name('profile.edit');
    Route::patch('/profile', [ProfileController::class, 'update'])->name('profile.update');
    Route::delete('/profile', [ProfileController::class, 'destroy'])->name('profile.destroy');
});

require __DIR__ . '/auth.php';
