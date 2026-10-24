import { describe, expect, it } from "vitest";
import {
  appliedCorrections,
  applyMinutes,
  estimateCostUsd,
  formatUsd,
  mmss,
  numbered,
  minutesMarkdown,
  paragraphText,
  speakerRuns,
  transcriptMarkdown,
  toParagraphs,
  type HaikuMinutes,
  MODEL,
  systemPrompt,
} from "./minutes";

const lines = [
  { idx: 0, start: 4.6, end: 6, text: "好我們開始今天的站會" },
  { idx: 1, start: 7, end: 9, text: "查了以後發現是 radius 的 lock 沒有釋放" },
  { idx: 2, start: 10, end: 12, text: "這個需求我會先開供單" },
  { idx: 3, start: 13, end: 15, text: "duplicate image duplicate image" },
];

const minutes: HaikuMinutes = {
  corrections: [
    { line: 2, original: "radius", corrected: "Redis", reason: "lock" },
    { line: 3, original: "供單", corrected: "工單", reason: "子單" },
    { line: 3, original: "不存在", corrected: "x", reason: "對不上原文" },
    { line: 9, original: "好", corrected: "x", reason: "行號超出範圍" },
    { line: 1, original: "站會", corrected: "站會", reason: "沒有改" },
  ],
  uncertain: [{ line: 4, text: "duplicate image duplicate image" }],
  paragraphs: [1, 3],
  summary: {
    kind: "meeting",
    title: "站會",
    overview: ["討論 Redis lock 問題。"],
    participants: ["小明"],
    topics: [{ title: "金流", points: ["Redis lock 沒釋放"] }],
    decisions: [],
    action_items: [{ task: "code review", owner: "小明", due: "今天" }],
  },
};

describe("applyMinutes", () => {
  it("applies only corrections that match their line", () => {
    const display = applyMinutes(lines, minutes);
    expect(display[1].spans).toEqual([
      { text: "查了以後發現是 " },
      { text: "Redis", original: "radius", reason: "lock", fix: 0 },
      { text: " 的 lock 沒有釋放" },
    ]);
    expect(appliedCorrections(display).map((s) => s.text)).toEqual(["Redis", "工單"]);
    expect(display[0].spans).toEqual([{ text: "好我們開始今天的站會" }]);
  });

  it("shows reverted corrections as the original text, still linked to the correction", () => {
    const display = applyMinutes(lines, { ...minutes, rejected: [0] });
    expect(display[1].spans[1]).toEqual({ text: "radius", reverted: true, corrected: "Redis", reason: "lock", fix: 0 });
    expect(appliedCorrections(display).map((s) => s.text)).toEqual(["工單"]);
    expect(paragraphText(toParagraphs(display, [1, 3])[0])).toContain("發現是 radius 的 lock");
  });

  it("marks uncertain text without changing it", () => {
    expect(applyMinutes(lines, minutes)[3].spans).toEqual([{ text: "duplicate image duplicate image", uncertain: true }]);
  });

  it("never marks inside a replacement", () => {
    const m = { ...minutes, corrections: [{ line: 3, original: "供單", corrected: "工單", reason: "" }], uncertain: [{ line: 3, text: "工單" }] };
    expect(applyMinutes(lines, m)[2].spans.filter((s) => s.uncertain)).toEqual([]);
  });

  it("leaves the transcript untouched without minutes", () => {
    expect(applyMinutes(lines, null).map((l) => l.spans)).toEqual(lines.map((l) => [{ text: l.text }]));
  });
});

describe("paragraphs", () => {
  it("splits at Haiku's line numbers and always starts at the first line", () => {
    const groups = toParagraphs(applyMinutes(lines, minutes), [3, 99, 3]);
    expect(groups.map((g) => g.map((l) => l.idx))).toEqual([[0, 1], [2, 3]]);
  });

  it("joins lines into a clean sentence with corrections applied", () => {
    const [first] = toParagraphs(applyMinutes(lines, minutes), minutes.paragraphs);
    expect(paragraphText(first)).toBe("好我們開始今天的站會，查了以後發現是 Redis 的 lock 沒有釋放。");
  });
});

describe("markdown", () => {
  const info = { title: "x", created_at: "2026-10-15T10:00:00+08:00", duration_s: 136.5 };
  const record = { ...minutes, model: "claude-haiku-5-5", generated_at: "", usage: { input_tokens: 0, output_tokens: 0 } };

  it("copies the minutes without the transcript", () => {
    const md = minutesMarkdown(info, record);
    expect(md).toContain("# 站會");
    expect(md).toContain("- 長度：02:16");
    expect(md).toContain("- 與會者：小明");
    expect(md).toContain("| code review | 小明 | 今天 |");
    expect(md).toContain("## 決議\n\n（無）");
    expect(md).not.toContain("逐字稿");
    expect(md).not.toContain("供單");
  });

  it("copies the transcript clean, keeping only the uncertain marks", () => {
    const md = transcriptMarkdown(info, lines, record);
    expect(md).toContain("# 站會（逐字稿）");
    expect(md).not.toContain("## 摘要");
    expect(md).toContain("[00:10] 這個需求我會先開工單，duplicate image duplicate image〔?〕。");
    expect(md).not.toContain("〔原");
  });

  it("copies the raw transcript before there are minutes", () => {
    const md = transcriptMarkdown({ title: "匯入的錄音", created_at: "2026-10-15T10:00:00+08:00", duration_s: 5 }, lines, null);
    expect(md).toContain("# 匯入的錄音（逐字稿）");
    expect(md).toContain("[00:04] 好我們開始今天的站會。");
  });
});

it("formats long durations with hours", () => {
  expect(mmss(3725)).toBe("1:02:05");
  expect(mmss(65.9)).toBe("01:05");
});

describe("cost", () => {
  it("prices Haiku per million tokens, with the long-prompt rate above 100K input", () => {
    // The script-02 meeting: 2,604 in, 3,493 out.
    expect(estimateCostUsd("claude-haiku-5-5", 2604, 3493)).toBeCloseTo(0.0020069, 7);
    expect(estimateCostUsd("claude-haiku-5-5", 100_000, 0)).toBeCloseTo(0.01, 9);
    expect(estimateCostUsd("claude-haiku-5-5", 100_001, 10_000)).toBeCloseTo(0.0750005, 7);
    expect(estimateCostUsd("unknown-model", 1, 1)).toBeNull();
  });

  it("formats small amounts with enough digits to be meaningful", () => {
    expect(formatUsd(0.0020069)).toBe("US$0.0020");
    expect(formatUsd(1.234)).toBe("US$1.23");
    expect(formatUsd(0)).toBe("US$0");
  });
});

describe("speakers", () => {
  const meeting = [
    { idx: 0, start: 0, end: 2, text: "大家好", speaker: "others" },
    { idx: 1, start: 2, end: 4, text: "今天先看報表", speaker: "others" },
    { idx: 2, start: 4, end: 6, text: "好我來處理", speaker: "me" },
    { idx: 3, start: 6, end: 7, text: "嗯" },
  ];

  it("labels lines for Claude only when the speaker is known", () => {
    expect(numbered(meeting)).toBe("L1（對方）: 大家好\nL2（對方）: 今天先看報表\nL3（我）: 好我來處理\nL4: 嗯");
    expect(numbered(lines)).toBe(lines.map((l, i) => `L${i + 1}: ${l.text}`).join("\n"));
  });

  it("uses the names the user labelled lines with", () => {
    expect(numbered([{ idx: 0, start: 0, end: 1, text: "我來處理", speaker: "阿彬" }])).toBe("L1（阿彬）: 我來處理");
  });

  it("splits paragraphs where the speaker changes", () => {
    expect(speakerRuns(meeting).map((r) => r.map((l) => l.idx))).toEqual([[0, 1], [2], [3]]);
    const mixed = [
      { idx: 0, start: 0, end: 1, text: "一", speaker: "RB" },
      { idx: 1, start: 1, end: 2, text: "二", speaker: "RB", speaker_score: 0.6 },
    ];
    expect(speakerRuns(mixed).length).toBe(1);
    expect(speakerRuns(mixed, true).length).toBe(2);
  });

  it("labels speaker runs in the Markdown transcript", () => {
    const record = (minutes: HaikuMinutes) => ({ ...minutes, model: MODEL, generated_at: "", usage: { input_tokens: 0, output_tokens: 0 } });
    const one = { ...minutes, corrections: [], uncertain: [], paragraphs: [1] };
    const md = transcriptMarkdown({ title: "線上會議", created_at: "2026-10-16T10:00:00+08:00", duration_s: 7 }, meeting, record(one));
    // One paragraph from Claude, split where the speaker changes.
    expect(md).toContain("[00:00] **對方**：大家好，今天先看報表。");
    expect(md).toContain("[00:04] **我**：好我來處理。");
    expect(md).toContain("[00:06] 嗯。");
  });
});

describe("systemPrompt", () => {
  it("lists the people speaking with their role", () => {
    const system = systemPrompt(["Hyperf"], [
      { name: "RB", role: "PHP" },
      { name: "Leah", role: "PM" },
      { name: "Toro", role: "" },
    ]);
    expect(system).toContain("Hyperf");
    expect(system).toContain("RB：PHP\nLeah：PM");
    expect(system).not.toContain("Toro");
  });

  it("adds nothing without a glossary or people", () => {
    expect(systemPrompt([], [])).not.toContain("講者名單（");
  });
});
