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
        // Schema::create('game_lists', function (Blueprint $table) {
        //     $table->id();
        //     $table->integer('gameID');
        //     $table->string('gameName');
        //     $table->string('imageUrl');
        //     $table->integer('gameType');
        //     $table->string('gameTypeName');
        //     $table->string('isPrePaidMode');
        //     $table->timestamps();
        // });
        Schema::create('game_lists', function (Blueprint $table) {
            $table->id();
            $table->integer('game_id');
            $table->string('game_name');
            $table->string('image_url');
            $table->integer('game_type');
            $table->string('game_type_name');
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
        Schema::dropIfExists('game_lists');
    }
};