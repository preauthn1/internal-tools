# 自動改進迴圈（loop engineering）

讓客服機器人的提示詞和內建意圖案例在固定題庫的保護下自動、持續變好（2026-10-09 起是 Go 版，只改提示詞和案例，不改程式邏輯）：

```
┌─> AI 讀紀錄與失敗明細 → 提出一個假設 → 做最小修改 → 寫 proposal.json   （claude -p，PROGRAM.md）
│                                  │
│   gate.py：靜態檢查 → 測試 → eval/verify.py 量開發＋驗證題組 → 依規則決定（程式決策，AI 不能干預）
│        ├─ 保留：git commit、更新基準
│        └─ 還原：git checkout、新檔移到 quarantine/
│                                  │
└──────── LEDGER.md／NOTES.md 記下結果與心得，下一輪會讀 ────────────────┘
```

## 檔案

| 檔案 | 用途 |
|---|---|
| `run.sh` | 外層迴圈：`loop/run.sh 10` 跑 10 次迭代。在主要工作目錄執行會自動開 worktree（`../aichatbot-loop-<時間>`，分支 `loop/<時間>`）去那裡跑；跑完推分支、開 MR。`touch <worktree>/loop/STOP` 讓它在這次跑完後停下 |
| `PROGRAM.md` | 每次迭代給 AI 的指令（選題方式、修改規則、快速檢查） |
| `gate.py` | 守門：檢查、跑測試、用 `eval/verify.py --out` 評分、保留或還原、寫紀錄 |
| `score.py` | **已停用**：Swift 版的評分，Go 版改用 `eval/verify.py` |
| `LEDGER.md`、`ledger.jsonl` | 每次嘗試的結果（人看的表格／程式讀的明細） |
| `NOTES.md` | 每次迭代留下的心得 |
| `baseline.json` | 目前基準分數（不進 git） |
| `runs/<時間>/` | 每次評分的輸出（`metrics.json`、`dev_report.md`、`holdout/`）；`runs/latest` 指向最新一次 |

## 題組

| | 開發（AI 看得到失敗明細） | 驗證（AI 只看到總分） |
|---|---|---|
| 完整流程 | `flow_tests.json` 104 步 | `holdout/flow_holdout.json` 26 步 |
| 分類 | `intent_set.json` 101 句 | `holdout/intent_holdout.json` 41 句 |
| 單題 | `testset.json`＋`faq_set.json` 118 題 | — |
| 誘答題、跨語言題 | `testset_extra.json` 16＋15 題 | — |

機器人用 greedy 取樣（2026-10-07 起），同一句話每次結果相同，所以每題只跑 1 次。

開發綜合分數＝分類、單題、流程題三個百分比的平均；驗證綜合分數＝驗證分類、驗證流程題的平均。

## 保留規則（`gate.py` 開頭的常數）

1. 只改內建意圖案例 `examples.jsonl`、`faq.txt`（只能改問法）、提示詞所在的 `go/internal/dialog/answer.go`、`classify.go`、`booking.go`（只能改字串內容：改動的行拿掉字串後要和原本一樣，整行新增的要是純字串）
2. 提示詞不能寫死測試句；案例庫不能和任何題庫重複；測試（`eval/verify.py --checks-only`）全部通過
3. 開發分數 ≥ 基準 +0.2（約等於淨修好 1 題，交換不算），
   **或**「答錯」加「誘答題 編造」少 1 題以上且開發分數沒有下降（寧可轉店員，也不要答錯）
   **或**跨語言多答對 1 題以上、開發分數沒降、答錯沒增加（跨語言不在開發分數裡）
4. 「答錯」加「誘答題 編造」不能增加，跨語言答對不能變少
5. 驗證題組**沒有任何一題從答對變答錯**（逐題比對，逐題基準在 `loop/items/`），驗證分數也最多掉 3.0（抓過度擬合）

greedy 取樣結果固定，任何一題的差異都是真的，所以門檻都是 1 題。守門規則的完整流程有測試：`python3 -m unittest loop/test_gate.py`（驗收指令會跑）。

每 5 次迭代重量一次基準：保留時的分數通常是「剛好跑得好」的那次，偏樂觀。

被還原的修改會存成 `loop/runs/<時間>/rejected.patch`（新增的檔案在同一資料夾的 `quarantine/`），可以用 `git apply` 重新套用。

## 執行流程（worktree → MR → 合回 main）

1. `loop/run.sh 10`：在主要工作目錄執行時，會從最新的 `origin/main` 開一個 worktree（`../aichatbot-loop-<時間>`、分支 `loop/<時間>`），改到那裡跑；主要工作目錄不會被動到，可以照常開發。在 worktree 裡執行就直接在那裡跑。
2. 新的 worktree 沒有 `baseline.json`（不進 git），會先量一次基準。
3. 跑完分支上有新的 commit，就推上去並開 MR（描述附這次的紀錄表和驗收步驟）。
4. MR 驗收：把最新的 main 合進分支 → `python3 eval/verify.py --holdout`（測試＋真模型量測，和 `eval/baseline-go.json` 比，沒有退步）。
5. 合併後**更新主要工作目錄**（`git pull --ff-only origin main`，有在跑的伺服器要重編重啟），並移除 worktree（`git worktree remove <路徑>`）。

## 常用指令

```bash
loop/run.sh 10
```

```bash
python3 loop/gate.py --baseline
```

```bash
python3 loop/gate.py --dry-run
```

```bash
touch ../aichatbot-loop-<時間>/loop/STOP
```

## 注意

- 本機模型的呼叫在常駐的 model-server 排隊：迴圈執行時另外跑量測不會失真，但兩邊都會變慢，盡量避開。
- Swift 版的舊 `baseline.json` 不能沿用（gate.py 會拒絕）：**恢復迴圈前先跑 `python3 loop/gate.py --baseline`**。
- 題庫的預約日期是相對標記 `{{date:+N:格式}}`，執行時依當天換算，不會過期（規則見 `eval/testdates.py`）。
- 驗證題組被看過或被拿來調整，就失去意義；累積很多次迭代後，應該再寫一批新的驗證題。
