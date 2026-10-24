#!/usr/bin/env python3
"""Meeting minutes with on-device Apple Foundation Models, rendered exactly like haiku.py for side-by-side comparison.

    python3 afm.py recordings/2026-10-15_1700.breeze.md

Runs ./afm-minutes (spike/AFM/main.swift) and writes <name>.afm.md next to the transcript.
Nothing leaves this Mac.
"""
import json
import subprocess
import sys
from pathlib import Path

import haiku

ROOT = Path(__file__).resolve().parent


def main() -> None:
    if len(sys.argv) != 2:
        sys.exit(__doc__)
    path = Path(sys.argv[1])
    result = subprocess.run(
        [str(ROOT / "afm-minutes"), str(path), str(ROOT / "glossary.txt")],
        stdout=subprocess.PIPE, check=True, text=True,
    )
    data = json.loads(result.stdout)
    header, lines = haiku.read_transcript(path)
    fixed, applied = haiku.apply(lines, data)
    haiku.MODEL = "AFM"
    out = path.with_name(path.name.removesuffix(path.suffix) + ".afm.md")
    out.write_text(haiku.render(header, fixed, data, applied), encoding="utf-8")
    skipped = len(data["corrections"]) - len(applied)
    print(f"AFM：耗時 {data['elapsed_seconds']:.0f} 秒，修正 {len(applied)} 處（略過 {skipped} 筆對不上原文的）→ {out}")


if __name__ == "__main__":
    main()
