<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('game_transfer_audits', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->comment('會員 ID')->constrained()->cascadeOnDelete();
            $table->string('provider')->comment('遊戲商代碼');
            $table->string('player_name')->comment('遊戲商玩家帳號');
            $table->foreignId('transfer_in_id')->nullable()->comment('本次遊戲期間起始轉入流水')->constrained('game_transfers')->nullOnDelete();
            $table->foreignId('transfer_out_id')->unique()->comment('本次接受稽核的轉出流水')->constrained('game_transfers')->cascadeOnDelete();
            $table->timestamp('period_started_at')->nullable()->comment('遊戲紀錄統計起始時間');
            $table->timestamp('period_ended_at')->comment('遊戲紀錄統計結束時間');
            $table->decimal('opening_balance', 16, 2)->nullable()->comment('轉入前遊戲錢包餘額');
            $table->decimal('total_transfer_in', 16, 2)->default(0)->comment('期間完成轉入合計');
            $table->decimal('prior_transfer_out', 16, 2)->default(0)->comment('本次以前期間完成轉出合計');
            $table->decimal('total_bet', 16, 2)->default(0)->comment('期間下注合計');
            $table->decimal('total_payout', 16, 2)->default(0)->comment('期間派彩合計');
            $table->decimal('total_adjustment', 16, 2)->default(0)->comment('已確認的其他資金調整合計');
            $table->decimal('expected_transfer_out', 16, 2)->nullable()->comment('依遊戲紀錄推算的應轉出金額');
            $table->decimal('actual_transfer_out', 16, 2)->comment('target 實際完成轉出金額');
            $table->decimal('difference', 16, 2)->nullable()->comment('實際轉出減應轉出差額');
            $table->unsignedInteger('record_count')->default(0)->comment('納入計算的遊戲紀錄筆數');
            $table->string('status')->default('pending_records')->comment('pending_records/matched/manual_review');
            $table->unsignedInteger('attempts')->default(0)->comment('稽核嘗試次數');
            $table->timestamp('next_retry_at')->nullable()->comment('下次延遲稽核時間');
            $table->timestamp('audited_at')->nullable()->comment('最後一次完成計算時間');
            $table->string('reason')->nullable()->comment('等待或人工處理原因');
            $table->json('meta')->nullable()->comment('稽核使用的補充資訊');
            $table->timestamps();

            $table->index(['status', 'next_retry_at']);
            $table->index(['user_id', 'provider', 'period_ended_at'], 'game_transfer_audits_user_period');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('game_transfer_audits');
    }
};
