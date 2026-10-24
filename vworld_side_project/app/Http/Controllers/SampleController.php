<?php

namespace App\Http\Controllers;

use App\Models\User;
use Hash;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\DB;
use Session;


require '../vendor/autoload.php';

class SampleController extends Controller
{
    function index()
    {
        return view('login');
    }

    function registration()
    {
        return view('registration');
    }

    function validate_registration(Request $request)
    {

        $request->validate([
            'name' => 'required|max:50|unique:users|alpha_dash',
            'email' => 'required|email|unique:users',
            'password' => 'required|regex:/[0-9a-zA-Z]{10}/',
        ]);

        $data = $request->all();
        if ($data['password'] == $data['re_password']) {
            User::create([
                'name' => $data['name'],
                'email' => $data['email'],
                'password' => Hash::make($data['password'])
            ]);
            return redirect('login')->with('success', 'Registration Completed, now you can login');
        } else
            return redirect('registration')->with('error_message', '密碼與確認密碼不一致');
    }
    function validate_login(Request $request)
    {
        $request->validate([
            'email' => 'required',
            'password' => 'required'
        ]);

        $credentials = $request->only('email', 'password');

        if (Auth::attempt($credentials)) {
            return redirect('dashboard');
        }

        return redirect('login')->with('success', 'Login details are not valid');
    }
    function logout()
    {
        Session::flush();
        Auth::logout();
        return Redirect('login');
    }

    function dashboard()
    {
        if (Auth::check()) {
            $user_id = Auth::user()->id;
            if (DB::table('players')->where('user_id', $user_id)->exists() && !Auth::user()->money) {
                $player_register_name = DB::table('players')->where('user_id', $user_id)->value('player');
                $para = ["playerName" => "$player_register_name"];
                $url = 'wallet/amount';
                $temp = guzzle_get($para, $url);

                if ($temp != null) {
                    $wallet_amount = $temp['amount'];
                    $transfer_type = 'transfer_out';
                    $this->withholding_money($transfer_type, $wallet_amount, $user_id);
                    $merchantOrderNo = $player_register_name . date("YmdHisu") . substr(microtime(FALSE), 2, 9);
                    $merchantOrderNo = str_replace(' ', '', $merchantOrderNo);
                    $para = ['merchantOrderNo' => "$merchantOrderNo", 'amount' => "$wallet_amount", 'playerName' => "$player_register_name"];
                    $url = 'wallet/transfer-out';
                    $transfer_out = guzzle_post($para, $url);
                    $this->transaction($transfer_out, $transfer_type, $user_id, $player_register_name, $wallet_amount, $merchantOrderNo);
                }
            }
            return view('dashboard', ['user' => Auth::user()->name, 'money' => DB::table('users')->where('id', Auth::user()->id)->value('money')]);
        }
        return redirect('login')->with('success', 'you are not allowed to access');
    }
    function games()
    {
        $this->create_players();
        $user_money = (double) Auth::user()->money;
        if ($user_money) {
            $user_id = Auth::user()->id;
            $transfer_type = 'transfer_in';
            $this->withholding_money($transfer_type, $user_money, $user_id);
            $player_register_name = DB::table('players')->where('user_id', $user_id)->value('player');
            $merchantOrderNo = $player_register_name . date("YmdHisu") . substr(microtime(FALSE), 2, 9);
            $merchantOrderNo = str_replace(' ', '', $merchantOrderNo);
            $para = ['amount' => $user_money, 'merchantOrderNo' => $merchantOrderNo, 'playerName' => $player_register_name];
            $url = 'wallet/transfer-in';
            $transfer_in = guzzle_post($para, $url);
            $this->transaction($transfer_in, $transfer_type, $user_id, $player_register_name, $user_money, $merchantOrderNo);
        }
        $game_lists = DB::table('game_lists')->get();
        return view('games', ['game_lists' => $game_lists]);
    }

    function withholding_money($transfer_type, $money, $user_id)
    {
        $user = User::find($user_id);
        if ($transfer_type == 'transfer_out')
            $user->money = $user->money + $money;
        else if ($transfer_type == 'transfer_in')
            $user->money = $user->money - $money;
        $user->save();

    }

    function transaction($transfer_data, $transfer_type, $user_id, $player_register_name, $money, $merchantOrderNo)
    {
        $username = Auth::user()->name;
        if ($transfer_data != null) {
            $sql = sprintf("INSERT INTO `money_books` (`user_id`, `status`, `user`, `player`, `transfer_type`, `merchant_order_no`, `transaction_no`, `amount`, `before_amount`, `after_amount`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%01.2f', '%01.2f', '%01.2f')", $user_id, 'success', $username, $player_register_name, $transfer_type, $merchantOrderNo, $transfer_data["transactionNo"], $transfer_data["amount"], $transfer_data["beforeAmount"], $transfer_data["afterAmount"]);
            DB::insert($sql);
        } else {
            $user = User::find($user_id);
            if ($transfer_type == 'transfer_in') {
                $user->money = $user->money - $money;
            } else if ($transfer_type == 'transfer_out') {
                $user->money = $user->money + $money;
            }
            $user->save();
            $sql = sprintf("INSERT INTO `money_books` (`user_id`, `status`, `user`, `player`, `transfer_type` `merchant_order_no`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s')", $user_id, 'fail', $username, $player_register_name, 'transfer-in', $merchantOrderNo);
            DB::insert($sql);
        }
    }
    function game_link(Request $request)
    {
        $data = $request->all();
        if ($data['gameID']) {
            $player_register_name = (string) DB::table('players')->where('user_id', Auth::user()->id)->value('player');
            $para = ['playerName' => "$player_register_name", 'lang' => 'zh-tw', 'gameID' => $data['gameID']];
            $url = 'game/link';
            $temp = guzzle_get($para, $url);
            if ($temp != null) {
                $game_link = $temp['link'];
                return view('game_link', ['game_link' => $game_link]);
            }
        } else
            return view('dashboard');
    }
    function create_players()
    {
        $user_id = Auth::user()->id;
        if (!DB::table('players')->where('user_id', $user_id)->exists()) {
            $serial = (int) DB::table('player_lists')->min('serial_NO');
            $player_register_name = DB::table('player_lists')->where('serial_NO', $serial)->value('player_account');
            $username = Auth::user()->name;
            $sql = sprintf("INSERT INTO `players` (`user_id`, `user`, `player`, `serial_NO`) VALUES ('%d', '%s', '%s', '%d')", $user_id, $username, $player_register_name, $serial);
            DB::insert($sql);
            DB::table('player_lists')->where('serial_NO', $serial)->delete();
        }
    }
    function game_screenshot(Request $request)
    {
        $data = $request->all();
        if (isset($data['transactionNo'])) {
            $para = ['transactionNo' => $data['transactionNo']];
            $url = 'game-records/screenshot-history/link';
            $link = guzzle_get($para, $url);
            if (!isset($link[0])) {
                return view('game_screenshot', ['game_screenshot' => $link]);
            }
        }
    }
    function transfer(Request $request)
    {
        $table_name = 'money_books';
        $view_name = 'transfer';
        $column_name = 'user_id';
        $column_value = Auth::user()->id;
        return $this->find_something($request, $table_name, $view_name, $column_name, $column_value);
    }

    function game_records(Request $request)
    {
        $table_name = 'game_records';
        $view_name = 'game_records';
        $column_name = 'player_name';
        $column_value = DB::table('players')->where('user_id', Auth::user()->id)->value('player');
        return $this->find_something($request, $table_name, $view_name, $column_name, $column_value);
    }
    function find_something($request, $table_name, $view_name, $column_name, $column_value)
    {
        $data = $request->all();
        if (!empty($data)) {
            $startTime = $data['startTime'];
            $endTime = $data['endTime'];
            $pageNumber = $data['pageNumber'];
            $temp = DB::table($table_name)->where($column_name, $column_value)->whereBetween('created_at', [$startTime, $endTime])->orderBy('created_at', 'DESC')->paginate($pageNumber);
            $compact_data = [
                'startTime' => $startTime,
                'endTime' => $endTime,
                'temp' => $temp
            ];
        } else
            $compact_data = [];
        return view($view_name, compact('compact_data'));
    }
}