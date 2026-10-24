<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateTransferPointsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('transfer_points', function (Blueprint $table) {
            $table->id();
            $table->float('amount');
            $table->string('merchantOrderNo');
            $table->string('name');
            $table->string('transactionNo');
            $table->float('beforeAmount');
            $table->float('afterAmount');
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
        Schema::dropIfExists('transfer_point');
    }
}
