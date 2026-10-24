# aichatbot

> 本目录原始 README 见文末，本节为快照补充信息。

原始路径：`aichatbot`　|　文件 **1220**　|　主语言：TypeScript 496, Vue 355, Go 121

## 目录构成

```
web/packages/  — 811
go/internal/  — 108
web/internal/  — 94
web/apps/  — 72
eval/  — 30
web/scripts/  — 20
Sources/BotCore/  — 14
go/cmd/  — 13
web/  — 10
loop/  — 8
Sources/Telegram/  — 5
eval/stores/  — 5
```

## 关键文件

- CI/构建：`.devcontainer/Dockerfile`、`web/scripts/deploy/Dockerfile`
- 前端入口：`web/apps/web-antd/index.html`、`web/apps/web-antd/src/main.ts`

---

## 原始说明

# 森日茶屋 AI 客服

這個專案要開發一個 LINE AI 客服，參考了「顧店」LINE AI 客服：用 macOS 內建的 Apple 本機模型（Foundation Models）擔任 LINE 官方帳號的客服，不需要 API、也不用付費。範例店家是「森日茶屋」。

訊息流程：Telegram（之後換成 LINE 官方帳號）→ Cloudflare Tunnel → 自己的 Mac 上的 webhook 伺服器 → Apple 本機 AI 回覆。

目前進度：主程式是一個 Swift Package（核心在 `Sources/BotCore/`）。有兩種用法：
- 執行檔 `bot`：在終端機扮演顧客對話，也是量測各題組的入口。
- 執行檔 `server`：Telegram webhook 伺服器，已經可以私訊對話，顧客狀態和對話紀錄存在 SQLite。

**還沒接上 LINE**，店員端（Telegram 群組）和管理介面也還沒做。研究過程與測試數據請看 [`docs/研究總結.md`](docs/研究總結.md)，開發順序看 [`docs/開發計畫.md`](docs/開發計畫.md)，待辦看 [`docs/TODO.md`](docs/TODO.md)。

## 設計原則

這些原則都是從 `eval/` 的測試結果歸納出來的，詳細說明在 `Sources/BotCore/Config.swift` 開頭的註解：

- **AI 只做選擇題**：意圖分類、挑 FAQ、讀訂單數量、抽預約欄位、抽名字。AI 不自由撰寫回覆。
- **規則和計算由程式處理**：價格、折扣、預約檢查都依照 `store.json` 計算。
- **對話狀態由程式記錄**：例如預約進行到哪一步、是否在等店員、顧客追問了幾次，不靠比對上一句回覆來判斷。
- **一次只跑一個 AI 呼叫**：同時跑多個呼叫品質會下降，所以所有 AI 呼叫（跨程序）都排隊執行。一般訊息只會呼叫 AI 1～2 次。
- **AI 用 greedy 取樣**：同一句話每次結果相同，量測每題只要跑 1 次，正式回覆也比較一致。
- **回覆只有三種來源**：`replies.json`、`faq.txt`，或程式計算出的結果。

## 環境需求

- Apple Silicon 的 Mac，macOS 27，並已啟用 Apple Intelligence（`bot` 啟動時會檢查本機模型能不能用，不能用就直接結束）。
- Swift 工具鏈：只裝 Command Line Tools 就能編譯（`swift build`，Swift 6.2 以上），但這樣無法使用 `@Generable` 巨集，所以程式改用 `DynamicGenerationSchema`。正式開發需要安裝 Xcode。
- Python 3：用來跑量測腳本，只需要標準函式庫。

### 開發容器（Go 版）

Go 版（見 [`docs/開發計畫.md`](docs/開發計畫.md)）的開發環境放在 `.devcontainer/`：Go 1.25、node 22＋pnpm（corepack）、python3、sqlite3。用 VS Code 的「Reopen in Container」，或用 CLI：

```bash
npx @devcontainers/cli up --workspace-folder .
npx @devcontainers/cli exec --workspace-folder . go version
```

- Apple 本機模型不能在容器裡跑：Swift 模型服務在 Mac 上執行，容器用 `http://host.docker.internal:<埠>` 呼叫。Swift 的編譯和測試也在 Mac 上做。
- 版本鎖在專案檔裡（`go/go.mod` 的 `toolchain`、`web/package.json` 的 `packageManager`、`.nvmrc`），直接在 Mac 上開發（包括 issue 迴圈）也會用同一套版本。
- 正式環境不用容器，用原生執行檔＋launchd。

## Go 服務（開發中）

`go/` 是 Go 版的服務（[ADR 0001](docs/adr/0001-go-service-with-vben-admin.md)）。目前有：SQLite、商家設定的匯入和匯出、Telegram webhook、對話層（分類意圖＋預約狀態機＋一般對話的挑 FAQ／AI 寫回覆，見下面「對話層」）、員工群組（轉人工、員工回覆、預約確認／拒絕，見下面「員工群組」）。

- Go 版本鎖在 `go/go.mod`（`go 1.25.0`＋`toolchain go1.25.14`）：本機的 Go 比較舊時會自動下載這一版（`GOTOOLCHAIN` 是 `auto` 時；設成 `local` 會不理 `toolchain`）；本機的 Go 比它新時用本機的，要完全一樣就設 `GOTOOLCHAIN=go1.25.14`。SQLite 用純 Go 的 `example.invalid/sqlite`，不需要 cgo，Mac 和開發容器編出來一樣。
- 測試：`go -C go test ./...`（webhook 層用假的 Telegram API，不需要 bot token）。

```bash
go -C go build -o ../.build/gobot ./cmd/bot     # 編譯成 .build/gobot
.build/gobot import eval                        # 把 eval/ 的森日茶屋匯入成商家設定（整份取代，商家在後台新增的意圖案例保留）
.build/gobot export merchant.json               # 匯出成一個 JSON 檔（- 是標準輸出）
.build/gobot import merchant.json               # 匯入匯出的 JSON
.build/gobot create-admin owner                 # 建立後台管理員（見下面「後台」）
.build/gobot bind-code                          # 印出員工群組的綁定碼和目前綁定的群組（見下面「員工群組」）
set -a; source .env; set +a; .build/gobot serve # 啟動 webhook 伺服器和後台
```

- 每個子指令都可以加 `-db <路徑>`，預設是 `BOT_DB`，沒設就是工作目錄的 `data/bot.sqlite`（和 Swift 版同一個預設位置，但資料表不同，兩版不要共用同一個檔案）。
- 環境變數和 Swift 版相同：`TELEGRAM_BOT_TOKEN`、`TELEGRAM_WEBHOOK_SECRET`、`TELEGRAM_PUBLIC_URL` 必填，`TELEGRAM_API_BASE`、`PORT` 選填（見下面「Telegram 伺服器」）。`serve` 啟動後會呼叫 `setWebhook`（`<TELEGRAM_PUBLIC_URL>/telegram/webhook`，失敗只記錄）。
- 伺服器的行為：`X-Telegram-Bot-Api-Secret-Token` 不符回 401；同一個 `update_id` 只處理一次（記在資料庫，重開後也記得）；先回 200 再在背景處理，同一位顧客的訊息依收到的順序一則一則處理（不同顧客同時處理）；私訊先送「正在輸入」再回覆（轉人工中不顯示、不回覆）；bot 指令（`/start`、`/help` 這類 `/` 接英文字母開頭的英數字）不送進分類，回打招呼（回覆種類 `greeting`，新顧客用商家慣用語言，不拿指令判斷語言，`go/bin/chat` 也一樣，#54）；進出各寫一筆對話紀錄（回覆送不出去時照樣記下回覆和 AI 判斷，標成送出失敗 `send_failed`，後台對話內容頁顯示「送出失敗」）；Telegram API 失敗等 1、2、4 秒重試（「正在輸入」不重試）；日誌不含 token 和 secret。收到 SIGINT／SIGTERM 時處理完已收下的訊息再結束。
- 啟動時和之後每 24 小時清一次舊資料：對話紀錄保留 90 天、update 紀錄保留 7 天、過期的後台 session 刪掉。
- 匯入 `eval/` 時的轉換：`store.json` 的規則寫成知識段落（營業時間、菜單與價格、付款、自取優惠、人數上限），商家名稱取自 `store.json`，介紹由菜單組成，聯絡方式取 `shop.txt` 的地址；FAQ 只存中文答案，第一個問法當問題、其餘當其他問法；`replies.json` 存成「key × 語言」（`*` 是所有語言共用），用到營業時間、人數上限這類結構化變數（`{open}`、`{close}`、`{max}`、`{min}`）的模板不匯入（ADR 0003 拿掉了這些欄位，對話層用內建的說法）；`examples.jsonl` 當內建案例，Swift 版的舊代號（greeting、thanks、goodbye、small_talk、store_question）改成 `general`，`state` 照存。預設預約、轉人工都開，慣用語言中文，顧客語言中英日。

### 模型層（#25）

`go/internal/model`：業務程式只透過 `ChatModel` 呼叫 AI（給指示和欄位定義：選擇、文字、整數，回傳各欄位的值；失敗或逾時回 nil）。組法是 `Switch`（模型開關、費用記錄、月上限）→ `Cache`（AI 回答快取，有開才有）→ `Queued`（本機模型一次一個）→ `Apple`（呼叫 Swift 模型服務）；Haiku 不排隊。對話層測試用照腳本回答的假模型 `Scripted`。

- **預設是 Apple 本機模型（AFM）**。Haiku 要設 `BOT_MODEL=haiku` 才用，而且要有 `ANTHROPIC_API_KEY`，沒有 key 時啟動失敗並說明（不會偷偷改用 AFM）。
- **Haiku**：直接用 `net/http` 呼叫 Messages API，`max_tokens: 300`、關閉思考（Haiku 5.5 預設會先思考，額度可能被用完）、不送 `temperature`（Haiku 5.5 只接受預設值，送 0 會回 400），用 structured outputs（`output_config.format` 的 `json_schema`）強制照欄位格式回答；逾時 15 秒；可以同時呼叫。
- **用量和費用**：每次實際呼叫在資料庫 `ai_usage` 表記一筆（模型、題目名稱、token 數、費用、成功與否；不存顧客的話）。價格表在 `go/internal/model/budget.go`（`claude-haiku-5-5`：每百萬 token 輸入 $0.10、輸出 $0.50）；不在表上的 model id 費用記 0 並警告。本機模型記 0 元，token 數照模型服務回的記；快取命中不記。
- **每月費用上限**：本月（台北時間）Haiku 累計費用到上限後，當月剩下的呼叫自動改用本機模型，記一次警告 log 並通知員工群組（每個月一次；還沒綁定群組時只記 log）；下個月自動恢復。
- **後台可以改模型開關和上限（#38）**：`BOT_MODEL`、`BOT_AI_MONTHLY_BUDGET_USD` 只是初始值；後台存過一次之後（資料庫 `model_setting` 表）以資料庫為準，啟動時讀出來套用，存檔後馬上生效、不用重開。資料庫存的是 haiku 但現在沒有 `ANTHROPIC_API_KEY` 時，記警告並先用 AFM（不會啟動失敗）。

| 環境變數 | 預設 | 說明 |
| --- | --- | --- |
| `BOT_MODEL` | `afm` | 模型開關：`afm` 或 `haiku` |
| `ANTHROPIC_API_KEY` | （無） | Haiku 的 API key；不寫進日誌 |
| `BOT_HAIKU_MODEL` | `claude-haiku-5-5` | Haiku 的 model id |
| `BOT_AI_MONTHLY_BUDGET_USD` | `10` | Haiku 每月費用上限（美元），`0` 表示不限 |
| `BOT_APPLE_MODEL_URL` | `http://0.0.0.0:8765` | Swift 模型服務；開發容器裡用 `http://host.docker.internal:8765` |
| `BOT_AI_CACHE`、`BOT_AI_CACHE_DIR`、`BOT_AI_CACHE_NS` | 關、`<使用者快取目錄>/aichatbot/ai`、`dev` | AI 回答快取，見下面 |
| `MODEL_SERVER_HOST`、`MODEL_SERVER_PORT` | `0.0.0.0`、`8765` | 模型服務（Swift）聽的位址和埠 |
| `MODEL_SERVER_QUEUE_MAX` | `4` | 模型服務最多幾個請求在排隊（不含正在跑的那個）；滿了新的請求回 503 |
| `MODEL_SERVER_ALLOW` | （無） | 模型服務另外允許的來源網段（IPv4 CIDR，逗號分隔）；本機（loopback）一律允許，其他來源回 403 |

**模型服務**（`Sources/model-server`）：只包 Apple 本機模型的 HTTP 服務。`POST /respond` 收 `{name, instructions, text, schema}`（schema 是 `model.Schema` 產生的 JSON schema），成功回 `{values, usage}`（`usage` 是這次的輸入、輸出 token 數 `{input_tokens, output_tokens}`，算不出來時沒有這個欄位），本機模型失敗或逾時回 502（提示詞超過本機模型上下文 8192 tokens 時多帶 `"reason":"context_exceeded"`，Go 端的 `Usage.Failure` 是 `too_long`，挑 FAQ、寫回覆的 trace 記 `too_long` 而不是 `ai_failed`，#65；被 Apple 的安全過濾擋下時帶 `"reason":"guardrail"`，`Usage.Failure` 是 `blocked_by_safety`，各步的 trace 記 `blocked_by_safety`，分類被擋時不回看不懂、當成一般對話往下走，`eval/verify.py` 另外列「被安全過濾擋下 N 次」、不算 AI 失敗，#68）。只改寫既有內容的翻譯（`Translate`、`TranslateCustomer`）用寬鬆的安全過濾（`permissiveContentTransformations`，只對純文字回答有效，所以這兩個請求不帶 schema、回答直接放進唯一的欄位）；`GET /version` 回模型版本和 macOS 版本。呼叫經過 Swift 版的 `AIQueue`，和 Swift 的量測程式共用同一把檔案鎖排隊。

- **排隊上限**：已經有 `MODEL_SERVER_QUEUE_MAX` 個在等時，新的請求直接回 **503**，不排隊。
- **請求方放棄就不執行**：Go 端送 `X-Start-By-Ms`（最晚開始時間＝送出時間＋Go 端逾時 20 秒－本機模型逾時 15 秒）。排到這個時間還沒輪到，模型服務就不執行、回 **504**，不白佔本機模型。
- Go 端收到 503、504 都當成 AI 失敗（分類失敗當看不懂、寫回覆失敗轉人工），不重試；log 會記 HTTP 狀態碼。
- 這兩個限制只在模型服務；Swift 的量測程式（`bot`、改進迴圈）照舊排再久都等到。

```bash
swift build -c release && .build/release/model-server       # 預設聽 0.0.0.0:8765
BOT_LIVE_MODEL_TEST=1 ANTHROPIC_API_KEY=[REDACTED] go -C go test -run TestLive -v ./internal/model/   # 實機：Haiku 和 AFM 回答同一個選擇題（沒有 key 時只跑 AFM）
```

- 只聽 `0.0.0.0`，不對外開放。Docker Desktop／Colima 會把容器連 `host.docker.internal` 轉到 Mac 的本機位址，所以開發容器用 `BOT_APPLE_MODEL_URL=http://host.docker.internal:8765` 連得到（來源是 loopback）。如果你的容器環境連不進來（例如來源變成 VM 的網段，模型服務回 403），把 `MODEL_SERVER_HOST` 設成 Mac 在容器網路那一側的位址（只聽那張網卡，不要用 `0.0.0.0`），再用 `MODEL_SERVER_ALLOW=<容器網段，例如 0.0.0.0/24>` 放行；其他來源照樣回 403。（這部分還沒在開發容器實測。）

**AI 回答快取（Go 版）**：照搬 Swift 版（見下面「量測」的說明），`BOT_AI_CACHE=1` 才開，正式伺服器不用。包在排隊外面（命中不用排隊）；key 是請求（題目名稱、指示、輸入文字、欄位定義）＋模型名稱與版本（Apple 再加 macOS 版本，由模型服務的 `/version` 提供）的 SHA256，開頭是 `go-v1`，和 Swift 版的快取不共用但放在同一個資料夾；只存欄位值；失敗或逾時不存；`BOT_AI_CACHE_NS=dev|holdout` 分開存。`bot serve` 結束時印「AI 快取：命中 N 次、實際呼叫 M 次」（題庫工具改接 Go 版後，量測指令也會印）。

- **Haiku 不是 greedy，同一個請求每次回答可能不一樣**（Haiku 5.5 不能調 temperature；Apple 用 greedy 是固定的）。快取讓重跑的結果固定，但做「改了 A 卻影響 B」的比對時要記得：沒命中快取、真的重新呼叫 Haiku 的請求，結果可能本來就會變，不一定是改動造成的。要確認時，對同一個請求多呼叫幾次看看是不是本來就不穩定。

### 對話層（#26）

`go/internal/dialog`：一則顧客訊息 → 分類意圖 → 依「狀態 × 意圖」分派 → 回覆。照搬 Swift 版（tag `swift-final`）的 `Handle`、`Intent`、`Booking`、`RelativeDate`、`State`，改成 CONTEXT.md 的 11 個內建意圖（`booking`、`cancel`、`follow_up`、`claim`、`identity`、`self_intro`、`human_request`、`complaint`、`off_topic`、`unclear`、`general`）。

- **分類**：一次選擇題，附目前狀態（等員工中、預約中問到哪個欄位）和意圖案例庫最像的 6 句（e5 向量，`query: ` 前綴；啟動時先算好案例庫的向量，純 Go 後端兩百句約 40 秒），同一次呼叫抽名字。案例庫有完全相同的句子（狀態也相同）就不問 AI。只有表情、標點不問 AI；1～2 個字被分成預約時改成看不懂（預約中不套用）；名字的防呆照 Swift 版。
- **預約**：狀態機照 Swift 版（相對日期由程式換算、短答不切換語言、人數／日期／時間和原文一致性檢查）；人數由程式從 AI 抄的出處讀出（中文數字會換算，「2大1小」算總數 3，日期、時間、電話裡的數字不算），AI 填的值要是出處裡的人數之一，否則改問人數；已經有人數時，不是在問人數的輪次要有人數出處才能覆蓋（#42）；同一個欄位問了 2 次還拿不到就轉人工（附上已收到的預約資料），但那一輪有收到其他新的預約資料（例如問人數時報了姓名）不算一次（#53）；**不檢查時段、公休日和人數**（ADR 0003），收齊後開預約申請（`ticket`＋`booking` 表）等員工確認。預約關閉時回固定回覆 `booking_off`（帶入聯絡方式 `{contact}`），不轉人工。
- **轉人工**：聲稱／無關累計次數（轉人工設定，預設 3 次）、想找真人時轉人工（之後 AI 不回這位顧客）；轉人工關閉時不轉，回固定回覆 `handoff_off`（帶入商家的聯絡方式 `{contact}`）。其他轉人工設定見下面「員工群組」。
- **語言**：程式判斷（先看文字系統，拉丁字母看常用字）；商家沒啟用的語言回「目前只支援…」並開待處理事項給員工。
- **固定回覆**：商家設定的（資料庫 `fixed_reply`）優先，沒有就用內建的中英日說法（`replies.go`）。
- **一般對話**（#27，ADR 0002 的 2026-10-09 修訂「混合式」，`answer.go`、`verify.go`）：
  1. **挑 FAQ**（`FAQ`）：照 Swift 版，一次選擇題。選項是 FAQ 代號加上問法和答案內容，另有 `price_calc`（算價）、`unclear`、`not_a_question`、`none`；先填 `ask`（逐字抄出問的事），最後填 `kind`（寒暄 `chat`／問商家 `question`；放最前面會影響挑 FAQ，「好」被挑成營業時間，#45）。商家沒有知識段落也沒有 FAQ 時不呼叫 AI。
  2. **挑到 FAQ**：顧客用商家慣用語言時回答案原文；其他語言用 `merchant.Translate` 翻譯，譯文要過下面第 5 步的驗證（FAQ 答案當原文片段），翻譯失敗或沒過就改走第 3 步。
  3. **挑知識**（`Pick`）：挑不到 FAQ 而且是問商家時，清單只列知識段落標題 `K1…`，三個選擇題欄位選最多 3 段，`none` 是不選。一段都沒挑到時，轉人工前再確認一次是不是在問店家（`Ask`，#45）：閒聊回 `general`、沒講出是哪件事回 `unclear`；AI 失敗當成問商家。
  4. **寫回覆**（`Answer`）：只給挑到的段落全文，依序填 `quote`（原文片段）、`source`、`answerable`（yes/no）、`reply`，用顧客的語言。
  5. **程式驗證**：原文片段要出現在出處裡（統一全形半形，空白、標點不算）。回覆裡的數字、時間、電話、網址要出自原文片段或顧客訊息；比對前統一全形半形、中文數字和千分位，時間統一成「時:分」，「晚上八點」對得上 20:00；「9 折」只對得上 9 折、10% off、1割引。回覆的語言要和顧客一致，長度上限 300 字。
  6. 回覆種類都是 `answer`（FAQ 原文、譯文、AI 寫的都是，從 traces 分辨）。答不出來時（沒挑到、`answerable: no`、AI 失敗或驗證沒過）：預約中正在問某個欄位就再問一次（`ask_<欄位>`，不轉人工）；寒暄回 `general`（程式依常見短詞分打招呼、道謝、道別、應答，各回一則固定回覆，種類記在 `chat` 這一步，#54），看不懂回 `unclear`，問商家的走「找不到答案時轉人工」（關著回看不懂）。每一步的結果（`intent`、`faq`、`translate`、`pick`、`ask`、`answer`、`verify`、`reask`、`chat`）放在 `Reply.Traces`，平台層寫進對話紀錄的 `traces`。
  7. **待補問題（#30）**：答不出來的訊息（問商家答不出來、`answerable: no` 是 `unanswerable`，原文片段或譯文驗證沒過是 `verify_failed`，分類或挑 FAQ 時看不懂是 `unclear`）在回覆後記進資料庫 `pending_question`（migration v10，每次出現另記在 `pending_question_hit`），回覆內容不變；有 AI 失敗的訊息不記（分類 AI 失敗時 `intent` 的 trace 多一個 `"error":"ai_failed"`）。相同問題合併計次：有 e5 時語意向量 cosine ≥ 0.92（`pendingMergeThreshold`），沒有 e5 時轉小寫、只留字母和數字後相同；只和待處理、已忽略的合併，已轉成 FAQ／知識／案例的又出現就開一筆新的。新的待補問題的顧客語言不是商家慣用語言時，`bot serve` 在背景翻成慣用語言（照原意翻顧客訊息的 `merchant.TranslateCustomer`，不加語氣詞和表情符號；後台對話紀錄的翻譯按鈕翻顧客訊息也用它，翻機器人回覆照舊用 `merchant.Translate`；`dialog.Bot.Async`；終端機和量測不翻，免得和下一題搶本機模型），翻不出來留空，後台詳細頁打開時補翻。轉人工的那次記下轉人工事項，員工在群組回覆那張（或它的緊急通知）時，員工原話接到 `staff_answer`（預填答案）。

**終端機對話工具**（`go/cmd/chat`，照搬 Swift 版的 `bot`）：每次啟動把 `eval/`（或 `BOT_CONFIG_DIR`）匯入暫存的資料庫，設了 `BOT_CONFIG_FILE=匯出檔` 時改讀那份匯出的設定；模型和快取的環境變數同 `bot serve`，要先開模型服務。`BOT_EXAMPLES_K` 可以改附幾句相似案例（量測用）。

```bash
go -C go build -o bin/chat ./cmd/chat
.build/release/model-server &                                                 # 模型服務（見上面「模型層」）
go/bin/chat                                                                   # 對話：/reset 換新顧客、/as 切換顧客、/pending 待處理、/quit
BOT_AI_CACHE=1 BOT_INTENT_EVAL=eval/intent_set.json go/bin/chat              # 分類題組：另外印「10 類換算」（不含想找真人）
BOT_AI_CACHE=1 BOT_TESTSET=eval/testset.json,eval/faq_set.json go/bin/chat   # 回答題庫：每題一位新顧客，列出答錯、該答沒答、需要計算的（calc-*）和被擋下的（不跑驗證題）
BOT_AI_CACHE=1 BOT_LANGUAGES=zh,en,ja,ko BOT_TESTSET=eval/testset_extra.json go/bin/chat  # 誘答題、跨語言題：另外印「誘答題 編造 N 題」「跨語言 答對 N/M」
BOT_AI_CACHE=1 python3 eval/run_flows.py --bot go-chat --exe go/bin/chat     # 流程題（對話模式＋FLOW_TRACE，可以接對話編號）
```

- 回答題庫每題的結果（ok／wrong／missed、回覆種類、挑到的 FAQ）另外寫到 `eval/results/testset-*.json`，最後印「結果：路徑」（`eval/snapshot.py` 讀）；`HOLDOUT=1` 只跑驗證題（給 `snapshot.py` 比對，輸出不要看）。
- 題目的 `group`：`trap`＝誘答題（資料裡沒有答案，轉人工或說不知道才對；回覆種類是 answer 而且判錯算「編造」），`xlang`＝跨語言題（英日韓文，要用顧客的語言回答、內容和 FAQ／知識段落一致）。
- 分類和回答題庫會印「每則平均 X 秒」（開快取時含命中）；每種模式結束時在 stderr 印「AI 用量：輸入 N token、輸出 M token、估算 US$X」（只算實際呼叫，費用照 `model.Cost` 的價格表，本機模型是 0）。
- `BOT_LANGUAGES=zh,en,ja,ko` 換掉店家啟用的語言（預設中英日；流程題要測「不支援的語言」，所以只在跨語言題開韓文）。
- **後台設定和檔案匯入的比對（#39）**：`go -C go run ./cmd/adminparity -o /tmp/admin.json` 從空的資料庫只透過後台 API 做出森日茶屋的設定（`admin.BuildViaAPI`，翻譯用假模型，不需要模型服務），匯出到檔案；再用 `eval/snapshot.py` 逐題比對兩份設定（同一天、機器空閒時跑）。設定層級的比對是 `go/internal/admin` 的 `TestBuildViaAPIMatchesEvalImport`。

  ```bash
  python3 eval/snapshot.py save import
  BOT_CONFIG_FILE=/tmp/admin.json python3 eval/snapshot.py save admin
  python3 eval/snapshot.py compare import admin
  ```

### 員工群組（#28）

照搬 Swift 版（tag `swift-final`）的店員群組，改用綁定碼綁定群組（不再讀 `TELEGRAM_STAFF_CHAT_ID`）。程式在 `go/internal/telegram/staff.go` 和 `go/internal/dialog/staff.go`。

- **綁定**：第一次啟動時產生綁定碼（存在資料庫 `handoff_setting`），`bot serve` 啟動時記在日誌，`bot bind-code` 也會印出來。把機器人加進員工群組，在群組裡傳「`/bind 綁定碼`」就綁定那個群組（蓋掉之前綁的），綁定後換一組新碼。群組開著隱私模式（Telegram 預設）時機器人只收得到指令和回覆機器人的訊息，所以要用 `/bind`；只傳綁定碼要關掉隱私模式或把機器人設成管理員。群組升級成超級群組（chat id 會變）時自動改綁新的 id。只認綁定的那個群組，其他群組的訊息和按鈕都不理會；還沒綁定時沒有人能接，**不轉人工**，照轉人工關閉回固定回覆 `handoff_off`（附聯絡方式），照樣開一張轉人工事項（註明「沒有員工群組，未轉人工」，不算在等員工；這位顧客已經有一張還沒處理的就不再開）給後台看，通知只記 log「沒有員工群組，未通知」（#55）。`go/bin/chat` 每次開新的資料庫，會綁一個假的群組，量測和已綁群組的商家一樣。
- **通知**：待處理事項（轉人工、聲稱、不支援的語言、追問 3 次的緊急通知、預約、非文字訊息）送到群組：編號、類型、通知文字、顧客名字＋`@username`＋使用者 ID。轉人工附最近 5 則對話（每則最多 100 字）和「結束轉人工」按鈕；預約附「確認／拒絕」按鈕。
- **員工回覆**：員工「回覆」通知（或轉送到群組的圖片），文字原樣加上「店員回覆」轉給顧客，記一筆員工的對話紀錄，群組回「已轉給顧客」（送不出去回「傳給顧客失敗」）。這件事和對應的緊急通知改成已處理、追問次數歸零；轉人工不會因此結束。
- **按鈕**：按「結束轉人工」恢復 AI 回覆（後台對話內容頁也能結束，見「後台」的結束轉人工）；按「確認」「拒絕」把固定回覆送給顧客。處理完通知改成「原通知＋結果（操作：員工名字 @username）」並拿掉按鈕，送不出去時註明；重複按只回「已處理」；舊通知上的「結束轉人工」不會結束之後另一次轉人工。拒絕的原因：回覆那則預約通知，內容會轉給顧客。
- **非文字訊息**：貼圖當閒聊（交給對話層的是貼圖的表情，不呼叫 AI）；圖片、語音、影片、檔案回「目前只看得懂文字」並開待處理事項，原始訊息用 `copyMessage` 轉送到群組、回覆在通知下面，轉送失敗時在通知後面註明。轉人工中不回覆顧客，照樣通知和轉送。
- **依顧客加鎖**：群組訊息和按鈕回呼排成一列依序處理；對話層每位顧客一把鎖，員工的操作和這位顧客的訊息不會插隊。
- **轉人工設定**（資料庫 `handoff_setting`，後台畫面見下面「後台」的轉人工設定）：總開關沿用功能開關 `handoff`；找不到答案時轉人工（預設開；一般對話答不出來的問商家訊息，預約中除外）；聲稱／無關累計次數（預設 3）；員工服務時段（依星期幾、台北時間，空的是全天；時段外照樣轉人工，回覆後面加上固定回覆 `off_hours`）；自動結束（分鐘，預設 30，0 是不自動結束；migration v12 把從沒存過、還是 0 的改成 30，#55；員工最後一次回覆、沒回覆過就從轉人工起算，超過就在顧客下一則訊息時恢復 AI）；提示文字（空的用固定回覆 `handoff`）。

### 後台前端（`web/`）

`web/` 是 [Vue Vben Admin](https://example.invalid/vbenjs/vue-vben-admin) **v5.7.0** 的骨架（UI 用 Ant Design Vue 的 `web-antd`），從官方 release 原始碼複製後精簡：只留 `apps/web-antd`，拿掉其他 UI 版本、`backend-mock`、docs、playground、lefthook（它會在 `pnpm install` 時把 git hook 裝進這個 repo）；npm registry 改成官方的 `example.invalid`；開發模式關掉假資料（`VITE_NITRO_MOCK=false`）。pnpm 版本由 `web/package.json` 的 `packageManager`（pnpm 10.33.4）鎖定，node 跟 `.nvmrc`（22）。

pnpm 要**在 `web/` 裡面**用 corepack 跑（從 repo 根目錄跑 `pnpm -C web` 會拿到錯的 pnpm 版本）：

```bash
cd web
corepack pnpm install --frozen-lockfile
corepack pnpm build:antd      # 輸出到 web/apps/web-antd/dist，Go 的後台伺服器提供這個資料夾
corepack pnpm dev:antd        # 開發模式（http://localhost:5666），/api 轉到 bot serve 的後台（0.0.0.0:8081）
```

目前有登入頁（#34）、商家設定的商家資訊、FAQ、語言（#35）、固定回覆、模型（#38）和預約設定、轉人工設定（#37）；其他頁面在後面的 issue。和 Vben 原本的差別：

- **登入**：接 Go 的 `/api/auth/*`，憑證是 HttpOnly 的 session cookie（前端拿不到 token，沒有 refresh token）。拿掉驗證碼、註冊、忘記密碼、QR code、手機登入和範例帳號。選單、路由寫死在前端（`accessMode: frontend`），路由用 history 模式。
- **介面語言**：繁中（`zh-TW`）、英（`en-US`）、日（`ja-JP`），拿掉簡體中文。順序是：管理員手動切換的（存在瀏覽器的 localStorage）＞登入後商家的慣用語言＞瀏覽器語言＞繁中。共用套件的語言檔在 `web/packages/locales/src/langs/`，後台自己的在 `web/apps/web-antd/src/locales/langs/`。手動選的語言跟目前一樣時也算手動選的：共用套件的語言切換選單每次點都在 window 發 `LOCALE_SELECTED_EVENT`（改了 `packages/locales`、`dropdown-radio-menu.vue`、`language-toggle.vue`，註解標了「aichatbot」，升級 Vben 時要手動合併）。

### 後台（#34）

`bot serve` 另外開一個 port 給後台（`BOT_ADMIN_PORT`，預設 8081，只聽 0.0.0.0）：`/api/` 是後台 API，其他路徑是建好的前端（`BOT_ADMIN_DIR`，預設是**工作目錄**底下的 `web/apps/web-antd/dist`，所以要在 repo 根目錄啟動，或設成絕對路徑）。還沒建置前端時頁面回 503 並提示建置指令。後台和 Telegram webhook 分開兩個 port，webhook 的公開網址碰不到後台。

```bash
.build/gobot create-admin owner            # 建立第一個管理員：密碼從終端機輸入兩次、不顯示（至少 10 個字）
(cd web && corepack pnpm build:antd)       # 建置前端
set -a; source .env; set +a; .build/gobot serve   # webhook 在 8080、後台在 8081
```

- **管理員**：只有一種角色。帳號 3～32 個英數字（可以用 `_`、`-`），密碼至少 10 個字、存 bcrypt；帳號已經有了 `create-admin` 會失敗，不會覆蓋。stdin 不是終端機時讀一行當密碼（給腳本用）。
- **API**：`POST /api/auth/login`（`{"username","password"}`）、`POST /api/auth/logout`、`GET /api/auth/me`（`{"username","language"}`，language 是商家慣用語言對應的介面語言，還沒匯入商家設定時是 `zh-TW`）。除了登入，`/api/` 底下一律要先登入，否則回 401。登入失敗一律回「帳號或密碼錯誤」。
- **意圖案例庫（#36）**：`GET /api/intents`（意圖代號和三語白話名稱 `{"code","labels":{"zh-TW","en-US","ja-JP"}}`）、`GET /api/intent-examples`、`POST /api/intent-examples`（`{"text","intent"}`）、`DELETE /api/intent-examples/{id}`。內建案例不能刪（403）；沒有修改的 API（PUT、PATCH 回 405），要改就刪掉再新增。新增的句子不能和題庫重複（比對規則和 `eval/check_overlap.py` 一樣：轉小寫、只留字母和數字），也不能和案例庫裡已有的句子重複，否則回 409（錯誤訊息寫出是哪個題庫，驗證題組一律只寫「驗證題組」）；題庫是**工作目錄**底下的 `eval/`（`BOT_EVAL_DIR` 可以改），每次新增都重新讀，讀不到就拒絕存檔（500）。
- **session**：32 bytes 隨機 token 放在 cookie（`HttpOnly; Secure; SameSite=Strict; Path=/`），資料庫只存它的 SHA-256；從登入算起 7 天過期，使用時不延長；過期的在每日清理時刪掉。會改資料的請求（POST、PUT、PATCH、DELETE）一律要 `Content-Type: application/json`，否則回 415（防 CSRF）。
- **cookie 有設 Secure，只能走 https**：本機的 `http://0.0.0.0:8081`、`http://localhost:5666` 在 Chrome、Firefox 會當成安全來源，可以登入；**Safari 會丟掉 http 上的 Secure cookie，本機開發請用 Chrome 或 Firefox**；手機要走下面的 Cloudflare 網址。不在本機（loopback）連線拿掉 Secure，是因為 cloudflared 轉進來的正式流量也是從 loopback 連進來，Go 分不出本機開發和正式流量，拿掉的話正式環境的 session cookie 也會沒有 Secure。
- **商家設定**（#35）：`GET`／`PUT /api/merchant`（`{"name","intro","contact","knowledge":[{"title","body"}]}`）、`/api/faq`（`{"faq":[{"key","question","alt_questions","answer"}]}`）、`/api/languages`（`{"preferred_language","languages":[{"code","enabled"}]}`；GET 另外回可勾選的 `supported` 和量過品質的 `measured`）。`/api/merchant`、`/api/faq` 的 GET、PUT 回應另外附上長度警告 `warnings`（#65，`merchant.KnowledgeWarnings`、`merchant.FAQWarnings`，門檻在 `go/internal/merchant/limits.go`，依 `go/cmd/ctxprobe` 的量測訂）：FAQ 超過 50 題時每則訊息先用句子向量挑最像的 20 題給 AI、一次送出的 FAQ 選項太長、知識段落單段或最長 3 段合計太長；只提醒不擋儲存，PUT 的 body 不能帶 `warnings`。每個區塊分開存，`PUT` 只換掉自己的區塊，但用整份設定檢查（必填欄位、FAQ 代號和知識段落標題不能重複、慣用語言要是有勾的語言、語言代碼要在可勾選清單裡），沒過回 400 和錯在哪裡。FAQ 的代號留空時程式給一個沒用過的 `faq-N`。還沒匯入商家設定時 GET 回空白欄位和預設語言（慣用中文、顧客語言中英日），要先存商家資訊，FAQ 和語言才能存。機器人和後台每次都直接讀資料庫，**存了馬上生效，不用重開服務**。
- **預約設定、轉人工設定（#37）**：`GET`／`PUT /api/booking-setting`（`{"enabled"}`，功能開關 `booking`）。`GET`／`PUT /api/handoff-setting`（`{"enabled","on_no_answer","streak_limit","service_hours":[{"weekday","start","end"}],"auto_end_minutes"}`；GET 另外回只讀的 `bind_code` 和 `bound`（有沒有綁好員工群組，不給 chat id））：總開關存到功能開關 `handoff`，其他存到 `handoff_setting`；提示文字 `prompt_text`、綁定碼、綁定的群組不動。檢查：累計次數 1～10、自動結束 0～1440 分鐘、服務時段星期 0～6、時間兩位數的 `HH:MM`、開始和結束不同（開始晚於結束是跨過午夜、延續到隔天）、時段之間不重疊（跨日的也算，週六可以跨到週日），沒過回 400；還沒有商家設定時 PUT 回 400「請先填商家資訊」。`POST /api/handoff-setting/bind-code` 換一組綁定碼（舊碼不能再用，已綁定的群組不變）。關閉時回的文字、轉人工提示不在這兩頁，連到固定回覆頁（`/settings/replies?key=…` 捲到那一則）。功能開關原本的 `off_text` 在 migration v9 拿掉，有填的搬到固定回覆 `booking_off`／`handoff_off` 的慣用語言版本；舊的 `bot export` 備份裡的 `off_text` 匯入時忽略。
- **固定回覆（#38、#37）**：只開放打招呼 `greeting`、預約關閉 `booking_off`（含聯絡方式 `{contact}`）、轉人工提示 `handoff`、轉人工關閉 `handoff_off`（含聯絡方式 `{contact}`）、時段外提示 `off_hours` 五則。`GET /api/fixed-replies` 回 `{"preferred_language","languages":[啟用的語言],"replies":[{"key","labels":{"zh-TW","en-US","ja-JP"},"texts":{"語言":["一句", ...]}}]}`（存的時候多句用 `|` 接起來，機器人隨機挑一句）；還沒設定的那則，慣用語言給預設文字（不含店名）、其他語言空白。`PUT /api/fixed-replies/{key}`（`{"texts":{...}}`）：慣用語言版本有改時，用目前開著的模型（AFM 或 Haiku，照常記用量）逐句翻成其他啟用的語言，取代送來的那些語言；沒改時照送來的存（修正翻譯），只有空著的語言才翻。譯文裡的 `{name}` 這類變數要和原文一樣，否則算翻譯失敗；任何一句翻譯失敗就整次不存、回 502。不能改的 key 回 404，內容沒過檢查（慣用語言至少一句、句子不能空白或含 `|`、只能存啟用的語言）回 400。機器人回覆時，商家改過的固定回覆優先（見上面「對話層」）。
- **模型（#38）**：`GET`／`PUT /api/model`（`{"model":"afm|haiku","monthly_budget_usd":10}`；GET 另外回 `haiku_available`，只說有沒有設定 API key，不回 key 本身）。沒有 key 時切到 haiku、上限是負數或不是數字回 400。
- **自訂意圖＋分類自我測試（#31）**：`GET`／`PUT /api/custom-intents`（`{"custom_intents":[{"code","name","description","examples":[...],"action":"fixed_text|handoff|knowledge","text","knowledge":[段落標題]}]}`；GET 另外回 `knowledge_titles`）整份取代，最多 10 個。代號留空時程式給一個沒用過的 `custom_N`（送來的代號只能是 `custom_1`～`custom_10`、不能重複，不能蓋過內建意圖），後台只顯示名稱。檢查：名稱、說明必填，名稱不能和內建意圖的白話名稱或其他自訂意圖重複；例句至少 3 句，不能和案例庫、其他自訂意圖的例句重複（沒過回 400），也不能和題庫重複（回 409，同新增意圖案例）；回固定文字要填文字；只根據知識回答要選存在的知識段落（之後改掉段落標題，商家資訊會存不進去，要先改自訂意圖）。內建意圖不能改，沒有對應的 API。`GET /api/intents` 接著列出自訂意圖，案例庫也能新增屬於自訂意圖的案例；刪掉自訂意圖時，屬於它的商家案例一起刪掉。對話層：分類時自訂意圖排在 `general` 前面（說明附例句），例句也加進完全相同比對和相似案例；**沒有自訂意圖時送給 AI 的提示和原本一模一樣**。回固定文字和轉人工提示只寫慣用語言，其他語言翻譯後驗證（同 FAQ 答案），沒過回原文；轉人工關閉時照轉人工的關閉文字；只根據知識回答時不挑 FAQ、只從指定段落挑知識，答不出來走「找不到答案時轉人工」。存好後在背景跑**分類自我測試**：`GET /api/custom-intents/self-test`（`{"status":"idle|running|done|failed","done","total","report":{"total","conflicts":[{"text","state","expected","got"}],"failed":[...]}}`）、`POST` 手動重跑（202）。同時間只跑一組，新的取消還沒跑完的；結果只放記憶體，重開服務後是 idle。
- **待補問題（#30）**：`GET /api/pending-questions?status=open|ignored|resolved`（預設 open，次數多的在前，同次數最近出現的在前）；`GET /api/pending-questions/{id}`（加上每次出現 `hits`，各附那位顧客前、後各 10 則對話 `conversation`；對話紀錄 90 天後清掉，待補問題和每次的原話一直保留；還沒有譯文時用目前的模型補翻）。動作（待處理、已忽略的都可以做；已處理的回 409）：`POST …/{id}/faq`（`{"question","alt_questions","answer"}`，只在 FAQ 最後新增一題，不整份覆蓋）、`POST …/{id}/knowledge`（`{"title","body"}`，標題是現有段落就接在後面，否則新增一段）、`POST …/{id}/example`（`{"text","intent"}`，檢查同新增意圖案例，重複回 409）、`POST …/{id}/ignore`（`{}`）。成功回這筆待補問題，轉成 FAQ／知識／案例的改成 `resolved`（`resolution` 記轉成什麼），忽略的改成 `ignored`（之後相同的問題照樣計次）。頁面上轉成 FAQ 時問題預填顧客原話、答案預填員工回覆。
- **結束轉人工（#55）**：`GET /api/conversations/{platform}/{userID}` 多回 `human`（顧客在轉人工中），對話內容頁這時顯示「轉人工中」和「結束轉人工」按鈕。`POST /api/conversations/{platform}/{userID}/end-handoff`（`{}`）效果同員工群組的按鈕：狀態歸零、關掉還沒處理的轉人工事項、對話紀錄記一筆員工的「（結束轉人工）」，和顧客訊息共用同一把顧客鎖（`bot serve` 把同一個 `dialog.Bot` 傳給後台）。不在轉人工中回 409，沒有接上對話層回 503。
- **後台頁面**：選單「商家設定」底下有商家資訊、FAQ、語言、固定回覆、預約設定、轉人工設定、模型七頁，另外有「意圖案例庫」「自訂意圖」「待補問題」，手機也能操作。中英日以外的語言會提示「品質沒量過」。

**對外：獨立網址＋Cloudflare Access**。後台用另一個 hostname（例如 `example.invalid`），和 webhook 的 hostname 分開，前面加 Cloudflare Access，只讓指定的人連得到：

1. cloudflared 的固定 tunnel 設定（`~/.cloudflared/config.yml`）加一條 ingress：

   ```yaml
   tunnel: <tunnel 名稱或 ID>
   credentials-file: /Users/<你>/.cloudflared/<tunnel ID>.json
   ingress:
     - hostname: example.invalid       # Telegram webhook（TELEGRAM_PUBLIC_URL）
       service: http://localhost:8080
     - hostname: example.invalid     # 後台
       service: http://localhost:8081
     - service: http_status:404
   ```

   再幫後台的 hostname 建 DNS：`cloudflared tunnel route dns <tunnel 名稱> example.invalid`，重啟 `cloudflared tunnel run <tunnel 名稱>`。
2. Cloudflare 的 Zero Trust 主控台 → Access → Applications → Add an application → **Self-hosted**：Application domain 填 `example.invalid`；Policy 選 Allow，Include 用 **Emails** 列出可以進後台的人（登入方式用預設的 One-time PIN 就可以）。**webhook 的 hostname 不要加 Access**（Telegram 過不了驗證）。
3. 用手機開 `https://example.invalid`：先過 Cloudflare Access 的 email 驗證，再用 `create-admin` 建的帳號登入。

### 多語言句子向量（技術驗證，#23）

`go/internal/embed` 用 hugot＋`intfloat/multilingual-e5-small` 算句子向量，給意圖案例庫找相似句（中英日可以互相比）。量測結果和結論見 [`docs/研究總結.md`](docs/研究總結.md) 第 15 節。

- **模型檔**（470 MB，不進 git）：Hugging Face 上 `intfloat/multilingual-e5-small` 的 `onnx/` 資料夾（`model.onnx`、`tokenizer.json`、`config.json` 等），放到 `~/.cache/aichatbot/models/multilingual-e5-small/`，或用 `BOT_EMBED_MODEL=<資料夾>` 指定。沒有模型檔時 `go test` 會跳過這部分。
- **macOS**：預設是純 Go 後端，不用 cgo、不用另外裝東西。
- **開發容器**：一樣是純 Go 後端，Dockerfile 不用多裝東西。`.devcontainer/devcontainer.json` 已經把 Mac 的 `~/.cache/aichatbot/models` 唯讀掛到容器的 `/opt/aichatbot/models`，並設好 `BOT_EMBED_MODEL`，模型檔放在 Mac 上就好（用 `docker run` 時自己加 `-v ~/.cache/aichatbot/models:<容器的 HOME>/.cache/aichatbot/models:ro`）。
- **ONNX Runtime 後端（選用，比較快，還沒量）**：`go build -tags ORT`，要 cgo 和 libonnxruntime。macOS 用 `brew install onnxruntime`；Linux arm64 從 onnxruntime 的 GitHub Releases 下載 `onnxruntime-linux-aarch64-<版本>.tgz`，把 `lib/libonnxruntime.so*` 放到 `/usr/lib/`（要寫進 Dockerfile 的話放在這一步）。函式庫不在系統預設路徑時，用 `ONNXRUNTIME_LIB=<.dylib／.so 路徑>` 指定。

```bash
go -C go test -v ./internal/embed/                          # 同義句比無關句像；日文、英文訊息找得到中文案例
go -C go test -bench . -run '^$' ./internal/embed/          # 一句多久
go -C go run ./cmd/xlang                                    # 跨語言抽查（也可以在後面接自己的句子）
go -C go run ./cmd/embed -o ../eval/results/e5.json         # 案例庫＋分類題組的句子先算成向量
BOT_AI_CACHE=1 BOT_EMBED_FILE=eval/results/e5.json BOT_INTENT_EVAL=eval/intent_set.json .build/release/bot  # Swift 版改用 e5 向量跑分類題組
```

`BOT_EMBED_FILE` 是 Swift 版的開關，只用來做 A/B 量測，預設不開：有設的時候，找相似案例改用這個檔裡的向量，而且不限同一種文字系統。檔案裡沒有的句子就不附案例，`BOT_DEBUG=1` 會印出「向量檔沒有這句」。程式在分類前改寫過的句子（例如相對日期換成「10/27」）要接在 `cmd/embed` 後面補算。

實機驗證（需要 bot token 和對外網址）：設好 `.env`、開 tunnel（見「Telegram 伺服器」）後啟動 `serve`，日誌出現「已設定 webhook」，再用 Telegram 私訊機器人，應該先看到「正在輸入」，再收到上面那句固定的話。

## 執行

先編譯，再執行 `.build/release/bot`（改了 `Sources/` 的程式要重新編譯）：

```bash
swift build -c release && .build/release/bot
```

對話中可以使用這些指令：

| 指令 | 作用 |
|---|---|
| `/reset` | 換一位新顧客（清空目前顧客的對話狀態） |
| `/as 使用者ID` | 切換成另一位顧客（預設 `user1`；每位顧客的狀態分開保存，可以模擬多位顧客交錯對話） |
| `/reload` | 重新讀取 `store.json`、`replies.json`、`faq.txt`、`examples.jsonl` |
| `/me 名字` | 手動設定目前顧客的名字 |
| `/pending` | 查看店員待處理清單（每項有編號） |
| `/reply 編號 文字` | 店員回覆：把文字（加上「店員回覆」前綴）轉給那位顧客，之後不再算「等店員中」 |
| `/confirm 編號` | 確認預約申請，通知顧客 |
| `/reject 編號 [原因]` | 拒絕預約申請，通知顧客（可附原因） |
| `/endhuman [使用者ID]` | 結束轉人工（不填就是目前顧客），之後恢復 AI 回覆 |
| `/quit` | 結束對話（也可以用 `/exit`） |

對話層測試（不呼叫本機模型，會呼叫 AI 的路徑用照腳本回答的假模型 `ScriptedModel`）：`swift build && .build/debug/bot-tests`。

想看每次 AI 判斷的結果，可以開除錯模式：

```bash
BOT_DEBUG=1 .build/release/bot
```

對話紀錄會存到 `eval/logs/`，也可以用環境變數 `CHAT_LOG_DIR` 指定其他位置。

### Telegram 伺服器

`server` 執行檔是 Telegram 的 webhook 伺服器：顧客私訊機器人，核心（`Bot.receive`）想好之後回覆，處理期間會顯示「正在輸入」。顧客狀態、已處理的 update、對話紀錄、待處理事項和預約申請存在 SQLite（`BOT_DB`），伺服器重開後可以接續（預約通知上的按鈕照樣有效）。其他群組、頻道一律忽略。

店員群組（`TELEGRAM_STAFF_CHAT_ID`）：

- 需要店員處理的事（知識庫沒有的問題、聲稱、不支援的語言、追問 3 次的緊急通知、轉人工、預約）開成待處理事項，通知送到店員群組：編號、類型、顧客名字（聊天中問到的名字，沒有就用 Telegram 名稱）＋`@username`＋使用者 ID、原話。緊急通知會引用原本那張。
- 店員在群組裡用「回覆」那則通知，文字原文（不經過 AI）加上依顧客語言的前綴（💬 店員回覆／Reply from our staff／スタッフより）轉給顧客；待處理事項改成已回覆，顧客不再算「等店員」。回覆緊急通知時對應回原本那張，兩張都改成已回覆；已回覆過的通知照樣可以再回。結果回在那則通知上：「✅ 已送給顧客」或「❌ 傳送失敗」。回覆照片、貼圖時不轉送，回「目前只能轉送文字」。
- 轉人工（連續 3 次無關或聲稱）的通知另外附上原因、最近 5 則對話（顧客／機器人／店員）和「結束轉人工」按鈕。轉人工期間 AI 不回覆（顧客的訊息照樣寫進對話紀錄），店員「回覆」那則通知就能和顧客來回對話；處理完按「結束轉人工」，AI 恢復回覆（不另外通知顧客），群組回「✅ 已結束轉人工」。重複按或按到舊通知只提示「已經結束轉人工了」。
- 預約通知附「✅ 確認／❌ 拒絕」按鈕。店員按下後，顧客收到依語言的「預約已確認」或「預約無法安排」，群組那則通知改成「原通知＋✅ 已確認／❌ 已拒絕（操作：店員名稱 @username）」並拿掉按鈕；預約狀態和操作的店員存進資料庫。已處理過的再按只回「已處理」；其他群組的按鈕回呼不處理。按鈕拒絕不帶原因，需要說明時用「回覆」那則通知補充。
- 沒有回覆通知的發言（店員之間聊天）只記錄，機器人不回話；其他群組的回覆、按鈕不處理。
- 顧客傳貼圖當閒聊；圖片、影片、語音、檔案回「目前只看得懂文字」，並開待處理事項（附說明文字）通知店員群組。轉人工期間不回覆顧客，但照樣通知店員群組。
- 同一位顧客的訊息和店員操作依序處理，店員的修改不會被正在等 AI 的那則訊息蓋掉。

設定都用環境變數（**不要寫進 git，也不要貼到 issue 或日誌**；可以放在 `.env`，已列入 `.gitignore`。範本是 `.env.example`：`cp .env.example .env` 再填值）：

| 變數 | 必填 | 內容 |
|---|---|---|
| `TELEGRAM_BOT_TOKEN` | ✅ | BotFather 給的 token |
| `TELEGRAM_WEBHOOK_SECRET` | ✅ | 自己產生的亂數（只能用英數字、`_`、`-`，最多 256 個字），例如 `openssl rand -hex 32` |
| `TELEGRAM_PUBLIC_URL` | ✅ | Tunnel 的對外網址（`https://` 開頭），webhook 會設成 `<網址>/telegram/webhook` |
| `TELEGRAM_STAFF_CHAT_ID` | | 店員群組的 chat ID（通常是負數）；沒設就不送通知，只記錄 |
| `PORT` | | 伺服器的 port，預設 8080（只聽 0.0.0.0） |
| `BOT_CONFIG_DIR` | | 店家設定檔的資料夾（見下面「設定檔」），預設是專案的 `eval/`；執行檔搬到別的資料夾或部署到測試機時要設 |
| `BOT_DB` | | SQLite 資料庫檔案，預設是**工作目錄**底下的 `data/bot.sqlite`（在專案根目錄執行時就是專案的 `data/`，已列入 `.gitignore`）；用 launchd 執行時要設 `BOT_DB` 或 `WorkingDirectory` |

資料庫：啟動時自動建表（GRDB migration）。對話紀錄每則訊息一筆（顧客、機器人、店員），機器人那筆附上每次 AI 呼叫的題目名稱、送進去的文字、各欄位回傳值、毫秒數、是否失敗（不存 instructions）。啟動時和之後每 24 小時清一次舊資料：對話紀錄保留 90 天、update 紀錄保留 7 天。終端機版 `bot` 用記憶體資料庫，不寫入正式資料。

啟動（會先檢查本機模型能不能用，再呼叫 `setWebhook`；`setWebhook` 失敗只記錄錯誤，伺服器照常執行）：

開發時可以用臨時 tunnel（不需要 Cloudflare 帳號，但每次重開網址都會變，要更新 `TELEGRAM_PUBLIC_URL` 再重啟伺服器；步驟見 `.env.example`）；測試機和正式版改用綁網域的固定 tunnel：

```bash
cloudflared tunnel --url http://localhost:8080   # 臨時 tunnel，另開一個終端機
cloudflared tunnel run <tunnel 名稱>              # 固定 tunnel：ingress 指到 http://localhost:8080
```

```bash
set -a; source .env; set +a
swift build -c release && .build/release/server
```

伺服器的行為：

- 只收一般訊息和按鈕回呼；`X-Telegram-Bot-Api-Secret-Token` 不符回 401；同一個 `update_id` 只處理一次；先回 200，再依序在背景處理。
- `/start` 回打招呼，不經過 AI。
- 呼叫 Telegram API 失敗時，等 1、2、4 秒重試 3 次，仍失敗就記錄錯誤，伺服器繼續服務（「正在輸入」不重試）。
- 日誌不會寫出 token、secret 和含 token 的 API 網址。

Webhook 層測試（對伺服器送真的 HTTP 請求，Telegram API 換成本機的假端點，不需要本機模型和 bot token）：`swift build && .build/debug/telegram-tests`。

### 設定檔（預設在 `eval/`）

| 檔案 | 內容 |
|---|---|
| `store.json` | 店家規則：菜單、價格、折扣、營業時間、公休日、人數上限、刷卡門檻 |
| `faq.txt` | FAQ：問法和中英日答案（一行一題，用「` \| `」分隔：id、問法、中文答案，後面可以再接英文、日文答案） |
| `replies.json` | 所有固定回覆（中英日，同一則可以放多種說法，隨機挑一句）；每個模板一定要有 `zh` 或 `*`，`en`、`ja` 缺的時候用 `zh` |
| `examples.jsonl` | 案例庫：已標好意圖的例句，給分類時參考（設 `NO_EXAMPLES=1` 可以關閉；沒有這個檔就是空的） |

- 位置：設 `BOT_CONFIG_DIR` 指定資料夾，沒設就用編譯時專案的 `eval/`（開發和量測照舊）。執行檔搬到別的資料夾或部署到測試機時一定要設。
- 驗證：`bot` 和 `server` 啟動時會檢查四個檔（檔案在不在、JSON 格式、必要欄位、`faq.txt` 和 `examples.jsonl` 每一行的格式，`#` 開頭和空行略過），也會檢查程式用到的回覆模板在 `replies.json` 都有（清單是 `Sources/BotCore/Config.swift` 的 `usedReplyKeys`，新增模板時要加進去，`bot-tests` 會掃原始碼檢查）。有錯就印出哪個檔案、哪裡錯（欄位或行號），然後結束（exit code 1）。
- 終端機版的 `/reload` 讀到有錯的設定時，會印出錯誤、繼續用舊的設定。

## 量測

量測分成三個題組。每個題組都有**開發題組**（可以看失敗明細、拿來改進）和**驗證題組**（只看總分）。各題組的題數請見 [`loop/README.md`](loop/README.md#題組)。

| 題組 | 開發題組 | 驗證題組 |
|---|---|---|
| 完整流程（多輪對話） | `eval/flow_tests.json` | `eval/holdout/flow_holdout.json` |
| 分類 | `eval/intent_set.json` | `eval/holdout/intent_holdout.json` |
| 單題（挑 FAQ） | `eval/testset.json`＋`eval/faq_set.json` | `eval/faq_set.json` 中標記 `holdout` 的題目 |

> ⚠️ 驗證題組（`eval/holdout/` 和 `faq_set.json` 的 `holdout` 題）**不要打開來看，也不要拿來調整程式**；一旦被看過，它就無法再用來檢查是否過度擬合。

以下指令都在專案根目錄執行，而且要先跑過 `swift build -c release`。AI 用 greedy 取樣，同一題每次結果相同，**每題跑 1 次就夠**（`--runs`、`ROUNDS` 設成大於 1 只會花更多時間，結果一樣）。

**完整流程**：把多輪對話餵給機器人，逐步比對回覆：

```bash
python3 eval/run_flows.py --bot bot --exe .build/release/bot
```

- 只跑指定的對話：在最後加上對話 ID，例如 `python3 eval/run_flows.py --bot bot --exe .build/release/bot F04 F12`。
- 換題庫用 `--tests`；另外輸出 JSON 總結用 `--json`。`--bot` 只是結果檔名和輸出上的標籤，實際執行的是 `--exe` 指定的執行檔。
- 沒加 `--bot` 和 `--exe` 時，預設測的是舊版 `chat.swift`。

**分類**（`ROUNDS` 是輪數，預設 1）：

```bash
BOT_INTENT_EVAL=eval/intent_set.json .build/release/bot
```

**單題**（`ROUNDS` 預設 1；設 `IDS=id1,id2…` 只跑指定題目，`eval/quick_ids.txt` 是一組快速題組）：

```bash
BOT_TESTSET=eval/testset.json,eval/faq_set.json .build/release/bot
```

```bash
IDS=$(cat eval/quick_ids.txt) BOT_TESTSET=eval/testset.json,eval/faq_set.json .build/release/bot
```

**檢查重複**：檢查案例庫、FAQ 問法和各題庫之間有沒有重複的句子。案例庫和題庫重複時，分數只是在測模型有沒有背下答案，所以兩者的重複數必須是 0：

```bash
python3 eval/check_overlap.py
```

**AI 排隊**：所有 AI 呼叫都經過同一個排隊，同一台機器一次只跑一個（本機模型同時跑兩個呼叫，品質會變差）。程序之間用檔案鎖排隊，鎖檔預設是 `~/Library/Caches/aichatbot/ai.lock`，可以用 `BOT_AI_LOCK` 指定。量測結束和終端機離開時會印出「等鎖總時間」。設 `BOT_AI_LOG=檔案` 會把每個 AI 呼叫的程序編號、開始和結束時間附加寫進那個檔案，再用 `eval/check_ai_log.py` 檢查有沒有重疊：

```bash
rm -f /tmp/ai.log
BOT_AI_LOG=/tmp/ai.log BOT_INTENT_EVAL=eval/intent_set.json ROUNDS=1 .build/release/bot &
BOT_AI_LOG=/tmp/ai.log BOT_INTENT_EVAL=eval/intent_set.json ROUNDS=1 .build/release/bot; wait
python3 eval/check_ai_log.py /tmp/ai.log
```

**一次跑完全部題組（Go 版驗收）**：`python3 eval/verify.py` 跑測試和真模型量測（分類、單題、誘答題、跨語言題、流程題），和 `eval/baseline-go.json` 比較，退步的標 ❌（結束碼 0 沒有退步、1 有退步、2 量測無效）。另外印每則平均時間（開快取時含命中，只顯示不比較），`--model haiku` 時再印 AI 用量（不含流程題）：

```bash
python3 eval/verify.py                      # 測試＋量測（--no-checks 跳過測試，--checks-only 只跑測試）
python3 eval/verify.py --no-cache           # 不用 AI 回答快取：量真實的每則平均時間
python3 eval/verify.py --model haiku        # 改用 Haiku 量（要 ANTHROPIC_API_KEY，沒有就結束碼 2）：只列分數、token 數和估算費用，不和基準比
python3 eval/verify.py --store salon        # 只量一家店（tea 森日茶屋、salon 霧光髮廊；預設 all 兩家都量）
```

新指標「誘答題 編造」（越少越好）和「跨語言 答對」來自 `eval/testset_extra.json`，另外跑（店家多開韓文），不算進單題題組的總計。`loop/score.py` 是 Swift 版的評分，已停用。

**兩家測試商家**：森日茶屋（茶飲店，`eval/` 的設定和題組）之外，還有虛構的美髮沙龍「霧光髮廊」（`eval/stores/salon/`）：`config.json` 是後台匯出檔格式的商家設定（`go/bin/chat` 用 `BOT_CONFIG_FILE` 讀；內建意圖案例照 `eval/examples.jsonl` 補上），`intent_set.json`、`testset.json`（含誘答題）、`flow_tests.json` 是它的開發題組。`verify.py` 兩家分開量，每家各一組指標（名稱前面標店名；沙龍的代號是 `salon.*`）和逐題比對，任何一家總數退步或有題目從答對變答錯都標 ❌。分數基準兩家都存在 `eval/baseline-go.json`，沙龍的逐題基準是 `eval/stores/salon/baseline-items.json`。沙龍的流程題照「沙龍應該怎麼做」寫（預約要問服務項目、不問人數），現在的預約流程做不到，基準記的是現況。沙龍還沒有驗證題組（之後由人另外寫）。

各項測試的詳細結果會存到 `eval/results/`（已列入 `.gitignore`）。

**逐題比對**：改寫程式但不想改變行為時（例如重構），先記下各題組（開發：流程、分類、單題、誘答與跨語言；驗證：流程、分類、`faq_set.json` 的驗證題）每一題的結果，改完再記一次，然後逐題比對。`save` 預設先編譯 `go/bin/chat` 再測它（`--exe` 可以指定別的執行檔）。結果存到 `eval/results/snapshots/<名稱>/`。開發題組會列出不同的題目；驗證題組只顯示「相同」或「不同 N 題」，不顯示題目。題目的預約日期是相對今天算的，前後兩次要在同一天跑：

`save` 預設兩家店都量（沙龍存在 `salon/` 底下），`--store tea|salon` 只量一家：

```bash
python3 eval/snapshot.py save before        # 改之前
python3 eval/snapshot.py save after         # 改之後
python3 eval/snapshot.py compare before after
```

**AI 回答快取**：`verify.py`、`snapshot.py` 預設開（`BOT_AI_CACHE=0`、`verify.py --no-cache` 關掉）；單獨執行 `go/bin/chat` 的量測、對話模式或 `run_flows.py` 時要自己設 `BOT_AI_CACHE=1`。Go 版（`go/internal/model/cache.go`）和 Swift 版用同一個快取目錄。AI 用 greedy，同一個請求每次回答都一樣，所以同一個請求（題目名稱、指示、輸入文字、欄位定義、模型版本、macOS 版本都相同）第二次直接用存好的回答，不排隊、不呼叫本機模型；行為不變的重構幾乎全部命中，改了某個指示也只有用到它的請求會真的呼叫模型。量測結束會印「AI 快取：命中 N 次、實際呼叫 M 次」。

- 快取放在 `~/Library/Caches/aichatbot/ai/`（`BOT_AI_CACHE_DIR` 可改），檔名是請求的 SHA256，內容只有 AI 回的欄位值（輸入本身不存；但欄位值可能含原文片段，例如抄出的日期原文、「顧客想做什麼」的描述）；AI 失敗或逾時不存
- 驗證題組分開存在 `holdout/`（`BOT_AI_CACHE_NS=holdout`，腳本跑驗證題組時自動帶），改進迴圈和 issue 迴圈的 AI 都不能讀快取
- macOS 更新（版本或 build 號變了）或改了 `AppleModel.version`（改生成參數時要一起改）會自動失效
- 清快取：`python3 eval/clear_ai_cache.py [--ns dev|holdout|all]`，或 `rm -rf ~/Library/Caches/aichatbot/ai`
- 關掉快取時行為完全不變；Telegram 伺服器不使用快取

## 自動改進迴圈

`loop/` 是一套自動改進流程。每次迭代時，AI 會先提出一個修改，再由 `loop/gate.py` 檢查、評分，並依照固定規則決定保留（commit）或還原。這個決定完全由程式做出，AI 無法干預。

```bash
loop/run.sh 10                    # 跑 10 次迭代
python3 loop/gate.py --baseline   # 重新量測目前的基準分數
python3 loop/gate.py --dry-run    # 只做檢查，不評分
touch <worktree>/loop/STOP        # 等這次迭代跑完後停止
```

`loop/run.sh` 會自動開 worktree（`../aichatbot-loop-<時間>`，分支 `loop/<時間>`）去那裡跑，跑完推分支開 MR，不直接改這裡的 `main`。迴圈會占用本機模型，執行期間在同一台 Mac 跑其他量測會互相排隊、兩邊都變慢。保留規則、紀錄檔和注意事項詳見 [`loop/README.md`](loop/README.md)。

## 目錄結構

| 路徑 | 內容 |
|---|---|
| `Package.swift` | Swift Package 設定（只需要 Command Line Tools） |
| `Sources/BotCore/` | **主程式的核心**（訊息進、動作出）：分類、挑 FAQ、讀訂單、預約、依狀態分派、模型層與 AI 排隊（自動改進迴圈改的是這裡） |
| `Sources/Storage/` | SQLite 資料庫（GRDB）：顧客狀態、已處理的 update、對話紀錄 |
| `Sources/bot/` | 執行檔 `bot`：終端機對話，也內建分類和單題量測模式 |
| `Sources/Telegram/` | Telegram 平台層：webhook 路由、update 分流、Telegram API 用戶端（含重試）、「正在輸入」 |
| `Sources/server/` | 執行檔 `server`：Telegram webhook 伺服器 |
| `Sources/model-server/` | 執行檔 `model-server`：模型服務，只包 Apple 本機模型的 HTTP 服務，給 Go 版呼叫（見上方「模型層」） |
| `Tests/` | `bot-tests`（對話層）、`telegram-tests`（webhook 層），都用假模型，不需要本機模型 |
| `go/` | Go 版服務（見上方「Go 服務」）：`cmd/bot`（執行檔）、`internal/store`（SQLite）、`internal/merchant`（商家設定、匯入匯出）、`internal/telegram`（webhook、API 用戶端）、`internal/model`（模型層：Haiku、Apple 模型服務、排隊、快取、費用） |
| `.env.example` | Telegram 伺服器的環境變數範本 |
| `data/` | 伺服器的 SQLite 資料庫（不進 git） |
| `eval/store.json`、`faq.txt`、`replies.json`、`examples.jsonl` | 設定檔（見上方說明） |
| `eval/flow_tests.json`、`intent_set.json`、`testset.json`、`faq_set.json`、`testset_extra.json` | 開發題組（`testset_extra.json` 是誘答題和跨語言題） |
| `eval/holdout/` | 驗證題組（不要看，也不要拿來調整） |
| `eval/stores/salon/` | 第二家測試商家（美髮沙龍「霧光髮廊」）：商家設定（後台匯出檔格式）、開發題組、逐題基準 |
| `eval/run_flows.py` | 完整流程測試 |
| `eval/verify.py` | Go 版驗收：測試＋真模型量測，和 `eval/baseline-go.json` 比較 |
| `eval/snapshot.py` | 逐題比對（Go 版）：記下各題組每一題的結果，比對改寫前後是否相同 |
| `eval/clear_ai_cache.py` | 清掉 AI 回答快取（`BOT_AI_CACHE=1` 時存的回答，Go 版和 Swift 版共用） |
| `eval/testdates.py`、`relativize_dates.py` | 題庫的相對日期標記：執行時依當天換算；把寫死的日期改成標記 |
| `eval/check_ai_log.py` | 檢查 AI 呼叫有沒有在時間上重疊（驗證跨程序排隊） |
| `eval/check_overlap.py` | 檢查案例庫和題庫是否重複 |
| `eval/quick_ids.txt` | 快速題組的題目代號 |
| `loop/` | 自動改進迴圈（`run.sh`、`gate.py`、`PROGRAM.md`、紀錄檔；`score.py` 已停用） |
| `docs/研究總結.md` | 研究過程、測試結果、目前架構、設計原則 |
| `docs/開發計畫.md` | 開發順序、平台與部署的決定（對應 GitLab issue） |
| `docs/TODO.md` | 待辦：暫緩、待決定、待修、已完成 |

研究用的舊程式（**非主線**，保留做對照）：

| 檔案 | 內容 |
|---|---|
| `eval/chat.swift`、`shop.txt` | 舊版互動客服和它的店家設定 |
| `eval/run_testset.swift` | 用單題題庫評估舊版的自由回答。設 `MODEL=haiku` 時會呼叫 Claude API，需要 `ANTHROPIC_API_KEY` 並**產生費用** |
| `eval/model_compare.swift` | 比較本機模型和 Apple 私有雲模型；私有雲需要另外申請權限 |
| `eval/count_tokens.swift` | 用 Anthropic token 計數 API 估算 Haiku 的費用。需要 `ANTHROPIC_API_KEY`，**會把資料送到外部 API** |
| `eval/fm_eval.swift`、`fm_eval_v2.swift`、`intent_eval.swift`、`router_eval.swift`、`name_eval.swift` | 早期的單項實驗 |
| `eval/past_conversations.txt` | 早期測試紀錄（事後重建） |

## 注意事項

- **同時跑多個測試只會變慢**：本機模型同時處理兩個 AI 呼叫時品質會明顯變差，所以所有 AI 呼叫都經檔案鎖排隊（等鎖的時間不算進 15 秒逾時），不會真的同時跑；但每個程式都要輪流等，時間會拉長。不同 worktree 編出的執行檔用同一把預設鎖（`~/Library/Caches/aichatbot/ai.lock`），改了 `BOT_AI_LOCK` 就不會一起排隊。
- **macOS 更新後要重量基準**：greedy 只保證同一個模型版本結果相同，系統更新換了模型，分數就可能改變（`python3 loop/gate.py --baseline`）。
- **題庫的日期用相對標記**：預約日期寫成 `{{date:+1:M/D}}` 這類標記，執行時依當天換算（規則見 `eval/testdates.py`），不會過期。新增題目時也要用標記，不要寫死日期；`BOT_TODAY=YYYY-MM-DD` 可以指定「今天」。測相對說法時，預期的日期用 `{{today:+1:MM-DD}}`（今天 + N 天）或 `{{weekday:next:2:MM-DD}}`（下週二；`this`＝這週X、`near`＝只說週X，1＝週一）。

