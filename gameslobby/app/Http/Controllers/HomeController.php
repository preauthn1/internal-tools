<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\User;
use Illuminate\Support\Facades\Auth;
use App\Models\Game;
use Exception;
use Illuminate\Support\Facades\DB;
use App\Repositories\TransferPointRepository;
use App\Services\PgsApiService;
use App\Http\Presenters\RecordPresenter;

class HomeController extends Controller
{
    protected $transferPointRepository;
    protected $pgs_api_service;
    protected $record_presenter;

    /**
     * Create a new controller instance.
     *
     * @return void
     */
    public function __construct(TransferPointRepository $transferPointRepository, PgsApiService $pgs_api_service, RecordPresenter $record_presenter)
    {
        // $this->middleware('auth');
        $this->transferPointRepository = $transferPointRepository;
        $this->PgsApiService = $pgs_api_service;
        $this->RecordPresenter = $record_presenter;
    }

    /**
     * Show the application dashboard.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function index()
    {
        $games = Game::all();
        if (Auth::check()) {
            try {
                // 檢查玩家是否存在
                $player_not_exist = $this->PgsApiService->player_exist();
                if ($player_not_exist) {
                    // 建立玩家帳號
                    $player = $this->PgsApiService->player();
                    if ($player) {
                        throw new Exception('player exist is error');
                    }
                }
                return view('home')->with('games', $games)->with('point', Auth::user()->point);
            } catch (Exception $e) {
                // return $e;
                return '首頁維護中';
            }
        } else {
            return view('home')->with('games', $games);
        }
    }

    public function game_link(Request $request)
    {
        $data = $request->all();
        try {
            $game_link = $this->PgsApiService->game_link($data);
            if ($game_link) {
                return response()->json(['state'=>0,'message'=>'success','data'=>$game_link]);
            } else {
                return response()->json(['state'=>1,'message'=>'game link is fail']);
            }
        } catch (Exception $e) {
            return response()->json(['state'=>2,'message'=>$e], 403);
        }
    }

    public function transfer_out($point = false)
    {
        try {
            //確認遊戲內的點數
            $wallet_amount = $this->PgsApiService->wallet_amount($point);
            if (!empty($wallet_amount['data'])) {
                //點數全部轉回
                if (!$point) {
                    $point = $wallet_amount['data'];
                //轉回一定數量之點數
                } else {
                    if ($point>$wallet_amount['data']) {
                        return response()->json(['state'=>2, 'message'=>'餘額不足，無法轉回此金額', 'data'=>null]);
                    }
                }
                $wallet = Auth::user()->point;
                $before_point = $wallet;
                $after_point = $wallet+$point;

                //點數記錄新增
                $transfer_point = $this->transferPointRepository->transferPointCreate($point, $before_point, $after_point);
                //打第三方api
                $transaction = $this->PgsApiService->wallet_transfer_out($point);
                if (!empty($transaction['data'])) {
                    //點數轉移成功
                    $this->transferPointRepository->transferPointUpdate($transfer_point, $transaction['data']);
                    //更新錢包點數
                    User::where('name', Auth::user()->name)->update(['point' => $after_point]);
                    return response()->json(['state'=>0,'message'=>'success', 'data'=>$after_point]);
                }
                $wallet_amount= $transaction;
            }
            return response()->json(['state'=>2,'message'=>$wallet_amount['message'], 'data'=>null]);
        } catch (Exception $e) {
            return response()->json(['state' => $e->getCode(), 'message' => $e->getMessage(), 'data'=>null]);
        }
    }

    public function transfer_in($point = false)
    {
        try {
            $wallet_amount = Auth::user()->point;
            $wallet['message']='錢包點數為零';
            if ($wallet_amount>0) {
                //點數全部轉回
                if (!$point) {
                    $point = $wallet_amount;
                //轉回一定數量之點數
                } else {
                    if ($point>$wallet_amount) {
                        return response()->json(['state'=>2,'message'=>'餘額不足，無法轉入此金額', 'data'=>null]);
                    }
                }
                $before_point = $wallet_amount;
                $after_point = $wallet_amount-$point;

                //點數記錄新增
                $transfer_point = $this->transferPointRepository->transferPointCreate($point, $before_point, $after_point);
                DB::beginTransaction();
                //先更新錢包為零
                User::where('name', Auth::user()->name)->update(['point' => $after_point]);
                //打第三方api
                $transaction = $this->PgsApiService->wallet_transfer_in($point);
                if (isset($transaction['data'])) {
                    //點數轉移成功
                    $this->transferPointRepository->transferPointUpdate($transfer_point, $transaction['data']);
                    DB::commit();
                    return response()->json(['state'=>0,'message'=>'success', 'data'=>$after_point]);
                }
                $wallet= $transaction;
                DB::rollback();
            }
            return response()->json(['state'=>2,'message'=>$wallet['message'], 'data'=>$point]);
        } catch (Exception $e) {
            DB::rollback();
            return response()->json(['code' => $e->getCode(), 'message' => $e->getMessage(), 'data'=>0]);
        }
    }

    public function game_records(Request $request)
    {
        $game_records = $this->PgsApiService->game_records();
        if (isset($game_records)) {
            $this->transferPointRepository->get_game_records($game_records);
        } else {
            return '遊戲記錄錯誤';
        }
        $data = $request->all();
        if ($data) {
            $sdata = $data['sdata'];
            $edata = $data['edata'];
            $name = $data['name'];
            $res = $this->RecordPresenter->game_records($sdata, $edata, $name);
        } else {
            $res = $this->RecordPresenter->game_records();
        }
        return $res;
    }

    public function transfer_points()
    {
        $res = $this->RecordPresenter->transfer_points();
        return $res;
    }
}
