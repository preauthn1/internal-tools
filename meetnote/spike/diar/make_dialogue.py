#!/usr/bin/env python3
"""Builds a synthetic 4-speaker meeting with known speaker turns, for diarization tests.

    python3 make_dialogue.py        # writes out/dialogue.wav (16 kHz mono) and out/dialogue.ref.tsv

Voices are macOS zh_TW `say` voices. Meijia and Flo are both female on purpose (similar voices).
ref.tsv: start<TAB>end<TAB>speaker, in seconds.
"""
import array, os, random, subprocess, sys, tempfile

TURNS = [
    ("Meijia", "好，我們開始今天的週會，先請大家輪流報告一下進度。"),
    ("Eddy", "我先說。後端的報表 API 已經做完了，昨天也部署到 staging。"),
    ("Meijia", "有測過效能嗎？"),
    ("Eddy", "有，一萬筆資料大概兩秒內回來。"),
    ("Flo", "前端這邊還在串接，預計明天可以給測試。"),
    ("Grandpa", "我這邊有一個問題，代理後台的權限要怎麼設定？是照角色分，還是照帳號分？"),
    ("Meijia", "照角色分比較好管理。"),
    ("Eddy", "對。"),
    ("Grandpa", "好，那我就照角色來設計，週五前把規格寫出來。"),
    ("Flo", "那規格出來之後，我再調整畫面。"),
    ("Meijia", "可以。另外上週的資安事件，修補都上線了嗎？"),
    ("Eddy", "都上線了，密碼重設的信也發出去了。"),
    ("Grandpa", "客服那邊還有收到一些詢問，我整理之後再轉給你們。"),
    ("Flo", "好。"),
    ("Meijia", "最後提醒一下，下週三要做季度回顧，請大家準備三到五分鐘的分享。"),
    ("Eddy", "要做簡報嗎？"),
    ("Meijia", "不用，簡單講就好。"),
    ("Flo", "了解。"),
    ("Grandpa", "那我沒有其他問題了。"),
    ("Meijia", "好，那今天就到這裡，謝謝大家。"),
]
RATE = 16_000

def read_wav(path):
    data = open(path, "rb").read()
    pos = 12
    while pos < len(data):
        tag, size = data[pos:pos + 4], int.from_bytes(data[pos + 4:pos + 8], "little")
        if tag == b"data":
            return array.array("h", data[pos + 8:pos + 8 + size])
        pos += 8 + size + (size & 1)

# `python3 make_dialogue.py natural`: the more natural voices instead. Meijia, Tingting and Sinji are
# all female (a hard case); Eddy, Flo and Grandpa share one older synthesis engine and sound alike.
VOICE_SETS = {"eloquence": {}, "natural": {"Eddy": "Rocko", "Flo": "Tingting", "Grandpa": "Sinji"}}


def main():
    variant = sys.argv[1] if len(sys.argv) > 1 else "eloquence"
    swap = VOICE_SETS[variant]
    global TURNS
    TURNS = [(swap.get(v, v), t) for v, t in TURNS]
    os.makedirs("out", exist_ok=True)
    rng = random.Random(7)
    audio = array.array("h", [0] * int(0.8 * RATE))
    refs = []
    with tempfile.TemporaryDirectory() as tmp:
        for i, (voice, text) in enumerate(TURNS):
            aiff, wav = f"{tmp}/{i}.aiff", f"{tmp}/{i}.wav"
            # Eddy, Flo, Grandpa and Rocko also exist as English and zh_CN voices; name the zh_TW one
            # exactly. Tingting (zh_CN) and Sinji (zh_HK, reads in Cantonese) have a single voice each.
            name = voice if voice in ("Tingting", "Sinji") else f"{voice} (中文（台灣）)"
            subprocess.run(["say", "-v", name, "-o", aiff, text], check=True)
            subprocess.run(["afconvert", "-f", "WAVE", "-d", "LEI16@16000", "-c", "1", aiff, wav], check=True)
            turn = read_wav(wav)
            start = len(audio) / RATE
            audio.extend(turn)
            refs.append((start, len(audio) / RATE, voice))
            # Pauses between turns vary from a quick reply to a longer gap.
            audio.extend([0] * int(rng.uniform(0.25, 1.2) * RATE))
    # A little room noise, so silence is not digital zero.
    audio = array.array("h", (s + rng.randint(-30, 30) for s in audio))
    import wave
    stem = "dialogue" if variant == "eloquence" else f"dialogue-{variant}"
    with wave.open(f"out/{stem}.wav", "wb") as w:
        w.setnchannels(1); w.setsampwidth(2); w.setframerate(RATE); w.writeframes(audio.tobytes())
    with open(f"out/{stem}.ref.tsv", "w") as f:
        for s, e, v in refs:
            f.write(f"{s:.2f}\t{e:.2f}\t{v}\n")
    print(f"{len(audio) / RATE:.1f}s, {len(TURNS)} turns, speakers: {sorted({v for *_, v in refs})}")

if __name__ == "__main__":
    main()
