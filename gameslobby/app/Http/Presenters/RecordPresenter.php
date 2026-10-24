<?php

namespace App\Http\Presenters;

use Illuminate\Support\Facades\Auth;
use App\Models\TransferPoint;
use App\Models\GameRecord;

class RecordPresenter
{
    public function game_records($sdata = null, $edata = null, $name = null)
    {
        $game_records = GameRecord::where('name', Auth::user()->name)
            ->orderBy('betTime', 'desc')
            ->join('games', 'game_records.gameID', 'games.gameID');
        if ($sdata) {
            $game_records = $game_records->where('betTime', '>=', $sdata);
        }
        if ($edata) {
            $game_records = $game_records->where('betTime', '<=', $edata.' 23:59:59');
        }
        if ($name) {
            $game_records = $game_records->where('gameName', 'like', '%'. $name.'%');
        }
        $game_records_count = $game_records->count();
        $game_records = $game_records->select('game_records.id as id', 'gameName', 'betAmount', 'winLossAmount', 'betTime')
            ->paginate(5);
        return view('game_record')->with('game_records', $game_records)->with('count', $game_records_count);
    }

    public function transfer_points()
    {
        $transfer_points = TransferPoint::where('name', Auth::user()->name)
            ->orderBy('created_at', 'desc')
            ->paginate(5);
        $transfer_points_count = TransferPoint::where('name', Auth::user()->name)->count();
        return view('transfer_point')->with('transfer_points', $transfer_points)->with('count', $transfer_points_count);
    }
}
