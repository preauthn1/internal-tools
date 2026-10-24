<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class DailyQuote extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'quote:daily';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'test log';

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {

        return Command::SUCCESS;
    }
}