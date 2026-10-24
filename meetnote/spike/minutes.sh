#!/bin/zsh
# Audio → Breeze ASR 25 transcript → Claude Haiku meeting minutes.
#   ./minutes.sh recordings/2026-10-15_0100.m4a
# Writes <name>.breeze.md (timestamped transcript) and <name>.breeze.haiku.md (cleaned text + minutes) next to the audio.
set -e
cd "${0:A:h}"
audio=$1
[[ -f $audio ]] || { echo "用法: ./minutes.sh <音檔>"; exit 1; }
base=${audio:r}
wav=out/tmp/${audio:t:r}.wav
mkdir -p out/tmp
afconvert -f WAVE -d LEI16@16000 -c 1 "$audio" "$wav"

words=$(grep -v '^#' glossary.txt | grep -v '^$' | paste -sd '、' -)
prompt="以下是繁體中文的會議逐字稿，內容可能夾雜英文術語，例如：$words。"
echo "Breeze ASR 25 辨識中…"
whisper-cli -m models/ggml-breeze-asr-25-q8_0.bin -f "$wav" -l zh --prompt "$prompt" --vad -vm models/ggml-silero-v5.1.2.bin -vt 0.35 -vp 300 -np 2>/dev/null | python3 -c '
import re, sys
for l in sys.stdin:
    m = re.match(r"\[(\d+):(\d+):(\d+)\.\d+ --> [^\]]+\]\s*(.*)", l)
    if m and m.group(4).strip():
        h, mi, s = map(int, m.group(1, 2, 3))
        print(f"[{h*60+mi:02d}:{s:02d}] {m.group(4).strip()}")' > "$wav.txt"

duration=$(afinfo "$audio" | awk '/estimated duration/ {printf "%02d:%02d", $3/60, $3%60}')
{
  echo "# 會議逐字稿 ${audio:t:r}"
  echo
  echo "- 長度：$duration"
  echo "- 音檔：[${audio:t}](${audio:t})"
  echo "- 引擎：Breeze ASR 25 q8_0（whisper.cpp，VAD）"
  echo
  echo "## 逐字稿"
  echo
  cat "$wav.txt"
} > "$base.breeze.md"
echo "逐字稿 → $base.breeze.md"
.venv/bin/python haiku.py "$base.breeze.md"
