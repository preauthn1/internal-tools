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
        // Schema::create('money_book', function (Blueprint $table) {
        //     $table->id();
        //     $table->string('user');
        //     $table->string('player');
        //     $table->string('merchantOrderNo');
        //     $table->string('transactionNo');
        //     $table->double('amount', 30, 2);
        //     $table->double('beforeAmount', 30, 2);
        //     $table->double('afterAmount', 30, 2);
        //     $table->timestamps();
        // });
        Schema::create('money_books', function (Blueprint $table) {
            $table->id();
            $table->integer('user_id');
            $table->string('status');
            $table->string('user');
            $table->string('player');
            $table->string('transfer_type');
            $table->string('merchant_order_no');
            $table->string('transaction_no');
            $table->double('amount', 30, 2);
            $table->double('before_amount', 30, 2);
            $table->double('after_amount', 30, 2);
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
        Schema::dropIfExists('money_book');
    }
};