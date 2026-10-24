<?php

namespace App\Console\Commands;

use App\Models\GameRecord;
use App\Models\User;
use App\Services\GameVendors\GameVendorFactory;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class FetchGameRecords extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'game-records:fetch
                            {--vendor= : The vendor ID (e.g., target)}
                            {--start_time= : Start time (format: YYYY-MM-DD HH:MM:SS)}
                            {--end_time= : End time (format: YYYY-MM-DD HH:MM:SS)}
                            {--page=1 : Page number}
                            {--page_size=10000 : Page size}
                            {--username= : Filter by username (optional)}
                            {--update : Update existing records instead of ignoring duplicates}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Fetch game records from a vendor service';

    /**
     * Execute the console command.
     */
    public function handle(): int
    {
        $vendorId = $this->option('vendor');
        $startTime = $this->option('start_time');
        $endTime = $this->option('end_time');
        $page = (int) $this->option('page');
        $pageSize = (int) $this->option('page_size');
        $username = $this->option('username') ?? '';

        // Validate required options
        if (empty($vendorId)) {
            $this->error("Error: --vendor option is required");
            $this->info("Usage: php artisan game-records:fetch --vendor=target --start_time=\"2025-12-01 00:00:00\" --end_time=\"2025-12-18 23:59:59\"");
            return Command::FAILURE;
        }

        if (empty($startTime)) {
            $this->error("Error: --start_time option is required");
            return Command::FAILURE;
        }

        if (empty($endTime)) {
            $this->error("Error: --end_time option is required");
            return Command::FAILURE;
        }

        try {
            // Validate vendorId
            if (!in_array($vendorId, GameVendorFactory::getAvailableVendors(), true)) {
                $this->error("Invalid vendor ID: {$vendorId}");
                $this->info("Available vendors: " . implode(', ', GameVendorFactory::getAvailableVendors()));
                return Command::FAILURE;
            }

            $this->info("Fetching game records for vendor: {$vendorId}");
            $this->info("Time range: {$startTime} to {$endTime}");

            $vendor = GameVendorFactory::make($vendorId);

            $filters = [
                'startTime' => $startTime,
                'endTime' => $endTime,
                'pageIndex' => $page,
                'pageSize' => $pageSize,
            ];

            if (!empty($username)) {
                $filters['username'] = $username;
            }

            // Check if game records exist in Redis cache using Laravel's built-in Cache facade
            $cacheKey = "game_records:{$vendorId}:" . md5($startTime . $endTime . $page . $pageSize . $username);
            $redisCache = Cache::store('redis');

            //OdinTodo : remove caching after testing
            //if ($redisCache->has($cacheKey)) {
                //// Use cached data if it exists
                //$gameRecords = $redisCache->get($cacheKey);
                //$this->info("Using cached game records from Redis (key: {$cacheKey})");
                //$this->info("Cached game records: " . print_r($gameRecords, true));
            //} else {
                // Fetch from vendor service if cache doesn't exist
                $this->info("Fetching game records from vendor service...");
                $gameRecords = $vendor->getGameRecords($filters);

                // Cache the fetched data for future use (Todo: remove this, for testing only)
                $redisCache->put($cacheKey, $gameRecords, now()->addHours(24)); // Cache for 24 hours
                $this->info("Game records cached in Redis with key: {$cacheKey}");
            //}

            // Display results
            $gameRecordsData = $gameRecords["items"];
            $count = $gameRecords["totalCount"];
            $this->info("Found {$count} game records");
            $vendorFieldUserName = $vendor->getVendorFieldUserName();

            if ($count === 0) {
                $this->warn("No game records found for the specified criteria");
                return Command::SUCCESS;
            }

            // Prepare data for batch insertion using Laravel's built-in methods
            $this->info("Preparing data for batch insertion...");
            $recordsToInsert = [];

            // Extract unique usernames from records (handle different key formats)
            $usernames_as_key = [];
            foreach ($gameRecordsData as $record) {
                if (!empty($record[$vendorFieldUserName])) {
                    $usernames_as_key[$record[$vendorFieldUserName]] = 0;
                }
            }
            $usernames = array_keys($usernames_as_key);
            $this->info("Unique usernames extracted: " . count($usernames));

            // Get user IDs mapping for batch lookup using Laravel's built-in pluck method
            $userMapping = [];
            if (!empty($usernames)) {
                $userMapping = User::whereIn('name', $usernames)
                    ->pluck('id', 'name')
                    ->toArray();
            }

            foreach ($gameRecordsData as $record) {
                // Map vendor response to database columns
                $recordsToInsert[] = [
                    'vendor' => $vendorId,
                    'username' => $record[$vendorFieldUserName],
                    'user_id' => $userMapping[$record[$vendorFieldUserName]] ?? null,
                    'game_id' => $record['gameID'],
                    'bet_id' => $record['betID'],
                    'bet_time' => date("Y-m-d H:i:s", strtotime($record['betTime'])),
                    'bet_result_amount' => (float) $record['winLossAmount'],
                    'bet_amount' => (float) $record['betAmount'],
                    'payout_amount' => (float) $record['payoutAmount'],
                    'valid_amount' => (float) $record['validAmount'],
                    'tax_amount' => (float) $record['taxAmount'],
                    'extra_info' => json_encode($record['extraInfo']),
                ];
            }

            // Batch insert using Laravel's built-in methods
            $this->info("Inserting {$count} records into database...");

            // Use Laravel's chunk method for large batches to avoid memory issues
            $chunks = array_chunk($recordsToInsert, 500); // Process 500 records at a time
            $inserted = 0;
            $updated = 0;

            foreach ($chunks as $chunk) {
                if ($this->option('update')) {
                    DB::table('game_records')->upsert(
                        $chunk,
                        ['bet_id', 'vendor'],
                        [
                            'username',
                            'user_id',
                            'game_id',
                            'bet_time',
                            'bet_result_amount',
                            'bet_amount',
                            'payout_amount',
                            'valid_amount',
                            'tax_amount',
                            'extra_info',
                        ]
                    );
                    $inserted += count($chunk);
                    $updated += count($chunk);
                } else {
                    $inserted += DB::table('game_records')->insertOrIgnore($chunk);
                }
            }

            if ($this->option('update')) {
                $this->info("Successfully processed {$inserted} game records (inserted or updated)");
            } else {
                $this->info("Successfully inserted {$inserted} game records");
            }

            // Display sample results
            $this->table(
                ['Bet ID', 'Username', 'Bet Time', 'Bet Amount', 'Payout Amount', 'Result Amount'],
                array_map(function ($record) {
                    return [
                        $record['bet_id'] ?? 'N/A',
                        $record['username'] ?? 'N/A',
                        $record['bet_time'] ?? 'N/A',
                        number_format($record['bet_amount'] ?? 0, 2),
                        number_format($record['payout_amount'] ?? 0, 2),
                        number_format($record['bet_result_amount'] ?? 0, 2),
                    ];
                }, array_slice($recordsToInsert, 0, 10)) // Show first 10 records
            );

            if ($count > 10) {
                $this->info("... and " . ($count - 10) . " more records");
            }

            $this->info("Command completed successfully");
            return Command::SUCCESS;

        } catch (\Exception $e) {
            $this->error("Error fetching game records: " . $e->getMessage());
            Log::error("FetchGameRecords command failed", [
                'vendorId' => $vendorId,
                'startTime' => $startTime,
                'endTime' => $endTime,
                'error' => $e->getMessage(),
                'trace' => $e->getTraceAsString(),
            ]);
            return Command::FAILURE;
        }
    }
}
