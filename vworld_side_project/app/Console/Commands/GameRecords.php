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
class GameRecords extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'game:records';

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
        $from = date('Y-m-d H:i:s', strtotime('-5 minute'));
        $to = date('Y-m-d H:i:s');
        $para = ['pageIndex' => '1', 'startTime' => $from, 'endTime' => $to];
        $url = 'game-records';
        $tmp = guzzle_get($para, $url);
        if ($tmp != null) {
            $totalPage = $tmp['totalPage'];
            for ($i = 1; $i <= $totalPage; $i++) {
                $para = ['pageIndex' => $i, 'startTime' => $from, 'endTime' => $to];
                $temp = guzzle_get($para, $url);
                if ($temp != null) {
                    $game_records = $temp['items'];
                    foreach ($game_records as $game_record) {
                        $betTime = date('Y-m-d H:i:s', strtotime(str_replace('Z', '', str_replace('T', ' ', $game_record["betTime"])) . '+8hours'));
                        $roundStartTime = date('Y-m-d H:i:s', strtotime(str_replace('Z', '', str_replace('T', ' ', $game_record["roundStartTime"])) . '+8hours'));
                        $roundEndTime = date('Y-m-d H:i:s', strtotime(str_replace('Z', '', str_replace('T', ' ', $game_record["roundEndTime"])) . '+8hours'));
                        $sql = sprintf("INSERT INTO `game_records` (`player_name`, `game_id`, `before_amount`, `win_loss_amount`, `round_id`, `bet_id`, `bet_time`, `round_start_time`, `round_end_time`, `transaction_no`, `bet_amount`, `payout_amount`, `tax_amount`, `valid_amount`) VALUES ('%s', '%d', '%f', '%f', '%s', '%s', '%s', '%s', '%s', '%s', '%f', '%f', '%f', '%f')", $game_record["playerName"], $game_record["gameID"], $game_record["beforeAmount"], $game_record["winLossAmount"], $game_record["roundID"], $game_record["betID"], $betTime, $roundStartTime, $roundEndTime, $game_record["transactionNo"], $game_record["betAmount"], $game_record["payoutAmount"], $game_record["taxAmount"], $game_record["validAmount"]);
                        DB::insert($sql);
                    }
                }
            }
        }
        return Command::SUCCESS;
    }
}