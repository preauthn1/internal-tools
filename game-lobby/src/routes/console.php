<?php

use App\Models\LogEntry;
use Illuminate\Foundation\Inspiring;
use Illuminate\Support\Facades\Artisan;
use Illuminate\Support\Facades\Schedule;

Artisan::command('inspire', function () {
    $this->comment(Inspiring::quote());
})->purpose('Display an inspiring quote');

Schedule::command('game:sync-records')->everyFiveMinutes()->withoutOverlapping();
Schedule::command('game:reconcile-transfers')->everyMinute()->withoutOverlapping();
Schedule::command('game:audit-wallet-sessions')->everyFiveMinutes()->withoutOverlapping();
Schedule::command('model:prune', ['--model' => [LogEntry::class]])
    ->daily()
    ->withoutOverlapping();
