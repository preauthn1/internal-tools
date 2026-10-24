<?php

namespace App\Services;

use App\Models\User;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Log;

class UserService
{
    public function __construct(protected User $user)
    {
    }

    public function subBalance(User $user, float $amount, bool $checkSufficientBalance = true): float
    {
        if ($amount <= 0) {
            throw new \InvalidArgumentException('Amount must be greater than zero');
        }

        if ($checkSufficientBalance && $user->balance < $amount) {
            throw new \RuntimeException(
                "Insufficient balance. Current balance: {$user->balance}, Required: {$amount}"
            );
        }

        // Use database-level decrement to avoid race conditions
        $user->decrement('balance', $amount);

        // Refresh to get the updated balance
        $user->refresh();

        Log::debug('User balance subtracted', [
            'user_id' => $user->id,
            'amount' => $amount,
            'new_balance' => $user->balance,
        ]);

        return (float) $user->balance;
    }

    /**
     * Add amount to user's balance
     *
     * @param User $user The user whose balance to update
     * @param float $amount The amount to add
     * @return float The updated balance after addition
     * @throws \InvalidArgumentException If amount is invalid
     */
    public function addBalance(User $user, float $amount): float
    {
        if ($amount <= 0) {
            throw new \InvalidArgumentException('Amount must be greater than zero');
        }

        // Use database-level increment to avoid race conditions
        $user->increment('balance', $amount);

        // Refresh to get the updated balance
        $user->refresh();

        Log::info('User balance added', [
            'user_id' => $user->id,
            'amount' => $amount,
            'new_balance' => $user->balance,
        ]);

        return (float) $user->balance;
    }

    public function hasSufficientBalance(User $user, float $amount): bool
    {
        return $user->balance >= $amount;
    }

    /**
     * Create a new user with optional signup bonus
     *
     * @param array $data User data (name, email, password, etc.)
     * @param float|null $signupBonus Optional signup bonus amount (defaults to config value)
     * @return User The created user
     */
    public function createUser(array $data, ?float $signupBonus = null): User
    {
        $userData = [
            'name' => $data['name'],
            'email' => $data['email'],
            'password' => Hash::make($data['password']),
        ];

        // Add signup bonus if provided or from config
        if ($signupBonus !== null) {
            $userData['balance'] = $signupBonus;
        } elseif (config('user.signup_bonus')) {
            $userData['balance'] = (float) config('user.signup_bonus');
        } else {
            $userData['balance'] = 0;
        }

        $user = $this->user->create($userData);

        Log::info('User created', [
            'user_id' => $user->id,
            'email' => $user->email,
            'signup_bonus' => $userData['balance'],
        ]);

        return $user;
    }

    /**
     * Update user profile information
     *
     * @param User $user The user to update
     * @param array $data Data to update (name, email, etc.)
     * @return User The updated user
     */
    //public function updateProfile(User $user, array $data): User
    //{
        //$user->update($data);
        //$user->refresh();

        //Log::info('User profile updated', [
            //'user_id' => $user->id,
            //'updated_fields' => array_keys($data),
        //]);

        //return $user;
    //}

    /**
     * Update user password
     *
     * @param User $user The user whose password to update
     * @param string $newPassword The new password (will be hashed)
     * @return User The updated user
     */
    //public function updatePassword(User $user, string $newPassword): User
    //{
        //$user->update([
            //'password' => Hash::make($newPassword),
        //]);

        //Log::info('User password updated', [
            //'user_id' => $user->id,
        //]);

        //return $user;
    //}
}
