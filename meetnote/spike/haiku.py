#!/usr/bin/env python3
"""Fix and summarize a transcript with Claude Haiku (cloud), as an alternative to on-device AFM.

    .venv/bin/python haiku.py recordings/2026-10-14_1332.breeze.md

Reads the "[mm:ss] ..." lines of a transcript .md (or every line of a .txt) and writes <name>.haiku.md
next to it: the transcript with corrections applied and marked, the list of corrections, and meeting
minutes. Haiku returns only line-numbered corrections, never the whole transcript, so output stays
small. Only the transcript text is sent; audio stays local.
Needs ANTHROPIC_API_KEY in the environment or in ../.env.
"""
import json
import os
import re
import sys
import time
from pathlib import Path

import anthropic

MODEL = "claude-haiku-5-5"
PRICE_IN, PRICE_OUT = 0.10, 0.50  # USD per 1M tokens (prompts up to 100K tokens)

SYSTEM = """你會收到一份語音辨識（ASR）產生的繁體中文逐字稿，內容可能夾雜英文術語，每一行開頭是行號（例如 L12:）。逐字稿常見的問題：
- 同音或近音錯字
- 英文術語被辨識成發音相近的中文（例如 deploy 被寫成「低破」）
- ASR 幻覺：出現和上下文完全無關的詞句，常見於靜音片段，例如「(音樂)」「請不吝點贊訂閱」

請完成四件事，不要重寫整份逐字稿：

1. corrections：需要修正的地方
   - 只有在上下文能明確判斷正確詞語時才修正
   - line 是行號；original 必須是該行裡原封不動出現的字串，盡量短，只包含要改的詞；corrected 是修正後的詞；reason 簡短說明依據
   - 英文術語保持英文

2. uncertain：無法確定意思的詞句，或疑似 ASR 幻覺的片段
   - line 是行號；text 必須是該行裡原封不動出現的字串

3. paragraphs：每個新段落（新的話題）開始的行號，依序排列，第一個必須是第一行的行號

4. summary：根據修正後的內容整理會議記錄（title 是一句話的主題）
   - 只能使用逐字稿裡的資訊，不要推論或補充；uncertain 裡的內容不要寫進摘要
   - topics 依講述順序排列，必須涵蓋從開頭到結尾的內容
   - decisions 只放會議中明確做出的決定；沒有就回傳空陣列
   - action_items 只放會議中明確指派的工作；演講、課程中對聽眾的一般建議不算。owner、due 逐字稿沒提到就填「未指定」。沒有就回傳空陣列"""


def glossary() -> str:
    path = Path(__file__).resolve().parent / "glossary.txt"
    if not path.exists():
        return ""
    words = [l.strip() for l in path.read_text(encoding="utf-8").splitlines() if l.strip() and not l.startswith("#")]
    return "\n\n講者常用的詞彙（辨識結果中發音相近的詞，很可能是這些）：" + "、".join(words)

SCHEMA = {
    "type": "object",
    "properties": {
        "corrections": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "line": {"type": "integer"},
                    "original": {"type": "string"},
                    "corrected": {"type": "string"},
                    "reason": {"type": "string"},
                },
                "required": ["line", "original", "corrected", "reason"],
                "additionalProperties": False,
            },
        },
        "uncertain": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {"line": {"type": "integer"}, "text": {"type": "string"}},
                "required": ["line", "text"],
                "additionalProperties": False,
            },
        },
        "paragraphs": {"type": "array", "items": {"type": "integer"}},
        "summary": {
            "type": "object",
            "properties": {
                "kind": {"type": "string", "enum": ["meeting", "talk", "memo", "other"]},
                "title": {"type": "string"},
                "topics": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "title": {"type": "string"},
                            "points": {"type": "array", "items": {"type": "string"}},
                        },
                        "required": ["title", "points"],
                        "additionalProperties": False,
                    },
                },
                "decisions": {"type": "array", "items": {"type": "string"}},
                "action_items": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "task": {"type": "string"},
                            "owner": {"type": "string"},
                            "due": {"type": "string"},
                        },
                        "required": ["task", "owner", "due"],
                        "additionalProperties": False,
                    },
                },
            },
            "required": ["kind", "title", "topics", "decisions", "action_items"],
            "additionalProperties": False,
        },
    },
    "required": ["corrections", "uncertain", "paragraphs", "summary"],
    "additionalProperties": False,
}

TIMESTAMP = re.compile(r"^\[\d+:\d+\]\s*")


def load_env() -> None:
    """Loads ANTHROPIC_API_KEY from the repo's .env without printing it."""
    if os.environ.get("ANTHROPIC_API_KEY"):
        return
    env = Path(__file__).resolve().parent.parent / ".env"
    if env.exists():
        for line in env.read_text().splitlines():
            key, _, value = line.partition("=")
            if key.strip() == "ANTHROPIC_API_KEY":
                os.environ["ANTHROPIC_API_KEY"] = value.strip().strip('"').strip("'")


def read_transcript(path: Path) -> tuple[str, list[str]]:
    """Returns (header, lines). For .md, the header is everything before the first post-processing section."""
    content = path.read_text(encoding="utf-8")
    lines = content.splitlines()
    timed = [l for l in lines if TIMESTAMP.match(l)]
    header = re.split(r"\n## (?:整理後|修正後)", content)[0].rstrip() if path.suffix == ".md" else ""
    return header, timed or [l for l in lines if l.strip()]


def apply(lines: list[str], data: dict) -> tuple[list[str], list[dict]]:
    """Applies corrections and uncertainty marks in place; returns the new lines and the corrections that applied."""
    lines = list(lines)
    applied = []
    for c in data["corrections"]:
        i = c["line"] - 1
        if 0 <= i < len(lines) and c["original"] and c["original"] != c["corrected"] and c["original"] in lines[i]:
            lines[i] = lines[i].replace(c["original"], f"{c['corrected']}〔原：{c['original']}〕", 1)
            applied.append(c)
    for u in data["uncertain"]:
        i = u["line"] - 1
        if 0 <= i < len(lines) and u["text"] and u["text"] in lines[i]:
            lines[i] = lines[i].replace(u["text"], u["text"] + "〔?〕", 1)
    return lines, applied


def paragraphs(lines: list[str], starts: list[int]) -> list[str]:
    """Joins short ASR lines into paragraphs, keeping the first line's timestamp."""
    cuts = sorted({s - 1 for s in starts if 0 < s <= len(lines)} | {0})
    out = []
    for begin, end in zip(cuts, cuts[1:] + [len(lines)]):
        stamp = TIMESTAMP.match(lines[begin])
        texts = [TIMESTAMP.sub("", l).strip() for l in lines[begin:end]]
        joined = ""
        for t in texts:
            if joined and not re.search(r"[。，、？！?!,.：:；;）)]$", joined):
                joined += "，"
            joined += t
        out.append((stamp.group(0) if stamp else "") + joined + ("" if re.search(r"[。？！?!]$", joined) else "。"))
    return out


def render(header: str, lines: list[str], data: dict, applied: list[dict]) -> str:
    summary = data["summary"]
    parts = [header, "", f"## 修正後逐字稿（{MODEL}）", ""]
    for p in paragraphs(lines, data["paragraphs"]):
        parts += [p, ""]
    parts += ["## 修正紀錄", ""]
    if applied:
        parts += ["| 行 | 原文 | 修正 | 依據 |", "|---|---|---|---|"]
        parts += [f"| {c['line']} | {c['original']} | {c['corrected']} | {c['reason']} |" for c in applied]
    else:
        parts.append("（無）")
    parts += ["", f"## 會議記錄（{MODEL}）", "", f"**{summary['title']}**（類型：{summary['kind']}）", ""]
    for topic in summary["topics"]:
        parts += [f"### {topic['title']}", ""] + [f"- {p}" for p in topic["points"]] + [""]
    if summary["decisions"]:
        parts += ["### 決議", ""] + [f"- {d}" for d in summary["decisions"]] + [""]
    if summary["action_items"]:
        parts += ["### 待辦事項", "", "| 事項 | 負責人 | 期限 |", "|---|---|---|"]
        parts += [f"| {a['task']} | {a['owner']} | {a['due']} |" for a in summary["action_items"]] + [""]
    return "\n".join(parts).lstrip() + "\n"


def process(path: Path) -> Path:
    header, lines = read_transcript(path)
    numbered = "\n".join(f"L{i}: {TIMESTAMP.sub('', l)}" for i, l in enumerate(lines, 1))
    load_env()
    client = anthropic.Anthropic()

    start = time.monotonic()
    with client.messages.stream(
        model=MODEL,
        max_tokens=64000,
        system=SYSTEM + glossary(),
        output_config={"effort": "medium", "format": {"type": "json_schema", "schema": SCHEMA}},
        messages=[{"role": "user", "content": numbered}],
    ) as stream:
        message = stream.get_final_message()
    elapsed = time.monotonic() - start

    if message.stop_reason != "end_turn":
        sys.exit(f"未完成：stop_reason={message.stop_reason}")
    data = json.loads(next(b.text for b in message.content if b.type == "text"))
    fixed, applied = apply(lines, data)

    out = path.with_name(path.name.removesuffix(path.suffix) + ".haiku.md")
    out.write_text(render(header, fixed, data, applied), encoding="utf-8")
    usage = message.usage
    cost = usage.input_tokens * PRICE_IN / 1e6 + usage.output_tokens * PRICE_OUT / 1e6
    skipped = len(data["corrections"]) - len(applied)
    print(
        f"{MODEL}：耗時 {elapsed:.1f} 秒，輸入 {usage.input_tokens} / 輸出 {usage.output_tokens} tokens，"
        f"約 US${cost:.4f}，修正 {len(applied)} 處（略過 {skipped} 筆對不上原文的）→ {out}"
    )
    return out


def main() -> None:
    if len(sys.argv) != 2:
        sys.exit(__doc__)
    process(Path(sys.argv[1]))


if __name__ == "__main__":
    main()
