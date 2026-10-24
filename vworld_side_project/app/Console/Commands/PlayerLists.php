<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use GuzzleHttp\Client;
use Illuminate\Http\Request;

use Hash;
use Session;
use App\Models\User;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\DB;
use Illuminate\Pagination\Paginator;

// include_once('default_parameter.php');

class PlayerLists extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'player:lists';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Command description';

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        $counter = (int) DB::table('player_lists')->count();
        while ($counter < 50) {
            if ($counter == 0)
                $serial = 1;
            else
                $serial = (int) DB::table('player_lists')->max('serial_NO') + 1;
            $counter = (int) DB::table('player_lists')->count();
            $player_register_name = "target-" . (string) str_pad($serial, 13, '0', STR_PAD_LEFT);
            $parameters = ["playerName" => "$player_register_name"];
            $url = 'player/exist';
            $is_exist = guzzle_get($parameters, $url);
            if ($is_exist != null) {
                if (isset($is_exist) && !$is_exist['isExist']) {
                    $para = ['playerName' => "$player_register_name"];
                    $url = 'player';
                    $obj = guzzle_post($para, $url);
                    if ($obj == null) {
                        $sql = sprintf("INSERT INTO `player_lists` (`player_account`, `serial_NO`) VALUES ('%s', '%d')", $player_register_name, $serial);
                        DB::insert($sql);
                    } else {
                        continue;
                    }
                } else if ($is_exist['isExist']) {
                    $serial++;
                    continue;
                } else {
                    continue;
                }
            }
        }
        return Command::SUCCESS;
    }
}