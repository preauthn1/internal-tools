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
        Schema::create('wallet_transactions', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->comment('所屬會員')->constrained()->cascadeOnDelete();
            $table->string('type')->comment('類型：initial 初始金 / transfer_in 轉入遊戲 / transfer_out 從遊戲轉回 / refund 轉入失敗退回');
            $table->decimal('amount', 16, 2)->comment('有號變動額（轉入為負、轉出/退回為正）');
            $table->decimal('balance_after', 16, 2)->comment('異動後結餘');
            $table->string('reference')->nullable()->comment('關聯 game_transfers.reference（轉帳類）');
            $table->json('meta')->nullable()->comment('附註（回填標記等）');
            $table->timestamps();

            $table->index(['user_id', 'id']);
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('wallet_transactions');
    }
};
