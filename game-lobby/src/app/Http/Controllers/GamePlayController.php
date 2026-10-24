<?php

namespace App\Http\Controllers;

use App\Services\Games\Exceptions\GameProviderException;
use App\Services\Games\Exceptions\GameTransferPendingException;
use App\Services\Games\PlayGameService;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Inertia\Inertia;
use InvalidArgumentException;

class GamePlayController extends Controller
{
    public function launch(Request $request, PlayGameService $service, string $provider, string $gameId): RedirectResponse
    {
        try {
            $link = $service->enter($request->user(), $provider, $gameId);
        } catch (GameTransferPendingException|InvalidArgumentException $exception) {
            return redirect()->route('home')->withErrors(['play' => $exception->getMessage()]);
        } catch (GameProviderException $exception) {
            Log::warning('遊戲啟動失敗', ['provider' => $provider, 'api_code' => $exception->apiCode, 'category' => 'provider']);

            return redirect()->route('home')->withErrors(['play' => '目前無法啟動遊戲，請稍後再試。']);
        }

        return redirect()->away($link);
    }

    public function settle(Request $request, PlayGameService $service): RedirectResponse
    {
        try {
            $result = $service->settle($request->user());
        } catch (GameProviderException|InvalidArgumentException $exception) {
            Log::warning('遊戲結算失敗', ['provider' => $exception instanceof GameProviderException ? $exception->provider : 'unknown', 'category' => 'settlement']);

            return back()->withErrors(['settlement' => '結算處理中，系統將自動重試。']);
        }

        if ($result->pending) {
            return back()->withErrors(['settlement' => '結算處理中，系統將自動重試。']);
        }
        if ($result->residual) {
            return back()->withErrors(['settlement' => '遊戲錢包殘值低於 0.10，暫時無法轉出。']);
        }
        if (bccomp($result->amount, '0', 2) > 0) {
            Inertia::flash('toast', ['type' => 'success', 'message' => "已轉回平台錢包 {$result->amount}"]);
        }

        return back();
    }
}
