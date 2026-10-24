<?php

namespace App\Services;

use App\Models\GameRecord;
use App\Models\User;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Support\Facades\Log;

class GameRecordService
{
    public function __construct(protected GameRecord $gameRecord)
    {
    }

    /**
     * Get game records with filters using Laravel's built-in Eloquent methods
     *
     * @param array $filters
     * @return LengthAwarePaginator
     */
    public function getRecords(array $filters = []): LengthAwarePaginator
    {
        $query = $this->gameRecord->newQuery();

        // Filter by vendor
        if (isset($filters['vendor'])) {
            $query->where('vendor', $filters['vendor']);
        }

        // Filter by username
        if (isset($filters['username'])) {
            $query->where('username', $filters['username']);
        }

        // Filter by user_id
        if (isset($filters['user_id'])) {
            $query->where('user_id', $filters['user_id']);
        }

        // Filter by bet_id
        if (isset($filters['bet_id'])) {
            $query->where('bet_id', $filters['bet_id']);
        }

        // Filter by game_id
        if (isset($filters['game_id'])) {
            $query->where('game_id', $filters['game_id']);
        }

        // Filter by date range (bet_time)
        if (isset($filters['start_time'])) {
            $query->where('bet_time', '>=', $filters['start_time']);
        }

        if (isset($filters['end_time'])) {
            $query->where('bet_time', '<=', $filters['end_time']);
        }

        // Filter by bet amount range
        if (isset($filters['min_bet_amount'])) {
            $query->where('bet_amount', '>=', $filters['min_bet_amount']);
        }

        if (isset($filters['max_bet_amount'])) {
            $query->where('bet_amount', '<=', $filters['max_bet_amount']);
        }

        // Order by bet_time (latest first)
        $query->orderBy('bet_time', 'desc');

        // Eager load user relationship using Laravel's built-in with() method
        $query->with('user');

        return $query->paginate($filters['per_page'] ?? 15);
    }

    /**
     * Get game records by vendor using Laravel's built-in Eloquent methods
     *
     * @param string $vendor
     * @param array $filters
     * @return LengthAwarePaginator
     */
    public function getByVendor(string $vendor, array $filters = []): LengthAwarePaginator
    {
        $filters['vendor'] = $vendor;
        return $this->getRecords($filters);
    }

    /**
     * Get game records by user using Laravel's built-in Eloquent methods
     * Filters by user_id if available, otherwise by username
     *
     * @param User $user
     * @param array $filters
     * @return LengthAwarePaginator
     */
    public function getByUser(User $user, array $filters = []): LengthAwarePaginator
    {
        // Filter by user_id if available, otherwise by username
        $query = $this->gameRecord->newQuery();

        // Apply vendor filter if provided
        if (isset($filters['vendor'])) {
            $query->where('vendor', $filters['vendor']);
        }

        // Filter by user_id or username
        $query->where(function ($q) use ($user) {
            $q->where('user_id', $user->id)
              ->orWhere('username', $user->name);
        });

        // Apply other filters
        if (isset($filters['start_time'])) {
            $query->where('bet_time', '>=', $filters['start_time']);
        }

        if (isset($filters['end_time'])) {
            $query->where('bet_time', '<=', $filters['end_time']);
        }

        if (isset($filters['min_bet_amount'])) {
            $query->where('bet_amount', '>=', $filters['min_bet_amount']);
        }

        if (isset($filters['max_bet_amount'])) {
            $query->where('bet_amount', '<=', $filters['max_bet_amount']);
        }

        // Order by bet_time (latest first)
        $query->orderBy('bet_time', 'desc');

        // Eager load user relationship if exists
        $query->with('user');

        return $query->paginate($filters['per_page'] ?? 15);
    }

    /**
     * Get game records by username using Laravel's built-in Eloquent methods
     *
     * @param string $username
     * @param array $filters
     * @return LengthAwarePaginator
     */
    public function getByUsername(string $username, array $filters = []): LengthAwarePaginator
    {
        $filters['username'] = $username;
        return $this->getRecords($filters);
    }

    /**
     * Get a single game record by ID using Laravel's built-in Eloquent methods
     *
     * @param int $id
     * @return GameRecord|null
     */
    public function getById(int $id): ?GameRecord
    {
        return $this->gameRecord->with('user')->find($id);
    }

    /**
     * Get a game record by bet_id and vendor using Laravel's built-in Eloquent methods
     *
     * @param string $betId
     * @param string $vendor
     * @return GameRecord|null
     */
    public function getByBetId(string $betId, string $vendor): ?GameRecord
    {
        return $this->gameRecord->where('bet_id', $betId)
            ->where('vendor', $vendor)
            ->with('user')
            ->first();
    }

    /**
     * Get game records statistics using Laravel's built-in Eloquent aggregation methods
     *
     * @param array $filters
     * @return array
     */
    public function getStatistics(array $filters = []): array
    {
        $query = $this->gameRecord->newQuery();

        // Apply same filters as getRecords
        if (isset($filters['vendor'])) {
            $query->where('vendor', $filters['vendor']);
        }

        if (isset($filters['username'])) {
            $query->where('username', $filters['username']);
        }

        if (isset($filters['user_id'])) {
            $query->where('user_id', $filters['user_id']);
        }

        if (isset($filters['start_time'])) {
            $query->where('bet_time', '>=', $filters['start_time']);
        }

        if (isset($filters['end_time'])) {
            $query->where('bet_time', '<=', $filters['end_time']);
        }

        // Use Laravel's built-in aggregation methods
        return [
            'total_records' => $query->count(),
            'total_bet_amount' => (float) $query->sum('bet_amount'),
            'total_payout_amount' => (float) $query->sum('payout_amount'),
            'total_valid_amount' => (float) $query->sum('valid_amount'),
            'total_tax_amount' => (float) $query->sum('tax_amount'),
            'total_result_amount' => (float) $query->sum('bet_result_amount'),
            'average_bet_amount' => (float) $query->avg('bet_amount'),
        ];
    }

    /**
     * Get recent game records using Laravel's built-in Eloquent methods
     *
     * @param int $limit
     * @param array $filters
     * @return Collection
     */
    public function getRecent(int $limit = 10, array $filters = []): Collection
    {
        $query = $this->gameRecord->newQuery();

        if (isset($filters['vendor'])) {
            $query->where('vendor', $filters['vendor']);
        }

        if (isset($filters['username'])) {
            $query->where('username', $filters['username']);
        }

        return $query->with('user')
            ->orderBy('bet_time', 'desc')
            ->limit($limit)
            ->get();
    }
}

