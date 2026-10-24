#!/usr/bin/env python3
"""MR 驗收（Go 版）：一個指令跑完測試和真模型量測，和基準 eval/baseline-go.json 比較。

執行（在 repo 根目錄）：
  python3 eval/verify.py                    # 測試＋真模型量測，和基準比較
  python3 eval/verify.py --holdout          # 另外量驗證題組（只有人類驗收時用；迴圈的 AI 不能用）
  python3 eval/verify.py --no-checks        # 跳過測試（迴圈已經跑過驗收指令時）
  python3 eval/verify.py --checks-only      # 只跑測試（迴圈的驗收指令 test_cmd 就是這個）
  python3 eval/verify.py --update-baseline  # 量完寫回基準（合併後分數有提升時用）
  python3 eval/verify.py --model haiku      # 改用 Haiku 量（要 ANTHROPIC_API_KEY；只印結果和費用，不和基準比）
  python3 eval/verify.py --no-cache         # 不用 AI 回答快取（量真實的每則平均時間、Haiku 費用）
  python3 eval/verify.py --out <資料夾>     # 另外寫 metrics.json、dev_report.md、holdout/ 明細（loop/gate.py 用）
  python3 eval/verify.py --no-items         # 不做逐題比對
  python3 eval/verify.py --baseline-from <資料夾>  # 不量測，用之前 --out 的結果更新基準
  python3 eval/verify.py --store salon      # 只量一家店（tea 森日茶屋、salon 霧光髮廊；預設 all 兩家都量）

- 兩家測試商家分開量（#66）：森日茶屋（eval/ 的題組）和霧光髮廊（美髮沙龍，eval/stores/salon/ 的商家設定和開發題組），
  每家各一組指標（名稱前面標店名；森日茶屋的代號沿用舊的，沙龍的是 salon.*）和逐題比對，任何一家退步都標 ❌。
  驗證題組只有森日茶屋有。
- 真模型量測一次只跑一組（檔案鎖 ~/Library/Caches/aichatbot/measure.lock）：有別的量測在跑就排隊等，不會同時搶本機模型。
- 逐題比對：總數之外，再和逐題基準（eval/baseline-items.json、eval/stores/salon/baseline-items.json、
  eval/holdout/baseline-items.json）比，有題目從答對變答錯就算退步（總數持平或進步也一樣，抓「修好一題、弄壞另一題」）。
  開發題組列出是哪幾題，驗證題組只給題數。

- 真模型量測預設用常駐的 model-server（0.0.0.0:8765，見 AGENTS.md），預設開 AI 回答快取（BOT_AI_CACHE=1）。
- 單題題組之外，另外跑 eval/testset_extra.json（誘答題、跨語言題，店家多開韓文），量「誘答題 編造」和「跨語言 答對」。
- 輸出一張表：指標｜基準｜這次｜差異，退步的標 ❌；另外印每則平均時間（開快取時含命中，只顯示不比較）和各店的量測時間；
  --model haiku 時再印 AI 用量（token、估算費用；流程題讀不到，不含）。
- 結束碼：0 沒有退步；1 有退步或測試失敗；2 量測無效（連不上模型服務、出現 AI 呼叫失敗）。
"""
import argparse
import contextlib
import datetime
import fcntl
import io
import json
import os
import re
import shutil
import subprocess
import sys
import time
import urllib.request
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
BASELINE = ROOT / "eval" / "baseline-go.json"
SALON = Path("eval") / "stores" / "salon"   # 第二家測試商家（美髮沙龍，#66）：商家設定、開發題組、逐題基準
# 店家：--store 的代號 → 指標名稱前面標的店名
STORES = {"tea": "森日茶屋", "salon": "霧光髮廊"}
# 逐題基準：每個題組答對的題目。驗證題組的放在 eval/holdout/ 底下（迴圈的 AI 不能讀）
ITEMS = {"dev": ROOT / "eval" / "baseline-items.json", "salon": ROOT / SALON / "baseline-items.json",
         "holdout": ROOT / "eval" / "holdout" / "baseline-items.json"}
# 逐題比對的題組 → (店名, 開發或驗證)
ITEM_GROUPS = {"dev": (STORES["tea"], "開發"), "salon": (STORES["salon"], "開發"), "holdout": (STORES["tea"], "驗證")}
CHAT = ROOT / "go" / "bin" / "chat"
MODEL_URL = os.environ.get("BOT_APPLE_MODEL_URL", "http://0.0.0.0:8765")
# 模型呼叫失敗時 Go 的模型層會記這兩種 log（stderr）
FAIL_MARKERS = ("本機模型呼叫失敗", "Haiku 呼叫失敗")
BLOCKED_MARKER = "被安全過濾擋下"   # go/internal/model 的 Apple 寫的 log（#68）

# 指標：名稱 → (說明, 越大越好?)。森日茶屋的代號沿用舊的（基準和 loop/gate.py 不用改），沙龍的前面加 salon.
METRICS = {
    "intent.ten": ("森日茶屋 分類題組 10 類換算", True),
    "intent.human": ("森日茶屋 分類題組 想找真人", True),
    "testset.ok": ("森日茶屋 單題題組 總計", True),
    "testset.wrong": ("森日茶屋 單題題組 答錯", False),
    "testset.missed": ("森日茶屋 單題題組 該答沒答", False),
    "testset.blocked": ("森日茶屋 單題題組 被擋下", False),
    "testset.trap_fabricated": ("森日茶屋 誘答題 編造", False),
    "testset.xlang_ok": ("森日茶屋 跨語言 答對", True),
    "flows.pass": ("森日茶屋 流程題 通過步數", True),
    "holdout.intent.ten": ("森日茶屋 驗證 分類 10 類換算", True),
    "holdout.flows.pass": ("森日茶屋 驗證 流程題 通過步數", True),
    "salon.intent.ten": ("霧光髮廊 分類題組 10 類換算", True),
    "salon.intent.human": ("霧光髮廊 分類題組 想找真人", True),
    "salon.testset.ok": ("霧光髮廊 單題題組 總計", True),
    "salon.testset.wrong": ("霧光髮廊 單題題組 答錯", False),
    "salon.testset.missed": ("霧光髮廊 單題題組 該答沒答", False),
    "salon.testset.blocked": ("霧光髮廊 單題題組 被擋下", False),
    "salon.testset.trap_fabricated": ("霧光髮廊 誘答題 編造", False),
    "salon.flows.pass": ("霧光髮廊 流程題 通過步數", True),
}


def run(cmd, env=None, timeout=3600):
    t0 = time.time()
    p = subprocess.run(cmd, cwd=ROOT, shell=isinstance(cmd, str), capture_output=True, text=True,
                       env={**os.environ, **(env or {})}, timeout=timeout)
    return p.returncode, p.stdout + p.stderr, time.time() - t0


def checks():
    """和迴圈的驗收指令相同的測試（Swift、Go、web、案例庫重複檢查）。"""
    steps = [
        ("swift build", "swift build"),
        ("bot-tests", ".build/debug/bot-tests"),
        ("telegram-tests", ".build/debug/telegram-tests"),
        ("gofmt", 'test -z "$(gofmt -l go)"'),
        ("go vet", "go -C go vet ./..."),
        ("go test", "go -C go test ./..."),
        ("web", "cd web && COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack pnpm install --frozen-lockfile "
                "&& COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack pnpm build:antd "
                "&& COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack pnpm --dir apps/web-antd typecheck"),   # vue-tsc 型別檢查
        ("check_overlap", "python3 eval/check_overlap.py"),
        ("gate-test", "python3 -m unittest loop/test_gate.py"),   # 改進迴圈守門規則的完整流程（假評分，幾秒）
    ]
    ok = True
    for name, cmd in steps:
        code, out, sec = run(cmd)
        print(f"  {'✅' if code == 0 else '❌'} {name}（{sec:.0f} 秒）", flush=True)
        if code != 0:
            print("\n".join("      " + l for l in out.strip().splitlines()[-15:]))
            ok = False
            break
    return ok


def ratio(pattern, text):
    m = re.search(pattern, text)
    return (int(m[1]), int(m[2])) if m else None


USAGE_RE = re.compile(r"AI 用量：輸入 (\d+) token、輸出 (\d+) token、估算 US\$([\d.]+)")


def intent_metrics(res, prefix, out):
    res[prefix + "intent.ten"] = ratio(r"10 類換算[^：]*：(\d+)/(\d+)", out)
    res[prefix + "intent.human"] = ratio(r"human_request\s+(\d+)/(\d+)", out)


def testset_metrics(res, prefix, out):
    total = ratio(r"總計：(\d+)/(\d+)", out)
    res[prefix + "testset.ok"] = total
    for key, word in (("testset.wrong", "答錯"), ("testset.missed", "該答沒答"), ("testset.blocked", "被擋下")):
        m = re.search(rf"{word} (\d+) 題", out)
        res[prefix + key] = (int(m[1]), total[1] if total else 0) if m else None


def flow_pass(out):
    m = re.search(r"平均通過 ([\d.]+)/(\d+) 步", out)
    return (float(m[1]), int(m[2])) if m else None


def measure(holdout, env, stores=tuple(STORES)):
    """真模型量測，回傳 ({指標: (分子, 分母)}, 有沒有 AI 呼叫失敗, {題組: 每則平均秒數}, [輸入 token, 輸出 token, 費用],
    {題組: 完整輸出}, {店名: 量測秒數})。stores 是要量的店（tea、salon）；驗證題組只有森日茶屋有。"""
    res, failed, avg, usage = {}, False, {}, [0, 0, 0.0]
    raw = {}   # 題組 → 完整輸出（--out 時寫成明細）
    seconds = {}   # 店名 → 量測秒數

    def go(label, cmd, extra_env=None):
        nonlocal failed
        code, out, sec = run(cmd, {**env, **(extra_env or {})})
        seconds[label.split()[0]] = seconds.get(label.split()[0], 0) + sec
        fails = sum(out.count(m) for m in FAIL_MARKERS)
        failed |= fails > 0
        # 被安全過濾擋下（#68）是本機模型固定的行為（greedy，每次一樣），不是機器忙，不算 AI 失敗、量測照樣有效，另外列次數
        blocked = out.count(BLOCKED_MARKER)
        cache =re.search(r"AI 快取[^：]*：命中 (\d+) 次、實際呼叫 (\d+) 次", out)
        for i, o, cost in USAGE_RE.findall(out):
            usage[0] += int(i); usage[1] += int(o); usage[2] += float(cost)
        if m := re.search(r"每則平均 ([\d.]+) 秒", out):
            avg[label] = float(m[1])
        print(f"  {'✅' if code == 0 and not fails else '⚠️'} {label}（{sec:.0f} 秒"
              f"{f'，快取命中 {cache[1]}、實際呼叫 {cache[2]}' if cache else ''}"
              f"{f'，AI 呼叫失敗 {fails} 次' if fails else ''}"
              f"{f'，被安全過濾擋下 {blocked} 次' if blocked else ''}）", flush=True)
        raw[label] = out
        return out

    if "tea" in stores:
        tea = STORES["tea"]
        intent_metrics(res, "", go(f"{tea} 分類題組", [str(CHAT)], {"BOT_INTENT_EVAL": "eval/intent_set.json"}))
        testset_metrics(res, "", go(f"{tea} 單題題組", [str(CHAT)], {"BOT_TESTSET": "eval/testset.json,eval/faq_set.json"}))

        # 誘答題、跨語言題另外跑，不動上面的分母；跨語言題有韓文，店家多開 ko（預設只開中英日，流程題要測不支援的語言）
        out = go(f"{tea} 誘答題、跨語言題", [str(CHAT)], {"BOT_TESTSET": "eval/testset_extra.json", "BOT_LANGUAGES": "zh,en,ja,ko"})
        res["testset.trap_fabricated"] = ratio(r"誘答題 編造 (\d+) 題（共 (\d+) 題）", out)
        res["testset.xlang_ok"] = ratio(r"跨語言 答對 (\d+)/(\d+)", out)

        res["flows.pass"] = flow_pass(go(f"{tea} 流程題", ["python3", "eval/run_flows.py", "--bot", "go-chat", "--exe", str(CHAT)]))

        if holdout:
            hold = {"BOT_AI_CACHE_NS": "holdout"}
            out = go(f"{tea} 驗證 分類題組", [str(CHAT)], {**hold, "BOT_INTENT_EVAL": "eval/holdout/intent_holdout.json"})
            res["holdout.intent.ten"] = ratio(r"10 類換算[^：]*：(\d+)/(\d+)", out)
            res["holdout.flows.pass"] = flow_pass(go(f"{tea} 驗證 流程題", [
                "python3", "eval/run_flows.py", "--bot", "go-chat", "--exe", str(CHAT),
                "--tests", "holdout/flow_holdout.json", "--cache-ns", "holdout"]))

    if "salon" in stores:
        # 第二家測試商家：商家設定改讀匯出檔；誘答題和其他單題在同一個檔，一次跑完。
        # AI 回答快取和森日茶屋同一個命名空間（dev）：快取的 key 含提示詞（店家資料）和訊息，兩家不會互相命中
        salon, senv = STORES["salon"], {"BOT_CONFIG_FILE": str(SALON / "config.json")}
        intent_metrics(res, "salon.", go(f"{salon} 分類題組", [str(CHAT)], {**senv, "BOT_INTENT_EVAL": str(SALON / "intent_set.json")}))
        out = go(f"{salon} 單題題組（含誘答題）", [str(CHAT)], {**senv, "BOT_TESTSET": str(SALON / "testset.json")})
        testset_metrics(res, "salon.", out)
        res["salon.testset.trap_fabricated"] = ratio(r"誘答題 編造 (\d+) 題（共 (\d+) 題）", out)
        res["salon.flows.pass"] = flow_pass(go(f"{salon} 流程題", [
            "python3", "eval/run_flows.py", "--bot", "go-chat", "--exe", str(CHAT),
            "--tests", str(ROOT / SALON / "flow_tests.json"), "--cache-ns", "dev"], senv))
    return res, failed, avg, usage, raw, seconds


def write_out(out_dir, res, failed, raw):
    """--out：寫 metrics.json（分數）、dev_report.md（開發題組的完整輸出，給改進迴圈看失敗明細）、
    holdout/（驗證題組的完整輸出，改進迴圈不能讀）。"""
    out_dir.mkdir(parents=True, exist_ok=True)
    (out_dir / "metrics.json").write_text(json.dumps(
        {"metrics": {k: list(v) if v else None for k, v in res.items()}, "failed": failed},
        ensure_ascii=False, indent=1) + "\n", encoding="utf-8")
    is_hold = lambda k: "驗證" in k.split()   # 題組名稱是「店名 驗證 …」
    dev = [f"## {k}\n\n```\n{v.strip()}\n```\n" for k, v in raw.items() if not is_hold(k)]
    (out_dir / "dev_report.md").write_text("# 開發題組明細\n\n" + "\n".join(dev), encoding="utf-8")
    hold = [(k, v) for k, v in raw.items() if is_hold(k)]
    if hold:
        (out_dir / "holdout").mkdir(exist_ok=True)
        for k, v in hold:
            (out_dir / "holdout" / f"{k}.txt").write_text(v, encoding="utf-8")


MEASURE_LOCK = Path.home() / "Library" / "Caches" / "aichatbot" / "measure.lock"
_lock = None


def lock_measurement():
    """真模型量測一次只跑一組：issue 迴圈的開發 AI、改進迴圈的守門程式、人的驗收都經過這裡，
    有別的量測在跑就排隊等它結束（同時跑會搶本機模型，量測失效）。鎖在程序結束時自動放掉。"""
    global _lock
    MEASURE_LOCK.parent.mkdir(parents=True, exist_ok=True)
    _lock = open(MEASURE_LOCK, "a+")
    try:
        fcntl.flock(_lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
    except BlockingIOError:
        _lock.seek(0)
        who = _lock.read().strip() or "另一組量測"
        print(f"⏳ 等其他量測結束（{who}）…", flush=True)
        t0 = time.time()
        fcntl.flock(_lock, fcntl.LOCK_EX)
        print(f"   等了 {time.time() - t0:.0f} 秒", flush=True)
    _lock.seek(0); _lock.truncate()
    _lock.write(f"pid {os.getpid()}，{ROOT.name}，{datetime.datetime.now():%H:%M:%S} 開始"); _lock.flush()


def write_baseline(res):
    """把量測結果寫進 eval/baseline-go.json（只更新有量到的指標）。"""
    commit = subprocess.run(["git", "rev-parse", "--short", "HEAD"], cwd=ROOT, capture_output=True, text=True).stdout.strip()
    old = json.loads(BASELINE.read_text(encoding="utf-8"))["metrics"] if BASELINE.exists() else {}
    old.update({k: list(v) for k, v in res.items() if v is not None})
    BASELINE.write_text(json.dumps({
        "說明": "Go 版的分數基準（eval/verify.py 產生，--update-baseline 更新）。[分子, 分母]；分母不同表示題組改過，要重量。",
        "updated": datetime.date.today().isoformat(), "commit": commit, "metrics": old,
    }, ensure_ascii=False, indent=1) + "\n", encoding="utf-8")
    print(f"已更新 {BASELINE.relative_to(ROOT)}（commit {commit}）")


def baseline_from(out_dir):
    """--baseline-from：用之前 --out 寫出的 metrics.json 和逐題結果更新基準，不量測。"""
    f = out_dir / "metrics.json"
    if not f.exists():
        print(f"❌ 找不到 {f}")
        return 2
    d = json.loads(f.read_text(encoding="utf-8"))
    if d.get("failed"):
        print("❌ 那次量測有 AI 呼叫失敗，不能拿來更新基準")
        return 2
    write_baseline({k: tuple(v) for k, v in d["metrics"].items() if v})
    for g, path in item_files(out_dir).items():
        if path.exists():
            ITEMS[g].write_text(path.read_text(encoding="utf-8"), encoding="utf-8")
            print(f"已更新 {ITEMS[g].relative_to(ROOT)}")
    return 0


def item_files(base_dir):
    """逐題基準的檔案：預設是 repo 的基準；--items-base 指定資料夾時（loop/gate.py 用）是那裡的 dev.json、salon.json、holdout/items.json。"""
    if base_dir is None:
        return ITEMS
    return {"dev": base_dir / "dev.json", "salon": base_dir / "salon.json", "holdout": base_dir / "holdout" / "items.json"}


def item_check(holdout, env, files=ITEMS, stores=tuple(STORES)):
    """逐題比對：用 eval/snapshot.py 再跑一次各題組（剛量過，幾乎全部命中快取），記下每題答對與否，
    和逐題基準比。回傳 (變化 {group: {"worse": [...], "better": [...]}}, 這次答對的題目, 沒有基準的題組)；
    跑不出來回傳 None。group 是 dev（森日茶屋）、salon（霧光髮廊）、holdout（森日茶屋的驗證題組）。"""
    sys.path.insert(0, str(ROOT / "eval"))
    import snapshot
    groups = (("dev",) if "tea" in stores else ()) + (("salon",) if "salon" in stores else ()) \
        + (("holdout",) if holdout and "tea" in stores else ())
    name = f"verify-{os.getpid()}"
    os.environ.update(env)
    try:
        with contextlib.redirect_stdout(io.StringIO()):
            snapshot.save(name, str(CHAT), groups)
    except SystemExit:
        return None
    try:
        cur = snapshot.pass_sets(name, groups)
    finally:
        shutil.rmtree(snapshot.SNAPSHOTS / name, ignore_errors=True)
    base = {g: json.loads(files[g].read_text(encoding="utf-8")) for g in groups if files[g].exists()}
    return snapshot.item_changes(base, {g: v for g, v in cur.items() if g in base}), cur, [g for g in groups if g not in base]


def print_items(changes, missing):
    """開發題組（兩家店）列出哪幾題變了；驗證題組只給題數，不給題目（迴圈的 AI 看得到這段輸出）。"""
    print("\n== 逐題比對（和逐題基準比，總數持平也可能有題目對調）")
    worse = False
    for group, ch in changes.items():
        label = " ".join(ITEM_GROUPS[group])
        w, b = ch["worse"], ch["better"]
        worse |= bool(w)
        print(f"  {'❌' if w else '✅'} {label}：從答對變答錯 {len(w)} 題、從答錯變答對 {len(b)} 題")
        if group != "holdout":
            for kind, k in w:
                print(f"      ❌ {kind} {k}")
            for kind, k in b:
                print(f"      ✅ {kind} {k}")
    for g in missing:
        print(f"  ⚠️ {' '.join(ITEM_GROUPS[g])}題組還沒有逐題基準（--update-baseline 會建立）")
    return worse


def fmt(v):
    if v is None:
        return "—"
    n, d = v
    return f"{n:g}/{d}"


def compare(res, base):
    """印出比較表，回傳有沒有退步。"""
    worse = False
    print("\n| 指標 | 基準 | 這次 | 差異 |\n|---|---|---|---|")
    for key, (name, higher) in METRICS.items():
        if key not in res:
            continue
        cur, old = res[key], base.get(key)
        mark = ""
        if cur is None:
            mark, worse = "❌ 讀不到結果", True
        elif old is not None:
            if cur[1] != old[1]:
                mark = "⚠️ 題數不同，基準要重量"
            else:
                diff = cur[0] - old[0]
                bad = diff < 0 if higher else diff > 0
                good = diff > 0 if higher else diff < 0
                mark = f"{'❌ 退步' if bad else '✅ 進步' if good else '持平'}（{diff:+g}）"
                worse |= bad
        print(f"| {name} | {fmt(old)} | {fmt(cur)} | {mark} |")
    return worse


def main():
    ap = argparse.ArgumentParser(description="MR 驗收（Go 版）")
    ap.add_argument("--holdout", action="store_true", help="另外量驗證題組（只有人類驗收時用）")
    ap.add_argument("--no-checks", action="store_true", help="跳過測試，只做真模型量測")
    ap.add_argument("--checks-only", action="store_true", help="只跑測試，不做真模型量測（迴圈的驗收指令）")
    ap.add_argument("--update-baseline", action="store_true", help="量完寫回 eval/baseline-go.json")
    ap.add_argument("--model", choices=["apple", "haiku"], default="apple",
                    help="量哪個模型（預設 apple；haiku 要 ANTHROPIC_API_KEY，只印結果和費用，不和基準比）")
    ap.add_argument("--no-cache", action="store_true", help="不用 AI 回答快取（量真實的每則平均時間、Haiku 費用）")
    ap.add_argument("--out", type=Path, help="量測結果和明細寫到這個資料夾（loop/gate.py 用）")
    ap.add_argument("--no-items", action="store_true", help="不做逐題比對（省約 2 分鐘）")
    ap.add_argument("--baseline-from", type=Path, help="不量測，直接用之前 --out 寫出的資料夾更新基準（accept_mr.py 合併時用，省一次量測）")
    ap.add_argument("--items-base", type=Path, help="逐題基準改讀這個資料夾（loop/gate.py 用）；搭配 --out 時，這次答對的題目也寫到 --out 底下同樣的檔名")
    ap.add_argument("--store", choices=[*STORES, "all"], default="all",
                    help="量哪家店：tea 森日茶屋、salon 霧光髮廊（美髮沙龍），預設 all 兩家都量；驗證題組只有森日茶屋有")
    args = ap.parse_args()
    stores = tuple(STORES) if args.store == "all" else (args.store,)
    if args.baseline_from:
        sys.exit(baseline_from(args.baseline_from))
    haiku = args.model == "haiku"
    if haiku and not args.checks_only:
        # 跑測試和量測之前就擋下：沒有 key 是量測無效，不是退步
        if not os.environ.get("ANTHROPIC_API_KEY"):
            print("❌ --model haiku 需要 ANTHROPIC_API_KEY（先 export ANTHROPIC_API_KEY=[REDACTED]")
            sys.exit(2)
        if args.update_baseline:
            print("❌ --update-baseline 只能搭配 --model apple：基準是本機模型的分數")
            sys.exit(2)

    if not args.no_checks:
        print("== 測試", flush=True)
        if not checks():
            sys.exit(1)
        if args.checks_only:
            print("✅ 測試全部通過")
            sys.exit(0)

    if not haiku:
        lock_measurement()
    print(f"== 真模型量測（{'Haiku' if haiku else 'Apple 本機模型'}，AI 回答快取{'關' if args.no_cache else '開'}）", flush=True)
    if not haiku:
        try:
            urllib.request.urlopen(MODEL_URL + "/version", timeout=5).read()
        except Exception as e:
            print(f"❌ 連不上模型服務 {MODEL_URL}（{e}）：model-server 要常駐，見 AGENTS.md")
            sys.exit(2)
    code, out, _ = run("go -C go build -o bin/chat ./cmd/chat")
    if code != 0:
        print(out)
        sys.exit(1)
    env = {"BOT_AI_CACHE": "0" if args.no_cache else "1", "BOT_MODEL": "haiku" if haiku else "afm"}
    res, failed, avg, usage, raw, seconds = measure(args.holdout, env, stores)
    if args.out:
        write_out(args.out, res, failed, raw)

    base = json.loads(BASELINE.read_text(encoding="utf-8"))["metrics"] if BASELINE.exists() else {}
    # Haiku 的分數不和本機模型的基準比：基準欄位留空，只列這次的數字
    base = {} if haiku else {k: tuple(v) for k, v in base.items()}
    worse = compare(res, base)
    items = None
    if not haiku and not failed and not args.no_items:
        t0 = time.time()
        items = item_check(args.holdout, env, item_files(args.items_base), stores)
        seconds["逐題比對"] = time.time() - t0
        if items is None:
            print("\n⚠️ 逐題比對跑不出來（snapshot.py 失敗），這次只比總數")
        else:
            worse |= print_items(items[0], items[2])
            if args.out:
                f = args.out / "metrics.json"
                d = json.loads(f.read_text(encoding="utf-8"))
                d["items"] = {g: {"worse": len(c["worse"]), "better": len(c["better"])} for g, c in items[0].items()}
                f.write_text(json.dumps(d, ensure_ascii=False, indent=1) + "\n", encoding="utf-8")
                for g, path in item_files(args.out).items():   # 這次答對的題目（gate 保留或量基準時拿來當新的逐題基準）
                    if g in items[1]:
                        path.parent.mkdir(parents=True, exist_ok=True)
                        path.write_text(json.dumps(items[1][g], ensure_ascii=False, indent=1) + "\n", encoding="utf-8")
    print("\n每則平均時間" + ("" if args.no_cache else "（含快取命中）") + "："
          + "、".join(f"{k} {v:.2f} 秒" for k, v in avg.items()))
    print("量測時間：" + "、".join(f"{k} {v / 60:.1f} 分鐘" for k, v in seconds.items())
          + f"，合計 {sum(seconds.values()) / 60:.1f} 分鐘")
    if haiku:
        # run_flows.py 不轉印 go/bin/chat 的 stderr，流程題的用量讀不到
        print(f"AI 用量（實際呼叫，不含流程題）：輸入 {usage[0]} token、輸出 {usage[1]} token、估算 US${usage[2]:.4f}")
    if failed:
        print("\n⚠️ 量測中有 AI 呼叫失敗，這次量測無效：在機器空閒時重跑（AGENTS.md）")
        sys.exit(2)

    if args.update_baseline:
        print()
        write_baseline(res)
        if items is not None:
            for g, sets in items[1].items():
                ITEMS[g].write_text(json.dumps(sets, ensure_ascii=False, indent=1) + "\n", encoding="utf-8")
                print(f"已更新 {ITEMS[g].relative_to(ROOT)}")
    if haiku:
        print("\nHaiku 的結果只列出，不和基準比較")
        sys.exit(0)
    print("\n❌ 有退步" if worse else "\n✅ 沒有退步")
    sys.exit(1 if worse else 0)


if __name__ == "__main__":
    main()
