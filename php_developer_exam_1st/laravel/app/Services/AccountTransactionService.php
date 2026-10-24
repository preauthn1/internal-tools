<?php

namespace App\Services;

use App\Models\AccountTransaction;
use App\Models\User;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class AccountTransactionService
{
    public function __construct( protected AccountTransaction $accountTransaction)
    {

    }

    /**
     * Get all transactions for a user
     */
    public function getByUser(User $user, array $filters = []): LengthAwarePaginator
    {
        $query = $this->accountTransaction->where('userid', $user->id);

        if (isset($filters['vendor'])) {
            $query->where('vendor', $filters['vendor']);
        }

        if (isset($filters['type'])) {
            $query->where('type', $filters['type']);
        }

        if (isset($filters['status'])) {
            $query->where('status', $filters['status']);
        }

        if (isset($filters['date_from'])) {
            $query->whereDate('created_at', '>=', $filters['date_from']);
        }

        if (isset($filters['date_to'])) {
            $query->whereDate('created_at', '<=', $filters['date_to']);
        }

        return $query->with('user')
            ->latest()
            ->paginate($filters['per_page'] ?? 15);
    }

    /**
     * Get a single transaction by ID
     */
    public function getById(int $id): ?AccountTransaction
    {
        return $this->accountTransaction->with('user')->find($id);
    }

    public function createPendingDepositTransaction(int $userId, string $vendor, float $amount): AccountTransaction
    {
        $transaction = new AccountTransaction();
        $transaction_no = $this->generateTransactionNo();
        $transaction->transaction_no = $this->generateTransactionNo();
        $transaction->userid = $userId;
        $transaction->vendor = $vendor;
        $transaction->type = AccountTransaction::TYPE_DEPOSIT;
        $transaction->amount = $amount;
        $transaction->status = AccountTransaction::STATUS_PENDING;
        $transaction->save();

        return $transaction;
    }

    /**
     * Create a pending withdraw transaction using Laravel's built-in Eloquent methods
     */
    public function createPendingWithdrawTransaction(int $userId, string $vendor, float $amount): AccountTransaction
    {
        $transaction = new AccountTransaction();
        $transaction->transaction_no = $this->generateTransactionNo();
        $transaction->userid = $userId;
        $transaction->vendor = $vendor;
        $transaction->type = AccountTransaction::TYPE_WITHDRAW;
        $transaction->amount = $amount;
        $transaction->status = AccountTransaction::STATUS_PENDING;
        $transaction->save();

        return $transaction;
    }

    /**
     * Generate transaction number based on the actual inserted ID
     * Uses the model's ID property (more reliable than PDO lastInsertId)
     */
    public function generateTransactionNo(): string
    {
        $lastInsertId = DB::getPdo()->lastInsertId();

        return sprintf('TN%s%04d', date("YmdHis"), $lastInsertId + 1);
    }

    /**
     * Get transaction by transaction ID (vendor's transaction ID)
     */
    public function getByTransactionId(string $transactionId): ?AccountTransaction
    {
        return $this->accountTransaction->with('user')
            ->where('transaction_id', $transactionId)
            ->first();
    }

    /**
     * Get transactions by vendor
     */
    //public function getByVendor(string $vendorId, array $filters = []): LengthAwarePaginator
    //{
        //$query = $this->accountTransaction->where('vendor_id', $vendorId);

        //if (isset($filters['user_id'])) {
            //$query->where('user_id', $filters['user_id']);
        //}

        //if (isset($filters['type'])) {
            //$query->where('type', $filters['type']);
        //}

        //if (isset($filters['status'])) {
            //$query->where('status', $filters['status']);
        //}

        //return $query->with('user')
            //->latest()
            //->paginate($filters['per_page'] ?? 15);
    //}

    /**
     * Get user's transaction summary
     */
    //public function getSummary(User $user, array $filters = []): array
    //{
        //$query = $this->accountTransaction->where('user_id', $user->id);

        //if (isset($filters['vendor_id'])) {
            //$query->where('vendor_id', $filters['vendor_id']);
        //}

        //if (isset($filters['date_from'])) {
            //$query->whereDate('created_at', '>=', $filters['date_from']);
        //}

        //if (isset($filters['date_to'])) {
            //$query->whereDate('created_at', '<=', $filters['date_to']);
        //}

        //$totalDeposits = (clone $query)->where('type', 'deposit')
            //->where('status', 'completed')
            //->sum('amount');

        //$totalWithdraws = (clone $query)->where('type', 'withdraw')
            //->where('status', 'completed')
            //->sum('amount');

        //$pendingCount = (clone $query)->where('status', 'pending')->count();
        //$failedCount = (clone $query)->where('status', 'failed')->count();

        //return [
            //'total_deposits' => (float) $totalDeposits,
            //'total_withdraws' => (float) $totalWithdraws,
            //'net_amount' => (float) ($totalDeposits - $totalWithdraws),
            //'pending_count' => $pendingCount,
            //'failed_count' => $failedCount,
            //'total_transactions' => $query->count(),
        //];
    //}

    /**
     * Get recent transactions for a user
     */
    //public function getRecent(User $user, int $limit = 10): Collection
    //{
        //return $this->accountTransaction->where('user_id', $user->id)
            //->with('user')
            //->latest()
            //->limit($limit)
            //->get();
    //}
}
