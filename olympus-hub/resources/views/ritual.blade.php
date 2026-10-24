<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sacred Ritual - 神聖祭典</title>
    <style>
        body {
            background: #0f172a;
            color: #e2e8f0;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            padding: 2rem;
            min-height: 100vh;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
        }
        .header {
            text-align: center;
            margin-bottom: 3rem;
        }
        h1 {
            color: #fbbf24;
            text-transform: uppercase;
            letter-spacing: 0.3rem;
            margin-bottom: 0.5rem;
        }
        .realm-info {
            color: #38bdf8;
            font-weight: bold;
            font-size: 1.2rem;
        }
        .ritual-form {
            background: rgba(30, 41, 59, 0.7);
            border: 2px solid #38bdf8;
            border-radius: 1rem;
            padding: 2rem;
            box-shadow: 0 0 20px rgba(56, 189, 248, 0.2);
        }
        .form-group {
            margin-bottom: 2rem;
        }
        label {
            display: block;
            margin-bottom: 0.8rem;
            color: #94a3b8;
            font-weight: bold;
            text-transform: uppercase;
            font-size: 0.9rem;
        }
        select, textarea {
            width: 100%;
            background: #0f172a;
            border: 1px solid #334155;
            border-radius: 0.5rem;
            color: #e2e8f0;
            padding: 1rem;
            font-family: 'Fira Code', 'Courier New', Courier, monospace;
            font-size: 1rem;
            box-sizing: border-box;
        }
        select:focus, textarea:focus {
            outline: none;
            border-color: #38bdf8;
            box-shadow: 0 0 10px rgba(56, 189, 248, 0.3);
        }
        optgroup {
            background: #1e293b;
            color: #fbbf24;
            font-style: normal;
        }
        option {
            background: #0f172a;
            color: #e2e8f0;
        }
        .unleash-btn {
            display: block;
            width: 100%;
            padding: 1.2rem;
            background: transparent;
            color: #fbbf24;
            border: 2px solid #fbbf24;
            border-radius: 0.5rem;
            font-size: 1.2rem;
            font-weight: bold;
            text-transform: uppercase;
            letter-spacing: 0.3rem;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        .unleash-btn:hover {
            background: #fbbf24;
            color: #0f172a;
            box-shadow: 0 0 20px rgba(251, 191, 36, 0.5);
        }
        .unleash-btn:disabled {
            border-color: #475569;
            color: #475569;
            cursor: not-allowed;
        }
        .chronicle-box {
            margin-top: 3rem;
            background: #000;
            border: 1px solid #334155;
            border-radius: 0.5rem;
            padding: 1.5rem;
            display: none;
        }
        .chronicle-header {
            color: #64748b;
            font-size: 0.8rem;
            margin-bottom: 1rem;
            border-bottom: 1px solid #1e293b;
            padding-bottom: 0.5rem;
            display: flex;
            justify-content: space-between;
        }
        #chronicle-content {
            white-space: pre-wrap;
            font-family: 'Fira Code', 'Courier New', Courier, monospace;
            color: #10b981;
            line-height: 1.5;
        }
        .back-link {
            display: block;
            text-align: center;
            margin-top: 2rem;
            color: #64748b;
            text-decoration: none;
        }
        /* 加載動畫 */
        .loading-overlay {
            display: none;
            position: fixed;
            top: 0; left: 0; right: 0; bottom: 0;
            background: rgba(15, 23, 42, 0.8);
            z-index: 100;
            justify-content: center;
            align-items: center;
            flex-direction: column;
        }
        .bolt {
            color: #fbbf24;
            font-size: 4rem;
            animation: flash 1s infinite;
        }
        @keyframes flash {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.5; transform: scale(1.1); }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>神聖祭典 (The Sacred Ritual)</h1>
            <p>正在與 <span class="realm-info">{{ $realm['name'] }}</span> 位面進行共鳴...</p>
        </div>

        <div class="ritual-form">
            <div class="form-group">
                <label for="command-select">選擇神諭指令 (Divine Command)</label>
                <select id="command-select">
                    <option value="" disabled selected>-- 請選擇要降下的神諭 --</option>
                    @foreach($commands as $namespace => $cmds)
                        <optgroup label="{{ $namespace }}">
                            @foreach($cmds as $signature => $description)
                                <option value="{{ $signature }}">{{ $signature }} - {{ $description }}</option>
                            @endforeach
                        </optgroup>
                    @endforeach
                </select>
            </div>

            <div class="form-group">
                <label for="arguments-input">輸入參數與選項 (Arguments & Options)</label>
                <textarea id="arguments-input" rows="4" placeholder="例如：--ids=1,2 --currency=MYR --dry-run"></textarea>
            </div>

            <button type="button" class="unleash-btn" id="unleash-btn">
                降下神諭 (Unleash Oracle)
            </button>
        </div>

        <div class="chronicle-box" id="chronicle-box">
            <div class="chronicle-header">
                <span>眾神編年史 (The Chronicles) - [{{ $realm['name'] }}]</span>
                <span id="chronicle-status"></span>
            </div>
            <div id="chronicle-content"></div>
        </div>

        <a href="{{ route('olympus.gates') }}" class="back-link">← 返回命運之門 (Back to Gates)</a>
    </div>

    <div class="loading-overlay" id="loading-overlay">
        <div class="bolt">⚡</div>
        <p>正在穿透位面障壁，降下神諭中...</p>
    </div>

    <script>
        document.getElementById('unleash-btn').addEventListener('click', async function() {
            const command = document.getElementById('command-select').value;
            const arguments = document.getElementById('arguments-input').value;
            const btn = this;
            const overlay = document.getElementById('loading-overlay');
            const box = document.getElementById('chronicle-box');
            const content = document.getElementById('chronicle-content');
            const status = document.getElementById('chronicle-status');

            if (!command) {
                alert('請先選擇要降下的神諭指令！');
                return;
            }

            if (!confirm('確定要在 [{{ $realm['name'] }}] 位面降下此神諭嗎？')) {
                return;
            }

            // 顯示加載中
            overlay.style.display = 'flex';
            btn.disabled = true;

            try {
                const response = await fetch("{{ route('olympus.invoke', $realm['id']) }}", {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-TOKEN': '{{ csrf_token() }}'
                    },
                    body: JSON.stringify({ command, arguments })
                });

                const data = await response.json();

                // 顯示結果
                overlay.style.display = 'none';
                btn.disabled = false;
                box.style.display = 'block';
                content.innerText = data.chronicle || data.message;
                status.innerText = data.status === 'success' ? '⚡ 成功' : '❌ 失敗';
                status.style.color = data.status === 'success' ? '#10b981' : '#ef4444';
                
                // 捲動到結果區
                box.scrollIntoView({ behavior: 'smooth' });

            } catch (error) {
                overlay.style.display = 'none';
                btn.disabled = false;
                alert('聖域連線中斷：' + error.message);
            }
        });
    </script>
</body>
</html>
