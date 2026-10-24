#!/usr/bin/env python3
"""Character error rate of one or more hypothesis transcripts against a hand-corrected reference.

    python3 cer.py samples/meeting.ref.txt out/meeting.speechanalyzer.txt out/meeting.whisper.txt

Timestamps, punctuation and whitespace are ignored; English is compared case-insensitively.
Note: whisper may output Simplified characters, which count as errors here.
"""
import re
import sys
import unicodedata


def normalize(text: str) -> str:
    # For a .md with a cleaned section (AFM or Haiku), score only that section.
    m = re.search(r"^## (?:整理後|修正後逐字稿)[^\n]*\n(.*?)(?=^## |\Z)", text, flags=re.M | re.S)
    if m:
        text = m.group(1)
    # For a raw transcript .md, score only the "[mm:ss] ..." lines, not the header.
    elif re.search(r"^\[\d+:\d+\]", text, flags=re.M):
        text = "\n".join(re.findall(r"^\[\d+:\d+\].*$", text, flags=re.M))
    # Drop Haiku's correction markers: 〔原：…〕 and 〔?〕.
    text = re.sub(r"〔[^〕]*〕", "", text)
    text = re.sub(r"^\[\d+:\d+\]\s*", "", text, flags=re.M)
    return "".join(
        ch.lower() for ch in text
        if not unicodedata.category(ch).startswith(("P", "Z", "C", "S"))
    )


def distance(a: str, b: str) -> int:
    prev = list(range(len(b) + 1))
    for i, ca in enumerate(a, 1):
        cur = [i]
        for j, cb in enumerate(b, 1):
            cur.append(min(prev[j] + 1, cur[j - 1] + 1, prev[j - 1] + (ca != cb)))
        prev = cur
    return prev[-1]


def main() -> None:
    if len(sys.argv) < 3:
        sys.exit(__doc__)
    ref = normalize(open(sys.argv[1], encoding="utf-8").read())
    for path in sys.argv[2:]:
        hyp = normalize(open(path, encoding="utf-8").read())
        print(f"{distance(ref, hyp) / max(len(ref), 1):6.1%}  {path}")


if __name__ == "__main__":
    main()
