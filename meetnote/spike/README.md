# Spike：Mac 本機錄音轉逐字稿

驗證 SpeechAnalyzer（轉錄）+ Apple Foundation Models（整理、摘要）在 Mac 上全本機可行，並與 whisper.cpp 對照。
只需 Command Line Tools，不需 Xcode。結論寫在 [FINDINGS.md](FINDINGS.md)。

## 準備

```bash
./build.sh                       # 編出 ./transcribe
swiftc -O -swift-version 5 -parse-as-library -o afm-minutes AFM/main.swift   # AFM 會議記錄
brew install whisper-cpp         # Breeze ASR 25 透過 whisper.cpp 執行
python3 -m venv .venv && .venv/bin/pip install anthropic   # Haiku；ANTHROPIC_API_KEY 放在 ../.env
curl -L -o models/ggml-breeze-asr-25-q8_0.bin https://example.invalid/shdennlin/breeze-asr-25-ggml/resolve/main/ggml-breeze-asr-25-q8_0.bin
curl -L -o models/ggml-silero-v5.1.2.bin https://example.invalid/ggml-org/whisper-vad/resolve/main/ggml-silero-v5.1.2.bin
curl -L -o models/ggml-large-v3-turbo.bin https://example.invalid/ggerganov/whisper.cpp/resolve/main/ggml-large-v3-turbo.bin   # 對照組
```

## 用法

```bash
.venv/bin/python meeting.py         # 錄會議：邊錄邊轉，Ctrl-C 停止後約 15 秒產生逐字稿與會議記錄
./minutes.sh recordings/x.m4a       # 現有音檔 → Breeze 逐字稿 → Haiku 會議記錄
.venv/bin/python afm.py recordings/x.breeze.md   # 改用本機 AFM 產生會議記錄（離線，需先 swiftc 編出 afm-minutes）
./transcribe --mic               # 錄音，邊錄邊出字，Ctrl-C 停止 → recordings/<時間>.m4a + .md
./transcribe samples/x.m4a       # 轉檔 → out/x.speechanalyzer.txt, x.cleaned.md, x.summary.md
./transcribe ai out/x.whisper.txt  # 只重跑 AFM 整理與摘要
./compare.sh                     # samples/ 全部跑 SpeechAnalyzer 與 whisper.cpp，耗時記在 out/timings.tsv
python3 cer.py samples/x.ref.txt out/x.speechanalyzer.txt out/x.whisper.txt   # 字錯率
```

共用參數：`--no-ai` 跳過 AFM、`--out DIR` 輸出目錄、`--locale zh_TW`。

第一次跑 `--mic` 時，macOS 會要求你的終端機 App 取得麥克風權限。
