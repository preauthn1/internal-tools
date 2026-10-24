#!/bin/zsh
# 自動改進迴圈：每次迭代 = AI 提出一個修改 → gate.py 評分 → 保留（commit）或還原
#
# 執行：loop/run.sh [迭代次數，預設 10]
#   在主要工作目錄執行時，會從最新的 origin/main 開一個 worktree（../<專案>-loop-<時間>，分支 loop/<時間>），
#   改到那裡跑，主要工作目錄不會被動到；在 worktree 裡執行就直接在那裡跑。
#   跑完分支上有新的 commit 就推上去並開 MR（驗收通過再合回 main，合併後記得更新主要工作目錄）。
# 停止：touch <worktree>/loop/STOP（目前這次迭代跑完後停止）
# 環境變數：LOOP_MODEL（預設 opus）、REBASELINE_EVERY（每幾次重量基準，預設 5）
set -u
cd "$(dirname "$0")/.."
N=${1:-10}
MODEL=${LOOP_MODEL:-opus}
REBASE=${REBASELINE_EVERY:-5}

# 主要工作目錄（git-dir 和 common-dir 相同）不直接跑：迴圈會 commit、還原檔案，曾在本機 main 累積 26 個沒推的 commit
if [[ "$(git rev-parse --path-format=absolute --git-dir)" == "$(git rev-parse --path-format=absolute --git-common-dir)" ]]; then
  stamp=$(date +%Y%m%d-%H%M)
  branch="loop/$stamp"
  wt="$(cd .. && pwd)/$(basename "$PWD")-loop-$stamp"
  git fetch -q origin || { echo "⚠️  git fetch 失敗"; exit 1; }
  git worktree add -q -b "$branch" "$wt" origin/main || exit 1
  echo "🌿 改到 worktree 跑迴圈：$wt（分支 $branch，從最新的 origin/main 開出）"
  echo "   停止：touch $wt/loop/STOP"
  exec "$wt/loop/run.sh" "$@"
fi

branch=$(git branch --show-current)
[[ -n "$branch" ]] || { echo "⚠️  worktree 沒有在分支上（detached HEAD），先切到分支再跑"; exit 1; }
mkdir -p loop/runs
rm -f loop/STOP
ledger_start=$(wc -l < loop/LEDGER.md)
notes_start=$(wc -l < loop/NOTES.md 2>/dev/null || echo 0)

# baseline.json 不進 git：新的 worktree 第一次跑會先量基準（eval/verify.py --holdout）
[[ -f loop/baseline.json ]] || python3 loop/gate.py --baseline || exit 1

for i in $(seq 1 $N); do
  [[ -f loop/STOP ]] && { echo "🛑 偵測到 loop/STOP，停止"; rm -f loop/STOP; break; }
  if [[ -n "$(git status --porcelain -- eval go Sources Package.swift)" ]]; then
    echo "⚠️  eval/、go/、Sources/ 或 Package.swift 有未提交的修改，先處理再跑迴圈"; git status --short -- eval go Sources Package.swift; exit 1
  fi
  ts=$(date +%Y%m%d-%H%M%S)
  echo "\n===== 第 $i/$N 次迭代（$ts）====="

  claude -p \
    --model "$MODEL" \
    --permission-mode acceptEdits \
    --allowedTools "Read,Edit,Write,Glob,Grep,Bash(go -C go build:*),Bash(BOT_AI_CACHE=1 python3 eval/run_flows.py:*),Bash(python3 eval/check_overlap.py:*),Bash(python3 loop/gate.py --dry-run:*),Bash(BOT_AI_CACHE=1 BOT_INTENT_EVAL=*),Bash(BOT_AI_CACHE=1 BOT_TESTSET=*),Bash(git diff:*),Bash(git log:*),Bash(git show:*)" \
    --disallowedTools "Read(eval/holdout/**),Read(loop/runs/*/holdout/**),Read(loop/items/holdout/**),Read(eval/results/snapshots/*/holdout/**),Read(~/Library/Caches/aichatbot/**),Read(//Users/*/Library/Caches/aichatbot/**),Edit(eval/holdout/**),Edit(loop/*.py),Edit(eval/*.py),Edit(eval/store.json),Edit(Package.swift),Edit(Sources/**),Edit(web/**)" \
    < loop/PROGRAM.md > "loop/runs/agent-$ts.log" 2>&1
  echo "🤖 AI 提案完成（紀錄：loop/runs/agent-$ts.log）"
  [[ -f loop/proposal.json ]] && cat loop/proposal.json && echo

  python3 loop/gate.py

  # 迭代心得算在迴圈本身的紀錄裡
  git add loop/NOTES.md loop/LEDGER.md loop/ledger.jsonl 2>/dev/null && git commit -q -m "loop: 紀錄第 $i 次迭代" 2>/dev/null

  # 定期重量基準：保留時的分數偏樂觀（贏家詛咒），重量一次避免門檻越墊越高
  if (( i % REBASE == 0 )); then
    echo "📏 重新量測基準"; python3 loop/gate.py --baseline
  fi
done

echo "\n===== 結束 ====="
tail -n 15 loop/LEDGER.md

# 推上分支並開 MR：保留的修改走 MR 驗收再合回 main，不直接進 main
git fetch -q origin
if (( $(git rev-list --count origin/main..HEAD) == 0 )); then
  echo "\n（分支上沒有新的 commit，不開 MR）"; exit 0
fi
kept=$(git log --format=%s origin/main..HEAD | grep -v '^loop: 紀錄第' | grep -c '^loop: ')
git push -q -u origin "$branch" || { echo "⚠️  推送失敗，請手動 git push -u origin $branch 再開 MR"; exit 1; }
export GLAB_PAGER=cat
if [[ -n "$(glab mr list --source-branch "$branch" 2>/dev/null | grep '!')" ]]; then
  echo "\n📬 分支已推上去，MR 已經存在（$branch）"; exit 0
fi
desc=$(mktemp)
{
  echo "改進迴圈在 worktree \`$(basename "$PWD")\` 跑的結果（分支 \`$branch\`），保留 $kept 次修改。"
  echo
  grep -m1 -A1 '^| 時間 |' loop/LEDGER.md   # 表頭和分隔線
  tail -n +$((ledger_start + 1)) loop/LEDGER.md
  echo
  stuck=$(tail -n +$((notes_start + 1)) loop/NOTES.md 2>/dev/null | grep "建議開 issue")   # 只列這一輪寫的
  if [[ -n "$stuck" ]]; then
    echo "## 建議開 issue（迴圈修不好、可能要改程式的題目）"
    echo
    echo "$stuck" | sed 's/^/- /'
    echo
  fi
  echo "## 合併前驗收"
  echo
  echo "1. 把最新的 main 合進這個分支（其他 MR 可能也改了提示詞或案例）"
  echo "2. \`python3 eval/verify.py --holdout\`：測試通過、和 \`eval/baseline-go.json\` 比沒有退步"
  echo
  echo "合併後記得更新主要工作目錄（\`git pull --ff-only origin main\`，有在跑的伺服器要重編重啟），並移除這個 worktree。"
} > "$desc"
glab mr create --source-branch "$branch" --target-branch main --title "loop: $branch 改進迴圈（保留 $kept 次修改）" \
  --description-file "$desc" --remove-source-branch --yes
rm -f "$desc"
echo "\n跑完後可以移除 worktree：git worktree remove $PWD"
