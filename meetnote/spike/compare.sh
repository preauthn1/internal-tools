#!/bin/zsh
# Runs every file in samples/ through SpeechAnalyzer (+AFM) and every whisper.cpp model in models/, writing results to out/.
#   ./compare.sh            all samples
#   ./compare.sh --no-ai    skip AFM post-processing
set -e
cd "${0:A:h}"
PROMPT="以下是繁體中文的會議逐字稿，內容可能夾雜英文術語。"
mkdir -p out/tmp

for f in samples/*.(m4a|mp3|wav|aiff|caf)(N); do
  base=${f:t:r}
  echo "=== $base ==="
  ./transcribe "$f" --out out "$@"

  wav=out/tmp/$base.wav
  afconvert -f WAVE -d LEI16@16000 -c 1 "$f" "$wav"
  audio=$(afinfo "$wav" | awk '/estimated duration/ {print $3}')
  for model in models/ggml-*.bin; do
    name=${${model:t:r}#ggml-}
    start=$(perl -MTime::HiRes=time -e 'printf "%.3f", time')
    whisper-cli -m $model -f "$wav" -l zh --prompt "$PROMPT" -otxt -of out/$base.$name -np > /dev/null 2>&1
    end=$(perl -MTime::HiRes=time -e 'printf "%.3f", time')
    elapsed=$(echo "$end - $start" | bc)
    printf "%s\t%s\t%.1f\t%.1f\t%.1f\n" $name ${f:t} $audio $elapsed $(echo "$audio / $elapsed" | bc -l) >> out/timings.tsv
    echo "$name：耗時 ${elapsed} 秒 → out/$base.$name.txt"
  done
done

echo
column -t -s $'\t' out/timings.tsv
