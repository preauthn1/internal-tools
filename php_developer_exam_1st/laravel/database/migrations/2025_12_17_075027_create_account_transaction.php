<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('account_transaction', function (Blueprint $table) {
            //?? remove id ?
            $table->integer('id')->autoIncrement()->primary();
            $table->string('transaction_no', 20)->unique()->index();
            $table->string('userid');
            $table->decimal('amount', 12, 2);
            $table->integer('type')->default(0)->comment('0: deposit to vendor, 1: withdraw from vendor');
            $table->integer('status')->default(0)->comment('0: pending - prepared, 1: success of status 0, 2: failed');
            //dont add vendor for exam1
            $table->string('vendor', 20);
            $table->string('vendor_transaction_no', 50)->unique()->nullable();
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('account_transaction');
    }
};
