<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::table('game_records', function (Blueprint $table) {
            $table->index(
                ['user_id', 'provider', 'player_name', 'bet_at'],
                'game_records_audit_lookup',
            );
            $table->index(
                ['user_id', 'bet_at', 'id'],
                'game_records_user_history',
            );
        });
    }

    public function down(): void
    {
        Schema::table('game_records', function (Blueprint $table) {
            $table->dropIndex('game_records_audit_lookup');
            $table->dropIndex('game_records_user_history');
        });
    }
};
