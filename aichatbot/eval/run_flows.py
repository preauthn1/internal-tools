#!/usr/bin/env python3
"""完整流程測試：把 flow_tests.json 的多輪對話餵給 chat.swift，比對每一步的回覆種類與內容。

執行：python3 eval/run_flows.py                         （測 chat.swift，全部對話）
      swift build -c release && python3 eval/run_flows.py --bot bot --exe .build/release/bot --runs 3
                                                       （測新版、跑 3 次取平均）
      python3 eval/run_flows.py F04 F12                （只跑指定的對話）
結果存到 eval/results/flows-<時間>.json
"""
import json, os, re, subprocess, sys, tempfile, time

EVAL = os.path.dirname(os.path.abspath(__file__))
import argparse
ap = argparse.ArgumentParser()
ap.add_argument("--bot", default="chat.swift")
ap.add_argument("--runs", type=int, default=1)
ap.add_argument("--tests", default="flow_tests.json", help="題庫檔（相對 eval/ 或絕對路徑）")
ap.add_argument("--exe", help="已編譯好的執行檔（省掉每次重新編譯；自動改進迴圈用）")
ap.add_argument("--json", help="另外把總結寫成 JSON（給 loop/score.py 讀）")
ap.add_argument("--cache-ns", choices=["dev", "holdout"], help="AI 回答快取的命名空間（BOT_AI_CACHE=1 時才有用）；預設依題庫位置判斷，eval/holdout/ 底下是 holdout")
ap.add_argument("only", nargs="*")
args = ap.parse_args()
if not args.exe and not os.path.exists(os.path.join(EVAL, args.bot)):
    # 新版已拆成 Swift Package，沒有 eval/bot.swift 可以直接 swift 執行
    sys.exit(f"找不到 eval/{args.bot}。新版請先 swift build -c release，再加 --exe .build/release/bot")
from testdates import resolve   # 題庫的相對日期 {{date:...}}
tests = json.loads(resolve(open(os.path.join(EVAL, args.tests), encoding="utf-8").read()))["conversations"]
if args.only:
    tests = [c for c in tests if c["id"] in args.only]


def lang_of(text):
    """跟 chat.swift 一樣先看文字系統；沒有中日韓文字就當英文。"""
    if re.search(r"[぀-ヿ]", text): return "ja"
    if re.search(r"[가-힯]", text): return "ko"
    if re.search(r"[一-鿿]", text): return "zh"
    return "en"


def judge(turn, got):
    """回傳 (是否通過, 原因)。"""
    if turn["user"] == "/pending":
        items = got.get("pending", [])
        if "pending_any" in turn and not any(w in i for i in items for w in turn["pending_any"]):
            return False, f"待處理清單沒有 {turn['pending_any']}（實際：{items}）"
        if "pending_none" in turn:
            bad = [i for i in items for w in turn["pending_none"] if w in i]
            if bad: return False, f"待處理清單不該有：{bad}"
        return True, ""
    key, reply = got.get("key", "?"), got.get("reply", "")
    if key not in turn["keys"]:
        return False, f"應為 {'/'.join(turn['keys'])}，實際是 {key}"
    if miss := [w for w in turn.get("all", []) if w not in reply]:
        return False, f"缺少 {miss}"
    if turn.get("any") and not any(w in reply for w in turn["any"]):
        return False, f"沒有提到 {'/'.join(turn['any'])}"
    if bad := [w for w in turn.get("none", []) if w in reply]:
        return False, f"不該出現 {bad}"
    # 判斷語言前先拿掉這一步預期會出現的字（顧客的名字「たなか」會讓中文回覆被當成日文）
    body = reply
    for w in turn.get("all", []): body = body.replace(w, "")
    if (want := turn.get("lang")) and reply and lang_of(body) != want:
        return False, f"語言應為 {want}，實際是 {lang_of(body)}"
    return True, ""


# 組出輸入：每段對話前 /reset，最後 /quit
lines = []
for conv in tests:
    lines.append("/reset")
    lines += [t["user"] for t in conv["turns"]]
lines.append("/quit")


# AI 回答快取：驗證題組分開存（見 Sources/BotCore/AICache.swift）
tests_path = os.path.realpath(os.path.join(EVAL, args.tests))
cache_ns = args.cache_ns or ("holdout" if tests_path.startswith(os.path.realpath(os.path.join(EVAL, "holdout")) + os.sep) else "dev")
CACHE_RE = re.compile(r"AI 快取：命中 (\d+) 次、實際呼叫 (\d+) 次")
cache_total = [0, 0]   # 命中、實際呼叫（全部執行加總）


def run_once():
    """跑一次全部對話，回傳每一步的結果。"""
    tmp = tempfile.mkdtemp()
    trace_path = os.path.join(tmp, "trace.jsonl")
    env = dict(os.environ, FLOW_TRACE=trace_path, CHAT_LOG_DIR=tmp, BOT_AI_CACHE_NS=cache_ns)
    proc = subprocess.run([args.exe] if args.exe else ["swift", os.path.join(EVAL, args.bot)], input="\n".join(lines) + "\n",
                          text=True, env=env, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE, timeout=3600)
    # 有開 AI 回答快取時，bot 離開前在 stderr 印命中幾次、實際呼叫幾次
    for hits, misses in CACHE_RE.findall(proc.stderr):
        cache_total[0] += int(hits); cache_total[1] += int(misses)
    entries = [json.loads(l) for l in open(trace_path, encoding="utf-8")] if os.path.exists(trace_path) else []
    # 程式沒產生任何回覆：多半是編譯失敗或一開始就當掉，直接顯示錯誤並停止，不要印一整頁「沒有回覆」
    if not any(e.get("key") not in ("reset", "pending") for e in entries):
        errors = [l for l in proc.stderr.splitlines() if "error" in l.lower() or "fatal" in l.lower()]
        print(f"❌ {args.bot} 沒有產生任何回覆（編譯失敗或執行就當掉）：")
        print("\n".join(errors[:10]) or proc.stderr[-2000:])
        sys.exit(1)
    groups, cur = [], None
    for e in entries:
        if e.get("key") == "reset":
            cur = []; groups.append(cur)
        elif cur is not None:
            cur.append(e)
    out = []
    for conv, got_list in zip(tests, groups + [[]] * (len(tests) - len(groups))):
        for i, turn in enumerate(conv["turns"]):
            got = got_list[i] if i < len(got_list) else {"key": "（沒有回覆）"}
            ok, reason = judge(turn, got)
            out.append({"conv": conv["id"], "title": conv["title"], "user": turn["user"], "pass": ok, "reason": reason, "got": got})
    return out


total_turns = sum(len(c["turns"]) for c in tests)
print(f"測試 {args.bot}：{len(tests)} 段對話、{total_turns} 步，跑 {args.runs} 次（約需幾分鐘）…")
start = time.time()
runs = []
for r in range(args.runs):
    res = run_once()
    runs.append(res)
    ok = sum(x["pass"] for x in res)
    convs_ok = sum(all(x["pass"] for x in res if x["conv"] == c["id"]) for c in tests)
    print(f"  第 {r + 1} 次：通過 {ok}/{len(res)} 步（{ok / max(len(res), 1) * 100:.0f}%），完全通過的對話 {convs_ok}/{len(tests)}")

# 逐步列出：每一步在幾次裡通過（只列有失敗的）
print("\n========== 有失敗的步驟（通過次數／總次數）==========")
last_conv = None
for i, x in enumerate(runs[0]):
    passes = sum(run[i]["pass"] for run in runs)
    if passes == len(runs):
        continue
    if x["conv"] != last_conv:
        print(f"\n【{x['conv']}】{x['title']}"); last_conv = x["conv"]
    print(f"  {passes}/{len(runs)}  {x['user']}")
    for run in runs:
        if not run[i]["pass"]:
            g = run[i]["got"]
            shown = "（待處理清單）" if x["user"] == "/pending" else f"[{g.get('key')}] {g.get('reply', '')[:70]}"
            print(f"        → {shown}\n          原因：{run[i]['reason']}")

avg = sum(sum(x["pass"] for x in run) for run in runs) / len(runs)
print("\n========== 總結 ==========")
print(f"{args.bot}：平均通過 {avg:.1f}/{total_turns} 步（{avg / total_turns * 100:.0f}%），跑 {len(runs)} 次")
print(f"耗時 {time.time() - start:.0f} 秒")
if os.environ.get("BOT_AI_CACHE") == "1":
    print(f"AI 快取：命中 {cache_total[0]} 次、實際呼叫 {cache_total[1]} 次")

os.makedirs(os.path.join(EVAL, "results"), exist_ok=True)
out = os.path.join(EVAL, "results", f"flows-{args.bot.replace('.swift', '')}-{int(time.time())}.json")
json.dump(runs, open(out, "w", encoding="utf-8"), ensure_ascii=False, indent=1)
print(f"詳細結果：{out}")

if args.json:
    json.dump({"avg_pass": avg, "total": total_turns, "pct": avg / total_turns * 100,
               "per_run": [sum(x["pass"] for x in run) for run in runs], "detail": out,
               "cache": {"hits": cache_total[0], "misses": cache_total[1]}},
              open(args.json, "w", encoding="utf-8"), ensure_ascii=False, indent=1)
