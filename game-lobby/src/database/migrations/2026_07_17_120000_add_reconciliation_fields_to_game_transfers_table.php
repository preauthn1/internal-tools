<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::table('game_transfers', function (Blueprint $table) {
            $table->unsignedInteger('attempts')->default(0)->after('status')->comment('已嘗試次數（重試/對帳用）');
            $table->timestamp('last_attempt_at')->nullable()->after('attempts')->comment('最後一次嘗試時間');
            $table->timestamp('next_retry_at')->nullable()->after('last_attempt_at')->comment('下次重試時間');
            $table->timestamp('resolved_at')->nullable()->after('next_retry_at')->comment('結案時間（完成/失敗/轉人工）');
            $table->string('error_code')->nullable()->after('resolved_at')->comment('遊戲商錯誤碼或內部原因');
            $table->string('provider_transaction_no')->nullable()->after('error_code')->comment('遊戲商交易流水號');

            $table->index(['status', 'next_retry_at']);
            $table->index(['user_id', 'provider', 'direction', 'status'], 'game_transfers_unresolved_lookup');
        });
    }

    public function down(): void
    {
        Schema::table('game_transfers', function (Blueprint $table) {
            $table->dropIndex(['status', 'next_retry_at']);
            $table->dropIndex('game_transfers_unresolved_lookup');
            $table->dropColumn([
                'attempts',
                'last_attempt_at',
                'next_retry_at',
                'resolved_at',
                'error_code',
                'provider_transaction_no',
            ]);
        });
    }
};
