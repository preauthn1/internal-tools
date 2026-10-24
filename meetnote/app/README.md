# meetnote 開發指南

Tauri v2 + Vue 3（TypeScript）+ Rust。專案總覽和決策理由在 [../README.md](../README.md)，規劃與進度在 [PLAN.md](PLAN.md)，開發規則在 [../AGENTS.md](../AGENTS.md)。

## 執行

```bash
./dev.sh          # 開發版
./bundle-dev.sh   # 打包成 debug 版 .app 再開：要錄電腦播放的聲音時用這個
```

- 需要 Command Line Tools、Homebrew 的 rustup（`/opt/homebrew/opt/rustup/bin`）、Node 22、cmake。
- 模型取自 `../spike/models/`，Claude API key 取自 `../.env`。開發版不讀鑰匙圈。
- `tauri dev` 會啟動 Vite（`http://localhost:1420`），App 視窗載入這個網址：改 Vue 會即時更新，改 Rust 會自動重新編譯並重開 App。
- 正式打包（`npm run tauri build`）會把前端編進執行檔，不啟動伺服器。
- **錄電腦播放的聲音**的權限（`NSAudioCaptureUsageDescription`）只會給 .app：從 `dev.sh`（終端機）啟動時錄到的是靜音，要用 `bundle-dev.sh`。

## 架構

```text
            ┌──────────── 錄音執行緒（recorder.rs，cpal）────────────┐
麥克風 ───▶ │ 轉單聲道 → 累計音量和時間（Meter）→ 交給辨識引擎         │
            └──────────────────────────┬────────────────────────────┘
電腦聲音 ─▶ system_audio.rs：meetnote-tap（Core Audio process tap）│ Msg::SystemAudio
                                       │ Msg::Audio
            ┌──────────── 辨識引擎執行緒（lib.rs spawn_engine）───────┐
匯入音檔 ─▶ │ 唯一持有模型的執行緒，匯入和錄音都在這裡處理            │
            │  live.rs：重新取樣成 16 kHz → 寫 WAV（有電腦聲音時是   │
            │           立體聲：左麥克風、右電腦）→ 每秒依 gap 規則   │
            │           切出講完的段落 → asr.rs（Breeze + VAD）      │
            │           → 依兩軌音量標「我／對方」                    │
            │  draft.rs：把 16 kHz 音訊轉給 meetnote-draft（灰字）    │
            └──────┬─────────────────────────────┬─────────────────┘
                   │ db.rs（SQLite）             │ Tauri 事件
                   ▼                             ▼
            會議、逐字稿、搜尋索引、        前端（Vue）
            API 用量、聲紋、人員名單       summarizer.ts：identify_speakers → Claude Haiku
                   ▲
                   │ voice.rs（identify_speakers 指令，spawn_blocking）
            每句算 3D-Speaker 聲紋（sherpa-onnx）→ 存進 line_voices →
            比對手動標的句子 → 標上講者（相似度 < 0.4 留空）
```

- **錄音的回呼函式只做最輕量的事**；計時和音量由獨立的小執行緒每 0.1 秒回報，不受辨識影響。
- **辨識引擎是唯一持有模型的地方**。錄音中匯入的檔案會排隊，錄音結束後才處理。
- **Claude 在前端呼叫**：Rust 把會議狀態設成 `summarizing`，前端看到後先認講者，再呼叫 Haiku，把結果存回去。
- **講者有兩種來源**：`lines.speaker_score` 是 NULL 表示手動標的（也只有這些會拿來建聲紋），有值表示依聲音判斷的（值是相似度）。`me`、`others` 是錄電腦聲音時依音量判斷的。
- **App 結束時**（`shutdown`）：先完成錄音、讓引擎卸載模型，再結束。否則 whisper.cpp 會在結束時中止程式。

## 模組

### Rust（`src-tauri/src/`）

| 檔案 | 職責 |
|---|---|
| `lib.rs` | Tauri 指令、辨識引擎執行緒、選單列、App 結束收尾、API key 與詞彙表 |
| `asr.rs` | Breeze（whisper-rs）和 Silero VAD；gap 規則 `group`、錄音時的切段 `live_cut` 與 `LiveRule` |
| `live.rs` | 一場錄音：串流重新取樣、寫 WAV、保留還沒切的音訊、送出講完的段落 |
| `recorder.rs` | cpal 麥克風擷取、裝置清單、`Meter`、錄音錯誤分類 |
| `draft.rs` | 啟動並管理 `meetnote-draft` 輔助程式（僅 macOS） |
| `audio.rs` | 匯入音檔解碼（symphonia）與重新取樣（rubato）成 16 kHz 單聲道 |
| `db.rs` | SQLite：會議、逐字稿、FTS5 trigram 搜尋（trigger 自動同步）、API 用量 |
| `models.rs` | 模型檔的網址、大小、SHA-256，下載與驗證 |
| `voice.rs` | 分辨講者：每句算聲紋（sherpa-onnx + 3D-Speaker），用手動標過的句子當聲紋，認出其他句；門檻 `MATCH` |
| `system_audio.rs` | 啟動並管理 `meetnote-tap` 輔助程式，錄電腦播放的聲音（僅 macOS） |
| `pdf.rs` | 用 WKWebView 的列印功能把會議記錄直接存成 PDF（僅 macOS） |
| `../draft/main.swift` | `meetnote-draft`：SpeechAnalyzer 即時辨識。`build.rs` 用 swiftc 編譯，只在 macOS 打包 |

### TypeScript（`src/`）

| 檔案 | 職責 |
|---|---|
| `lib/api.ts` | 所有 Rust 指令與事件的型別包裝 |
| `lib/minutes.ts` | Haiku 提示、輸出 schema、套用修正、複製用的 Markdown（會議記錄、逐字稿）、費用估算（`PRICES`） |
| `lib/summarizer.ts` | 產生會議記錄的流程、錯誤訊息、記錄每次 API 用量 |
| `lib/draft.ts` | 灰字草稿的狀態：哪些草稿已經被 Breeze 取代 |
| `App.vue` | 主畫面、側邊欄、錄音控制、事件訂閱、載入畫面 |
| `components/MeetingView.vue` | 會議詳情：錄音中畫面、播放器、會議記錄、逐字稿、修正、本場詞彙 |
| `components/MinutesEditor.vue` | 會議記錄編輯器 |
| `components/SettingsView.vue` | API key、Claude 用量、錄音來源、麥克風、詞彙表、人員名單、模型下載 |
| `components/MeetingList.vue`、`SearchResults.vue`、`SplashScreen.vue` | 會議列表、搜尋結果、載入畫面 |

### 事件（Rust → 前端）

| 事件 | 時機 |
|---|---|
| `meeting-updated` | 會議狀態改變、匯入進度、錄音中多了新的逐字稿 |
| `recording-level` | 錄音中每 0.1 秒：音量和已錄時間 |
| `recording-draft` / `recording-covered` | 灰字草稿；Breeze 已經辨識到哪個時間點 |
| `recording-changed` / `recording-error` | 開始、停止、暫停；無法繼續錄音的錯誤 |
| `model-load` / `model-download` / `models-changed` | 模型載入狀態、下載進度、下載完成 |
| `tray-error` | 選單列操作失敗 |

## 資料位置

`~/Library/Application Support/com.twrory.meetnote/`（Windows 是 `%APPDATA%\com.twrory.meetnote\`）

| 檔案 | 內容 |
|---|---|
| `meetnote.db` | 會議、逐字稿（含講者）、會議記錄（JSON）、搜尋索引、API 用量、每句的聲紋（`line_voices`）、人員名單（`people`） |
| `audio/` | 錄音（WAV）和匯入音檔的複本 |
| `glossary.txt` | 全域詞彙表 |
| `models/` | 正式版下載的模型（開發版改用 `MEETNOTE_MODEL_DIR`） |
| `api_key_hint.txt` | 遮罩後的 API key，用來判斷有沒有設定，不含完整 key |
| `input_device.txt` | 選擇的麥克風 |
| `capture_system` | 錄音來源：有這個檔案表示同時錄電腦播放的聲音 |

## 測試

```bash
cd src-tauri && cargo test --lib       # Rust 單元測試
npx vitest run                         # TypeScript 單元測試
npx vue-tsc --noEmit                   # 型別檢查
```

`src-tauri/examples/` 裡的驗證工具：

| 工具 | 用途 |
|---|---|
| `asr_check` | 用 App 的匯入流程辨識一個音檔，輸出可以交給 `../spike/cer.py` 算 CER |
| `live_check` | 模擬錄音（48 kHz、每 10 毫秒一次），量測延遲和 CER；`LIVE_MIN_S`、`LIVE_PAUSE_S` 可調切段參數 |
| `mic_check` | 用 App 的程式實際錄 3 秒，確認麥克風和取樣率 |
| `download_check` | 下載小的 VAD 模型，驗證下載和 SHA-256 檢查 |
| `exit_check` | 重現「模型沒卸載就結束會中止」的問題 |
| `voice_check` | 用 App 的流程在**資料庫複本**上認講者，印出耗時和每人句數 |

用法和回歸測試流程見 [../AGENTS.md](../AGENTS.md)。
