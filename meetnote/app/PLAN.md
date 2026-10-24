# meetnote（會議筆記）規劃

會議錄音（實體或線上）或匯入音檔，產生繁體中文（中英夾雜）逐字稿與會議記錄。
各項選擇的實測依據在 [../spike/FINDINGS.md](../spike/FINDINGS.md)。

- 識別碼：`meetnote`，bundle id `com.twrory.meetnote`（原本是 `com.target.meetnote`，2026-10-20 改；沒有寫搬資料的程式；鑰匙圈名稱 `KEYRING_SERVICE` 也跟著改，舊的 API key 要在設定頁重新輸入）
- 平台：macOS 先做，Windows 在 M4
- 使用對象：先自己用，之後發給團隊內部

## 辨識與 AI

- **辨識**：在本機用 Breeze ASR 25 q8_0（1.66 GB），透過 whisper-rs 執行，Metal 加速。模型在第一次開啟 App 時下載。
- **邊錄邊轉（gap 規則）**：Rust 端每秒對還沒切的錄音跑一次 Silero VAD（`WhisperVadContext`）。停頓超過 2 秒、累積 15 秒後遇到短暫停頓、或段落長度達到 60 秒，就送出一段，停頓的部分丟掉。段落開頭保留 1 秒、結尾保留 0.3 秒，而且開頭不能早於上一段的結尾。送進 whisper 時**不開 VAD**。每一句的時間戳 = 段落在整場錄音裡的起始時間 + whisper 回傳的時間。
- **灰字即時草稿（只有 macOS）**：錄音時把 16 kHz 音訊同時交給 `meetnote-draft`（SpeechAnalyzer），即時顯示灰字；Breeze 辨識完的段落會取代對應的灰字。草稿不會被儲存。
- **詞彙表**：一份全域清單，加上每場會議臨時補的詞。給 whisper 的 initial prompt 限制在約 200 個 token 以內（上限是 224，超過的部分會被截掉）；給 Haiku 的是完整清單。
- **會議記錄**：用 Claude Haiku 5.5（`claude-haiku-5-5`），透過 TypeScript 的 `@anthropic-ai/sdk` 呼叫，用 JSON schema 結構化輸出。停止錄音後自動執行，也有「重新產生」按鈕。
  - Haiku 只回傳會議記錄和**修正清單**（以行號標示 原文 → 修正，附判斷依據），不重寫全文。修正由 App 在本機套用，每一處都可以接受或還原。
  - 只有逐字稿的文字會上傳，音檔不離開本機。
- **API key**：使用者自己提供，存在系統鑰匙圈（macOS 的 Keychain、Windows 的 Credential Manager）。只在產生會議記錄時讀取，每次開啟 App 最多讀一次。
- **離線或呼叫失敗**：保留原始逐字稿，顯示失敗原因，使用者按「重新產生」重試。App 重新開啟時，會自動重跑上次沒完成的會議記錄。不使用 AFM。
- **用量紀錄**：每次呼叫 Claude 都記下 token 數和估算費用（含失敗的呼叫），在會議記錄下方和設定頁顯示。
- **線上會議（只有 macOS）**：`meetnote-tap`（Core Audio process tap）錄電腦播放的聲音，和麥克風存成立體聲 WAV，依兩軌音量標「我／對方」。設定頁可以選「只有麥克風」。
- **分辨講者**：使用者在逐字稿幫每個人標幾句；每句用 3D-Speaker（sherpa-onnx）算聲紋，存在 SQLite，其他句標成聲音最像的人（相似度低於 0.4 留空）。只有手動標的句子會拿來建聲紋。產生會議記錄前先跑一次，講者名字、職位會交給 Haiku。

## 會議記錄格式（第一版只有這一種，可以編輯）

1. 會議資訊：日期、時間長度、主題（Haiku 推斷）、與會者（使用者填寫或 Haiku 推斷）
2. 摘要：3 到 5 句
3. 討論重點：主題 → 要點
4. 決議
5. 待辦事項：內容、負責人、期限（沒提到就留空）
6. 有時間戳的完整逐字稿

## 畫面

- **會議列表**：可以全文搜尋（SQLite FTS5）
- **錄音中**：計時、音量、每段辨識完就更新逐字稿；可以選麥克風、暫停或繼續
- **會議詳情**：「會議記錄」和「逐字稿」兩個分頁加播放器。點時間戳就從那裡播放；修正可以接受或還原，修正過的詞可以一鍵加進詞彙表
- **設定**：API key、Claude 用量、錄音來源、錄音裝置、詞彙表、人員名單（職位、聲紋）、模型下載
- 選單列（macOS）或系統匣（Windows）圖示，可以開始和停止錄音

## 資料

都存在 App 的資料夾（`~/Library/Application Support/com.twrory.meetnote/`）：SQLite（會議、逐字稿、會議記錄、全文索引、API 用量）、錄音（WAV）和匯入音檔的複本、詞彙表。PDF 只在使用者匯出時寫到選擇的位置。

## 技術

- Tauri v2 搭配 Vue 3（TypeScript）
- **Rust 盡量少寫**：只負責錄音（cpal）、辨識（whisper-rs + VAD）、SQLite、讀寫檔案、存取鑰匙圈
- **TypeScript 負責**：呼叫 Haiku、組裝會議記錄、套用修正清單、介面
- **發佈前不裝 Xcode**：用 Command Line Tools 就能開發和編譯。麥克風權限寫在 `src-tauri/Info.plist`（`NSMicrophoneUsageDescription`）
- 開發時要注意：`tauri dev` 的麥克風權限會算在終端機或 VS Code 身上；重新簽章後權限可能要重新授予（`tccutil reset Microphone com.twrory.meetnote`）
- 開發版用 `app/dev.sh` 啟動：模型取自 `spike/models`，API key 取自 `.env`，不碰鑰匙圈（每次重新編譯簽章都會變，macOS 會一直詢問鑰匙圈權限）

## 第一版不做

範本、全域快捷鍵、發佈到 Outline、GitLab 或 Telegram、AFM。

## 開發順序

- [x] **切段規則驗證**：`spike/segtest.py`，見 FINDINGS「App 的切段規則」
- [x] **M1**：匯入音檔 → 逐字稿 → 會議記錄 → 會議列表（只做 Mac）；模型下載在 M3 補上
- [x] **M2**：麥克風錄音、邊錄邊轉、選單列圖示
- [x] **M3**：編輯、點時間戳播放、接受或還原修正、詞彙表、全文搜尋、在 App 裡下載模型
- [x] **M3 之後**：錄線上會議的電腦聲音、存成 PDF、複製 Markdown、分辨講者與人員名單
- [ ] **M4**：Windows 版（要先確認 sherpa-onnx 的靜態 C 執行環境和 whisper.cpp 能不能連結在一起）
- [ ] **發佈**：給團隊內部。屆時裝 Xcode，申請 Apple Developer 帳號做簽章和公證；Windows 版也要簽章

## 已知限制

- **錄音存成 WAV**：1 小時約 115 MB，還沒有壓縮。
- **App 在錄音中當掉時**，WAV 檔頭沒有寫完整，那段錄音可能無法播放；已辨識的逐字稿會保留。正常結束（⌘Q、關視窗、選單列「結束」）不會有這個問題。
- **選單列圖示**是彩色的 App 圖示，macOS 慣例是單色的 template 圖示，發佈前要換。
- **正式版第一次開啟**：之前存在鑰匙圈的 key 沒有提示檔，會顯示「尚未設定」，要重新輸入一次。
- **音檔保存天數**還沒做，目前一律永久保存。
- **舊的 Claude 呼叫**（用量紀錄功能加入之前）沒有算進設定頁的累計。
- **Windows 版的灰字草稿**做法還沒決定：有 GPU 的電腦用 Breeze 3 秒 / 0.7 秒，沒有 GPU 的用小型 whisper 模型，要在 Windows 上量測後再定。

## 回歸測試

`spike/` 保留下來：`samples/script-01`、`samples/script-02`（照稿念，有正確答案）、`cer.py`、`segtest.py`，以及 `src-tauri/examples/` 的 `asr_check`、`live_check`。之後換模型、改切段規則或改 prompt 時重跑一次，確認字錯率沒有變差。做法見 [../AGENTS.md](../AGENTS.md)。
