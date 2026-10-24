<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateGameRecordsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('game_records', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->integer('gameID');
            $table->float('beforeAmount');
            $table->float('winLossAmount');
            $table->integer('roundID');
            $table->integer('betID');
            $table->dateTime('roundStartTime');
            $table->dateTime('roundEndTime');
            $table->string('transactionNo');
            $table->float('betAmount');
            $table->dateTime('betTime');
            $table->float('payoutAmount');
            $table->string('taxAmount');
            $table->float('validAmount');
            $table->timestamps();
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
}
