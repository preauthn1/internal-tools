<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::table('wallet_transactions', function (Blueprint $table) {
            // 外鍵關聯 game_providers（轉帳類帳本才有值；初始金為 null）
            $table->foreignId('game_provider_id')->nullable()->after('type')
                ->comment('對應的遊戲商（轉帳類才有；初始金為 null）')
                ->constrained()->nullOnDelete();
        });

        // 回填既有帳本的遊戲商：reference → game_transfers.provider（key）→ game_providers.id
        DB::table('wallet_transactions')
            ->whereNull('game_provider_id')
            ->whereNotNull('reference')
            ->update([
                'game_provider_id' => DB::raw(
                    '(select gp.id from game_transfers gt'
                    .' join game_providers gp on gp.`key` = gt.provider'
                    .' where gt.reference = wallet_transactions.reference limit 1)'
                ),
            ]);
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('wallet_transactions', function (Blueprint $table) {
            $table->dropConstrainedForeignId('game_provider_id');
        });
    }
};
