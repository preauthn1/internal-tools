#!/usr/bin/env python3
"""守門：檢查改進迴圈的修改、評分、決定保留（git commit）或還原。決策完全由這支程式做，不交給 AI。

執行：python3 loop/gate.py --baseline            量目前 HEAD 的分數，存成 loop/baseline.json
      python3 loop/gate.py                       評估工作目錄的修改（說明讀 loop/proposal.json）
      python3 loop/gate.py --dry-run             只做檢查、不評分

Go 版（2026-10-09，#33）：只允許改提示詞和內建意圖案例，評分用 eval/verify.py。

保留條件（全部成立）：
  1. 只改了允許的檔案：內建意圖案例 eval/examples.jsonl、faq.txt（只能改問法）、提示詞所在的 Go 檔
     （PROMPT_FILES，只能改字串內容：改動的行拿掉字串後要和原本的一樣，整行新增或刪除的要是純字串）；
     提示詞不能寫死測試句（意圖說明的例句逐句比對）
  2. 案例庫和所有題庫（含驗證題）沒有重複
  2b. 測試全部通過（python3 eval/verify.py --checks-only，不用本機模型；在評分前擋下改壞的行為）
  3. 開發綜合分數 ≥ 基準 + MIN_GAIN，
     或「答錯少了 MIN_WRONG_DROP 次以上且開發分數沒有下降」（寧可轉店員，也不要答錯），
     或「跨語言多答對、開發分數沒降、答錯沒增加」（跨語言不在開發綜合分數裡）
  4. 單題「答錯」和「誘答題 編造」（會誤導顧客的）不能比基準多超過 MAX_WRONG_UP；跨語言答對不能變少
  5. 驗證題組沒有任何一題從答對變答錯（逐題比對，eval/verify.py 的逐題基準放在 loop/items/），
     驗證綜合分數也不能比基準低超過 MAX_HOLDOUT_DROP（偵測過度擬合）
  6. 第二家測試商家（霧光髮廊，#66）不能退步：沒有題目從答對變答錯，分類、單題、流程題的通過數不能少，
     答錯、誘答題編造不能多（分母不同或基準沒有沙龍指標時略過那一項）。沙龍不算進開發綜合分數，只當守門條件
  greedy 取樣結果固定，任何一題的差異都是真的，所以門檻都是 1 題。
  開發綜合分數＝分類、單題、流程題三個通過率的平均（×100）；驗證綜合分數＝驗證分類、驗證流程題的平均。
"""
import argparse, collections, json, re, shutil, subprocess, sys, time
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
LOOP = ROOT / "loop"
EVAL = ROOT / "eval"
BASELINE = LOOP / "baseline.json"
ITEMS_DIR = LOOP / "items"   # 逐題基準（dev.json、salon.json、holdout/items.json），量基準和保留時更新；不進 git，迴圈的 AI 不能讀 holdout/
SALON = EVAL / "stores" / "salon"
# 沙龍的守門指標：代號 → 越大越好?
SALON_GUARDS = {"salon.intent.ten": True, "salon.testset.ok": True, "salon.flows.pass": True,
                "salon.testset.wrong": False, "salon.testset.trap_fabricated": False}
LEDGER = LOOP / "ledger.jsonl"
LEDGER_MD = LOOP / "LEDGER.md"
PROPOSAL = LOOP / "proposal.json"

# 提示詞寫在這幾個 Go 檔的 Instructions 字串裡：檔案可以改，但改動的行只能是字串（見 prompt_only）
PROMPT_FILES = {"go/internal/dialog/answer.go", "go/internal/dialog/classify.go", "go/internal/dialog/booking.go"}
ALLOWED = {"eval/examples.jsonl", "eval/faq.txt"} | PROMPT_FILES
allowed = lambda f: f in ALLOWED
MIN_GAIN = 0.2          # greedy 結果固定，任何差異都是真的。修好 1 題約 +0.28（單題 118）～+0.33（分類 101、流程 104 步），
                        # 0.2 = 至少淨修好 1 題，且擋下「修好一題、改壞另一題」的交換。
                        # 題庫變大時要重算：例如單題 300 題時修好 1 題只有 +0.11
MAX_WRONG_UP = 0        # greedy 結果固定、每題 1 輪：多一個答錯就是真的退步
MIN_WRONG_DROP = 1      # 答錯少這麼多題、開發分數沒降，就算綜合分數進步不到 MIN_GAIN 也保留（greedy 結果固定，少 1 題就是真的少）
MAX_HOLDOUT_DROP = 3.0


def sh(*cmd, check=True):
    p = subprocess.run(cmd, cwd=ROOT, capture_output=True, text=True)
    if check and p.returncode != 0:
        raise SystemExit(f"指令失敗：{' '.join(cmd)}\n{p.stderr}")
    return p.stdout


LOOP_STATE = {"loop/NOTES.md", "loop/LEDGER.md", "loop/ledger.jsonl", "loop/proposal.json"}   # 迴圈自己的紀錄，不算修改


SCOPE = ("eval", "go", "web", "Sources", "Tests", "Package.swift")   # 迴圈可能改到的地方（含不允許的，才抓得到）


def changed_files():
    # 只看 SCOPE：專案其他地方的修改（使用者自己在改的文件等）不是迴圈改的，不能還原（曾把使用者改的 docs/TODO.md 還原掉）
    tracked = [l for l in sh("git", "diff", "--name-only", "HEAD", "--", *SCOPE).splitlines() if l and l not in LOOP_STATE]
    # 只看 SCOPE 底下的新檔：專案其他地方沒被追蹤的檔案（使用者自己的筆記、金鑰等）不是迴圈改的，不能動
    untracked = [l for l in sh("git", "ls-files", "--others", "--exclude-standard", "--", *SCOPE).splitlines() if l and l not in LOOP_STATE]
    return tracked, untracked


norm = lambda t: re.sub(r"[\W_]+", "", t.lower())


def test_sentences():
    s = set()
    for f in ["intent_set.json", "holdout/intent_holdout.json"]:
        s |= {norm(i["text"]) for i in json.loads((EVAL / f).read_text())["items"]}
    for f in ["flow_tests.json", "holdout/flow_holdout.json"]:
        s |= {norm(t["user"]) for c in json.loads((EVAL / f).read_text())["conversations"] for t in c["turns"]}
    for f in ["testset.json", "faq_set.json", "testset_extra.json"]:
        s |= {norm(c["q"]) for c in json.loads((EVAL / f).read_text())["cases"]}
    if SALON.exists():   # 第二家測試商家的開發題組（#66）
        s |= {norm(i["text"]) for i in json.loads((SALON / "intent_set.json").read_text())["items"]}
        s |= {norm(t["user"]) for c in json.loads((SALON / "flow_tests.json").read_text())["conversations"] for t in c["turns"]}
        s |= {norm(c["q"]) for c in json.loads((SALON / "testset.json").read_text())["cases"]}
    return {x for x in s if len(x) >= 3 and not x.startswith("pending")}


def faq_answers(text):
    out = {}
    for line in text.splitlines():
        if line.strip() and not line.startswith("#"):
            parts = [p.strip() for p in line.split(" | ")]
            out[parts[0]] = parts[2:]
    return out


def static_checks(tracked, untracked):
    """回傳問題清單（空的代表通過）。"""
    problems = []
    for f in tracked + untracked:
        if not allowed(f):
            problems.append(f"改了不允許的檔案：{f}")
    if "eval/faq.txt" in tracked:
        old, new = faq_answers(sh("git", "show", "HEAD:eval/faq.txt")), faq_answers((EVAL / "faq.txt").read_text())
        for k in new.keys() - old.keys(): problems.append(f"faq.txt 新增了題目 {k}（答案是店家事實，只能由店家新增）")
        for k in old.keys() & new.keys():
            if old[k] != new[k]: problems.append(f"faq.txt 改了 {k} 的答案（只能改問法）")
    prompt_changed = [f for f in tracked if f in PROMPT_FILES]
    if prompt_changed:
        diff = sh("git", "diff", "-U0", "HEAD", "--", *prompt_changed).splitlines()
        removed = [l[1:] for l in diff if l.startswith("-") and not l.startswith("---")]
        added = [l[1:] for l in diff if l.startswith("+") and not l.startswith("+++")]
        for line in code_changes(removed, added):
            problems.append(f"提示詞檔改了字串以外的程式：{line.strip()[:100]!r}")
        tests = test_sentences()
        for line in added:
            for lit in re.findall(r'"([^"\\]{2,})"', line):
                # 意圖說明的例句是一串用頓號、逗號隔開的句子，逐句比對
                for piece in [lit] + re.split(r"[、，,。：:；;／/（）()「」]", lit):
                    if norm(piece) in tests:
                        problems.append(f"提示詞寫死了測試句：{piece.strip()!r}")
    p = subprocess.run([sys.executable, str(EVAL / "check_overlap.py")], cwd=ROOT, capture_output=True, text=True)
    if p.returncode != 0:
        problems.append("案例庫和題庫重複：\n" + p.stdout.splitlines()[0] + "\n" + "\n".join(p.stdout.splitlines()[1:6]))
    return problems


STRING_RE = re.compile(r'"(?:[^"\\]|\\.)*"|`[^`]*`')


def prompt_only(line):
    """提示詞檔的一行改動是不是只有字串：去掉字串後只剩 Instructions:、+、逗號、冒號（提示用的 map）、括號、
    變數或欄位名稱（提示詞會串接店家名稱、語言這類變數）。空行和註解也算。"""
    t = line.strip()
    if not t or t.startswith("//"):
        return True
    if not STRING_RE.search(t):
        return False
    rest = STRING_RE.sub("", t)
    rest = re.sub(r"^Instructions:", "", rest.strip())
    return (re.fullmatch(r"[\w.\s+,():]*", rest) is not None and ":=" not in rest
            and not re.search(r"\b(if|for|return|func|var)\b", rest))


def skeleton(line):
    """把字串內容拿掉後剩下的程式（字串換成 ""），用來判斷一行改動是不是只改了字串。"""
    return STRING_RE.sub('""', line).strip()


def code_changes(removed, added):
    """提示詞檔的改動裡，改到程式的行。刪掉和新增的行先依骨架配對：骨架相同就是只改了字串內容
    （例如意圖說明 {"booking", "…"}、欄位說明 Description: "…"），放行；配不上的行（整行新增或刪除）
    要是純字串（prompt_only），否則算改到程式。"""
    left = collections.Counter(skeleton(l) for l in removed)
    out = []
    for line in added:
        k = skeleton(line)
        if left[k] > 0:
            left[k] -= 1
        elif not prompt_only(line):
            out.append(line)
    for line in removed:
        k = skeleton(line)
        if left[k] > 0:
            left[k] -= 1
            if not prompt_only(line):
                out.append(line)
    return out


def checks():
    """跑 eval/verify.py --checks-only（Swift、Go、web 的測試，不用本機模型），回傳問題說明（None 代表全部通過）。"""
    p = subprocess.run([sys.executable, str(EVAL / "verify.py"), "--checks-only"], cwd=ROOT, capture_output=True, text=True)
    if p.returncode == 0:
        return None
    lines = [l.strip() for l in p.stdout.splitlines() if l.strip()]
    failed = next((l for l in lines if l.startswith("❌")), "測試失敗")
    return f"{failed}：" + " ".join(lines[-5:])[:300]


def ratio(m, key):
    v = m.get(key)
    return v[0] / v[1] * 100 if v and v[1] else 0.0


def score(out):
    """跑 eval/verify.py 量開發和驗證題組，回傳 ({dev, holdout, metrics}, 錯誤說明)。"""
    t0 = time.time()
    cmd = [sys.executable, str(EVAL / "verify.py"), "--no-checks", "--holdout", "--out", str(out)]
    if ITEMS_DIR.exists():
        cmd += ["--items-base", str(ITEMS_DIR)]
    p = subprocess.run(cmd, cwd=ROOT, capture_output=True, text=True)
    seconds = round(time.time() - t0)
    if p.returncode == 2 or not (out / "metrics.json").exists():
        return None, "量測無效（連不上模型服務或 AI 呼叫失敗）：" + (p.stdout + p.stderr)[-1500:]
    data = json.loads((out / "metrics.json").read_text())
    m = data["metrics"]
    dev = round(sum(ratio(m, k) for k in ("intent.ten", "testset.ok", "flows.pass")) / 3, 2)
    hold = round(sum(ratio(m, k) for k in ("holdout.intent.ten", "holdout.flows.pass")) / 2, 2)
    wrong = (m.get("testset.wrong") or [0])[0] + (m.get("testset.trap_fabricated") or [0])[0]
    return {"dev": {"score": dev, "testset": {"wrong": wrong}}, "holdout": {"score": hold}, "metrics": m,
            "items": data.get("items") or {}, "seconds": seconds}, ""


def save_items(out):
    """把這次答對的題目存成新的逐題基準（量基準、保留時）。"""
    for rel in ("dev.json", "salon.json", "holdout/items.json"):
        src = out / rel
        if src.exists():
            dst = ITEMS_DIR / rel
            dst.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy(src, dst)


def salon_regression(s, base):
    """第二家測試商家有沒有退步（#66）：回傳說明，沒退步是空字串。分母不同（題組改過）或基準沒有那一項就略過"""
    worse = s["items"].get("salon", {}).get("worse", 0)
    if worse:
        return f"霧光髮廊有 {worse} 題從答對變答錯"
    bm = base.get("metrics", {})
    for key, higher in SALON_GUARDS.items():
        cur, old = s["metrics"].get(key), bm.get(key)
        if not cur or not old or cur[1] != old[1]:
            continue
        if (cur[0] < old[0]) if higher else (cur[0] > old[0]):
            return f"霧光髮廊 {key} 退步：{old[0]}/{old[1]} → {cur[0]}/{cur[1]}"
    return ""


def revert(tracked, untracked, run_dir):
    run_dir.mkdir(parents=True, exist_ok=True)
    if tracked:   # 還原前先存下修改，之後還能重新檢視或重試
        (run_dir / "rejected.patch").write_text(sh("git", "diff", "HEAD", "--", *tracked), encoding="utf-8")
        sh("git", "checkout", "HEAD", "--", *tracked)
    if untracked:   # 不直接刪除，移到這次紀錄的 quarantine/
        q = run_dir / "quarantine"; q.mkdir(parents=True, exist_ok=True)
        for f in untracked:
            dest = q / f; dest.parent.mkdir(parents=True, exist_ok=True)
            shutil.move(str(ROOT / f), str(dest))


def cell(text, n):
    return text.replace("|", "／").replace("\n", " ")[:n]


def record(entry):
    with LEDGER.open("a", encoding="utf-8") as f:
        f.write(json.dumps(entry, ensure_ascii=False) + "\n")
    if not LEDGER_MD.exists():
        LEDGER_MD.write_text("# 改進迴圈紀錄\n\n每一行是一次嘗試（新的在下面）。⛔ = 被擋下或還原，✅ = 保留。\n"
                             "開發分數＝流程／分類／單題的平均；驗證分數只看總數。\n\n"
                             "| 時間 | 結果 | 開發 | 驗證 | 答錯 | 假設 | 原因 |\n|---|---|---|---|---|---|---|\n", encoding="utf-8")
    d, h = entry.get("dev") or {}, entry.get("holdout") or {}
    fmt = lambda x, dx: f"{x}（{dx:+.1f}）" if x is not None and dx is not None else "-"
    row = (f"| {entry['time'][5:16]} | {'✅' if entry['decision'] == 'accept' else '⛔'} "
           f"| {fmt(d.get('score'), entry.get('d_dev'))} | {fmt(h.get('score'), entry.get('d_hold'))} "
           f"| {d.get('testset', {}).get('wrong', '-')} | {cell(entry['hypothesis'], 80)} | {cell(entry['reason'], 120)} |\n")
    with LEDGER_MD.open("a", encoding="utf-8") as f:
        f.write(row)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--baseline", action="store_true")
    ap.add_argument("--dry-run", action="store_true")
    a = ap.parse_args()
    ts = time.strftime("%Y%m%d-%H%M%S")
    run_dir = LOOP / "runs" / ts

    if a.baseline:
        tracked, untracked = changed_files()
        if tracked:
            raise SystemExit(f"工作目錄有未提交的修改，先處理再量基準：{tracked}")
        s, err = score(run_dir)
        if not s: raise SystemExit(err)
        s["commit"] = sh("git", "rev-parse", "--short", "HEAD").strip()
        save_items(run_dir)
        BASELINE.write_text(json.dumps(s, ensure_ascii=False, indent=1))
        (LOOP / "runs" / "latest").unlink(missing_ok=True); (LOOP / "runs" / "latest").symlink_to(ts)
        print(f"基準：開發 {s['dev']['score']}　驗證 {s['holdout']['score']}　（{s['seconds']} 秒）")
        return

    if not a.dry_run and (not BASELINE.exists() or "metrics" not in json.loads(BASELINE.read_text())):
        # Swift 版的舊基準（沒有 metrics）不能和 Go 版的分數比
        raise SystemExit("還沒有 Go 版的基準，先跑 python3 loop/gate.py --baseline")
    base = json.loads(BASELINE.read_text()) if BASELINE.exists() else {}
    prop = json.loads(PROPOSAL.read_text()) if PROPOSAL.exists() else {}
    hypothesis = prop.get("hypothesis", "（沒有寫假設）")
    tracked, untracked = changed_files()
    entry = {"time": time.strftime("%Y-%m-%d %H:%M:%S"), "hypothesis": hypothesis, "category": prop.get("category"),
             "files": tracked + untracked, "base_commit": base.get("commit"), "run_dir": str(run_dir.relative_to(ROOT))}

    def finish(decision, reason, s=None):
        entry.update(decision=decision, reason=reason)
        if s:
            entry.update(dev=s["dev"], holdout=s["holdout"],
                         d_dev=round(s["dev"]["score"] - base["dev"]["score"], 2),
                         d_hold=round(s["holdout"]["score"] - base["holdout"]["score"], 2))
        if decision == "accept":
            sh("git", "add", *[f for f in tracked + untracked])
            msg = f"loop: {hypothesis}\n\n開發 {s['dev']['score']}（{entry['d_dev']:+.1f}）驗證 {s['holdout']['score']}（{entry['d_hold']:+.1f}）\n{reason}"
            sh("git", "commit", "-q", "-m", msg)
            entry["commit"] = sh("git", "rev-parse", "--short", "HEAD").strip()
            s["commit"] = entry["commit"]
            save_items(run_dir)
            BASELINE.write_text(json.dumps(s, ensure_ascii=False, indent=1))
        elif not a.dry_run:
            revert(tracked, untracked, run_dir)
        if not a.dry_run:
            record(entry)
        if not a.dry_run:
            PROPOSAL.unlink(missing_ok=True)
        if (run_dir / "dev_report.md").exists():
            (LOOP / "runs" / "latest").unlink(missing_ok=True); (LOOP / "runs" / "latest").symlink_to(ts)
        label = {"accept": "✅ 保留", "skip": "⏭️ 略過"}.get(decision, "⛔ 擋下（dry run，未還原）" if a.dry_run else "⛔ 還原")
        print(f"{label}：{reason}")
        print(json.dumps({k: entry.get(k) for k in ("d_dev", "d_hold", "commit")}, ensure_ascii=False))

    if not tracked and not untracked:
        return finish("skip", "沒有任何修改")
    problems = static_checks(tracked, untracked)
    if problems:
        return finish("reject", "檢查未通過：" + "；".join(problems))
    if (problem := checks()):
        return finish("reject", "檢查未通過：" + problem)
    if a.dry_run:
        print("✅ 檢查通過、測試全部通過（dry run，不評分）"); return

    s, err = score(run_dir)
    if not s:
        return finish("reject", err.strip().splitlines()[0][:200] if err.strip() else "評分失敗")
    d_dev = s["dev"]["score"] - base["dev"]["score"]
    d_hold = s["holdout"]["score"] - base["holdout"]["score"]
    d_wrong = s["dev"]["testset"]["wrong"] - base["dev"]["testset"]["wrong"]
    d_xlang = ((s["metrics"].get("testset.xlang_ok") or [0])[0]
               - (base.get("metrics", {}).get("testset.xlang_ok") or [0])[0])
    if d_wrong > MAX_WRONG_UP:
        return finish("reject", f"答錯或誘答題編造增加 {d_wrong} 次（會誤導顧客）", s)
    if d_xlang < 0:
        return finish("reject", f"跨語言少答對 {-d_xlang} 題", s)
    hold_worse = s["items"].get("holdout", {}).get("worse", 0)
    if hold_worse:
        return finish("reject", f"驗證題組有 {hold_worse} 題從答對變答錯（可能過度擬合開發題組）", s)
    if d_hold < -MAX_HOLDOUT_DROP:
        return finish("reject", f"驗證分數下降 {d_hold:.1f}（可能過度擬合開發題組）", s)
    if (why := salon_regression(s, base)):
        return finish("reject", why + "（改動只對森日茶屋有效）", s)
    if d_dev >= MIN_GAIN:
        return finish("accept", f"開發 {d_dev:+.1f}、驗證 {d_hold:+.1f}、答錯 {d_wrong:+d}", s)
    if d_xlang >= 1 and d_dev >= 0 and d_wrong <= 0:
        return finish("accept", f"跨語言多答對 {d_xlang} 題、開發分數沒降、答錯沒增加：開發 {d_dev:+.1f}、驗證 {d_hold:+.1f}", s)
    if d_wrong <= -MIN_WRONG_DROP and d_dev >= 0:
        return finish("accept", f"答錯少 {-d_wrong} 次且開發分數沒降：開發 {d_dev:+.1f}、驗證 {d_hold:+.1f}", s)
    return finish("reject", f"進步不夠：開發 {d_dev:+.1f}（需 ≥ +{MIN_GAIN}）、驗證 {d_hold:+.1f}", s)


if __name__ == "__main__":
    main()
