<?php

namespace App\Services;

use Illuminate\Support\Facades\Auth;
use Exception;

class PgsApiService
{
    public function player_exist()
    {
        $params = [];
        $params['playerName']        = Auth::user()->name;
        $params['merchantID']        = env('merchantID');
        $params['version']           = 1;
        $params['merchantSecretKey'] = env('merchantSecretKey');
        $res = $this->guzzleHttp('player/exist?', $params, 'GET');

        $player_not_exist = $this->confirmDataFormat('player_exist', $res);
        return $player_not_exist;
    }

    public function player()
    {
        $params = [];
        $params['playerName']        = Auth::user()->name;
        $params['merchantSecretKey'] = env('merchantSecretKey');
        $params['merchantID']        = env('merchantID');
        $res = $this->guzzleHttp('player?', $params, 'POST');

        $player = $this->confirmDataFormat('player', $res);
        return $player;
    }

    public function game_link($data)
    {
        $params = [];
        $params['playerName']        = Auth::user()->name;
        $params['merchantID']        = env('merchantID');
        $params['version']           = 1;
        $params['merchantSecretKey'] = env('merchantSecretKey');
        $params['lang']              = 'zh-tw';
        $params['gameID']            = $data['gameID'];
        $res = $this->guzzleHttp('game/link?', $params, 'GET');

        $game_link = $this->confirmDataFormat('game_link', $res);
        return $game_link;
    }

    public function wallet_amount()
    {
        try {
            $params = [];
            $params['playerName']        = Auth::user()->name;
            $params['merchantID']        = env('merchantID');
            $params['version']           = 1;
            $params['merchantSecretKey'] = env('merchantSecretKey');
            $res = $this->guzzleHttp('wallet/amount?', $params, 'GET');

            $wallet = $this->confirmDataFormat('wallet_amount', $res);
            if ($wallet === false) {
                return ['state'=>2, 'message'=>'第三方無法查詢遊戲點數', 'data'=>null];
            } elseif ($wallet == 0) {
                return ['state'=>2, 'message'=>'遊戲內已無點數可轉出', 'data'=>null];
            } else {
                return ['state'=>0, 'message'=>'success', 'data'=>$wallet];
            }
        } catch (Exception $e) {
            return $this->error($e);
        }
    }

    public function wallet_transfer_out($wallet)
    {
        try {
            $params = [];
            $params['amount']            = $wallet;
            $params['merchantOrderNo']   = 'charley'.date('Y/m/d H:i:s');
            $params['playerName']        = Auth::user()->name;
            $params['merchantSecretKey'] = env('merchantSecretKey');
            $params['merchantID']        = env('merchantID');
            $res = $this->guzzleHttp('wallet/transfer-out?', $params, 'POST');

            $transaction = $this->confirmDataFormat('wallet_transfer_out', $res);
            if ($transaction) {
                return ['state'=>0,'message'=>'success', 'data'=>$transaction];
            }
            return ['state'=>2,'message'=>'轉回失敗，還原遊戲內點數', 'data'=>null];
        } catch (Exception $e) {
            return $this->error($e);
        }
    }

    public function wallet_transfer_in($point)
    {
        try {
            $params = [];
            $params['amount']            = $point;
            $params['merchantOrderNo']   = 'charley'.date('Y/m/d H:i:s');
            $params['playerName']        = Auth::user()->name;
            $params['merchantSecretKey'] = env('merchantSecretKey');
            $params['merchantID']        = env('merchantID');
            $res = $this->guzzleHttp('wallet/transfer-in?', $params, 'POST');

            $transaction = $this->confirmDataFormat('wallet_transfer_in', $res);
            if ($transaction) {
                return ['state'=>0,'message'=>'success', 'data'=>$transaction];
            }
            return ['state'=>2,'message'=>'轉入失敗，還原錢包內點數', 'data'=>null];
        } catch (Exception $e) {
            return $this->error($e);
        }
    }

    public function game_records()
    {
        $params = [];
        $params['playerName']        = Auth::user()->name;
        $params['merchantSecretKey'] = env('merchantSecretKey');
        $params['merchantID']        = env('merchantID');
        $params['version']           = 1;
        $params['startTime']         = date('Y-m-d H:i:s', strtotime('-1 day'));
        $params['endTime']           = date('Y-m-d H:i:s');
        $params['pageIndex']         = 1;
        $params['pageSize']          = 100;
        $res = $this->guzzleHttp('game-records?', $params, 'GET');

        $game_records = $this->confirmDataFormat('game_records', $res);
        return $game_records;
    }

    private function guzzleHttp($url, $params, $method)
    {
        $client = new \GuzzleHttp\Client();
        $url = env('PGS_API_Vendor_Development').$url;
        switch ($method) {
            case 'GET':
                $url .= http_build_query($params);
                $res = $client->request($method, $url);
                break;
            case 'POST':
                $res = $client->request($method, $url, ['json' => $params]);
                break;
            default:
                return ['state'=>2,'message'=>'api失敗', 'data'=>null];
        }
        return $res;
    }

    private function confirmDataFormat($type, $res)
    {
        $data = json_decode($res->getBody());
        $result = false;
        if (isset($data) && $data->code === 0) {
            switch($type) {
                case 'player_exist':
                    //玩家是否已存在
                    if (isset($data->data->isExist)) {
                        if ($data->data->isExist === false) {
                            $result = true;
                        }
                    }
                    break;
                case 'player':

                    break;
                case 'game_link':
                    if (isset($data->data->link)) {
                        $result = $data->data->link;
                    }
                    break;
                case 'wallet_amount':
                    // return $data;
                    if (isset($data->data->amount)) {
                        if ($data->data->amount >= 0) {
                            $result = $data->data->amount;
                        }
                    }
                    break;
                case 'wallet_transfer_out':
                    if (isset($data->data->afterAmount)) {
                        if ($data->data->afterAmount >= 0) {
                            $result = [];
                            $result['transactionNo'] = $data->data->transactionNo;
                            $result['amount']        = $data->data->amount;
                            $result['beforeAmount']  = $data->data->beforeAmount;
                            $result['afterAmount']   = $data->data->afterAmount;
                        }
                    }
                    break;
                case 'wallet_transfer_in':
                    if (isset($data->data->afterAmount)) {
                        if ($data->data->afterAmount >= 0) {
                            $result = [];
                            $result['transactionNo'] = $data->data->transactionNo;
                            $result['amount']        = $data->data->amount;
                            $result['beforeAmount']  = $data->data->beforeAmount;
                            $result['afterAmount']   = $data->data->afterAmount;
                        }
                    }
                    break;
                case 'game_records':
                    if (isset($data->data->items)) {
                        $result = [];
                        foreach ($data->data->items as $k => $v) {
                            $result[$k]['name']           = $v->playerName;
                            $result[$k]['gameID']         = $v->gameID;
                            $result[$k]['beforeAmount']   = $v->beforeAmount;
                            $result[$k]['winLossAmount']  = $v->winLossAmount;
                            $result[$k]['roundID']        = $v->roundID;
                            $result[$k]['betID']          = $v->betID;
                            $result[$k]['roundStartTime'] = $v->roundStartTime;
                            $result[$k]['roundEndTime']   = $v->roundEndTime;
                            $result[$k]['transactionNo']  = $v->transactionNo;
                            $result[$k]['betAmount']      = $v->betAmount;
                            $result[$k]['betTime']        = $v->betTime;
                            $result[$k]['payoutAmount']   = $v->payoutAmount;
                            $result[$k]['taxAmount']      = $v->taxAmount;
                            $result[$k]['validAmount']    = $v->validAmount;
                        }
                    }
                    break;
            }
        }
        return $result;
    }

    private function error($e)
    {
        switch($e->getCode()) {
            case '400':
                return ['state'=>2,'message'=>'缺參數', 'data'=>null];
            case '404':
                return ['state'=>2,'message'=>'url不對，無回傳值', 'data'=>null];
            default:
                return ['state'=>2,'message'=>'其他錯誤', 'data'=>null];
        }
    }
}
