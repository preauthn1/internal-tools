#!/usr/bin/env python3
"""Compares two ways of cutting a recording into pieces for whisper, offline and reproducibly.

    .venv/bin/python segtest.py

Builds out/segtest/test.wav from samples/script-01.m4a:
    65 s room noise | script-01 (136 s) | 45 s room noise | first 25 s of script-01 at 30% volume
and simulates the live recorder: cut at a speech pause once a piece is >= 20 s, always cut at 60 s.

    old: every piece goes to whisper-server started with --vad (what meeting.py does today)
    new: Silero VAD only decides whether a piece contains speech at all; pieces with speech go to
         whisper-server without --vad, and timestamps are piece offset + whisper's own times.
    gap: pieces are built from VAD speech segments: a pause longer than GAP_S ends the piece and the
         silence is dropped, each piece spans first onset - PAD_BEFORE to last offset + PAD_AFTER
         (never before the previous piece's end), max MAX_S long.
         Sent without --vad; timestamps are piece start + whisper's own times.

Prints each transcript and its CER against script-01 plus the start of the quiet replay.
"""
import array
import json
import random
import re
import subprocess
import sys
import time
import urllib.request
import wave
from pathlib import Path

ROOT = Path(__file__).resolve().parent
OUT = ROOT / "out/segtest"
RATE = 16_000
MODEL = ROOT / "models/ggml-breeze-asr-25-q8_0.bin"
VAD_MODEL = ROOT / "models/ggml-silero-v5.1.2.bin"
MIN_S, MAX_S = 20, 60
# Lead padding catches a soft first word the VAD starts late on ("大家好" in recordings/2026-10-15_1326);
# trailing padding stays short so neighbouring pieces never overlap and repeat words.
GAP_S, PAD_BEFORE, PAD_AFTER = 2.0, 1.0, 0.3
VAD_ARGS = ["-vt", "0.35", "-vp", "300"]
# Where speech starts in test.wav, for checking timestamps (script-01's first word is at 4.6 s).
EXPECTED = {"script-01 開頭「好，我們開始今天的站會」": 65 + 4.6, "小聲重播開頭": 65 + 136 + 45 + 4.6}


def read_wav(path: Path) -> array.array:
    """Reads 16-bit mono PCM. Parses chunks by hand: afconvert writes WAVE_FORMAT_EXTENSIBLE, which `wave` rejects."""
    data = path.read_bytes()
    pos = 12
    while pos < len(data):
        tag, size = data[pos:pos + 4], int.from_bytes(data[pos + 4:pos + 8], "little")
        if tag == b"data":
            return array.array("h", data[pos + 8:pos + 8 + size])
        pos += 8 + size + (size & 1)
    raise ValueError(f"{path} 沒有 data 區塊")


def write_wav(path: Path, samples: array.array) -> None:
    with wave.open(str(path), "wb") as w:
        w.setnchannels(1)
        w.setsampwidth(2)
        w.setframerate(RATE)
        w.writeframes(samples.tobytes())


def noise(seconds: float) -> array.array:
    rng = random.Random(seconds)
    return array.array("h", (int(rng.gauss(0, 40)) for _ in range(int(seconds * RATE))))


def build_test() -> Path:
    OUT.mkdir(parents=True, exist_ok=True)
    src = OUT / "script-01.wav"
    subprocess.run(["afconvert", "-f", "WAVE", "-d", "LEI16@16000", "-c", "1",
                    str(ROOT / "samples/script-01.m4a"), str(src)], check=True)
    speech = read_wav(src)
    quiet = array.array("h", (int(s * 0.3) for s in speech[: 25 * RATE]))
    test = noise(65) + speech + noise(45) + quiet
    path = OUT / "test.wav"
    write_wav(path, test)
    return path


def vad_segments(path: Path) -> list[tuple[float, float]]:
    """Speech segments in seconds, from whisper.cpp's Silero VAD (it prints centiseconds)."""
    out = subprocess.run(["whisper-vad-speech-segments", "-vm", str(VAD_MODEL), "-f", str(path), *VAD_ARGS],
                         capture_output=True, text=True, check=True).stdout
    return [(float(a) / 100, float(b) / 100)
            for a, b in re.findall(r"start = ([\d.]+), end = ([\d.]+)", out)]


def cut(total: float, speech: list[tuple[float, float]]) -> list[tuple[float, float]]:
    """Simulates the recorder: a pause is the end of a VAD speech segment."""
    pieces, start = [], 0.0
    pauses = [end for _, end in speech]
    t = 0.0
    while t < total:
        t = min(t + 0.1, total)
        at_pause = any(abs(t - p) < 0.05 for p in pauses)
        if (at_pause and t - start >= MIN_S) or t - start >= MAX_S:
            pieces.append((start, t))
            start = t
    if total - start > 0.1:
        pieces.append((start, total))
    return pieces


def group(total: float, speech: list[tuple[float, float]]) -> list[tuple[float, float]]:
    """The "gap" rule: speech-only pieces, split at pauses > GAP_S or before exceeding MAX_S."""
    pieces: list[tuple[float, float]] = []
    start = end = None

    def flush() -> None:
        floor = pieces[-1][1] if pieces else 0.0
        pieces.append((max(floor, start - PAD_BEFORE), min(total, end + PAD_AFTER)))

    for s, e in speech:
        if start is not None and (s - end > GAP_S or e - start > MAX_S):
            flush()
            start = None
        if start is None:
            start = s
        end = e
    if start is not None:
        flush()
    return pieces


def start_server(port: int, vad: bool) -> subprocess.Popen:
    cmd = ["whisper-server", "-m", str(MODEL), "-l", "zh", "--port", str(port),
           "--prompt", "以下是繁體中文的會議逐字稿，內容可能夾雜英文術語。"]
    if vad:
        cmd += ["--vad", "-vm", str(VAD_MODEL), *VAD_ARGS]
    server = subprocess.Popen(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    for _ in range(120):
        try:
            urllib.request.urlopen(f"http://0.0.0.0:{port}/", timeout=1)
            return server
        except OSError:
            time.sleep(0.5)
    sys.exit("whisper-server 啟動逾時")


def transcribe(port: int, wav: Path, offset: float) -> list[tuple[float, str]]:
    result = subprocess.run(
        ["curl", "-s", f"http://0.0.0.0:{port}/inference", "-F", f"file=@{wav}",
         "-F", "response_format=verbose_json", "-F", "temperature=0"],
        capture_output=True, check=True,
    )
    return [(offset + s["start"], s["text"].strip())
            for s in json.loads(result.stdout).get("segments", []) if s["text"].strip()]


def mmss(t: float) -> str:
    return f"{int(t) // 60:02d}:{int(t) % 60:02d}"


def run_gap_on(audio: Path) -> None:
    """`segtest.py <audio>`: apply only the gap rule to a real recording and print the transcript."""
    OUT.mkdir(parents=True, exist_ok=True)
    wav = OUT / f"{audio.stem}.wav"
    subprocess.run(["afconvert", "-f", "WAVE", "-d", "LEI16@16000", "-c", "1", str(audio), str(wav)], check=True)
    samples = read_wav(wav)
    total = len(samples) / RATE
    speech = vad_segments(wav)
    pieces = group(total, speech)
    print(f"{audio.name}：{mmss(total)}，{len(pieces)} 段：" + "、".join(f"{mmss(a)}–{mmss(b)}" for a, b in pieces))
    server = start_server(8184, vad=False)
    for i, (a, b) in enumerate(pieces):
        piece = OUT / f"{audio.stem}_{i:02d}.wav"
        write_wav(piece, samples[int(a * RATE): int(b * RATE)])
        for t, text in transcribe(8184, piece, a):
            print(f"[{mmss(t)}] {text}")
    server.terminate()


def main() -> None:
    if len(sys.argv) > 1:
        return run_gap_on(Path(sys.argv[1]))
    sys.path.insert(0, str(ROOT))
    from cer import distance, normalize

    test = build_test()
    samples = read_wav(test)
    total = len(samples) / RATE
    speech = vad_segments(test)
    timed = cut(total, speech)
    print(f"測試音檔 {mmss(total)}，VAD 找到 {len(speech)} 段語音")
    script = (ROOT / "samples/script-01.ref.txt").read_text(encoding="utf-8")
    replay = "好，我們開始今天的站會。先講上週的進度，Hyperf 的 API 服務已經 deploy 到 staging 環境，" \
             "目前跑了三天都很穩定，沒有看到 error log。第二件事，金流商的 callback 上"
    ref = normalize(script + replay)

    modes = [("old", 8181, True, timed, False), ("new", 8182, False, timed, True),
             ("gap", 8183, False, group(total, speech), True)]
    for mode, port, vad, pieces, skip_silent in modes:
        server = start_server(port, vad)
        lines: list[tuple[float, str]] = []
        skipped = []
        for i, (a, b) in enumerate(pieces):
            if skip_silent and not any(s < b and e > a for s, e in speech):
                skipped.append(f"{mmss(a)}–{mmss(b)}")
                continue
            wav = OUT / f"{mode}_{i:02d}.wav"
            write_wav(wav, samples[int(a * RATE): int(b * RATE)])
            try:
                lines += transcribe(port, wav, a)
            except subprocess.CalledProcessError:
                # whisper-server --vad crashes on a piece with no speech; note it and restart.
                lines.append((a, f"⚠️ whisper-server 當掉（{mmss(a)}–{mmss(b)}），重新啟動"))
                server.wait()
                server = start_server(port, vad)
        server.terminate()
        server.wait()
        transcript = "\n".join(f"[{mmss(t)}] {text}" for t, text in lines)
        (OUT / f"{mode}.txt").write_text(transcript + "\n", encoding="utf-8")
        hyp = normalize("\n".join(text for _, text in lines if not text.startswith("⚠️")))
        print(f"\n===== {mode}：{len(pieces)} 段，CER {distance(ref, hyp) / len(ref):.1%}，"
              f"略過無語音段：{', '.join(skipped) or '無'} =====")
        print(transcript)

    print("\n預期位置：" + "；".join(f"{k} 在 {mmss(v)}" for k, v in EXPECTED.items()))


if __name__ == "__main__":
    main()
