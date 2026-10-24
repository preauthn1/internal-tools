"""loop/gate.py 的完整流程測試：評分換成假的（不用本機模型），在暫時的 git repo 跑
量基準 → 保留 → 還原，確認規則和紀錄。執行：python3 -m unittest loop/test_gate.py"""
import importlib.util, json, os, shutil, subprocess, sys, tempfile, unittest
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent


def load_gate(root):
    spec = importlib.util.spec_from_file_location("gate_under_test", root / "loop" / "gate.py")
    gate = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(gate)
    return gate


class GateFlow(unittest.TestCase):
    def setUp(self):
        self.tmp = Path(tempfile.mkdtemp())
        self.root = self.tmp / "repo"
        # 只複製 gate 會讀的檔案，跑在獨立的 git repo（不會動到真的 repo）
        for rel in ["loop/gate.py", "eval/check_overlap.py", "eval/examples.jsonl", "eval/faq.txt",
                    "eval/intent_set.json", "eval/testset.json", "eval/faq_set.json", "eval/testset_extra.json",
                    "eval/flow_tests.json", "eval/holdout/intent_holdout.json", "eval/holdout/flow_holdout.json",
                    "eval/stores/salon/config.json", "eval/stores/salon/intent_set.json",
                    "eval/stores/salon/testset.json", "eval/stores/salon/flow_tests.json"]:
            (self.root / rel).parent.mkdir(parents=True, exist_ok=True)
            shutil.copy(REPO / rel, self.root / rel)
        for f in ("answer.go", "classify.go", "booking.go"):
            (self.root / "go/internal/dialog").mkdir(parents=True, exist_ok=True)
            shutil.copy(REPO / "go/internal/dialog" / f, self.root / "go/internal/dialog" / f)
        (self.root / ".gitignore").write_text("loop/runs/\nloop/baseline.json\nloop/proposal.json\nloop/items/\n")
        git = lambda *a: subprocess.run(["git", *a], cwd=self.root, check=True, capture_output=True)
        git("init", "-q"); git("add", "-A")
        git("-c", "user.name=t", "-c", "user.email=t@t", "commit", "-qm", "init")
        os.environ.setdefault("GIT_AUTHOR_NAME", "t"); os.environ.setdefault("GIT_AUTHOR_EMAIL", "t@t")
        os.environ.setdefault("GIT_COMMITTER_NAME", "t"); os.environ.setdefault("GIT_COMMITTER_EMAIL", "t@t")
        self.gate = load_gate(self.root)
        self.gate.checks = lambda: None
        self.scores = []
        self.gate.score = self.fake_score

    def tearDown(self):
        shutil.rmtree(self.tmp)

    def fake_score(self, out):
        # 第六個值（選填）是沙龍：(分類答對, 答錯, 逐題變差)；沒給就是沒有沙龍指標（舊基準）
        dev, hold, wrong, xlang, hold_worse, *salon = self.scores.pop(0)
        out.mkdir(parents=True, exist_ok=True)
        (out / "dev_report.md").write_text("x")
        (out / "dev.json").write_text("{}")
        metrics, items = {"testset.xlang_ok": [xlang, 15]}, {"holdout": {"worse": hold_worse, "better": 0}}
        if salon:
            ten, s_wrong, s_worse = salon[0]
            metrics.update({"salon.intent.ten": [ten, 37], "salon.testset.wrong": [s_wrong, 41]})
            items["salon"] = {"worse": s_worse, "better": 0}
            (out / "salon.json").write_text("{}")
        return {"dev": {"score": dev, "testset": {"wrong": wrong}}, "holdout": {"score": hold},
                "metrics": metrics, "items": items, "seconds": 1}, ""

    def run_gate(self, *args):
        cwd = os.getcwd(); os.chdir(self.root)
        try:
            sys.argv = ["gate.py", *args]
            self.gate.main()
        finally:
            os.chdir(cwd)

    def propose(self, text):
        with open(self.root / "eval/examples.jsonl", "a", encoding="utf-8") as f:
            f.write(json.dumps({"text": text, "intent": "general"}, ensure_ascii=False) + "\n")
        (self.root / "loop/proposal.json").write_text(json.dumps({"hypothesis": text}, ensure_ascii=False))

    def last(self):
        return json.loads((self.root / "loop/ledger.jsonl").read_text(encoding="utf-8").splitlines()[-1])

    def test_flow(self):
        self.scores = [(80.0, 90.0, 2, 10, 0)]
        self.run_gate("--baseline")
        self.assertTrue((self.root / "loop/baseline.json").exists())
        self.assertTrue((self.root / "loop/items/dev.json").exists(), "量基準時要存逐題基準")

        cases = [
            ("開發進步 → 保留", (81.0, 90.0, 2, 10, 0), "accept"),
            ("答錯少 1 次 → 保留", (81.0, 90.0, 1, 10, 0), "accept"),
            ("跨語言多 1 題 → 保留", (81.0, 90.0, 1, 11, 0), "accept"),
            ("驗證題組有 1 題變差 → 還原", (85.0, 95.0, 1, 11, 1), "reject"),
            ("跨語言少 1 題 → 還原", (85.0, 90.0, 1, 10, 0), "reject"),
            ("答錯多 1 次 → 還原", (85.0, 90.0, 2, 11, 0), "reject"),
            ("沒進步 → 還原", (81.0, 90.0, 1, 11, 0), "reject"),
        ]
        for i, (name, score, want) in enumerate(cases):
            self.scores = [score]
            self.propose(f"測試用的案例句子{i}zzq")
            self.run_gate()
            self.assertEqual(self.last()["decision"], want, name)
        examples = (self.root / "eval/examples.jsonl").read_text(encoding="utf-8")
        self.assertIn("測試用的案例句子2zzq", examples, "保留的修改要留在檔案裡")
        self.assertNotIn("測試用的案例句子3zzq", examples, "還原的修改要拿掉")

    def test_salon_guard(self):
        """第二家測試商家（#66）：沙龍退步就還原；舊基準沒有沙龍指標時照舊判斷"""
        self.scores = [(80.0, 90.0, 2, 10, 0, (30, 5, 0))]
        self.run_gate("--baseline")
        self.assertTrue((self.root / "loop/items/salon.json").exists(), "量基準時要存沙龍的逐題基準")
        cases = [
            ("沙龍有 1 題變差 → 還原", (85.0, 90.0, 2, 10, 0, (30, 5, 1)), "reject"),
            ("沙龍分類少 1 題 → 還原", (85.0, 90.0, 2, 10, 0, (29, 5, 0)), "reject"),
            ("沙龍答錯多 1 題 → 還原", (85.0, 90.0, 2, 10, 0, (30, 6, 0)), "reject"),
            ("沙龍持平、開發進步 → 保留", (85.0, 90.0, 2, 10, 0, (30, 5, 0)), "accept"),
        ]
        for i, (name, score, want) in enumerate(cases):
            self.scores = [score]
            self.propose(f"沙龍守門的案例句子{i}zzq")
            self.run_gate()
            self.assertEqual(self.last()["decision"], want, name)

    def test_salon_guard_old_baseline(self):
        """基準是加沙龍之前量的（沒有沙龍指標和逐題基準）：不報錯，照森日茶屋判斷"""
        self.scores = [(80.0, 90.0, 2, 10, 0)]
        self.run_gate("--baseline")
        self.scores = [(81.0, 90.0, 2, 10, 0, (20, 9, 0))]
        self.propose("舊基準的案例句子zzq")
        self.run_gate()
        self.assertEqual(self.last()["decision"], "accept")


if __name__ == "__main__":
    unittest.main()
