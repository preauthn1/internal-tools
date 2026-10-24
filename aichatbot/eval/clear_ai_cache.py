#!/usr/bin/env python3
"""清掉 AI 回答快取（BOT_AI_CACHE=1 時存的回答）。Go 版（go/internal/model/cache.go）和 Swift 版（tag swift-final）用同一個目錄，一起清。

執行：python3 eval/clear_ai_cache.py [--ns dev|holdout|all]   （預設 all）
快取目錄：BOT_AI_CACHE_DIR（Go 版和 Swift 版都認這個環境變數），預設 ~/Library/Caches/aichatbot/ai/{dev,holdout}；
  命名空間是 BOT_AI_CACHE_NS：dev＝開發題組、holdout＝驗證題組。也可以直接 rm -rf ~/Library/Caches/aichatbot/ai
系統更新、模型服務版本或 Haiku 的 model id 改了時快取會自動失效（都在 key 裡），不必手動清；舊的檔案留著只佔空間
"""
import argparse, os, shutil
from pathlib import Path

ROOT = Path(os.environ.get("BOT_AI_CACHE_DIR") or Path.home() / "Library" / "Caches" / "aichatbot" / "ai")

ap = argparse.ArgumentParser()
ap.add_argument("--ns", choices=["dev", "holdout", "all"], default="all")
args = ap.parse_args()
for ns in (["dev", "holdout"] if args.ns == "all" else [args.ns]):
    d = ROOT / ns
    # 只數檔案數，不列內容（驗證題組的快取不給看）
    n = sum(1 for _ in d.rglob("*.json")) if d.exists() else 0
    shutil.rmtree(d, ignore_errors=True)
    print(f"已清掉 {ns}：{n} 筆（{d}）")
