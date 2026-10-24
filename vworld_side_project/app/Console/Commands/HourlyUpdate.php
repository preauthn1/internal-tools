<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

use Illuminate\Support\Facades\DB;

class HourlyUpdate extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'hour:update';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Add money';

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        // $user = User::all();
        // foreach ($user as $a) {
        //     Mail::raw("This is automatically generated Hourly Update", function ($message) use ($a) {
        //         $message->from('user@example.invalid');
        //         $message->to($a->email)->subject('Hourly Update');
        //     });
        // }
        // $this->info('Hourly Update has been send successfully');
        DB::update("UPDATE `users` SET `money` = `money` +'10000' WHERE `users`.`id` = 27");
        return Command::SUCCESS;
        // return 'Successfully sent daily quote to everyone.';
        // $this->info('Successfully sent daily quote to everyone.');
    }
}