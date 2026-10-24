#!/usr/bin/env python3
"""Record a meeting and produce minutes, transcribing while recording so the result is ready soon after stopping.

    .venv/bin/python meeting.py [--model breeze-asr-25-q8_0]

Starts whisper-server (Breeze ASR 25 + VAD, model kept in memory) and `transcribe --mic --chunks --events`,
which writes 16 kHz pieces cut at speech pauses and reports SpeechAnalyzer results as JSON lines. Each piece
is sent to whisper-server as soon as it appears. The terminal shows Breeze lines as they are ready, with a
dimmed SpeechAnalyzer draft of the not-yet-transcribed tail on the last line. Press Ctrl-C (or send SIGINT to this process) to stop; then the
last piece is transcribed and Haiku writes the minutes. Output goes to recordings/<time>.breeze.md and
recordings/<time>.breeze.haiku.md next to the .m4a.
"""
import argparse
import json
import shutil
import signal
import subprocess
import sys
import threading
import time
import unicodedata
import urllib.request
import wave
from pathlib import Path

import haiku

ROOT = Path(__file__).resolve().parent
PORT = 8178


def log(msg: str) -> None:
    print(msg, flush=True)


def glossary_prompt() -> str:
    path = ROOT / "glossary.txt"
    words = [l.strip() for l in path.read_text(encoding="utf-8").splitlines() if l.strip() and not l.startswith("#")]
    return f"以下是繁體中文的會議逐字稿，內容可能夾雜英文術語，例如：{'、'.join(words)}。"


def start_server(model: str) -> subprocess.Popen:
    server = subprocess.Popen(
        [
            "whisper-server", "-m", str(ROOT / f"models/ggml-{model}.bin"), "-l", "zh",
            "--prompt", glossary_prompt(), "--vad", "-vm", str(ROOT / "models/ggml-silero-v5.1.2.bin"),
            # The default threshold (0.5) dropped quiet speech in testing; 0.35 with 300 ms padding kept it.
            "-vt", "0.35", "-vp", "300",
            "--port", str(PORT),
        ],
        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
        start_new_session=True,  # Ctrl-C in the terminal must not kill it before the last piece is done
    )
    for _ in range(120):
        try:
            urllib.request.urlopen(f"http://0.0.0.0:{PORT}/", timeout=1)
            return server
        except OSError:
            if server.poll() is not None:
                sys.exit("whisper-server 啟動失敗")
            time.sleep(0.5)
    sys.exit("whisper-server 啟動逾時")


def width(text: str) -> int:
    return sum(2 if unicodedata.east_asian_width(c) in "WF" else 1 for c in text)


class Display:
    """Breeze lines scroll up as they arrive; the bottom line is a dimmed draft of what Breeze hasn't covered yet."""

    def __init__(self):
        self.lock = threading.Lock()
        self.drafts: list[tuple[float, str]] = []  # (end seconds, text) of SpeechAnalyzer finals
        self.volatile = ""
        self.covered = 0.0  # Breeze has transcribed audio up to here
        self.shown = False

    def _clear(self) -> None:
        if self.shown:
            sys.stdout.write("\r\x1b[2K")
            self.shown = False

    def _draw(self) -> None:
        text = "".join(t for end, t in self.drafts if end > self.covered) + self.volatile
        if not text:
            return
        room = max(shutil.get_terminal_size().columns - 3, 10)
        tail = ""
        for c in reversed(text):
            if width(c + tail) > room:
                break
            tail = c + tail
        sys.stdout.write(f"\x1b[2m… {tail}\x1b[0m")
        sys.stdout.flush()
        self.shown = True

    def draft_final(self, end: float, text: str) -> None:
        with self.lock:
            self.drafts.append((end, text))
            self.volatile = ""
            self._clear()
            self._draw()

    def draft_volatile(self, text: str) -> None:
        with self.lock:
            self.volatile = text
            self._clear()
            self._draw()

    def breeze(self, lines: list[str], covered: float) -> None:
        with self.lock:
            self._clear()
            for line in lines:
                print(line)
            self.covered = max(self.covered, covered)
            self.drafts = [d for d in self.drafts if d[0] > self.covered]
            self._draw()

    def message(self, text: str) -> None:
        with self.lock:
            self._clear()
            print(text, flush=True)  # same stream as the draft line, so they never interleave
            self._draw()

    def close(self) -> None:
        with self.lock:
            self._clear()
            sys.stdout.flush()


def transcribe_piece(wav: Path) -> list[str]:
    """Returns "[mm:ss] text" lines, offset by the piece's start sample encoded in its name."""
    offset = int(wav.stem.split("_")[1]) / 16_000
    result = subprocess.run(
        ["curl", "-s", f"http://0.0.0.0:{PORT}/inference", "-F", f"file=@{wav}",
         "-F", "response_format=verbose_json", "-F", "temperature=0"],
        capture_output=True, check=True,
    )
    lines = []
    for seg in json.loads(result.stdout).get("segments", []):
        text = seg["text"].strip()
        if text:
            t = int(offset + seg["start"])
            lines.append(f"[{t // 60:02d}:{t % 60:02d}] {text}")
    return lines


class Worker(threading.Thread):
    """Transcribes pieces in order as they appear, until the recorder writes `done`."""

    def __init__(self, chunks: Path, display: Display):
        super().__init__(daemon=True)
        self.chunks = chunks
        self.display = display
        self.lines: list[str] = []
        self.done_count = 0
        self.last_finished = 0.0

    def run(self) -> None:
        seen: set[str] = set()
        while True:
            finished = (self.chunks / "done").exists()
            for wav in sorted(self.chunks.glob("*.wav")):
                if wav.name not in seen:
                    seen.add(wav.name)
                    with wave.open(str(wav)) as w:
                        end = int(wav.stem.split("_")[1]) / 16_000 + w.getnframes() / w.getframerate()
                    try:
                        lines = transcribe_piece(wav)
                    except (subprocess.CalledProcessError, ValueError) as e:
                        lines = []
                        self.display.message(f"第 {wav.stem} 段辨識失敗，略過：{e}")
                    self.lines += lines
                    self.display.breeze(lines, end)
                    self.done_count += 1
                    self.last_finished = time.monotonic()
            if finished:
                return
            time.sleep(0.5)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--model", default="breeze-asr-25-q8_0")
    args = parser.parse_args()

    log("載入 Breeze 模型…")
    server = start_server(args.model)
    chunks = ROOT / "out/chunks" / time.strftime("%Y%m%d_%H%M%S")
    recorder = subprocess.Popen(
        [str(ROOT / "transcribe"), "--mic", "--no-ai", "--events", "--out", str(ROOT / "recordings"), "--chunks", str(chunks)],
        cwd=ROOT, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
    )
    # Ctrl-C reaches the recorder through the terminal; a SIGINT sent only to this process is forwarded.
    signal.signal(signal.SIGINT, lambda *_: recorder.send_signal(signal.SIGINT))
    display = Display()
    worker = Worker(chunks, display)
    worker.start()

    def read_events() -> None:
        for line in recorder.stdout:
            try:
                event = json.loads(line)
            except ValueError:
                if line.strip():
                    display.message(line.rstrip())
                continue
            if event["type"] == "final":
                display.draft_final(event["end"], event["text"])
            else:
                display.draft_volatile(event["text"])

    def read_messages() -> None:
        for line in recorder.stderr:
            if line.strip():
                display.message(line.strip())

    readers = [threading.Thread(target=f, daemon=True) for f in (read_events, read_messages)]
    for r in readers:
        r.start()
    recorder.wait()
    for r in readers:
        r.join()
    stopped = time.monotonic()
    worker.join()
    display.close()
    server.terminate()
    log(f"逐字稿完成：停止後 {time.monotonic() - stopped:.1f} 秒，共 {worker.done_count} 段")

    audio = max((ROOT / "recordings").glob("*.m4a"), key=lambda p: p.stat().st_mtime)
    duration = float(subprocess.run(["afinfo", str(audio)], capture_output=True, text=True).stdout
                     .split("estimated duration:")[1].split()[0])
    transcript = audio.with_name(audio.stem + ".breeze.md")
    transcript.write_text(
        "\n".join([
            f"# 會議逐字稿 {audio.stem}", "",
            f"- 長度：{int(duration) // 60:02d}:{int(duration) % 60:02d}",
            f"- 音檔：[{audio.name}]({audio.name})",
            f"- 引擎：{args.model}（whisper.cpp，錄音中分段辨識，VAD）", "",
            "## 逐字稿", "", *worker.lines, "",
        ]),
        encoding="utf-8",
    )
    log(f"逐字稿 → {transcript}")
    haiku.process(transcript)
    log(f"全部完成：停止後 {time.monotonic() - stopped:.1f} 秒")


if __name__ == "__main__":
    main()
