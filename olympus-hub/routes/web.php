<?php

use Illuminate\Support\Facades\Route;

Route::get('/', 'HubController@index')->name('olympus.hub');
Route::get('/gates', 'HubController@gates')->name('olympus.gates'); // 金流專用
Route::get('/game-gates', 'HubController@gameGates')->name('olympus.game_gates'); // 遊戲專用
Route::get('/ritual/{realm}', 'HubController@ritual')->name('olympus.ritual');
Route::post('/invoke/{realm}', 'HubController@invoke')->name('olympus.invoke');
