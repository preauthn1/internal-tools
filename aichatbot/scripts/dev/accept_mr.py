#!/usr/bin/env python3
"""MR 驗收（人類用，AGENTS.md 開發流程第 3～5 步）：開 worktree、合最新的 main、跑完整驗收，合併時更新基準、清理。

執行（在主要工作目錄或任何 worktree）：
  python3 scripts/dev/accept_mr.py 71                 # 驗收：worktree → 合 main → verify.py --holdout，印結果
  python3 scripts/dev/accept_mr.py 71 --merge         # 合併：基準有變就更新並 commit → push → 合併 → 主要工作目錄 pull → 清理
  python3 scripts/dev/accept_mr.py 71 --merge --note "驗收說明"   # 合併前先在 MR 留言

- 要不要合併由人判斷（自動合併規則），這支只做重複的步驟；驗收和合併分兩次執行，中間可以在 worktree 補修改再 commit。
- worktree 在主要工作目錄旁邊：../<專案>-mr-<編號>，分支 review/<編號>。
- 基準：驗收結果有進步、或題數變了，--merge 會更新基準再 commit。驗收時的結果存在 worktree 的 eval/results/accept/，
  合併時分支沒有新的 commit 就直接拿那次的結果更新（verify.py --baseline-from），不再量一次；有新 commit 才重量。
- 合併後測試伺服器有在跑，就用 scripts/dev/test-server.sh restart 換成新版。
"""
import argparse
import json
import os
import re
import subprocess
import sys
import time
from pathlib import Path

ENV = {**os.environ, "GLAB_PAGER": "cat"}


def sh(cmd, cwd=None, check=True, capture=True):
    p = subprocess.run(cmd, cwd=cwd, env=ENV, text=True, capture_output=capture)
    if check and p.returncode != 0:
        sys.exit(f"指令失敗：{' '.join(cmd)}\n{(p.stdout or '')[-2000:]}{(p.stderr or '')[-2000:]}")
    return p


def main_dir():
    common = sh(["git", "rev-parse", "--path-format=absolute", "--git-common-dir"]).stdout.strip()
    return Path(common).parent


def mr_info(n, cwd):
    data = json.loads(sh(["glab", "api", f"projects/:id/merge_requests/{n}"], cwd=cwd).stdout)
    return data["source_branch"], data["state"], data["title"]


def worktree(n, branch, root):
    wt = root.parent / f"{root.name}-mr-{n}"
    sh(["git", "fetch", "-q", "origin"], cwd=root)
    if not wt.exists():
        sh(["git", "worktree", "add", "-q", "-B", f"review/{n}", str(wt), f"origin/{branch}"], cwd=root)
    return wt


def merge_main(wt):
    if sh(["git", "merge-base", "--is-ancestor", "origin/main", "HEAD"], cwd=wt, check=False).returncode == 0:
        print("分支已包含最新的 main")
        return
    p = sh(["git", "merge", "--no-edit", "origin/main"], cwd=wt, check=False)
    if p.returncode != 0:
        sys.exit(f"合併 main 有衝突，請在 {wt} 解決後再執行：\n{p.stdout[-1500:]}")
    print("已合併最新的 main")


ACCEPT = Path("eval/results/accept")   # worktree 裡存驗收結果的地方（eval/results/ 不進 git）
# verify.py 的基準：分數（兩家店）、逐題基準（森日茶屋、霧光髮廊、驗證題組）
BASELINES = ["eval/baseline-go.json", "eval/baseline-items.json", "eval/stores/salon/baseline-items.json",
             "eval/holdout/baseline-items.json"]


def head(wt):
    return sh(["git", "rev-parse", "HEAD"], cwd=wt).stdout.strip()


def verify(wt, *extra, show=True):
    p = subprocess.run([sys.executable, "eval/verify.py", "--holdout", *extra], cwd=wt, env=ENV, text=True,
                       capture_output=True)
    out = p.stdout + p.stderr
    if show:
        keep = [l for l in out.splitlines()
                if l.startswith(("| ", "  ❌", "  ✅", "  ⚠️", "      ", "⚠️", "❌", "✅", "⏳", "   等了"))
                or "逐題比對" in l]
        print("\n".join(keep))
    return p.returncode, out


def changed(out):
    """量測結果有進步或題數不同（要更新基準）。"""
    return "✅ 進步" in out or "題數不同" in out or bool(re.search(r"從答錯變答對 [1-9]", out))


def merge(n, root, retries=4):
    for i in range(retries):
        p = sh(["glab", "mr", "merge", str(n), "--yes"], cwd=root, check=False)
        if "Merged" in p.stdout + p.stderr:
            return True
        time.sleep(4)
    return False


def main():
    ap = argparse.ArgumentParser(description="MR 驗收")
    ap.add_argument("mr", type=int)
    ap.add_argument("--merge", action="store_true", help="更新基準、push、合併、清理")
    ap.add_argument("--note", help="合併前在 MR 留的驗收說明")
    a = ap.parse_args()
    root = main_dir()
    branch, state, title = mr_info(a.mr, root)
    print(f"!{a.mr} {title}（{branch}，{state}）")
    if state != "opened":
        sys.exit("MR 不是開啟狀態")
    wt = worktree(a.mr, branch, root)
    print(f"worktree：{wt}")
    merge_main(wt)

    if not a.merge:
        print(sh(["git", "diff", "--stat", "origin/main"], cwd=wt).stdout.strip().splitlines()[-1])
        acc = wt / ACCEPT
        sh(["rm", "-rf", str(acc)])
        code, out = verify(wt, "--out", str(acc))
        acc.mkdir(parents=True, exist_ok=True)
        (acc / "run.json").write_text(json.dumps({"head": head(wt), "code": code, "output": out}, ensure_ascii=False))
        print(f"\n結束碼 {code}（0 沒有退步、1 有退步、2 量測無效）；要合併：python3 scripts/dev/accept_mr.py {a.mr} --merge")
        return

    acc = wt / ACCEPT
    run = json.loads((acc / "run.json").read_text()) if (acc / "run.json").exists() else None
    if run and run["head"] == head(wt) and run["code"] != 2:
        print("沿用驗收時的量測結果（分支沒有新的 commit）")
        code, out = run["code"], run["output"]
        update = lambda: sh([sys.executable, "eval/verify.py", "--baseline-from", str(acc)], cwd=wt)
    else:
        print("分支有新的 commit（或沒有驗收紀錄），重新量測")
        code, out = verify(wt, "--no-checks", show=False)
        update = lambda: verify(wt, "--no-checks", "--update-baseline", show=False)
    if code == 2:
        sys.exit("量測無效（AI 呼叫失敗或連不上模型服務），先不要合併")
    if changed(out):
        update()
        if sh(["git", "status", "--porcelain", *BASELINES], cwd=wt).stdout.strip():
            sh(["git", "add", *BASELINES], cwd=wt)
            sh(["git", "commit", "-q", "-m", f"更新 Go 版基準（!{a.mr} 驗收）\n\nCo-Authored-By: Claude Opus 5.5 <user@example.invalid>"], cwd=wt)
            print("已更新基準並 commit")
    sh(["git", "push", "-q", "origin", f"HEAD:{branch}"], cwd=wt)
    if a.note:
        sh(["glab", "mr", "note", "create", str(a.mr), "-m", a.note], cwd=root)
    time.sleep(2)
    if not merge(a.mr, root):
        sys.exit(f"合併失敗，worktree 留著：{wt}")
    print(f"✅ 已合併 !{a.mr}")
    sh(["git", "pull", "-q", "--ff-only", "origin", "main"], cwd=root)
    sh(["git", "worktree", "remove", "--force", str(wt)], cwd=root)
    sh(["git", "branch", "-D", "-q", f"review/{a.mr}"], cwd=root, check=False)
    print(f"主要工作目錄已更新：{sh(['git', 'log', '--oneline', '-1'], cwd=root).stdout.strip()}")
    if sh(["git", "diff", "--name-only", "HEAD@{1}", "HEAD", "--", "Sources/"], cwd=root, check=False).stdout.strip():
        print("⚠️ 改到 Sources/：在主要工作目錄跑 scripts/dev/install-model-server.sh")
    if subprocess.run(["pgrep", "-f", ".build/gobot serve"], capture_output=True).returncode == 0:
        print("測試伺服器有在跑，換成新版：")
        subprocess.run(["zsh", "scripts/dev/test-server.sh", "restart"], cwd=root)


if __name__ == "__main__":
    main()
