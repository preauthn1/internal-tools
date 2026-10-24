#!/usr/bin/env python3
"""檢查 AI 呼叫日誌：所有 AI 呼叫（不分程序）的時間有沒有重疊。
執行：python3 eval/check_ai_log.py <日誌>（沒有重疊 exit 0，有重疊 exit 1）
日誌由 bot 在設了 BOT_AI_LOG 時寫入，每行「程序編號<TAB>開始<TAB>結束」（Unix 秒，到毫秒）。
例：同時開兩個量測程序，確認跨程序排隊有效
    rm -f /tmp/ai.log
    BOT_AI_LOG=/tmp/ai.log BOT_INTENT_EVAL=eval/intent_set.json ROUNDS=1 .build/release/bot &
    BOT_AI_LOG=/tmp/ai.log BOT_INTENT_EVAL=eval/intent_set.json ROUNDS=1 .build/release/bot; wait
    python3 eval/check_ai_log.py /tmp/ai.log"""
import sys

if len(sys.argv) != 2: raise SystemExit(__doc__)
calls = []
for line in open(sys.argv[1], encoding="utf-8"):
    if line.strip():
        pid, start, end = line.split("\t")
        calls.append((float(start), float(end), pid))
calls.sort()

# 依開始時間排序後，每一筆都要在前面所有呼叫中最晚結束的那一筆之後才開始
overlaps, last = [], None
for c in calls:
    if last and last[1] > c[0]: overlaps.append((last, c))
    if not last or c[1] > last[1]: last = c
print(f"AI 呼叫 {len(calls)} 次，來自 {len({c[2] for c in calls})} 個程序，重疊 {len(overlaps)} 處{' ✅' if not overlaps else ' ❌'}")
for a, b in overlaps[:20]:
    print(f"   程序 {a[2]} {a[0]:.3f}～{a[1]:.3f} 和程序 {b[2]} {b[0]:.3f}～{b[1]:.3f}")
sys.exit(1 if overlaps else 0)
