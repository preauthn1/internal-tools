<?php

use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/

Auth::routes();

Route::get('/', [App\Http\Controllers\HomeController::class, 'index']);
Route::post('/game_link', [App\Http\Controllers\HomeController::class, 'game_link'])->name('game_link');
Route::post('/transfer_in', [App\Http\Controllers\HomeController::class, 'transfer_in'])->name('transfer_in');
Route::post('/transfer_out', [App\Http\Controllers\HomeController::class, 'transfer_out'])->name('transfer_out');
Route::get('/game_records', [App\Http\Controllers\HomeController::class, 'game_records'])->name('game_records');
Route::get('/transfer_points', [App\Http\Controllers\HomeController::class, 'transfer_points'])->name('transfer_points');
