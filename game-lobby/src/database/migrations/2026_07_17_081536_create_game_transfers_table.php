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
        Schema::create('game_transfers', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->comment('所屬會員')->constrained()->cascadeOnDelete();
            $table->string('provider')->comment('遊戲商代號，如 target');
            $table->string('direction')->comment('轉帳方向：in 轉入遊戲 / out 轉回平台');
            $table->decimal('amount', 16, 2)->comment('轉帳金額');
            $table->string('reference')->unique()->comment('送遊戲商的流水號 merchantOrderNo（唯一、冪等）');
            $table->string('status')->default('pending')->comment('狀態：pending/unknown/completed/failed/manual_review');
            $table->json('meta')->nullable()->comment('遊戲商回應與 player_name 等附註');
            $table->timestamps();

            $table->index(['user_id', 'provider']);
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('game_transfers');
    }
};
