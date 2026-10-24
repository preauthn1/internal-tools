# aichatbot

開發一個 LINE AI 客服（參考「顧店」LINE AI 客服），用 macOS 內建的 Apple 本機模型回覆。接手前先讀 `docs/研究總結.md`、`docs/開發計畫.md`、`docs/TODO.md`。

## 開發流程

每個開發迴圈（自動改進迴圈 `loop/run.sh`、issue→MR 迴圈、人或 AI 手動處理 issue）都照這個流程：

1. **在 git worktree 開發**：從最新的 `origin/main` 開分支和 worktree，不在主要工作目錄（`~/Workspace/rory/aichatbot`）的 `main` 上直接改或 commit。`loop/run.sh` 在主要工作目錄執行時會自動開 worktree。
2. **開發完沒問題就推分支、開 MR**。
3. **MR 驗收也在 worktree 做**（人類驗收可以用 `python3 scripts/dev/accept_mr.py <MR 編號>`：開 worktree、合 main、跑 `verify.py --holdout`；確定要合併再加 `--merge`，會更新基準（分支沒有新 commit 就沿用驗收時的結果，不再量一次）、push、合併、pull 主要工作目錄、清 worktree，測試伺服器有在跑就換成新版）：先把最新的 main 合進分支，再跑 `python3 eval/verify.py --checks-only`（Swift 的 `swift build`、`bot-tests`、`telegram-tests`，Go 的 gofmt、vet、test，web 編譯和 vue-tsc 型別檢查，案例庫重複檢查，改進迴圈守門規則的測試）。改到 Go 對話層、提示詞或題庫時，用 `python3 eval/verify.py` 和基準比分數（見下面）；行為不變的重構用 `python3 eval/snapshot.py save`、`compare` 逐題比對 Go 版（`go/bin/chat`）前後的結果。這兩支預設開 AI 回答快取（放在 `~/Library/Caches/aichatbot/ai/`，不能讀，見 README），`BOT_AI_CACHE=0`（`verify.py` 是 `--no-cache`）關掉。`loop/score.py` 是 Swift 版的評分，已停用。
   - **自己跑 `go/bin/chat` 量測（`BOT_INTENT_EVAL`、`BOT_TESTSET`、對話模式、`run_flows.py --bot go-chat`）一律加 `BOT_AI_CACHE=1`**：只有 `verify.py`、`snapshot.py` 預設開快取，單獨跑 `go/bin/chat` 不設就每題重新呼叫本機模型（分類題組一次 101 次呼叫）。
   - **`verify.py` 的真模型量測一次只跑一組**（檔案鎖 `~/Library/Caches/aichatbot/measure.lock`）：issue 迴圈、改進迴圈、人的驗收同時要量時會自動排隊，輸出「⏳ 等其他量測結束」。直接跑 `go/bin/chat`、`run_flows.py` 不經過這把鎖。
   - **用本機模型量測時不要同時跑吃 CPU 的工作**（Docker、編譯、embedding、另一組量測）：本機模型會逾時，結果出現「AI 失敗，用預設值」。只要出現 AI 失敗，那次量測就無效，要在機器空閒時重跑，不能拿來下結論。
   - **Go 版用真模型量測**：`model-server`（Apple 本機模型服務）由 launchd 常駐在 `0.0.0.0:8765`（`com.rory.aichatbot.model-server`），**不要自己啟動或停掉它**，`go/bin/chat`、`bot serve`、`eval/verify.py` 預設就連它；量測輸出出現連不上模型服務、或全部都是 AI 失敗時，回報 blocked（迴圈的 AI 不能用 curl，不用另外檢查）。它和 Swift 的量測共用同一把檔案鎖排隊。
   - **Go 版的驗收一律用 `python3 eval/verify.py`**：測試＋真模型量測（分類、單題、誘答題編造數、跨語言答對數、流程題），和 `eval/baseline-go.json` 的基準比較，另外**逐題比對**（逐題基準在 `eval/baseline-items.json`、`eval/stores/salon/baseline-items.json`、`eval/holdout/baseline-items.json`）：總數退步、或有任何一題從答對變答錯（總數持平也算）都標 ❌（結束碼 1；量測無效是 2）。開發題組列出是哪幾題，驗證題組只給題數。另外印每則平均時間。迴圈的 AI 已經跑過驗收指令時加 `--no-checks`；**只有人類驗收可以加 `--holdout`**（驗證題組）。合併後分數有提升時，人類在主要工作目錄用 `--update-baseline` 更新基準。`--model haiku`（要 `ANTHROPIC_API_KEY`，沒有就結束碼 2）改用 Haiku 量，只列分數、token 數和估算費用，不和基準比；量真實延遲或費用時加 `--no-cache`。
   - **兩家測試商家分開量**：森日茶屋（`eval/`）和美髮沙龍「霧光髮廊」（`eval/stores/salon/`：後台匯出檔格式的商家設定和開發題組），每家各一組指標（名稱前面標店名）和逐題比對，**改提示詞要兩家都不退步**。`--store tea|salon` 只量一家（預設兩家）。沙龍沒有驗證題組（之後由人另外寫）；沙龍的題組句子不能和森日茶屋的題組、案例庫、沙龍自己的案例重複（`check_overlap.py` 會擋）。
   - **數字門檻**：會改到回答或分類行為的 issue（`go/internal/dialog`、`go/internal/model`、提示詞、意圖案例、題組），驗收條件要有數字門檻；沒特別寫的話是「`eval/verify.py` 沒有退步」。
4. **驗收通過才合回 main**。
5. **合併後更新主要工作目錄**（實機測試環境用 `zsh scripts/dev/test-server.sh start|restart|stop|status`：Go 版伺服器＋Cloudflare 臨時 tunnel，資料庫 `data/go-test.sqlite`）：`git pull --ff-only origin main`；有在跑的伺服器要重新編譯、重啟（改到 `Sources/` 時，在主要工作目錄跑 `scripts/dev/install-model-server.sh`：編譯、等沒有人在用本機模型時才重啟常駐的 model-server）；移除驗收用的 worktree。

## Agent skills

### Issue tracker

Issues 放在自架 GitLab（Rory/aichatbot），用 glab 操作。See `docs/agents/issue-tracker.md`.

### Triage labels

使用預設的五個 triage label（needs-triage、needs-info、ready-for-agent、ready-for-human、wontfix）。See `docs/agents/triage-labels.md`.

### Domain docs

Single-context：根目錄一份 CONTEXT.md，ADR 放在 docs/adr/。See `docs/agents/domain.md`.
