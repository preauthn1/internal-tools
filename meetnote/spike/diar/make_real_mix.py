#!/usr/bin/env python3
"""Builds a "meeting" from four real people's recordings, interleaved, for diarization tests.

    python3 make_real_mix.py      # writes out/real.wav and out/real.ref.tsv

Each source is cut into speech pieces of 1.5-8 s with whisper.cpp's Silero VAD, and the pieces are
interleaved in a shuffled turn order. Caveat: every source has its own microphone and room, which
makes speakers easier to tell apart than in a real meeting on one microphone.
"""
import array, random, re, subprocess, wave
from pathlib import Path

SPIKE = Path(__file__).resolve().parent.parent
SOURCES = {
    "Rory": SPIKE / "samples/script-01.m4a",
    "talk-1014": SPIKE / "recordings/2026-10-14_1332.m4a",
    "talk-1015a": SPIKE / "recordings/2026-10-15_1330.m4a",
    "talk-1015b": SPIKE / "recordings/2026-10-15_1423.m4a",
}
TURNS_PER_SPEAKER = 6
RATE = 16_000


def read_wav(path):
    data = path.read_bytes()
    pos = 12
    while pos < len(data):
        tag, size = data[pos:pos + 4], int.from_bytes(data[pos + 4:pos + 8], "little")
        if tag == b"data":
            return array.array("h", data[pos + 8:pos + 8 + size])
        pos += 8 + size + (size & 1)


def pieces(name, src, rng):
    wav = Path("out") / f"src-{name}.wav"
    subprocess.run(["afconvert", "-f", "WAVE", "-d", "LEI16@16000", "-c", "1", str(src), str(wav)], check=True)
    out = subprocess.run(
        ["whisper-vad-speech-segments", "-vm", str(SPIKE / "models/ggml-silero-v5.1.2.bin"), "-f", str(wav), "-vt", "0.5"],
        capture_output=True, text=True, check=True,
    ).stdout
    segs = [(float(a) / 100, float(b) / 100) for a, b in re.findall(r"start = ([\d.]+), end = ([\d.]+)", out)]
    segs = [(a, b) for a, b in segs if 1.5 <= b - a <= 8.0]
    audio = read_wav(wav)
    chosen = rng.sample(segs, TURNS_PER_SPEAKER)
    return [audio[int(a * RATE):int(b * RATE)] for a, b in chosen]


def main():
    Path("out").mkdir(exist_ok=True)
    rng = random.Random(11)
    pool = {name: pieces(name, src, rng) for name, src in SOURCES.items()}
    order = [name for name in SOURCES for _ in range(TURNS_PER_SPEAKER)]
    rng.shuffle(order)
    audio = array.array("h", [0] * int(0.5 * RATE))
    with open("out/real.ref.tsv", "w") as ref:
        for name in order:
            piece = pool[name].pop()
            start = len(audio) / RATE
            audio.extend(piece)
            ref.write(f"{start:.2f}\t{len(audio) / RATE:.2f}\t{name}\n")
            audio.extend([0] * int(rng.uniform(0.3, 1.0) * RATE))
    with wave.open("out/real.wav", "wb") as w:
        w.setnchannels(1); w.setsampwidth(2); w.setframerate(RATE); w.writeframes(audio.tobytes())
    print(f"{len(audio) / RATE:.1f}s, {len(order)} turns from {len(SOURCES)} people")


if __name__ == "__main__":
    main()
