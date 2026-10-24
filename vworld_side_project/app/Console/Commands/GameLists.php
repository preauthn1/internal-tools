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

class GameLists extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'game:lists';

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
        $data = [];
        $url = 'games';
        $game_lists = guzzle_get($data, $url);
        if ($game_lists != null) {
            foreach ($game_lists as $game_list) {
                $gameID = $game_list['gameID'];
                $gameName = $game_list["gameName"];
                $imageUrl = $game_list["imageUrl"];
                $gameType = $game_list["gameType"];
                $gameTypeName = $game_list["gameTypeName"];
                if (!DB::table('game_lists')->where('game_id', $gameID)->exists()) {
                    $sql = sprintf("INSERT INTO `game_lists` (`game_id`, `game_name`, `image_url`, `game_type`, `game_type_name`) VALUES ('%d', '%s', '%s', '%d', '%s')", $gameID, $gameName, $imageUrl, $gameType, $gameTypeName);
                    DB::insert($sql);
                } else {
                    $sql = sprintf("UPDATE `game_lists` SET `game_name`='%s',`image_url`='%s',`game_type`='%d',`game_type_name`='%s' where `game_id` = '%d'", $gameName, $imageUrl, $gameType, $gameTypeName, $gameID);
                    DB::update($sql);
                }
            }
        }
        return Command::SUCCESS;
    }
}