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
        Schema::create('logs', function (Blueprint $table) {
            $table->id();
            $table->unsignedSmallInteger('level')->comment('等級數值（Monolog：100 DEBUG ~ 600 EMERGENCY）');
            $table->string('level_name', 20)->comment('等級名稱（DEBUG/INFO/WARNING/ERROR…）');
            $table->string('channel')->comment('來源頻道');
            $table->text('message')->comment('訊息內容');
            $table->json('context')->nullable()->comment('上下文資料');
            $table->json('extra')->nullable()->comment('Monolog 額外資訊（processor 產生）');
            $table->timestamp('logged_at')->comment('記錄時間');
            $table->timestamps();

            $table->index('level_name');
            $table->index('channel');
            $table->index('logged_at');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('logs');
    }
};
