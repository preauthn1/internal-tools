<?php

namespace Tests\Feature;

use App\Models\LogEntry;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Log;
use Tests\TestCase;

class DatabaseLoggingTest extends TestCase
{
    use RefreshDatabase;

    public function test_log_is_written_to_the_database(): void
    {
        Log::channel('database')->warning('遊戲轉帳遭拒', ['reference' => 'abc-123']);

        $this->assertDatabaseCount('logs', 1);

        $entry = LogEntry::firstOrFail();
        $this->assertSame('WARNING', $entry->level_name);
        $this->assertSame(300, $entry->level);
        $this->assertSame('遊戲轉帳遭拒', $entry->message);
        $this->assertSame('abc-123', $entry->context['reference']);
        $this->assertNotNull($entry->logged_at);
    }

    public function test_error_level_is_written(): void
    {
        Log::channel('database')->error('boom');

        $this->assertDatabaseHas('logs', ['level_name' => 'ERROR', 'message' => 'boom']);
    }

    public function test_model_prune_removes_only_expired_database_logs(): void
    {
        config()->set('logging.channels.database.retention_days', 30);
        LogEntry::create([
            'level' => 300,
            'level_name' => 'WARNING',
            'channel' => 'testing',
            'message' => 'expired',
            'logged_at' => now()->subDays(31),
        ]);
        LogEntry::create([
            'level' => 300,
            'level_name' => 'WARNING',
            'channel' => 'testing',
            'message' => 'active',
            'logged_at' => now()->subDays(29),
        ]);

        $this->artisan('model:prune', ['--model' => [LogEntry::class]])
            ->assertSuccessful();

        $this->assertDatabaseMissing('logs', ['message' => 'expired']);
        $this->assertDatabaseHas('logs', ['message' => 'active']);
    }
}
