<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\MassPrunable;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Carbon;

/**
 * 應用程式日誌（由 App\Logging\DatabaseLogHandler 寫入的 database log channel）。
 *
 * @property int $id
 * @property int $level
 * @property string $level_name
 * @property string $channel
 * @property string $message
 * @property array<string, mixed>|null $context
 * @property array<string, mixed>|null $extra
 * @property Carbon $logged_at
 */
class LogEntry extends Model
{
    use MassPrunable;

    protected $table = 'logs';

    protected $fillable = [
        'level',
        'level_name',
        'channel',
        'message',
        'context',
        'extra',
        'logged_at',
    ];

    protected function casts(): array
    {
        return [
            'level' => 'integer',
            'context' => 'array',
            'extra' => 'array',
            'logged_at' => 'datetime',
        ];
    }

    /** @return Builder<static> */
    public function prunable(): Builder
    {
        $retentionDays = max(1, (int) config('logging.channels.database.retention_days', 30));

        return static::query()->where('logged_at', '<', now()->subDays($retentionDays));
    }
}
