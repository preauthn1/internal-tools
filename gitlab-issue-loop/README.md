# gitlab-issue-loop

把 GitLab issue 指派給 bot 帳號，AI 就會：讀 issue → 有疑問先在 issue 上提問、等到全部釐清 → 開發 → 跑測試驗證 → 開 MR → 回覆 issue 請人類審閱。

## 流程

```
指派給 bot ──► 釐清（唯讀讀碼）
                 ├─ 有疑問 ─► 留言提問，label needs-info ─► 人類回覆 ─► 再釐清（可多輪）
                 ├─ 不適合 ─► 留言說明，label ready-for-human（人類改回 ready-for-agent 才會再接）
                 └─ 都清楚 ─► 留言「開始處理」，label ready-for-agent + ai::doing
                                 ▼
                        共用 worktree 切到分支 ai/issue-N 開發（tdd skill → 驗證 → code-review skill，審閱在前景等結果）
                                 ▼
                 腳本檢查流程痕跡（summary 有「## 自我審閱」、test_notes 有「## TDD」）── 缺 ─► 退回補做一次
                                 ▼                                                      └─ 仍缺 ─► 留言，label ai::blocked
                 腳本獨立跑 test_cmd ── 失敗 ─► 把輸出丟回 AI 修（最多 max_fix_rounds 輪）
                                 ▼ 通過                       └─ 仍失敗 ─► 留言，label ai::blocked
                 commit、push、開/更新 MR（Closes #N，reviewer = issue 作者）
                                 ▼
                 留言 @作者 請審閱 MR，label ai::review
```

- label 分兩組：**triage label**（`needs-triage`、`needs-info`、`ready-for-agent`、`ready-for-human`、`wontfix`，與專案的 triage 流程共用，同時只留一個）表示「這張 issue 該誰做」；**執行狀態**（`ai::doing`、`ai::review`、`ai::blocked`）表示迴圈做到哪。
- 觸發條件只有「指派給 bot」；`ready-for-agent` 本身不會讓迴圈接單。指派給 bot 但標了 `ready-for-human`、`wontfix` 的不處理。
- 人類可以手動改 label：拿掉執行狀態、改回 `ready-for-agent` 就會從頭再來。
- **順序**：沒有等待的 issue 從編號小的開始做。issue 描述最上面可以寫兩種等待（只認同專案的 `#N`，「## Blocked by」小節不讀）：
  - `Blocked by: #26, #34`：**需求依賴**，沒有那張的功能這張做不了。
  - `排在後面：#27（都改 go/internal/dialog，避免衝突）`：**排程順序**，功能上不需要，只是怕同時改同一塊程式；括號寫原因，可以隨時拿掉。
  - 兩種迴圈都一樣處理：列的 issue 還沒關閉（MR 還沒合併）就先跳過，log 印「⏸ #N 等 #27（需求）關閉」。
- `needs-info`、`ai::blocked` 狀態下，只要最後一則留言不是 bot 的，就會重新釐清。
- 開發時 AI 用 `mattpocock-skills:tdd` 先寫測試再實作，完成後用 `mattpocock-skills:code-review` 自我審閱，結果寫在 MR 說明的「自我審閱」一節；需要安裝 mattpocock-skills plugin。
- AI 只負責讀碼與改檔；留言、push、開 MR 都由 `loop.py` 執行。AI 拿不到 GitLab token，且預設禁止 `git push`、`curl`、`WebFetch` 等對外動作（issue 內容是不可信輸入）。

## 設定

1. 建一個 bot 帳號，加入專案（至少 Developer），用它建立 Personal Access Token（scope：`api`，用來留言、改 label、開 MR）。
   git 走 SSH：替 bot 產生金鑰（`ssh-keygen -t ed25519 -f ~/.ssh/my_ai_bot -N ''`），把公鑰加到 bot 帳號的 SSH Keys，`config.json` 的 `ssh_key` 指向私鑰。clone 位址預設用 API 回傳的 `ssh_url_to_repo`，可用專案的 `ssh_url` 覆寫（SSH 不是 22 埠時用 `ssh://git@host:port/group/proj.git`）。
2. 複製設定檔並修改：
   ```bash
   cp config.example.json config.json
   ```
   每個專案設定 `test_cmd`（驗證關卡）、`setup_cmd`（裝依賴）、`agent_bash`（允許 AI 執行的指令）、`keep_dirs`（換工單時保留的建置快取，預設 `.build`、`node_modules`）。
3. 設定 token（存成 `token_file` 指定的檔案，或用環境變數）：
   ```bash
   export GITLAB_BOT_TOKEN=[REDACTED]
   ```

## 常駐（launchd）

`ops/install.sh` 把迴圈和看門狗裝成 LaunchAgent（你登入後才跑，Mac 可以鎖螢幕）：

- **迴圈**（`com.rory.gitlab-issue-loop`）：開機登入後自動啟動、結束或當掉自動重啟。紀錄在 `~/Library/Logs/gitlab-issue-loop/loop.log`。
- **改了 `loop.py`、提示詞或 `config.json` 之後**：`touch STOP`。迴圈做完手上那張就結束、自己刪掉 STOP，launchd 會用新程式重開，不會打斷正在做的單。
- **看門狗**（`com.rory.gitlab-issue-loop.watchdog`）：每 5 分鐘檢查迴圈程序和 `watchdog.urls`（預設 model-server），停掉超過 10 分鐘用 Telegram 通知一次，恢復時再通知。
- 要真的停掉：`ops/install.sh uninstall`。

## 通知（Telegram）

`config.json` 的 `notify`：`telegram_token_file`（bot token 存在檔案裡，權限 600；AI 不能讀這個資料夾）和 `chat_id`。會通知：釐清時有問題要回答（needs-info）、不適合自動處理（ready-for-human）、開發中斷（ai::blocked）、開了 MR 但**沒達到數字門檻**、看門狗發現服務停掉。沒設定時不通知。

## MR 上的執行紀錄

每張 issue 每次呼叫 claude 的輪數、時間、估算費用（API 牌價，用訂閱登入不另外收費）、被權限擋下的指令，附在 MR 說明和中斷時的 issue 留言。擋下超過 5 次或某一步超過 200 輪時，最上面加警告。

## 數字門檻

- 釐清時：會改到回答或分類行為的 issue 沒寫數字門檻，就回 needs-info 問門檻（建議「不低於基準」）。
- 開發時：沒達標最多試 3 種做法；還是沒達標就照樣開 MR，標題前面加「未達標：」、說明最上面寫門檻和實際數字，並通知人類決定。

## 執行

```bash
./loop.py plan                     # 只列出會處理哪些 issue
./loop.py issue my-group/proj#12   # 只處理一個 issue（第一次建議用這個試）
./loop.py once                     # 處理所有待辦一輪
./loop.py watch --interval 300     # 每 5 分鐘輪詢；touch STOP 停止
```

工作目錄在 `work/`（clone 與所有 issue 共用的 `worktree/`：換工單時切分支、清掉未追蹤檔案，但保留 `keep_dirs` 讓編譯只做增量），每次 AI 執行的完整輸出在 `runs/`。
