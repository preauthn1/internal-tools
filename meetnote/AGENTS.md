# AGENTS.md

給在這個 repo 工作的 AI 工具（Claude Code、Codex、Cursor、Copilot…）的開發須知。專案總覽在 [README.md](README.md)，架構和模組在 [app/README.md](app/README.md)。

## 溝通

- 用**繁體中文**回覆。
- 說到「做好了」「沒問題」之前，要有實際執行的結果（測試、量測、截圖）。沒辦法驗證的部分（例如需要對著麥克風講話、需要點 App 視窗）要明說還沒驗證。

## Repo 結構

- `app/`：產品。Tauri v2 + Vue 3（TypeScript）+ Rust。
- `spike/`：早期實測程式和**回歸測試資料**（`samples/script-01`、`samples/script-02` 有正確稿件）。實測結果記在 `spike/FINDINGS.md`。
- `spike/` 裡可能有其他人進行中、尚未 commit 的修改（例如 `spike/Sources/main.swift`、`spike/meeting.py`），**不要順手 commit**。

## 環境

- macOS、Apple Silicon，**只有 Command Line Tools，沒有 Xcode**。正式發佈（簽章、公證）之前不要裝 Xcode，也不要依賴 Xcode 才有的東西（例如 Swift 的 `@Generable` 巨集、`metal` 編譯器）。
- Rust 用 Homebrew 的 rustup，路徑不在預設 PATH：`export PATH=/opt/homebrew/opt/rustup/bin:$PATH`。
- 模型檔在 `spike/models/`（不進 git）；開發版透過 `MEETNOTE_MODEL_DIR` 使用。
- Claude API key 在 repo 根目錄的 `.env`（不進 git）。**不要印出、貼出或 commit 它。**

## 執行

```bash
app/dev.sh            # 開發版：自動讀 .env 的 key、用 spike/models 的模型；改 Rust 會自動重新編譯
```

- 開發版**不讀鑰匙圈**（每次重新編譯簽章都會變，macOS 會一直詢問）。不要把開發流程改回讀鑰匙圈。
- 不要直接執行 `target/debug/meetnote`，那樣讀不到 key 和模型。
- 要錄**電腦播放的聲音**時用 `app/bundle-dev.sh`（打包成 .app 再開）。從 `dev.sh` 啟動時 macOS 只會把權限給終端機，錄到的是靜音。

## 測試

改完一定要跑：

```bash
cd app/src-tauri && cargo test --lib     # Rust 單元測試
cd app && npx vitest run && npx vue-tsc --noEmit
```

**改到辨識、切段、重新取樣或 Haiku 提示時**，還要跑回歸測試，並把 CER 和改之前比較：

```bash
cd app/src-tauri
cargo run --release --example asr_check -- ../../spike/samples/script-01.m4a ../../spike/models/ggml-breeze-asr-25-q8_0.bin ../../spike/models/ggml-silero-v5.1.2.bin > /tmp/out.txt
python3 ../../spike/cer.py ../../spike/samples/script-01.ref.txt /tmp/out.txt     # 目前約 2.5%
cargo run --release --example live_check -- <音檔> <模型> <VAD 模型> <輸出.wav>   # 模擬錄音，看延遲和 CER
```

**改到分辨講者（`voice.rs` 的 `MATCH`、聲紋模型、取音訊的方式）時**，用已標好講者的會議重新量測，結果寫進 `spike/FINDINGS.md`「從會議裡建聲紋再認人」：

```bash
cd spike/diar
cargo run --release --bin identify -- <會議的單聲道 wav> <聲紋模型> <手動標的句子.tsv> [整場逐字稿.tsv]   # 留一句測一句；GAP=60 排除前後 60 秒
cd app/src-tauri
cargo run --release --example voice_check -- <meetnote.db 的複本> <會議 id> <聲紋模型>   # App 的流程；不要直接用正式的資料庫
```

呼叫真實 Claude API 的測試要自己設定 key：`ANTHROPIC_API_KEY=… npx vitest run src/lib/minutes.live.test.ts`（會花一點錢）。

## 寫程式的原則

- **Rust 盡量少寫**：只做錄音、辨識、SQLite、讀寫檔案、鑰匙圈這類系統層的事。Claude 呼叫、會議記錄、介面邏輯放在 TypeScript。團隊不熟 Rust。
- 重要參數旁邊要寫**為什麼是這個值**和量測依據（參考 `asr.rs` 的 `LiveRule`、gap 規則常數）。
- 改了有量測依據的參數，要重新量測，結果寫進 `spike/FINDINGS.md`，必要時更新 `README.md` 的決策紀錄。
- Claude 模型或價格改了，要一起更新 `app/src/lib/minutes.ts` 的 `MODEL` 和 `PRICES`。
- 不要在錄音的回呼函式（`recorder.rs`）裡做會阻塞的事，它只負責轉成單聲道和交出資料。
- 長時間的工作不要放在 Tauri 的同步指令裡（會卡住主執行緒）；用 `async` 指令加 `spawn_blocking`。

## Git

- **commit 前先問使用者**，除非使用者已經說要 commit。
- `git add` 要**明確列出路徑**，不要用 `git add -A` 或 `git add .`，避免把別人進行中的修改一起 commit。
- 不 commit：`.env`、`spike/models/`、`spike/recordings/`、`spike/out/`、`app/src-tauri/binaries/`（建置時產生）。
- commit 訊息寫清楚改了什麼、為什麼，有量測數據就附上。
