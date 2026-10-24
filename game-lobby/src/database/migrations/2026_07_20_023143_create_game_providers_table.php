<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('game_providers', function (Blueprint $table) {
            $table->id();
            $table->string('key')->unique()->comment('遊戲商代號，如 target（唯一）');
            $table->string('name')->comment('顯示名稱');
            $table->string('driver')->comment('實作 GameProvider 的 driver 類別 FQCN');
            $table->json('settings')->comment('該遊戲商設定（merchant 憑證、base_url、player_prefix、快取/紀錄參數…）');
            $table->boolean('enabled')->default(true)->comment('是否啟用');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('game_providers');
    }
};
