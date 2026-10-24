# meetnote

> 本目录原始 README 见文末，本节为快照补充信息。

原始路径：`meetnote`　|　文件 **58**　|　主语言：TypeScript 10, JSON 9, Python 8

## 目录构成

```
app/src/  — 16
app/  — 12
spike/  — 11
app/src-tauri/  — 8
spike/diar/  — 4
spike/samples/  — 2
spike/Sources/  — 1
spike/AFM/  — 1
app/.vscode/  — 1
```

## 关键文件

- CI/构建：`spike/build.sh`
- 前端入口：`app/index.html`、`app/src/main.ts`

---

## 原始说明

# meetnote（會議筆記）

在 Mac 上錄會議或匯入錄音，**在本機**轉成繁體中文（可中英夾雜）逐字稿，再用 Claude 產生會議記錄。錄音不會離開電腦，只有逐字稿文字會送到 Claude。

| 項目 | 值 |
|---|---|
| App 名稱 | `meetnote`（Finder、Dock、權限提示顯示這個）；視窗標題「會議筆記」 |
| bundle id | `com.twrory.meetnote`（資料夾位置、系統權限都綁這個，不要隨便改；改過一次，舊資料在 `com.target.meetnote` 資料夾，見 [app/WINDOWS.md](app/WINDOWS.md)） |
| 資料位置 | `~/Library/Application Support/com.twrory.meetnote/`：SQLite（會議、逐字稿、聲紋）、錄音、詞彙表、模型 |
| 平台 | macOS（Apple Silicon）；Windows 版開發中 |

## 功能

- **錄音**：實體會議用麥克風；線上會議（Google Meet、Slack…）可以同時錄電腦播放的聲音。邊錄邊轉，灰字即時草稿。
- **匯入音檔**：m4a、mp3、wav 等。
- **會議記錄**：摘要、討論重點、決議、待辦事項（負責人、期限）。Claude 修正的錯字都標出來，可以逐處還原。
- **分辨講者**：在逐字稿幫每個人標幾句，其他句會依聲音自動認出來；聲紋跨會議共用。設定頁的人員名單可以填職位，Claude 用來判斷待辦事項的負責人。
- **整理和分享**：全文搜尋、點時間戳播放、編輯逐字稿和會議記錄、詞彙表；複製會議記錄或逐字稿（Markdown），會議記錄可以存成 PDF。
- **用量**：每次呼叫 Claude 的 token 和估算費用。

## 文件

| 文件 | 內容 |
|---|---|
| [app/PLAN.md](app/PLAN.md) | 設計規劃、**目前進度**、已知限制 |
| [app/README.md](app/README.md) | 開發指南：架構、模組、執行、測試 |
| [spike/FINDINGS.md](spike/FINDINGS.md) | 所有實測紀錄（按時間順序） |
| [spike/README.md](spike/README.md) | 實測與回歸測試工具的用法 |
| [AGENTS.md](AGENTS.md) | 給 AI 工具的開發須知 |

## 開啟

```bash
app/dev.sh            # 開發版：改 Vue 即時更新，改 Rust 自動重新編譯
app/bundle-dev.sh     # 打包成 .app 再開：要錄電腦播放的聲音時用這個
```

第一次會先編譯（1 到 2 分鐘），開啟後顯示載入畫面，模型載入完成才進主畫面。錄電腦聲音的權限只會給 .app，從 `dev.sh` 啟動錄到的會是靜音。詳見 [app/README.md](app/README.md)。

## 怎麼運作

```text
麥克風（＋電腦播放的聲音）/ 匯入的音檔
  → Silero VAD 找出有人說話的片段，依 gap 規則切段
  → Breeze ASR 25（本機、GPU）辨識成正式逐字稿
     （錄音時，macOS SpeechAnalyzer 先即時顯示灰字草稿）
  → 3D-Speaker 聲紋（本機）：依手動標過的句子認出每句是誰講的
  → Claude Haiku 5.5 回傳修正清單和會議記錄（只上傳文字和講者名字、職位）
  → SQLite 儲存
```

## 決策紀錄

每個選擇都有實測依據，細節在 [spike/FINDINGS.md](spike/FINDINGS.md) 對應的章節。要推翻之前，先看那一節的數據。

### 辨識

| 決策 | 理由 | 依據（FINDINGS 章節） |
|---|---|---|
| 正式逐字稿用 **Breeze ASR 25 q8_0** | 照稿念的測試 CER 3.3%（q8_0 為 3.5%），whisper large-v3-turbo 9.8%、SpeechAnalyzer 12.9%；英文術語和台灣用語明顯較好。q8_0 比 f16 小一半（1.66 GB），準確度幾乎一樣 | 照稿念的真人錄音；會議記錄流程加速 |
| 切段用 **gap 規則**（停頓超過 2 秒就切並丟掉停頓、開頭留 1 秒、結尾留 0.3 秒），whisper 不開 VAD | 送給 whisper 的音訊夾著長時間雜音會產生大量幻覺；whisper-server 開 VAD 遇到整段無聲會當掉；開頭要留 1 秒才不會切掉「大家好」 | App 的切段規則 |
| 錄音時正式版用 **15 秒 / 2 秒**切段，即時顯示交給 **SpeechAnalyzer 灰字草稿** | 改成 3 秒 / 0.7 秒能把延遲從約 10 秒降到約 4 秒，但 CER 從 3.9% 升到 4.9%（段落邊界的字會錯）。用草稿補足即時性，正式版保留完整上下文 | `asr.rs` 的 `LiveRule` 註解；`examples/live_check.rs` |

### 會議記錄

| 決策 | 理由 | 依據（FINDINGS 章節） |
|---|---|---|
| 會議記錄用 **Claude Haiku 5.5**，不用 Apple Foundation Models | AFM 修不了辨識錯字，還會把錯字寫成摘要裡的「事實」；Haiku 能依上下文修正並標出不確定的片段。2 分鐘的會議約 US$0.002，25 分鐘約 US$0.007 | Claude Haiku 5.5 取代 AFM |
| Haiku **只回傳行號修正清單**，不重寫全文 | 25 分鐘的錄音：145 秒 → 51 秒，輸出 40K → 11K tokens；原始逐字稿和時間戳都保留，每處修正可以還原 | 會議記錄流程加速 |

### 講者與錄音來源

| 決策 | 理由 | 依據（FINDINGS 章節） |
|---|---|---|
| 線上會議用 **Core Audio process tap** 錄電腦聲音，存成立體聲（左麥克風、右電腦聲音） | 不用裝虛擬音訊驅動程式；兩軌分開，可以依音量分出「我／對方」。電腦聲音那一軌一直沒有聲音時（實體會議）就不標 | `live.rs` 的 `attribute` |
| 分辨講者用**手動標幾句建聲紋再認人**，不用自動分群 | 真實的 96 分鐘、7 人會議：自動分群完全失效（分出 361 群）；每人標 3 到 12 句後認人，留一句測一句 100%、排除前後 60 秒 88%，低於相似度 0.4 留空。只用手動標的句子建聲紋，認錯的不會回頭污染聲紋 | 分辨講者；從會議裡建聲紋再認人 |
| 聲紋用 **3D-Speaker（中文）**，透過 sherpa-onnx | 真人語音上比英文的 TitaNet 準（100% 對 79%）；sherpa-onnx 靜態連結，不用另外附 DLL；96 分鐘的會議算聲紋約 2 分鐘 | 分辨講者 |

### 技術

| 決策 | 理由 | 依據（FINDINGS 章節） |
|---|---|---|
| **Tauri + Rust**，不換成 Go（Wails） | Tauri v2 已穩定（選單列、附帶輔助程式、簽章都有支援），Wails v3 還在 beta；whisper.cpp 是 C++，換 Go 也解決不了跨平台編譯的難點 | 2026-10-16 的討論 |
| **發佈前不裝 Xcode** | Command Line Tools 就能編譯 Tauri、whisper.cpp（Metal 著色器在執行時編譯）和 Swift 輔助程式；簽章和公證時才需要 | 可行性 |
| 開發版的 API key 從 `.env` 讀，**不碰鑰匙圈** | 開發版沒有固定簽章，每次重新編譯 macOS 都會再問一次鑰匙圈權限 | `lib.rs` 的 `dev_api_key` |
| PDF 用 **WKWebView 的列印功能**直接存檔 | 不用另外帶 PDF 函式庫，版面沿用網頁的列印樣式（`@media print`），文字可以選取 | `pdf.rs` |

## 進度

以 [app/PLAN.md](app/PLAN.md) 的「開發順序」為準：Mac 版已完成匯入、錄音（含線上會議）、即時草稿、編輯、搜尋、分辨講者，下一步是 M4 Windows 版，交接和待辦見 [app/WINDOWS.md](app/WINDOWS.md)。

