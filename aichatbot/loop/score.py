#!/usr/bin/env python3
"""⚠️ 已停用（2026-10-09，#33）：Go 版的評分改用 python3 eval/verify.py，loop/gate.py 也改呼叫它。
這支只量 Swift 版（tag swift-final），留著是因為舊的紀錄（ledger.jsonl、NOTES.md）會提到它。

評分：用 swift build 把 bot 執行檔（Sources/BotCore＋Sources/bot）編譯一次，依序跑開發題組與驗證題組，輸出一份 JSON 分數。

執行：python3 loop/score.py [--out loop/runs/<名稱>] [--flow-runs 1] [--rounds 1]
（bot 用 greedy 取樣，同一句話每次結果相同，所以預設各跑 1 次就夠）
- 開發題組（dev）：flow_tests.json、intent_set.json、testset.json＋faq_set.json（非驗證句）
  失敗明細寫到 <out>/dev_report.md，給改進迴圈看
- 驗證題組（holdout）：eval/holdout/*、faq_set.json 的 holdout 題
  明細只放在 <out>/holdout/，改進迴圈不能看，只看總分
- 一律依序執行：本機模型同時跑兩個 AI 呼叫，品質會明顯變差
"""
import argparse, json, os, re, shutil, subprocess, sys, time

os.environ.setdefault("BOT_AI_CACHE", "1")   # 預設開 AI 回答快取（同樣的請求直接用存好的回答）；BOT_AI_CACHE=0 關掉
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
EVAL = ROOT / "eval"
BIN = ROOT / "loop" / "bin" / "bot"


def compile_bot():
    BIN.parent.mkdir(parents=True, exist_ok=True)
    p = subprocess.run(["swift", "build", "-c", "release", "--product", "bot"], cwd=ROOT, capture_output=True, text=True)
    if p.returncode != 0:
        out = p.stdout + p.stderr
        errs = [l for l in out.splitlines() if "error" in l]
        raise SystemExit("編譯失敗：\n" + "\n".join(errs[:20] or out.splitlines()[-20:]))
    bin_dir = subprocess.run(["swift", "build", "-c", "release", "--show-bin-path"], cwd=ROOT, capture_output=True, text=True).stdout.strip()
    # 複製一份到 loop/bin/：評分途中 .build/ 被重新編譯也不會換掉正在測的執行檔
    # 先刪再複製：直接覆寫已簽章的執行檔，macOS 會用舊的簽章快取把它擋掉（killed: 9）
    BIN.unlink(missing_ok=True)
    shutil.copy2(Path(bin_dir) / "bot", BIN)


# AI 回答快取（預設開，BOT_AI_CACHE=0 關掉）：bot、run_flows.py 會印命中幾次、實際呼叫幾次，這裡加總寫進 score.json
CACHE_RE = re.compile(r"AI 快取：命中 (\d+) 次、實際呼叫 (\d+) 次")
cache_total = [0, 0]


def count_cache(text):
    for hits, misses in CACHE_RE.findall(text):
        cache_total[0] += int(hits); cache_total[1] += int(misses)


def run_flows(tests, runs, out_json, log):
    cmd = [sys.executable, str(EVAL / "run_flows.py"), "--bot", "bot", "--exe", str(BIN),
           "--tests", tests, "--runs", str(runs), "--json", str(out_json)]
    p = subprocess.run(cmd, capture_output=True, text=True, timeout=7200)
    log.write_text(p.stdout + p.stderr, encoding="utf-8")
    count_cache(p.stdout)   # run_flows.py 會依題庫位置自己決定快取命名空間
    if not out_json.exists():
        raise SystemExit(f"流程測試失敗：{tests}\n{p.stdout[-1500:]}{p.stderr[-1500:]}")
    d = json.loads(out_json.read_text(encoding="utf-8"))
    return {"pct": round(d["pct"], 2), "pass": d["avg_pass"], "total": d["total"], "per_run": d["per_run"]}


def run_env(env_extra, log):
    p = subprocess.run([str(BIN)], env=dict(os.environ, **env_extra), capture_output=True, text=True,
                       stdin=subprocess.DEVNULL, timeout=7200)
    log.write_text(p.stdout + p.stderr, encoding="utf-8")
    count_cache(p.stdout)
    return p.stdout


def run_intent(path, rounds, log, holdout=False):
    out = run_env({"BOT_INTENT_EVAL": str(path), "ROUNDS": str(rounds), "BOT_AI_CACHE_NS": "holdout" if holdout else "dev"}, log)
    m = re.search(r"總計：(\d+)/(\d+)", out)
    if not m: raise SystemExit(f"分類測試沒有結果：{path}\n{out[-1500:]}")
    ok, n = int(m[1]), int(m[2])
    unstable = int(re.search(r"不穩定：(\d+)", out)[1])
    return {"pct": round(ok / n * 100, 2), "ok": ok, "total": n, "unstable": unstable}


def run_testset(paths, rounds, holdout, log):
    env = {"BOT_TESTSET": ",".join(str(p) for p in paths), "ROUNDS": str(rounds)}
    if holdout: env["HOLDOUT"] = "1"
    env["BOT_AI_CACHE_NS"] = "holdout" if holdout else "dev"   # 驗證題組的快取分開存
    out = run_env(env, log)
    m = re.search(r"全部：(\d+)/(\d+).*?答錯 (\d+).*?該答沒答 (\d+)", out)
    if not m: raise SystemExit(f"單題測試沒有結果\n{out[-1500:]}")
    ok, n, wrong, missed = map(int, m.groups())
    return {"pct": round(ok / n * 100, 2), "ok": ok, "total": n, "wrong": wrong, "missed": missed}


def section(text, start):
    """取出輸出中某個標題之後到下一個 ===== 標題之前的內容。"""
    i = text.find(start)
    if i < 0: return ""
    rest = text[i + len(start):]
    j = rest.find("==========")
    return rest[:j if j >= 0 else None].strip()


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--out", default=str(ROOT / "loop" / "runs" / time.strftime("%Y%m%d-%H%M%S")))
    ap.add_argument("--flow-runs", type=int, default=1)
    ap.add_argument("--rounds", type=int, default=1)
    ap.add_argument("--skip-holdout", action="store_true")
    a = ap.parse_args()
    out = Path(a.out); (out / "holdout").mkdir(parents=True, exist_ok=True)
    t0 = time.time()
    compile_bot()

    dev = {
        "flows": run_flows("flow_tests.json", a.flow_runs, out / "flows.json", out / "flows.txt"),
        "intent": run_intent(EVAL / "intent_set.json", a.rounds, out / "intent.txt"),
        "testset": run_testset([EVAL / "testset.json", EVAL / "faq_set.json"], a.rounds, False, out / "testset.txt"),
    }
    dev["score"] = round(sum(dev[k]["pct"] for k in ("flows", "intent", "testset")) / 3, 2)

    hold = None
    if not a.skip_holdout:
        h = out / "holdout"
        hold = {
            "flows": run_flows(str(EVAL / "holdout" / "flow_holdout.json"), a.flow_runs, h / "flows.json", h / "flows.txt"),
            "intent": run_intent(EVAL / "holdout" / "intent_holdout.json", a.rounds, h / "intent.txt", holdout=True),
            "faq": run_testset([EVAL / "faq_set.json"], a.rounds, True, h / "faq.txt"),
        }
        hold["score"] = round(sum(hold[k]["pct"] for k in ("flows", "intent", "faq")) / 3, 2)

    # 開發題組的失敗明細（改進迴圈靠這份找問題）
    flows_txt = (out / "flows.txt").read_text(encoding="utf-8")
    intent_txt = (out / "intent.txt").read_text(encoding="utf-8")
    testset_txt = (out / "testset.txt").read_text(encoding="utf-8")
    report = [f"# 開發題組失敗明細（{time.strftime('%Y-%m-%d %H:%M')}）", "",
              f"綜合分數 {dev['score']}：流程 {dev['flows']['pct']}%、分類 {dev['intent']['pct']}%（不穩定 {dev['intent']['unstable']} 題）、"
              f"單題 {dev['testset']['pct']}%（答錯 {dev['testset']['wrong']}、該答沒答 {dev['testset']['missed']}）", "",
              "## 完整流程（通過次數／總次數）", "```", section(flows_txt, "有失敗的步驟（通過次數／總次數）=========="), "```", "",
              "## 分類答錯", "```", section(intent_txt, "答錯的（各輪結果）=========="), "```", "",
              "## 單題沒通過（✗答錯 最嚴重；△沒答 是該答沒答）", "```", section(testset_txt, "沒通過的（各輪：結果／挑 FAQ／回覆種類）=========="), "```"]
    (out / "dev_report.md").write_text("\n".join(report) + "\n", encoding="utf-8")

    result = {"time": time.strftime("%Y-%m-%d %H:%M:%S"), "seconds": round(time.time() - t0), "dev": dev, "holdout": hold, "dir": str(out)}
    if os.environ.get("BOT_AI_CACHE") == "1": result["cache"] = {"hits": cache_total[0], "misses": cache_total[1]}
    (out / "score.json").write_text(json.dumps(result, ensure_ascii=False, indent=1), encoding="utf-8")
    print(json.dumps(result, ensure_ascii=False))


if __name__ == "__main__":
    main()
