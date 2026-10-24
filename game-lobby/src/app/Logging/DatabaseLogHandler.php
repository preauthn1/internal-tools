<?php

namespace App\Logging;

use App\Models\LogEntry;
use Monolog\Handler\AbstractProcessingHandler;
use Monolog\Level;
use Monolog\LogRecord;
use Throwable;

/**
 * 將 log 記錄寫入資料庫（logs 資料表）的 Monolog handler。
 * 於 config/logging.php 的 database channel 使用。
 */
class DatabaseLogHandler extends AbstractProcessingHandler
{
    public function __construct(int|string|Level $level = Level::Debug, bool $bubble = true)
    {
        parent::__construct($level, $bubble);
    }

    protected function write(LogRecord $record): void
    {
        // 寫入失敗時靜默略過，避免「記錄失敗又觸發記錄」造成遞迴。
        try {
            LogEntry::create([
                'level' => $record->level->value,
                'level_name' => $record->level->getName(),
                'channel' => $record->channel,
                'message' => $record->message,
                'context' => $record->context ?: null,
                'extra' => $record->extra ?: null,
                'logged_at' => $record->datetime,
            ]);
        } catch (Throwable) {
            // 忽略：資料庫不可用時仍應讓其他 log channel（檔案）正常運作。
        }
    }
}
