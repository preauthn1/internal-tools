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
        Schema::create('game_provider_accounts', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->comment('所屬會員')->constrained()->cascadeOnDelete();
            $table->string('provider')->comment('遊戲商代號，如 target');
            $table->string('player_name')->comment('送遊戲商的玩家識別（非平台 account）');
            $table->json('meta')->nullable()->comment('createPlayer 回傳的遊戲商參數');
            $table->timestamps();

            $table->unique(['user_id', 'provider']);    // 一位會員在每個遊戲商僅一組帳號
            $table->index(['provider', 'player_name']);
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('game_provider_accounts');
    }
};
