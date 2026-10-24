#!/usr/bin/env python3
"""看門狗：檢查 issue 迴圈和 model-server 有沒有在跑，停掉超過 10 分鐘就用 Telegram 通知一次，恢復時再通知一次。

由 launchd 每 5 分鐘執行（ops/com.rory.gitlab-issue-loop.watchdog.plist）。通知設定沿用 config.json 的 notify。
檢查的服務寫在 config.json 的 watchdog.services（預設：迴圈程序、model-server 的 /version）。
"""
import json
import os
import subprocess
import sys
import time
import urllib.request

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, ROOT)
from loop import notify  # noqa: E402

STATE = os.path.expanduser("~/Library/Caches/gitlab-issue-loop/watchdog.json")
GRACE = 600  # 停掉超過 10 分鐘才通知（launchd 自己重啟的空檔不吵人）


def loop_alive():
    return subprocess.run(["pgrep", "-f", "loop.py watch"], capture_output=True).returncode == 0


def url_alive(url):
    try:
        urllib.request.urlopen(url, timeout=5).read()
        return True
    except Exception:
        return False


def main():
    with open(os.path.join(ROOT, "config.json")) as f:
        cfg = json.load(f)
    services = {"issue 迴圈": loop_alive}
    for name, url in (cfg.get("watchdog", {}).get("urls") or {"model-server": "http://0.0.0.0:8765/version"}).items():
        services[name] = lambda u=url: url_alive(u)
    try:
        with open(STATE) as f:
            state = json.load(f)
    except Exception:
        state = {}
    now = time.time()
    for name, check in services.items():
        s = state.setdefault(name, {"down_since": None, "alerted": False})
        if check():
            if s["alerted"]:
                notify(cfg, f"✅ {name} 恢復了（停了 {(now - s['down_since']) / 60:.0f} 分鐘）")
            state[name] = {"down_since": None, "alerted": False}
            continue
        s["down_since"] = s["down_since"] or now
        if not s["alerted"] and now - s["down_since"] >= GRACE:
            notify(cfg, f"🛑 {name} 停了 {(now - s['down_since']) / 60:.0f} 分鐘，launchd 沒能讓它恢復，要人看一下")
            s["alerted"] = True
        print(time.strftime("%H:%M:%S"), name, "停掉中", flush=True)
    os.makedirs(os.path.dirname(STATE), exist_ok=True)
    with open(STATE, "w") as f:
        json.dump(state, f)


if __name__ == "__main__":
    main()
