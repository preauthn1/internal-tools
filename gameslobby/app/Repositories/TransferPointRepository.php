<?php

namespace App\Repositories;

use App\Models\TransferPoint;
use Illuminate\Support\Facades\Auth;
use App\Models\GameRecord;

class TransferPointRepository
{
    public function transferPointCreate($point, $before_point, $after_point)
    {
        $transfer_point = TransferPoint::create([
            'amount'          => $point,
            'merchantOrderNo' => 'charley'.date('Ymd His'),
            'name'            => Auth::user()->name,
            'transactionNo'   => '',
            'beforeAmount'    => $before_point,
            'afterAmount'     => $after_point,
            'state'    => 0,
        ]);
        return $transfer_point;
    }

    public function transferPointUpdate($transfer_point, $transaction)
    {
        $transfer_point = TransferPoint::where('merchantOrderNo', $transfer_point['merchantOrderNo'])
            ->update([
                'transactionNo'   => $transaction['transactionNo'],
                'state'    => 1,
            ]);
        return $transfer_point;
    }

    public function get_game_records($game_records)
    {
        foreach ($game_records as $k => $v) {
            GameRecord::updateOrCreate([
                'transactionNo' => $v['transactionNo'],
            ], [
                'name'          => $v['name'],
                'gameID'        => $v['gameID'],
                'beforeAmount'  => $v['beforeAmount'],
                'winLossAmount' => $v['winLossAmount'],
                'roundID'       => $v['roundID'],
                'betID'         => $v['betID'],
                'roundStartTime'=> substr($v['roundStartTime'], 0, 10).' '.substr($v['roundStartTime'], 11, 8),
                'roundEndTime'  => substr($v['roundEndTime'], 0, 10).' '.substr($v['roundEndTime'], 11, 8),
                'betAmount'     => $v['betAmount'],
                'betTime'       => substr($v['betTime'], 0, 10).' '.substr($v['betTime'], 11, 8),
                'payoutAmount'  => $v['payoutAmount'],
                'taxAmount'     => $v['taxAmount'],
                'validAmount'   => $v['validAmount'],
            ]);
        }
    }
}
