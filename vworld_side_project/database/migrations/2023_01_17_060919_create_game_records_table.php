<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration {
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        // Schema::create('game_records', function (Blueprint $table) {
        //     $table->id();
        //     $table->string('playerName');
        //     $table->integer('gameID');
        //     $table->double('beforeAmount', 30, 2);
        //     $table->double('winLossAmount', 30, 2);
        //     $table->string('roundID');
        //     $table->string('betID');
        //     $table->string('betTime');
        //     $table->string('roundStartTime');
        //     $table->string('roundEndTime');
        //     $table->string('transactionNo');
        //     $table->double('betAmount', 30, 2);
        //     $table->double('payoutAmount', 30, 2);
        //     $table->double('taxAmount', 30, 2);
        //     $table->double('validAmount', 30, 2);
        //     $table->timestamps();
        // });
        Schema::create('game_records', function (Blueprint $table) {
            $table->id();
            $table->string('player_name');
            $table->integer('game_id');
            $table->double('before_amount', 30, 2);
            $table->double('win_loss_amount', 30, 2);
            $table->string('round_id');
            $table->string('bet_id');
            $table->dateTime('bet_time');
            $table->dateTime('round_start_time');
            $table->dateTime('round_end_time');
            $table->string('transaction_no');
            $table->double('bet_amount', 30, 2);
            $table->double('payout_amount', 30, 2);
            $table->double('tax_amount', 30, 2);
            $table->double('valid_amount', 30, 2);
            $table->dateTime('created_at')->useCurrent();
            $table->dateTime('updated_at')->useCurrent();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('game_records');
    }
};