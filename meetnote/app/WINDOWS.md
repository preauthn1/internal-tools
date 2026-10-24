# Windows 版（M4）交接

在 Windows 電腦上接續開發 meetnote 時從這裡開始。專案總覽在 [../README.md](../README.md)，架構在 [README.md](README.md)，開發規則在 [../AGENTS.md](../AGENTS.md)（AGENTS.md 裡的指令是 macOS 版，Windows 的對應寫在下面）。

## 目前狀態（2026-10-19，commit `e0ed5b7`）

Mac 版已完成：匯入音檔、麥克風錄音邊錄邊轉、線上會議（錄電腦聲音）、即時灰字草稿、Haiku 會議記錄、編輯、點時間播放、詞彙表、搜尋、分辨講者（手動標幾句再依聲音認人）、人員名單、匯出 PDF、在 App 裡下載模型。

**2026-10-19 到 20 在 Windows 11 上編譯並量測過一輪**（連結、單元測試、CPU／Vulkan／CUDA 的速度與 CER、沒有 Vulkan 裝置時的行為），數字和做法見 [../spike/FINDINGS.md](../spike/FINDINGS.md) 的「Windows 版」。第 3 節每一項後面標了目前進度；沒標的就是還沒做。`dev.ps1` 已實際執行過，App 能啟動、匯入音檔走 GPU 辨識（見 FINDINGS）。**沒用麥克風錄過音**（開發用的這台沒有麥克風）。

## 1. 準備開發環境

| 項目 | 說明 |
|---|---|
| Git | 用 SSH 存取 GitLab：`git clone ssh://user@example.invalid:1022/Rory/meetnote.git`（網頁 http://example.invalid:1688/Rory/meetnote ，private），Windows 上要另外設定 SSH key |
| Node.js 22 | Mac 上是 v22.18.0 |
| Rust | rustup，用 `x86_64-pc-windows-msvc` toolchain |
| Visual Studio Build Tools | 勾「使用 C++ 的桌面開發」（MSVC + Windows SDK）。Tauri 和 whisper.cpp 都需要 |
| CMake | `whisper-rs` 會用 CMake 編譯 whisper.cpp |
| LLVM（libclang） | `whisper-rs-sys` 用 bindgen 產生綁定，需要 libclang；裝完設定 `LIBCLANG_PATH` |
| WebView2 | Windows 11 內建；Windows 10 要確認有裝 |
| Python 3 | 只有跑 `spike/cer.py`（算字錯率）時需要 |
| GPU 相關（看第 3 節的決定） | NVIDIA：CUDA Toolkit；AMD／Intel 顯示卡：Vulkan SDK。目前 `Cargo.toml` 的 Windows 設定是 Vulkan |

實際編譯過的版本：MSVC 14.44、Rust 1.99.1、CMake 4.4.4、LLVM 23.1.3、Node 22.23.2、Vulkan SDK 0.0.0.0、CUDA 13.4。注意事項：

- `LIBCLANG_PATH` 要設成 LLVM 的 `bin` 資料夾（例如 `C:\Program Files\LLVM\bin`）。
- 編 Vulkan 版要設 `CARGO_TARGET_DIR` 到很短的路徑（例如 `C:\t`），預設的 `app\src-tauri\target` 編不過（原因見 FINDINGS）。`app\dev.ps1` 已處理（`tauri dev` 有認得 `CARGO_TARGET_DIR`，輸出在 `C:\t\debug`）。
- `.cargo/config.toml` 的設定（`+crt-static` 和 CMake 變數）是必要的，沒有它們連結會失敗。改了那兩個 CMake 變數要 `cargo clean -p whisper-rs-sys`。
- 裝完 Vulkan SDK 或 CUDA 後，**已經開著的終端機要關掉重開**，否則讀不到 `VULKAN_SDK`、`CUDA_PATH_V13_4`。
- 非互動的 shell 找不到 `winget`，要用 `%LOCALAPPDATA%\Microsoft\WindowsApps\winget.exe` 的完整路徑。
- 用 PowerShell 5.1 的 `>` 轉導辨識結果會加 BOM，`cer.py` 會漏算第一行；改用 `Start-Process -RedirectStandardOutput`。

## 2. 不在 git 裡、要另外帶過去的東西

| 東西 | 在 Mac 上的位置 | Windows 上怎麼處理 |
|---|---|---|
| Claude API key | repo 根目錄的 `.env`（`ANTHROPIC_API_KEY=[REDACTED] | 在 Windows 的 repo 根目錄建立同樣的 `.env`。**用安全的方式傳，不要貼進聊天、不要 commit** |
| 模型檔 | `spike/models/` | 照下表放到 Windows 的 `spike/models/`，或在 App 裡下載（設定頁） |
| 會議資料 | `~/Library/Application Support/com.target.meetnote/meetnote.db` 和錄音（**這是改識別碼為 `com.twrory.meetnote` 之前的舊資料夾**；新版 App 讀的是 `com.twrory.meetnote`，舊資料要自己複製過去） | 開發 Windows 版不需要。要拿真實會議測試的話，再另外複製資料庫和音檔 |

App 需要的模型（`app/src-tauri/src/models.rs` 也有下載網址）：

| 檔案 | 大小 | 用途 |
|---|---|---|
| `ggml-breeze-asr-25-q8_0.bin` | 1.66 GB | 辨識（Breeze ASR 25） |
| `ggml-silero-v5.1.2.bin` | 0.9 MB | VAD（切段） |
| `3dspeaker_speech_eres2net_base_sv_zh-cn_3dspeaker_16k.onnx` | 40 MB | 聲紋（分辨講者） |

`spike/models/` 裡另外兩個（`ggml-breeze-asr-25-f16.bin`、`ggml-large-v3-turbo.bin`）只用於早期實測，App 不需要。

回歸測試用的錄音和正確稿件已經在 git 裡：`spike/samples/script-01.m4a`、`script-01.ref.txt`、`script-02.ref.txt`。

## 3. M4 要做的事

依建議順序。第 1 項沒過的話後面都做不了，要最先確認。

1. **確認 sherpa-onnx 和 whisper.cpp 能連結在一起**（PLAN 列為 M4 的前提）。sherpa-onnx 是靜態連結，用 MT（靜態 C 執行環境），whisper.cpp 通常用 MD，兩者混用可能連結失敗。先 `cargo build` 看看；失敗的話選一個統一（例如設 `RUSTFLAGS=-C target-feature=+crt-static` 全部改 MT，或找 sherpa-onnx 的 MD 版）。
   **✅ 已解決**：真的撞到 MT／MD 衝突（`LNK2038`）。用 `app/src-tauri/.cargo/config.toml` 設 `+crt-static`，debug 建置另外要兩個 CMake 變數，release 和 debug 都能連結，`cargo test --lib` 34 項全過。細節和 debug 為什麼還要額外設定（原因不明）見 FINDINGS。**沒在 Mac 上驗證這些設定。**
2. **`whisper-rs` 的 `metal` feature 要改成只限 macOS**。`app/src-tauri/Cargo.toml` 現在是 `whisper-rs = { version = "0.16", features = ["metal"] }`，不分平台。改成依 target 分開設定，Windows 選 `cuda`、`vulkan` 或不開 GPU。
   **✅ 已做**：`Cargo.toml` 現在 macOS 用 `metal`、Windows 用 `vulkan`，其他平台不開 GPU。CUDA 只做過實驗（暫時換 feature 測完改回），沒有進 repo。
3. **量測 Windows 上的辨識速度**，決定 GPU 方案。邊錄邊轉需要至少 1 倍即時、最好 2 倍以上（Mac M4 約 6.5 倍）。用 `examples/asr_check` 和 `examples/live_check` 量（指令見第 5 節），把結果寫進 `spike/FINDINGS.md`。
   **✅ 已量一輪（GTX 1650，各跑一次）**：CPU 0.16 倍、Vulkan 2.5 到 3.2 倍（單次變動大，重跑過一次）、CUDA 5.2 倍即時，CER 2.5%／2.5%／2.9%。**GPU 方案（Vulkan、CUDA 或兩者並存）還沒決定**；目標筆電（Intel 內顯）量過 3 次約 1.5 倍即時（最慢一步約 14.5 秒），見 FINDINGS。
4. **灰字草稿**（`src/draft.rs`）：macOS 用 SpeechAnalyzer 的 Swift 程式，其他平台 `Draft::start` 直接回傳 `None`，所以 Windows 現在沒有草稿，但不會壞。PLAN 的方向：有 GPU 用 Breeze 每 3 秒 / 0.7 秒跑一次，沒有 GPU 用小型 whisper 模型，依第 3 項的量測決定。
5. **線上會議錄電腦聲音**（`src/system_audio.rs`）：macOS 用 Core Audio process tap 的 Swift 程式，其他平台 `available()` 回傳 `false`，所以 Windows 現在不會顯示這個選項。Windows 可以用 WASAPI loopback；`cpal` 支援 loopback（Cargo.toml 已經有 `cpal`）。錄到之後要和麥克風一起存成立體聲，接到現有的「我／對方」標註流程。
6. **匯出 PDF**（`src/pdf.rs`）：只有 macOS 實作（WKWebView 列印），其他平台會回傳錯誤「這個平台還不支援直接存成 PDF」。WebView2 有 `PrintToPdf` 可以用。
7. **API key 存放**：`keyring` crate（4.2.0）在 Windows 對應 Credential Manager，要確認不需要另外開 feature 就能編譯和讀寫。
8. **開發用的啟動腳本**：`app/dev.sh` 是 zsh。Windows 要做一個對應的 `app/dev.ps1`：設定 `MEETNOTE_MODEL_DIR` 指向 `spike\models`、從 `..\.env` 讀 `ANTHROPIC_API_KEY`，然後 `npm run tauri dev`。開發版不讀系統鑰匙圈這條規則（見 AGENTS.md）也照用。
   **✅ 已寫好 `app/dev.ps1` 並實際執行過**：第一次完整 debug 編譯約數分鐘，之後 App 視窗正常開啟。匯入後 Claude 會議記錄也成功產生（key 在 App 設定頁輸入，不需要 `.env`；見 FINDINGS）。
9. **系統匣圖示**：`TrayIconBuilder` 在 Windows 會變成系統匣圖示，確認選單、開始／停止錄音都正常。
   **部分驗證**：系統匣圖示有出現，右鍵選單有「開始錄音」「打開會議筆記」「結束」。「開始錄音」在沒有麥克風時顯示找不到麥克風、App 沒當掉；「打開會議筆記」視窗有被叫到前景；「結束」的結束碼是 0、沒有殘留行程。**沒測**：左鍵行為、有麥克風時的開始／停止錄音、錄音或匯入進行中按「結束」。
10. **打包和簽章**：Tauri 產生 MSI 或 NSIS 安裝檔；正式發佈要簽章（PLAN 的「發佈」）。
   **已打出 NSIS 安裝檔和免安裝 zip**（`app/src-tauri/tauri.windows.conf.json` 只做 nsis；安裝檔 13.77 MiB、`meetnote.exe` 85.27 MB，都未簽章），單檔 exe 可穩定開著。**安裝檔本身還沒在任何機器上安裝測試過，簽章沒做**；MSI 沒做。

不用處理的：`build.rs` 只在 macOS 編譯 Swift 輔助程式，`tauri.macos.conf.json` 的 `externalBin` 也只在 macOS 生效，Windows 編譯不會碰到 `swiftc`。

## 4. 從 Mac 帶過去、還沒驗證的事

這些和 Windows 無關，但還沒做，別忘了：

- Haiku 提示詞新增了「負責人要對應人員名單」的規則，還沒用真的 API 跑過 `src/lib/minutes.live.test.ts`（會花一點錢）。
- 加入 sherpa-onnx 之後，Mac 的打包版（`bundle-dev.sh`、開了 LTO 的 release 建置）還沒重新驗證。
- 分辨講者：換會議室或電腦位置後還認不認得出來，還沒驗證；標註講者、人員名單的畫面還沒人實際操作過。
- Mac 資料庫第 10 場會議的講者標註曾經直接用 SQL 清除，前提是「那場是實體會議」，還沒和使用者確認。

### Windows 出貨前要驗證（2026-10-20 整理，都還沒做）

- **缺 `vulkan-1.dll` 的機器**：Vulkan 版的 `vulkan-1.dll` 是載入時相依，缺檔時程式應該連啟動都不能（從 `dumpbin` 匯入表推論，沒實測）。要在乾淨的 VM 或 Windows Sandbox 驗證；解法候選是安裝檔附上 `vulkan-1.dll` 或改成延遲載入。
- **沒有 Vulkan 裝置但載入器存在**：模擬過，會退回 CPU（很慢），沒當掉；沒等它跑完。舊的 Intel／AMD 驅動會不會在建立 Vulkan 實例時當掉，沒測。
- **CUDA 版的散佈**：CUDA 版需要 `cublas64_13.dll` 等 CUDA Toolkit 的檔案，驅動不附；體積大、只對 NVIDIA 有用。
- **目標機型的速度**：已量（Acer SFL14-54M，Intel 內顯，約 1.5 倍即時，見 FINDINGS）；其他機型沒量。
- **實際使用路徑**：麥克風（WASAPI）錄音、`voice_check`、`exit_check`、`download_check`、1 小時以上錄音的記憶體。（`dev.ps1` 啟動 App 和匯入音檔已驗證。）
- **Mac 回歸**：✅ 2026-10-20 已重跑，`cargo test --lib` 34 項全過、vitest 和 `vue-tsc` 通過、`asr_check` script-01 CER 2.5%（見 FINDINGS「Mac 回歸」）。沒測 script-02、`live_check`、打包版和 LTO release 建置。
- **安裝檔實測**：SmartScreen、安裝、首次啟動、App 內下載模型（約 1.7 GB）在 Windows 上都沒測過（`meetnote.exe` 單檔執行有測）。

## 5. 測試指令（PowerShell）

```powershell
cd app\src-tauri; cargo test --lib
cd app; npx vitest run; npx vue-tsc --noEmit
```

改到辨識、切段、重新取樣或 Haiku 提示時的回歸測試（Mac 上 script-01 約 2.5%）：

```powershell
cd app\src-tauri
cargo run --release --example asr_check -- ..\..\spike\samples\script-01.m4a ..\..\spike\models\ggml-breeze-asr-25-q8_0.bin ..\..\spike\models\ggml-silero-v5.1.2.bin > $env:TEMP\out.txt
python ..\..\spike\cer.py ..\..\spike\samples\script-01.ref.txt $env:TEMP\out.txt
```

## 6. 給 Windows 上的 AI 工具

Windows 上用 Claude Code 等工具接續時，先讀 `AGENTS.md` 和這份文件。AGENTS.md 的「環境」一節寫的是 Mac（Homebrew rustup、只有 Command Line Tools），在 Windows 上以這份文件第 1 節為準；其餘規則（繁體中文、說做好之前要有實測結果、commit 前先問、`git add` 要列出路徑、不 commit `.env` 和模型檔）照樣適用。
