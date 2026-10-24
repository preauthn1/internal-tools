<?php

namespace App\Actions\Fortify;

use App\Concerns\PasswordValidationRules;
use App\Concerns\ProfileValidationRules;
use App\Models\User;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Validator;
use Laravel\Fortify\Contracts\CreatesNewUsers;

class CreateNewUser implements CreatesNewUsers
{
    use PasswordValidationRules, ProfileValidationRules;

    /**
     * Validate and create a newly registered user.
     *
     * @param  array<string, string>  $input
     */
    public function create(array $input): User
    {
        Validator::make($input, [
            'email' => $this->emailRules(),
            'account' => $this->accountRules(),
            'password' => $this->passwordRules(),
        ])->validate();

        return DB::transaction(fn (): User => User::create([
            // 顯示名稱預設為帳號，之後可於個人資料頁修改。
            // UserObserver 建立的錢包與初始帳本會一併納入此交易。
            'name' => $input['account'],
            'account' => $input['account'],
            'email' => $input['email'],
            'password' => $input['password'],
        ]));
    }
}
