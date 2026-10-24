#!/bin/sh
# 編譯 model-server、安裝到 ~/.local/aichatbot/bin，並（重新）啟動 launchd 常駐。
# 合併改到 Sources/ 的 MR 之後在主要工作目錄執行（AGENTS.md 開發流程第 5 步）。
# 重啟前會等到沒有人在用本機模型（迴圈的開發 AI、量測程式），避免量測出現 AI 失敗；--now 不等。
# 用法：scripts/dev/install-model-server.sh [--now]
#       scripts/dev/install-model-server.sh uninstall
set -e
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
LABEL=com.rory.aichatbot.model-server
PLIST="$HOME/Library/LaunchAgents/$LABEL.plist"
BIN="$HOME/.local/aichatbot/bin/model-server"
if [ "$1" = uninstall ]; then
  launchctl bootout "gui/$(id -u)/$LABEL" 2>/dev/null || true
  rm -f "$PLIST"; echo "已移除 $LABEL"; exit 0
fi
cd "$ROOT"
swift build -c release --product model-server
if [ "$1" != --now ]; then
  # 有開發 AI（claude -p）、Go 或 Swift 的量測程式在跑時，等它們結束（最多 2 小時）
  for i in $(seq 1 240); do
    pgrep -f "claude -p|go/bin/chat|run_flows.py|\.build/release/bot|snapshot.py|score.py" >/dev/null || break
    [ "$i" = 1 ] && echo "有人在用本機模型，等它結束再重啟…"
    sleep 30
  done
fi
mkdir -p "$(dirname "$BIN")" "$HOME/Library/Logs/aichatbot" "$(dirname "$PLIST")"
cp .build/release/model-server "$BIN.new" && mv "$BIN.new" "$BIN"
sed -e "s|__BIN__|$BIN|g" -e "s|__HOME__|$HOME|g" "$ROOT/scripts/dev/$LABEL.plist" > "$PLIST"
plutil -lint "$PLIST" >/dev/null
# 每次都重新註冊，不用 kickstart -k：launchd 記著執行檔的簽章需求（LWCR），換了新執行檔只重啟會一直 spawn failed（EX_CONFIG）
launchctl bootout "gui/$(id -u)/$LABEL" 2>/dev/null || true
for i in $(seq 1 10); do launchctl print "gui/$(id -u)/$LABEL" >/dev/null 2>&1 || break; sleep 1; done
launchctl bootstrap "gui/$(id -u)" "$PLIST"
# 啟動後偶爾會被 macOS 結束（例如簽章原因），launchd 會自動重啟，所以等到 60 秒；中途被結束時印一次原因
reported=
for i in $(seq 1 60); do
  curl -s -m 2 http://0.0.0.0:8765/version && echo && echo "✅ model-server 已啟動" && exit 0
  info=$(launchctl print "gui/$(id -u)/$LABEL" 2>/dev/null)
  reason=$(echo "$info" | sed -n 's/.*last exit reason = //p' | head -1)
  echo "$info" | grep -q "state = running" && reason=   # 在跑時 last exit reason 是上一次留下的
  [ -n "$reason" ] && [ -z "$reported" ] && echo "model-server 被結束（$reason），等 launchd 重啟…" && reported=1
  sleep 1
done
echo "❌ model-server 60 秒內沒有回應，看 ~/Library/Logs/aichatbot/model-server.log"; exit 1
