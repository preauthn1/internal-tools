<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gates of Destiny - 命運之門</title>
    <style>
        body {
            background: #0f172a;
            color: #e2e8f0;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            padding: 2rem;
            min-height: 100vh;
        }
        .header {
            text-align: center;
            margin-bottom: 3rem;
        }
        h1 {
            color: #fbbf24;
            text-transform: uppercase;
            letter-spacing: 0.5rem;
            font-size: 2.5rem;
            text-shadow: 0 0 15px rgba(251, 191, 36, 0.5);
        }
        .realm-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 2rem;
            max-width: 1200px;
            margin: 0 auto;
        }
        .realm-card {
            background: rgba(30, 41, 59, 0.7);
            border: 1px solid #334155;
            border-radius: 1rem;
            padding: 2rem;
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
        }
        .realm-card:hover {
            border-color: #38bdf8;
            transform: translateY(-5px);
            box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
        }
        .realm-name {
            font-size: 1.5rem;
            font-weight: bold;
            color: #38bdf8;
            margin-bottom: 1rem;
        }
        .realm-desc {
            color: #94a3b8;
            font-size: 0.95rem;
            line-height: 1.6;
            margin-bottom: 1.5rem;
            min-height: 3rem;
        }
        .realm-meta {
            display: flex;
            justify-content: space-between;
            font-size: 0.8rem;
            color: #64748b;
            margin-bottom: 2rem;
            border-top: 1px solid #334155;
            padding-top: 1rem;
        }
        .invoke-btn {
            display: block;
            width: 100%;
            padding: 0.8rem;
            background: #1e293b;
            color: #38bdf8;
            border: 1px solid #38bdf8;
            border-radius: 0.5rem;
            text-align: center;
            text-decoration: none;
            font-weight: bold;
            text-transform: uppercase;
            letter-spacing: 0.1rem;
            text-indent: 0.1rem; /* 補償 letter-spacing 造成的偏移 */
            box-sizing: border-box; /* 確保寬度包含邊框與內距 */
            transition: all 0.3s ease;
            cursor: pointer;
        }
        .invoke-btn:hover {
            background: #38bdf8;
            color: #0f172a;
            box-shadow: 0 0 15px rgba(56, 189, 248, 0.5);
        }
        .back-link {
            display: block;
            text-align: center;
            margin-top: 3rem;
            color: #64748b;
            text-decoration: none;
            transition: color 0.3s ease;
        }
        .back-link:hover {
            color: #e2e8f0;
        }
        /* 簡易動畫 */
        @keyframes pulse {
            0% { opacity: 0.6; }
            50% { opacity: 1; }
            100% { opacity: 0.6; }
        }
        .status-dot {
            display: inline-block;
            width: 8px;
            height: 8px;
            background: #10b981;
            border-radius: 50%;
            margin-right: 5px;
            animation: pulse 2s infinite;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{ $title }}</h1>
        <p>觀測平行位面，降下金流神諭</p>
    </div>

    <div class="realm-grid">
        @foreach($realms as $realm)
        <div class="realm-card">
            <div class="realm-name">{{ $realm['name'] }}</div>
            <div class="realm-desc">{{ $realm['description'] }}</div>
            <div class="realm-meta">
                <span>狀態: <span class="status-dot"></span> {{ $realm['status'] }}</span>
                <span>神力密度: {{ $realm['power'] }}</span>
            </div>
            
            <a href="{{ route('olympus.ritual', $realm['id']) }}" class="invoke-btn">run artisan</a>
        </div>
        @endforeach
    </div>

    <a href="{{ route('olympus.hub') }}" class="back-link">← 返回聖所 (Return to Sanctum)</a>
</body>
</html>
