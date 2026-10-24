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
            $table->string('vendor', 20);
            $table->string('username', 20);
            $table->integer('user_id')->nullable();
            $table->string('game_id', 50);
            $table->string('bet_id', 50);
            $table->string('bet_time', 45);
            $table->float('bet_result_amount');
            $table->float('bet_amount');
            $table->float('payout_amount');
            $table->float('valid_amount');
            $table->float('tax_amount');
            $table->string('extra_info', 255)->nullable();
            $table->timestamps();

            // Indexes
            $table->unique(['bet_id', 'vendor'], 'game_records_betid_vendor');
            $table->index(['vendor', 'username', 'bet_time'], 'game_records_vendor_user');
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
