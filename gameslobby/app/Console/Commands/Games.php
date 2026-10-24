<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Models\Game;
use Exception;

class Games extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'command:games';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Command description';

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        $client = new \GuzzleHttp\Client();
        $res = $client->request('GET', env('PGS_API_Vendor_Development').'games?merchantID='.env('merchantID').'&version='.env('version').'&merchantSecretKey='.env('merchantSecretKey'));
        $games = json_decode($res->getBody());
        if(isset($games->data)){
            foreach ($games->data as $k => $v) {
                Game::updateOrCreate ([
                    'gameID' => $v->gameID,
                ],[
                    'gameCode' => $v->gameCode,
                    'gameName' => $v->gameName,
                    'imageUrl' => $v->imageUrl,
                ]);
            }
            return 1;
        } else {
            throw new Exception('game is null');
        }
    }
}
