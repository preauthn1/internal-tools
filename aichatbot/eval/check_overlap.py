#!/usr/bin/env python3
"""檢查案例庫、FAQ 問法和各題庫之間有沒有重複的句子（案例和考題重複，分數就只是背答案）。
執行：python3 eval/check_overlap.py
- 案例庫 × 任何題庫：必須是 0（有就要改案例）
- 第二家測試商家（eval/stores/salon/，#66）：沙龍自己的案例 × 任何題庫、沙龍題組 × 森日茶屋題組，都必須是 0
- FAQ 問法 × 題庫：只列出來（問法本來就該寫常見說法；BOT_TESTSET 會另外統計排除後的分數）
- 題庫之間：只列出來（同一句在不同題庫測不同步驟是正常的）"""
import json, re, sys
from pathlib import Path

here = Path(__file__).parent
salon = here / "stores" / "salon"
norm = lambda t: re.sub(r"[\W_]+", "", t.lower())

tests = {}
for i in json.load(open(here / "intent_set.json"))["items"]: tests.setdefault(norm(i["text"]), set()).add("intent_set")
for c in json.load(open(here / "testset.json"))["cases"]: tests.setdefault(norm(c["q"]), set()).add("testset")
for c in json.load(open(here / "faq_set.json"))["cases"]: tests.setdefault(norm(c["q"]), set()).add("faq_set")
for c in json.load(open(here / "testset_extra.json"))["cases"]: tests.setdefault(norm(c["q"]), set()).add("testset_extra")
for conv in json.load(open(here / "flow_tests.json"))["conversations"]:
    for t in conv["turns"]: tests.setdefault(norm(t["user"]), set()).add("flow_tests")
for i in json.load(open(here / "holdout/intent_holdout.json"))["items"]: tests.setdefault(norm(i["text"]), set()).add("intent_holdout")
for conv in json.load(open(here / "holdout/flow_holdout.json"))["conversations"]:
    for t in conv["turns"]: tests.setdefault(norm(t["user"]), set()).add("flow_holdout")
# loop/test_gate.py 的暫時 repo 只複製森日茶屋的檔案，沒有沙龍的資料夾就略過沙龍的檢查
has_salon = salon.exists()
if has_salon:
    for i in json.load(open(salon / "intent_set.json"))["items"]: tests.setdefault(norm(i["text"]), set()).add("salon/intent_set")
    for c in json.load(open(salon / "testset.json"))["cases"]: tests.setdefault(norm(c["q"]), set()).add("salon/testset")
    for conv in json.load(open(salon / "flow_tests.json"))["conversations"]:
        for t in conv["turns"]:
            if not t["user"].startswith("/"): tests.setdefault(norm(t["user"]), set()).add("salon/flow_tests")

examples = [json.loads(l) for l in open(here / "examples.jsonl") if l.strip()]
faq_phrases = {norm(q) for line in open(here / "faq.txt") if line.strip() and not line.startswith("#")
               for q in line.split(" | ")[1].split("/")}

bad = [(e["text"], sorted(tests[norm(e["text"])])) for e in examples if norm(e["text"]) in tests]
print(f"案例庫 {len(examples)} 句 × 題庫：{'無重複 ✅' if not bad else f'{len(bad)} 句重複 ❌'}")
for t, where in bad: print(f"   {t}　（{'、'.join(where)}）")

bad_faq = [(e["text"], e["faq"]) for e in examples if e.get("faq") and norm(e["text"]) in faq_phrases]
if bad_faq: print(f"FAQ 案例和 FAQ 問法重複（多餘）：{bad_faq}")

leaked = sorted(t for t in tests if t in faq_phrases)
print(f"FAQ 問法 × 題庫：{len(leaked)} 句（只列出）：{'、'.join(leaked)}")

# 第二家測試商家：沙龍自己的案例（商家新增的，內建案例就是上面的案例庫）、沙龍題組和森日茶屋題組
bad_salon, cross = [], []
if has_salon:
    cfg = json.load(open(salon / "config.json"))
    salon_examples = [e["text"] for e in cfg["intent_examples"] if not e["builtin"]]
    bad_salon = [(t, sorted(tests[norm(t)])) for t in salon_examples if norm(t) in tests]
    print(f"沙龍案例 {len(salon_examples)} 句 × 題庫：{'無重複 ✅' if not bad_salon else f'{len(bad_salon)} 句重複 ❌'}")
    for t, where in bad_salon: print(f"   {t}　（{'、'.join(where)}）")
    cross = sorted((t, sorted(w)) for t, w in tests.items()
                   if any(x.startswith("salon/") for x in w) and any(not x.startswith("salon/") for x in w))
    print(f"沙龍題組 × 森日茶屋題組：{'無重複 ✅' if not cross else f'{len(cross)} 句重複 ❌'}")
    for t, where in cross:
        # 和驗證題組重複時不列句子（改進迴圈的 AI 看得到這段輸出，驗證題組不能被看出內容）
        hold = any(x.endswith("_holdout") for x in where)
        print(f"   {'（和驗證題組重複，不列句子）' if hold else t}　（{'、'.join(x for x in where if not x.endswith('_holdout'))}）")
    salon_faq = {norm(q) for f in cfg["faq"] for q in [f["question"], *f["alt_questions"]]}
    leaked = sorted(t for t, w in tests.items() if t in salon_faq and any(x.startswith("salon/") for x in w))
    print(f"沙龍 FAQ 問法 × 沙龍題組：{len(leaked)} 句（只列出）：{'、'.join(leaked)}")
sys.exit(1 if bad or bad_salon or cross else 0)
