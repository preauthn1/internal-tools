#!/bin/sh
# 安裝或更新 issue 迴圈和看門狗的 launchd 設定（LaunchAgent，登入後才跑）。
# 用法：ops/install.sh            安裝並啟動（已經在跑的會重載）
#       ops/install.sh uninstall  停掉並移除
# 迴圈的紀錄在 ~/Library/Logs/gitlab-issue-loop/loop.log
set -e
ROOT=$(cd "$(dirname "$0")/.." && pwd)
DEST="$HOME/Library/LaunchAgents"
UID_=$(id -u)
# launchd 的 PATH 很少：claude、go、node（corepack、pnpm）、glab 都要找得到
NODE_BIN=$(dirname "$(command -v node)")
PATH_="$HOME/.local/bin:/opt/homebrew/bin:$NODE_BIN:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
for name in com.rory.gitlab-issue-loop com.rory.gitlab-issue-loop.watchdog; do
  launchctl bootout "gui/$UID_/$name" 2>/dev/null || true
  if [ "$1" = uninstall ]; then rm -f "$DEST/$name.plist"; echo "已移除 $name"; continue; fi
  mkdir -p "$DEST" "$HOME/Library/Logs/gitlab-issue-loop"
  sed -e "s|__ROOT__|$ROOT|g" -e "s|__HOME__|$HOME|g" -e "s|__PATH__|$PATH_|g" "$ROOT/ops/$name.plist" > "$DEST/$name.plist"
  plutil -lint "$DEST/$name.plist" >/dev/null
  launchctl bootstrap "gui/$UID_" "$DEST/$name.plist"
  echo "已啟動 $name"
done
