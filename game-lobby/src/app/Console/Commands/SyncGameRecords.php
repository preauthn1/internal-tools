<?php

namespace App\Console\Commands;

use App\Models\GameProviderAccount;
use App\Models\GameRecord;
use App\Services\Games\Contracts\GameProvider;
use App\Services\Games\Data\GameRecordData;
use App\Services\Games\Exceptions\GameProviderException;
use App\Services\Games\GameProviderManager;
use Carbon\CarbonImmutable;
use Illuminate\Console\Attributes\Description;
use Illuminate\Console\Attributes\Signature;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Log;

#[Signature('game:sync-records {--provider= : 指定 provider} {--since= : 起始 YYYY-MM-DD} {--until= : 結束 YYYY-MM-D }')]
#[Description('同步遊戲紀錄')]
class SyncGameRecords extends Command
{
    public function handle(GameProviderManager $manager): int
    {
        $providers = $this->option('provider')
            ? [$manager->provider((string) $this->option('provider'))]
            : $manager->enabled();

        foreach ($providers as $provider) {
            $this->syncProvider($provider);
        }

        return self::SUCCESS;
    }

    private function syncProvider(GameProvider $provider): void
    {
        $key = $provider->key();
        $until = $this->option('until')
            ? CarbonImmutable::parse((string) $this->option('until'))->utc()
            : CarbonImmutable::now('UTC');
        $lookback = (int) $provider->setting('records_lookback_minutes', 120);
        $since = $this->option('since')
            ? CarbonImmutable::parse((string) $this->option('since'))->utc()
            : CarbonImmutable::parse($this->cursor($key) ?? $until)->utc()->subMinutes($lookback);

        if ($since->greaterThanOrEqualTo($until)) {
            $this->info("[{$key}] 無需同步：since >= until");

            return;
        }

        $stats = ['fetched' => 0, 'matched' => 0, 'upserted' => 0, 'skipped' => 0];
        try {
            foreach ($this->windows($provider, $since, $until) as [$start, $end]) {
                $this->syncWindow($provider, $start, $end, $stats);
            }
            $this->setCursor($key, $until);
        } catch (GameProviderException $exception) {
            Log::warning('遊戲紀錄同步失敗', ['provider' => $key, 'api_code' => $exception->apiCode, 'category' => 'provider']);
            $this->error("[{$key}] 同步失敗，游標未更新");
        }

        $this->info("[{$key}] fetched={$stats['fetched']} matched={$stats['matched']} upserted={$stats['upserted']} skipped={$stats['skipped']}");
    }

    /** @return iterable<array{0: CarbonImmutable, 1: CarbonImmutable}> */
    private function windows(GameProvider $provider, CarbonImmutable $since, CarbonImmutable $until): iterable
    {
        $window = (int) $provider->setting('records_window_minutes', 60);
        $start = $since;
        while ($start->lessThan($until)) {
            $end = $start->addMinutes($window);
            $end = $end->greaterThan($until) ? $until : $end;
            yield [$start, $end];
            $start = $end;
        }
    }

    /** @param array<string, int> $stats */
    private function syncWindow(GameProvider $provider, CarbonImmutable $start, CarbonImmutable $end, array &$stats): void
    {
        $key = $provider->key();
        $pageSize = (int) $provider->setting('records_page_size', 200);
        $page = 1;
        do {
            $result = $provider->gameRecords([
                'startTime' => $start->format('Y-m-d\TH:i:s\Z'),
                'endTime' => $end->format('Y-m-d\TH:i:s\Z'),
                'pageIndex' => $page,
                'pageSize' => $pageSize,
            ]);
            $this->storeRecords($key, $result['items'], $stats);
            $page++;
        } while ($page <= $result['totalPage']);
    }

    /**
     * @param  array<int, GameRecordData>  $records
     * @param  array<string, int>  $stats
     */
    private function storeRecords(string $key, array $records, array &$stats): void
    {
        $stats['fetched'] += count($records);
        $playerNames = array_values(array_unique(array_map(fn (GameRecordData $record): string => $record->playerName, $records)));
        $bindings = GameProviderAccount::query()->where('provider', $key)->whereIn('player_name', $playerNames)
            ->pluck('user_id', 'player_name');
        $now = now();
        $rows = [];

        foreach ($records as $record) {
            $userId = $bindings->get($record->playerName);
            if ($record->recordId === '' || $userId === null) {
                $stats['skipped']++;

                continue;
            }
            $attributes = $record->toAttributes();
            $attributes['raw'] = json_encode($attributes['raw'], JSON_THROW_ON_ERROR);
            // upsert 不會套用 Eloquent datetime cast，需自行把 ISO8601（帶 T/Z）正規化為 DB datetime。
            $attributes['bet_at'] = $attributes['bet_at'] !== null
                ? CarbonImmutable::parse($attributes['bet_at'])->utc()->format('Y-m-d H:i:s')
                : null;
            $rows[] = ['user_id' => $userId, ...$attributes, 'created_at' => $now, 'updated_at' => $now];
            $stats['matched']++;
        }

        if ($rows === []) {
            return;
        }

        GameRecord::upsert($rows, ['provider', 'record_id'], [
            'user_id',
            'player_name',
            'game_id',
            'game_code',
            'game_name',
            'bet_amount',
            'payout_amount',
            'valid_amount',
            'currency',
            'bet_at',
            'raw',
            'updated_at',
        ]);
        $stats['upserted'] += count($rows);
    }

    private function cursor(string $key): ?string
    {
        $value = Cache::get("game-records:cursor:{$key}");

        return is_string($value) ? $value : null;
    }

    private function setCursor(string $key, CarbonImmutable $until): void
    {
        Cache::forever("game-records:cursor:{$key}", $until->format('Y-m-d\TH:i:s\Z'));
    }
}
