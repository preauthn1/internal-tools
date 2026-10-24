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
        Schema::create('game_records', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->comment('所屬會員')->constrained()->cascadeOnDelete();
            $table->string('provider')->comment('遊戲商代號，如 target');
            $table->string('player_name')->comment('送遊戲商的玩家識別');
            $table->string('record_id')->comment('遊戲商單筆唯一識別（供冪等 upsert）');
            $table->string('game_id')->nullable()->comment('遊戲 ID（數字）');
            $table->string('game_code')->nullable()->comment('遊戲代碼（字串）');
            $table->string('game_name')->nullable()->comment('遊戲名稱');
            $table->decimal('bet_amount', 16, 2)->default(0)->comment('下注金額');
            $table->decimal('payout_amount', 16, 2)->default(0)->comment('派彩金額');
            $table->decimal('valid_amount', 16, 2)->default(0)->comment('有效投注額');
            $table->string('currency')->nullable()->comment('幣別');
            $table->timestamp('bet_at')->nullable()->comment('下注時間');
            $table->json('raw')->nullable()->comment('完整原始 item（保留供日後校正）');
            $table->timestamps();

            $table->unique(['provider', 'record_id']);      // 冪等鍵
            $table->index(['user_id', 'provider']);
            $table->index('bet_at');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('game_records');
    }
};
