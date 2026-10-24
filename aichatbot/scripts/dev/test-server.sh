#!/bin/zsh
# 開發機的實機測試環境：Go 版伺服器（Telegram webhook＋後台）＋ Cloudflare 臨時 tunnel，都放在 screen 裡背景跑。
#
#   zsh scripts/dev/test-server.sh start     # 編譯、建前端、開 tunnel（沒開的話）、啟動伺服器
#   zsh scripts/dev/test-server.sh restart   # 換成目前的程式碼重啟伺服器（tunnel 還連得到就不動，網址不變）
#   zsh scripts/dev/test-server.sh stop      # 停伺服器和 tunnel
#   zsh scripts/dev/test-server.sh status    # 另外從外面連 tunnel，連不到（臨時 tunnel 被 Cloudflare 收掉）就換新的、重啟伺服器
#
# - 在主要工作目錄執行；Telegram 設定讀 .env，資料庫是 data/go-test.sqlite（BOT_DB 可以改；第一次會匯入 eval/ 的森日茶屋）。
# - 後台：http://0.0.0.0:8081/ ；管理員帳號：.build/gobot create-admin <帳號>（BOT_DB 同上，密碼自己輸入）。
# - 紀錄：~/Library/Logs/aichatbot/bot-serve.log、tunnel.log。tunnel 網址是臨時的，tunnel 重開就會換（start 會重設 webhook）。
set -u
cd "$(dirname "$0")/../.."
LOGS=~/Library/Logs/aichatbot; mkdir -p $LOGS
DB=${BOT_DB:-data/go-test.sqlite}

tunnel_url() { grep -ao "https://[a-z0-9-]*\.trycloudflare\.com" $LOGS/tunnel.log 2>/dev/null | tail -1; }
serving() { lsof -iTCP:8080 -sTCP:LISTEN -P >/dev/null 2>&1; }
# tunnel_reachable：從外面打 webhook。伺服器回的 401（沒帶 secret）、伺服器沒開時 Cloudflare 回的 502 都代表 tunnel 還在；
# 連不到（000）或 530（Cloudflare 找不到 tunnel）就是失效了。
# DNS 走 Cloudflare 的 DoH：新網址剛建好時系統 DNS 查不到會被快取住，之後一段時間都查不到
tunnel_reachable() {
  local url=$(tunnel_url) code
  [[ -n $url ]] || return 1
  code=$(curl --doh-url https://0.0.0.0/dns-query -s -o /dev/null -w "%{http_code}" -m 10 -X POST "$url/telegram/webhook" -d '{}')
  [[ $code != 000 && $code != 530 ]]
}
# tunnel_alive：在跑而且外面連得到；連不到時隔 5 秒再試一次，剛建好或偶爾一次失敗不算失效
tunnel_alive() {
  screen -ls | grep -q bot-tunnel || return 1
  tunnel_reachable || { sleep 5; tunnel_reachable; }
}

build() {
  echo "編譯伺服器、建前端…"
  go -C go build -o ../.build/gobot ./cmd/bot || exit 1
  (cd web && COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack pnpm install --frozen-lockfile >/dev/null 2>&1 \
    && COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack pnpm build:antd >/dev/null 2>&1) || { echo "❌ 前端建置失敗"; exit 1; }
  [[ -f $DB ]] || BOT_DB=$DB .build/gobot import eval
}

stop_server() {
  screen -S bot-serve -X quit >/dev/null 2>&1
  pkill -f ".build/gobot serve" 2>/dev/null
  for i in {1..20}; do serving || return 0; sleep 0.5; done
  echo "❌ 8080 還被占用：$(lsof -iTCP:8080 -sTCP:LISTEN -P | tail -1)"; exit 1
}

# stop_tunnel：screen quit 不會結束裡面的 cloudflared，要直接停掉，不然舊的 tunnel 會一直留著、寫進同一個紀錄檔
stop_tunnel() {
  screen -S bot-tunnel -X quit >/dev/null 2>&1
  pkill -f "cloudflared tunnel --no-autoupdate --url http://localhost:8080" 2>/dev/null
  for i in {1..20}; do pgrep -f "cloudflared tunnel --no-autoupdate --url http://localhost:8080" >/dev/null || return 0; sleep 0.5; done
  pkill -9 -f "cloudflared tunnel --no-autoupdate --url http://localhost:8080" 2>/dev/null
}

start_tunnel() {
  if tunnel_alive; then return; fi
  screen -ls | grep -q bot-tunnel && echo "⚠️ tunnel $(tunnel_url) 外面連不到，換一個新的"
  stop_tunnel
  : > $LOGS/tunnel.log
  screen -dmS bot-tunnel zsh -c "cloudflared tunnel --no-autoupdate --url http://localhost:8080 >> $LOGS/tunnel.log 2>&1"
  for i in {1..30}; do [[ -n "$(tunnel_url)" ]] && break; sleep 1; done
  [[ -n "$(tunnel_url)" ]] || { echo "❌ tunnel 沒有拿到網址（看 $LOGS/tunnel.log）"; exit 1; }
  # 新網址的 DNS 要幾秒才查得到，等外面連得到再設 webhook
  for i in {1..30}; do tunnel_reachable && return; sleep 2; done
  echo "⚠️ 新的 tunnel $(tunnel_url) 一分鐘內外面還連不到，照樣設 webhook（Telegram 會重試）"
}

start_server() {
  local url=$(tunnel_url)
  screen -dmS bot-serve zsh -c "set -a && source .env && set +a && TELEGRAM_PUBLIC_URL=$url BOT_DB=$DB exec .build/gobot serve >> $LOGS/bot-serve.log 2>&1"
  for i in {1..30}; do serving && break; sleep 0.5; done
  serving || { echo "❌ 伺服器沒有起來："; tail -5 $LOGS/bot-serve.log; exit 1; }
  echo "✅ 伺服器已啟動（$DB）：webhook $url ／ 後台 http://0.0.0.0:8081/"
}

case "${1:-status}" in
  start)   build; stop_server; start_tunnel; start_server ;;
  restart) build; stop_server; start_tunnel; start_server ;;
  stop)    stop_server; stop_tunnel; echo "已停止伺服器和 tunnel" ;;
  status)  serving && echo "伺服器：執行中（後台 http://0.0.0.0:8081/）" || echo "伺服器：沒有在跑"
           if tunnel_alive; then echo "tunnel：$(tunnel_url)（外面連得到）"
           elif serving; then start_tunnel; stop_server; start_server   # 換了 tunnel 要重啟伺服器，才會用新網址設 webhook
           else echo "tunnel：沒有在跑或外面連不到（start 時會換新的）"
           fi ;;
  *) echo "用法：zsh scripts/dev/test-server.sh start|restart|stop|status"; exit 1 ;;
esac
