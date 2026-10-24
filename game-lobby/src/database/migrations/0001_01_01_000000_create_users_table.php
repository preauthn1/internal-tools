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
        Schema::create('users', function (Blueprint $table) {
            $table->id();
            $table->string('name')->comment('顯示名稱（註冊時預設為 account）');
            $table->string('account')->unique()->comment('登入帳號（唯一）');
            $table->string('email')->unique()->comment('電子郵件（唯一）');
            $table->timestamp('email_verified_at')->nullable()->comment('信箱驗證時間');
            $table->string('password')->comment('密碼雜湊');
            $table->rememberToken();
            $table->timestamps();
        });

        Schema::create('password_reset_tokens', function (Blueprint $table) {
            $table->string('email')->primary()->comment('電子郵件（主鍵）');
            $table->string('token')->comment('重設密碼權杖');
            $table->timestamp('created_at')->nullable()->comment('建立時間');
        });

        // session / cache / queue 皆改用 Redis，不再建立對應資料表。
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('users');
        Schema::dropIfExists('password_reset_tokens');
    }
};
