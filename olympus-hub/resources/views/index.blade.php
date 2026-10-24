<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Olympus Hub - 奧林帕斯樞紐</title>
    <style>
        body {
            background: #0f172a;
            color: #e2e8f0;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            overflow: hidden;
        }
        .container {
            text-align: center;
            border: 2px solid #38bdf8;
            padding: 3rem;
            border-radius: 1rem;
            box-shadow: 0 0 30px rgba(56, 189, 248, 0.4);
            background: rgba(15, 23, 42, 0.9);
            position: relative;
        }
        h1 {
            font-size: 3.5rem;
            margin-bottom: 1rem;
            color: #38bdf8;
            text-transform: uppercase;
            letter-spacing: 0.8rem;
            text-shadow: 0 0 10px rgba(56, 189, 248, 0.8);
        }
        p {
            font-style: italic;
            color: #94a3b8;
            margin-bottom: 2rem;
            font-size: 1.2rem;
        }
        .status {
            margin-top: 2rem;
            font-weight: bold;
            color: #10b981;
            text-shadow: 0 0 5px rgba(16, 185, 129, 0.5);
        }
        .enter-btn {
            display: inline-block;
            margin: 1rem;
            padding: 1.2rem 2rem;
            background: transparent;
            text-decoration: none;
            text-transform: uppercase;
            letter-spacing: 0.15rem;
            font-weight: bold;
            border-radius: 0.5rem;
            transition: all 0.3s ease;
        }
        .enter-btn.cashflow {
            color: #38bdf8;
            border: 1px solid #38bdf8;
            box-shadow: 0 0 10px rgba(56, 189, 248, 0.2);
        }
        .enter-btn.cashflow:hover {
            background: #38bdf8;
            color: #0f172a;
            box-shadow: 0 0 20px rgba(56, 189, 248, 0.6);
            transform: scale(1.05);
        }
        .enter-btn.game {
            color: #fbbf24;
            border: 1px solid #fbbf24;
            box-shadow: 0 0 10px rgba(251, 191, 36, 0.2);
        }
        .enter-btn.game:hover {
            background: #fbbf24;
            color: #0f172a;
            box-shadow: 0 0 20px rgba(251, 191, 36, 0.6);
            transform: scale(1.05);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>⚡ Olympus Hub ⚡</h1>
        <p>「歡迎來到數位世界的頂端。從這裡俯瞰，所有的 Production 領域皆為臣民。」</p>
        
        <div class="portal-container">
            @if($showCashFlow)
            <a href="{{ route('olympus.gates') }}" class="enter-btn cashflow">
                進入金流命運之門 (Cash Flow Realm)
            </a>
            @endif

            @if($showGame)
            <a href="{{ route('olympus.game_gates') }}" class="enter-btn game">
                進入遊戲核心聖域 (Game Core Realm)
            </a>
            @endif

            @if(!$showCashFlow && !$showGame)
            <p style="color: #ef4444; font-weight: bold;">⚠️ 偵測到未授權的位面接入。請聯繫大魔法師。</p>
            @endif
        </div>

        <div class="status">● 神聖聖所已就緒 (Sanctum Online)</div>
    </div>
</body>
</html>
