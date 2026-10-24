#!/usr/bin/env python3
"""逐題比對（Go 版）：記下各題組（開發＋驗證，兩家測試商家）每一題的結果，確認行為不變的重構前後結果完全相同。

執行：python3 eval/snapshot.py save <名稱> [--exe <執行檔>] [--store tea|salon|all]   跑完各題組，存下每一題的結果（預設先編譯 go/bin/chat）
      python3 eval/snapshot.py compare <名稱A> <名稱B>         逐題比對兩次的結果（全部相同 exit 0，有不同 exit 1）
結果存到 eval/results/snapshots/<名稱>/{dev,salon,holdout}/<題組>.json（eval/results/ 不進 git）
- 題組：開發是流程題、分類、單題（testset.json＋faq_set.json）、誘答與跨語言（testset_extra.json）；驗證是流程題、分類、faq_set.json 的驗證題
- 第二家測試商家（霧光髮廊，美髮沙龍，#66）是 salon：商家設定讀 eval/stores/salon/config.json（BOT_CONFIG_FILE），
  題組是同一個資料夾的流程題、分類、單題（含誘答題）；--store salon 只量沙龍，tea 只量森日茶屋
- 跑法和 eval/verify.py 相同：每題 1 輪、依序執行（本機模型同時跑兩個 AI 呼叫，品質會變差）
- 比的是每一題的通過與否、分類標籤、回覆種類、挑到的 FAQ；同一種回覆有多種說法會隨機挑，所以不比回覆文字
- 驗證題組只由程式比對：save 不印驗證題組的任何輸出，compare 只印「相同」或「不同 N 題」，不印題目和代號
- 題目的預約日期是相對今天算的，前後兩次要在同一天跑
- 預設開 AI 回答快取（BOT_AI_CACHE=1，會傳給 go/bin/chat；設 BOT_AI_CACHE=0 關掉）：同樣的 AI 請求直接用存好的回答，驗證題組的快取分開存（BOT_AI_CACHE_NS=holdout）
"""
import argparse, json, os, re, subprocess, sys, tempfile, time
from testdates import resolve
from pathlib import Path

EVAL = Path(__file__).resolve().parent
ROOT = EVAL.parent
os.environ.setdefault("BOT_AI_CACHE", "1")   # 預設開 AI 回答快取；BOT_AI_CACHE=0 關掉
SNAPSHOTS = EVAL / "results" / "snapshots"
SALON = EVAL / "stores" / "salon"   # 第二家測試商家（美髮沙龍，#66）：商家設定和開發題組
KINDS = {"dev": ("flows", "intent", "testset", "extra"), "salon": ("flows", "intent", "testset"),
         "holdout": ("flows", "intent", "testset")}
GROUPS = {"dev": "開發", "salon": "沙龍", "holdout": "驗證"}
STORE_GROUPS = {"tea": ("dev", "holdout"), "salon": ("salon",), "all": ("dev", "salon", "holdout")}
CACHE_RE = re.compile(r"AI 快取[^：\n]*：命中 (\d+) 次、實際呼叫 (\d+) 次")
cache_total = [0, 0]   # 命中、實際呼叫（全部題組加總）
NAMES = {"flows": "完整流程", "intent": "分類", "testset": "單題", "extra": "誘答與跨語言"}


def run(cmd, env_extra, log):
    p = subprocess.run(cmd, env=dict(os.environ, **env_extra), capture_output=True, text=True,
                       stdin=subprocess.DEVNULL, timeout=7200)
    log.write_text(p.stdout + p.stderr, encoding="utf-8")
    # 有開快取（BOT_AI_CACHE=1）時，go/bin/chat 和 run_flows.py 會印命中幾次、實際呼叫幾次
    for hits, misses in CACHE_RE.findall(p.stdout):
        cache_total[0] += int(hits); cache_total[1] += int(misses)
    return p


def flows(exe, tests, log, holdout=False, env=None):
    """每一步：通過與否、回覆種類。key = 對話 id#第幾步。env 是另外給 go/bin/chat 的環境變數（沙龍的 BOT_CONFIG_FILE）"""
    with tempfile.TemporaryDirectory() as tmp:
        summary = Path(tmp) / "summary.json"
        # 驗證題組的 AI 回答快取分開存（run_flows.py 用 --cache-ns 傳給 go/bin/chat）
        run([sys.executable, str(EVAL / "run_flows.py"), "--bot", "go-chat", "--exe", exe,
             "--tests", str(tests), "--runs", "1", "--json", str(summary), "--cache-ns", "holdout" if holdout else "dev"],
            env or {}, log)
        if not summary.exists(): return None
        detail = json.loads(summary.read_text(encoding="utf-8"))["detail"]
    out, step = {}, {}
    for x in json.loads(Path(detail).read_text(encoding="utf-8"))[0]:
        step[x["conv"]] = step.get(x["conv"], 0) + 1
        out[f"{x['conv']}#{step[x['conv']]}"] = {"pass": x["pass"], "key": x["got"].get("key")}
    return out


def intent(exe, path, log, holdout=False, env=None):
    """每一句：答對記 ok，答錯記分類結果。stdout 只列答錯的句子（依題目順序），其餘就是答對。key = 第幾句。"""
    env = {**(env or {}), "BOT_INTENT_EVAL": str(path), "ROUNDS": "1"}
    if holdout: env["BOT_AI_CACHE_NS"] = "holdout"
    out_text = run([exe], env, log).stdout
    m = re.search(r"總計：(\d+)/(\d+)", out_text)
    start = out_text.find("答錯的（各輪結果）==========")
    if not m or start < 0: return None
    wrong = [l for l in out_text[start:m.start()].splitlines()[1:] if l.strip()]
    items = json.loads(resolve(Path(path).read_text(encoding="utf-8")))["items"]   # 和 bot 一樣換掉相對日期，才對得上輸出
    out, j = {}, 0
    for i, it in enumerate(items, 1):
        prefix = (f"［{it['state']}］" if it.get("state") else "") + f"{it['text']}　應為 {'/'.join(it['accept'])}　→ "
        if j < len(wrong) and wrong[j].startswith(prefix):
            out[str(i)] = wrong[j][len(prefix):]; j += 1
        else:
            out[str(i)] = "ok"
    # 對不上（題目和輸出不一致）就當作沒有結果，不要存一份錯的快照
    if j != len(wrong) or int(m[2]) != len(items) or int(m[1]) != len(items) - len(wrong): return None
    return out


def testset(exe, paths, holdout, log, langs=None, env=None):
    """每一題：結果（ok／wrong／missed）、回覆種類、挑到的 FAQ。key = 題目 id。
    go/bin/chat 把每題結果寫到 eval/results/testset-*.json，印「結果：路徑」；HOLDOUT=1 只跑驗證題"""
    env = {**(env or {}), "BOT_TESTSET": ",".join(str(p) for p in paths)}
    if langs: env["BOT_LANGUAGES"] = langs   # 跨語言題有韓文，店家多開 ko（和 eval/verify.py 相同）
    if holdout: env.update(HOLDOUT="1", BOT_AI_CACHE_NS="holdout")
    m = re.search(r"^結果：(.+)$", run([exe], env, log).stdout, re.M)
    if not m or not Path(m[1]).exists(): return None
    return {r["id"]: {"result": r["result"], "key": r["key"], "pick": r["pick"]}
            for r in json.loads(Path(m[1]).read_text(encoding="utf-8"))}


def save(name, exe, groups=STORE_GROUPS["all"]):
    if exe is None:
        p = subprocess.run(["go", "-C", str(ROOT / "go"), "build", "-o", "bin/chat", "./cmd/chat"], capture_output=True, text=True)
        if p.returncode != 0: raise SystemExit("編譯 go/bin/chat 失敗：\n" + p.stdout + p.stderr)
        exe = ROOT / "go" / "bin" / "chat"
    exe = os.path.abspath(exe)
    if not os.access(exe, os.X_OK): raise SystemExit(f"找不到執行檔：{exe}")
    t0 = time.time()
    salon_env = {"BOT_CONFIG_FILE": str(SALON / "config.json")}
    jobs = {
        "dev": {
            "flows": lambda log: flows(exe, EVAL / "flow_tests.json", log),
            "intent": lambda log: intent(exe, EVAL / "intent_set.json", log),
            "testset": lambda log: testset(exe, [EVAL / "testset.json", EVAL / "faq_set.json"], False, log),
            "extra": lambda log: testset(exe, [EVAL / "testset_extra.json"], False, log, langs="zh,en,ja,ko"),
        },
        # 沙龍：商家設定改讀 eval/stores/salon/config.json；誘答題和其他單題放在同一個檔
        "salon": {
            "flows": lambda log: flows(exe, SALON / "flow_tests.json", log, env=salon_env),
            "intent": lambda log: intent(exe, SALON / "intent_set.json", log, env=salon_env),
            "testset": lambda log: testset(exe, [SALON / "testset.json"], False, log, env=salon_env),
        },
        "holdout": {
            "flows": lambda log: flows(exe, EVAL / "holdout" / "flow_holdout.json", log, holdout=True),
            "intent": lambda log: intent(exe, EVAL / "holdout" / "intent_holdout.json", log, holdout=True),
            "testset": lambda log: testset(exe, [EVAL / "faq_set.json"], True, log),
        },
    }
    jobs = {g: k for g, k in jobs.items() if g in groups}
    failed = False
    for group, kinds in jobs.items():
        d = SNAPSHOTS / name / group; d.mkdir(parents=True, exist_ok=True)
        for kind, job in kinds.items():
            label = f"{GROUPS[group]} {NAMES[kind]}"
            print(f"{label}：執行中…", flush=True)
            res = job(d / f"{kind}.log")
            if res is None:
                failed = True
                # 驗證題組的輸出可能含題目，不印出來
                print(f"{label}：❌ 沒有結果（輸出在 {d / f'{kind}.log'}）" if group != "holdout" else f"{label}：❌ 沒有結果")
                continue
            (d / f"{kind}.json").write_text(json.dumps(res, ensure_ascii=False, indent=1), encoding="utf-8")
            print(f"{label}：已記錄 {len(res)} 題")
    print(f"快照：{SNAPSHOTS / name}")
    if os.environ.get("BOT_AI_CACHE") == "1":
        print(f"AI 快取（全部題組）：命中 {cache_total[0]} 次、實際呼叫 {cache_total[1]} 次")
    print(f"耗時 {time.time() - t0:.0f} 秒")
    if failed: sys.exit(1)


def passed(kind, v):
    """一題算不算答對：流程題看 pass、分類是 "ok"、單題和誘答與跨語言看 result 是 ok。"""
    if kind == "flows":
        return bool(v) and bool(v.get("pass"))
    if kind == "intent":
        return v == "ok"
    return bool(v) and v.get("result") == "ok"


ITEM_SOURCES = {("dev", "intent"): EVAL / "intent_set.json", ("salon", "intent"): SALON / "intent_set.json",
                ("holdout", "intent"): EVAL / "holdout" / "intent_holdout.json"}


def pass_sets(name, groups=STORE_GROUPS["all"]):
    """快照裡每個題組答對的題目代號 {group: {kind: [代號, ...]}}。分類題組的代號是第幾句（從 1 起算），換成句子本身
    （有狀態的前面加「［狀態］」，同一句在不同狀態是不同題），題組中間插入新句子時才不會整批對錯位。沒有快照的題組不列。"""
    out = {}
    for group in groups:
        for kind in KINDS[group]:
            f = SNAPSHOTS / name / group / f"{kind}.json"
            if not f.exists():
                continue
            res = json.loads(f.read_text(encoding="utf-8"))
            src = ITEM_SOURCES.get((group, kind))
            texts = [(f"［{i['state']}］" if i.get("state") else "") + i["text"]
                     for i in json.loads(src.read_text(encoding="utf-8"))["items"]] if src else None
            # intent() 的代號從 1 起算（2026-10-13 前誤用成從 0 起算，列出的句子都是答錯那句的下一句）
            key = (lambda k: texts[int(k) - 1] if texts and k.isdigit() and 1 <= int(k) <= len(texts) else k)
            out.setdefault(group, {})[kind] = sorted(key(k) for k, v in res.items() if passed(kind, v))
            out[group].setdefault("_all", {})[kind] = sorted(key(k) for k in res)
    return out


def item_changes(base, cur):
    """逐題比較兩份 pass_sets：{group: {"worse": [(kind, 代號)], "better": [...]}}。只比兩邊都有的題目。"""
    out = {}
    for group, kinds in cur.items():
        b = base.get(group, {})
        worse, better = [], []
        for kind, ok in kinds.items():
            if kind == "_all" or kind not in b:
                continue
            ok, was = set(ok), set(b[kind])
            both = set(kinds["_all"][kind]) & set(b.get("_all", {}).get(kind, []))
            worse += [(kind, k) for k in sorted((was - ok) & both)]
            better += [(kind, k) for k in sorted((ok - was) & both)]
        out[group] = {"worse": worse, "better": better}
    return out


def compare(a, b):
    """逐題比對兩次快照；只比任一邊有存的題組（save --store 只跑一家店時，另一家不算缺少）。"""
    same = True
    for group, kinds in KINDS.items():
        if not any((SNAPSHOTS / n / group).exists() for n in (a, b)):
            continue
        for kind in kinds:
            label = f"{GROUPS[group]} {NAMES[kind]}"
            fa, fb = (SNAPSHOTS / n / group / f"{kind}.json" for n in (a, b))
            if not fa.exists() or not fb.exists():
                same = False; print(f"{label}：❌ 缺少快照（{fa if not fa.exists() else fb}）"); continue
            ra, rb = (json.loads(f.read_text(encoding="utf-8")) for f in (fa, fb))
            diff = [k for k in list(ra) + [k for k in rb if k not in ra] if ra.get(k) != rb.get(k)]
            same = same and not diff
            print(f"{label}：相同（{len(ra)} 題）" if not diff else f"{label}：不同 {len(diff)} 題")
            if group != "holdout":   # 驗證題組只給數字，不給題目和代號
                for k in diff:
                    print(f"  {k}：{json.dumps(ra.get(k), ensure_ascii=False)} → {json.dumps(rb.get(k), ensure_ascii=False)}")
    print("\n✅ 全部相同" if same else "\n❌ 有不同")
    sys.exit(0 if same else 1)


def main():
    ap = argparse.ArgumentParser()
    sub = ap.add_subparsers(dest="cmd", required=True)
    s = sub.add_parser("save", help="跑完各題組，存下每一題的結果")
    s.add_argument("name")
    s.add_argument("--exe", help="要測的執行檔（預設先編譯 go/bin/chat 再測它）")
    s.add_argument("--store", choices=list(STORE_GROUPS), default="all",
                   help="量哪家店：tea 森日茶屋（開發＋驗證）、salon 霧光髮廊（開發），預設 all 兩家都量")
    c = sub.add_parser("compare", help="逐題比對兩次快照")
    c.add_argument("a"); c.add_argument("b")
    args = ap.parse_args()
    save(args.name, args.exe, STORE_GROUPS[args.store]) if args.cmd == "save" else compare(args.a, args.b)


if __name__ == "__main__":
    main()
