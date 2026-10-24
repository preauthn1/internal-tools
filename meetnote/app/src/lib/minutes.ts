// Meeting minutes from Claude Haiku. Haiku never rewrites the transcript: it returns line-numbered
// corrections, uncertain spans, paragraph starts and a structured summary, and the app applies them
// locally. Ported from spike/haiku.py; see "會議記錄流程加速" in spike/FINDINGS.md for why.

import Anthropic from "@anthropic-ai/sdk";
import { zodOutputFormat } from "@anthropic-ai/sdk/helpers/zod";
import { z } from "zod/v4";

export const MODEL = "claude-haiku-5-5";

/**
 * USD per million tokens. Claude Haiku 5.5 bills a request at the higher rate once its prompt is
 * over 100K tokens (roughly a meeting of three hours or more). Thinking counts as output.
 */
const PRICES: Record<string, { input: number; output: number; longInput: number; longOutput: number; longAbove: number }> = {
  "claude-haiku-5-5": { input: 0.1, output: 0.5, longInput: 0.5, longOutput: 2.5, longAbove: 100_000 },
};

/** Estimated cost of one call, or null for a model without a known price. */
export function estimateCostUsd(model: string, inputTokens: number, outputTokens: number): number | null {
  const p = PRICES[model];
  if (!p) return null;
  const long = inputTokens > p.longAbove;
  return (inputTokens * (long ? p.longInput : p.input) + outputTokens * (long ? p.longOutput : p.output)) / 1e6;
}

export function formatUsd(usd: number): string {
  if (usd === 0) return "US$0";
  return usd < 0.01 ? `US$${usd.toFixed(4)}` : `US$${usd.toFixed(2)}`;
}

export function formatTokens(n: number): string {
  return n >= 1000 ? `${(n / 1000).toFixed(1)}K` : String(n);
}

/** Thrown when Claude answered but the minutes could not be used; the tokens were still spent. */
export class MinutesError extends Error {
  constructor(
    message: string,
    readonly usage: { input_tokens: number; output_tokens: number },
  ) {
    super(message);
  }
}

export interface TranscriptLine {
  idx: number;
  start: number;
  end: number;
  text: string;
  /**
   * "me" (microphone) or "others" (the meeting's system audio) for online-meeting recordings, or a
   * person's name: labelled by hand, or matched by voice.
   */
  speaker?: string | null;
  /** Set when the name was matched by voice (the similarity, 0.4-1); absent when labelled by hand. */
  speaker_score?: number | null;
}

export const SPEAKER_LABEL: Record<string, string> = { me: "我", others: "對方" };

export const speakerLabel = (speaker: string) => SPEAKER_LABEL[speaker] ?? speaker;

const Correction = z.object({
  line: z.number().int(),
  original: z.string(),
  corrected: z.string(),
  reason: z.string(),
});

const Uncertain = z.object({ line: z.number().int(), text: z.string() });

const ActionItem = z.object({ task: z.string(), owner: z.string(), due: z.string() });

export const MinutesSchema = z.object({
  corrections: z.array(Correction),
  uncertain: z.array(Uncertain),
  paragraphs: z.array(z.number().int()),
  summary: z.object({
    kind: z.enum(["meeting", "talk", "memo", "other"]),
    title: z.string(),
    overview: z.array(z.string()),
    participants: z.array(z.string()),
    topics: z.array(z.object({ title: z.string(), points: z.array(z.string()) })),
    decisions: z.array(z.string()),
    action_items: z.array(ActionItem),
  }),
});

export type HaikuMinutes = z.infer<typeof MinutesSchema>;
export type Correction = z.infer<typeof Correction>;

/** What the app stores for a meeting: Haiku's answer plus when and how it was made. */
export interface MinutesRecord extends HaikuMinutes {
  model: string;
  generated_at: string;
  usage: { input_tokens: number; output_tokens: number };
  /** Indexes into `corrections` the user reverted. */
  rejected?: number[];
  /** Set once the user edits the summary by hand. */
  edited_at?: string;
}

const SYSTEM = `你會收到一份語音辨識（ASR）產生的繁體中文逐字稿，內容可能夾雜英文術語，每一行開頭是行號（例如 L12:）。線上會議的逐字稿在行號後面會標出講者：（我）是錄音的使用者本人，（對方）是會議另一端的其他與會者，其他名稱是使用者標註的講者姓名，沒標示的是無法判斷。逐字稿常見的問題：
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

4. summary：根據修正後的內容整理會議記錄
   - title 是一句話的主題；overview 是 3 到 5 句的摘要
   - participants 是逐字稿中出現的與會者名字；沒有就回傳空陣列
   - 判斷待辦事項的負責人時可以參考講者標示，例如（我）說「這個我來處理」，負責人就是「我」；有附講者名單時，負責人寫名單上的名字，提到稱呼或職位（例如「PM 確認後再上」）也對應到名單上的人
   - 只能使用逐字稿裡的資訊，不要推論或補充；uncertain 裡的內容不要寫進摘要
   - topics 依講述順序排列，必須涵蓋從開頭到結尾的內容
   - decisions 只放會議中明確做出的決定；沒有就回傳空陣列
   - action_items 只放會議中明確指派的工作；演講、課程中對聽眾的一般建議不算。owner、due 逐字稿沒提到就填「未指定」。沒有就回傳空陣列`;

export function numbered(lines: TranscriptLine[]): string {
  return lines
    .map((l, i) => {
      const who = l.speaker ? speakerLabel(l.speaker) : undefined;
      return who ? `L${i + 1}（${who}）: ${l.text}` : `L${i + 1}: ${l.text}`;
    })
    .join("\n");
}

/**
 * Splits a paragraph where the speaker changes, so each run can be labelled. With `byVoice`, also
 * where hand labels and voice-matched ones meet, so the app can show which is which.
 */
export function speakerRuns<T extends TranscriptLine>(lines: T[], byVoice = false): T[][] {
  const key = (l: TranscriptLine) => `${l.speaker ?? ""}${byVoice && l.speaker_score != null ? "\u0000voice" : ""}`;
  const runs: T[][] = [];
  for (const line of lines) {
    const last = runs[runs.length - 1];
    if (last && key(last[0]) === key(line)) last.push(line);
    else runs.push([line]);
  }
  return runs;
}

/** Who the labelled speakers are, as the user described them in 人員名單. */
export interface PersonInfo {
  name: string;
  role: string;
}

/** The instructions, plus the glossary and the people speaking in this meeting when there are any. */
export function systemPrompt(glossary: string[], people: PersonInfo[]): string {
  let system = SYSTEM;
  if (glossary.length) system += `\n\n講者常用的詞彙（辨識結果中發音相近的詞，很可能是這些）：${glossary.join("、")}`;
  const described = people.filter((p) => p.role);
  if (described.length) {
    const who = described.map((p) => `${p.name}：${p.role}`);
    system += `\n\n講者名單（逐字稿標示的名字：職位）：\n${who.join("\n")}`;
  }
  return system;
}

export async function generateMinutes(
  apiKey: [REDACTED],
  lines: TranscriptLine[],
  glossary: string[] = [],
  people: PersonInfo[] = [],
): Promise<MinutesRecord> {
  // The key lives in the OS keychain and only reaches this local webview, never a third-party page.
  const client = new Anthropic({ apiKey, dangerouslyAllowBrowser: true });
  const system = systemPrompt(glossary, people);

  const message = await client.messages
    .stream({
      model: MODEL,
      max_tokens: 64000,
      system,
      output_config: { effort: "medium", format: zodOutputFormat(MinutesSchema) },
      messages: [{ role: "user", content: numbered(lines) }],
    })
    .finalMessage();

  if (message.stop_reason !== "end_turn" || !message.parsed_output) {
    throw new MinutesError(`會議記錄沒有完成（${message.stop_reason}）`, {
      input_tokens: message.usage.input_tokens,
      output_tokens: message.usage.output_tokens,
    });
  }
  return {
    ...message.parsed_output,
    model: MODEL,
    generated_at: new Date().toISOString(),
    usage: { input_tokens: message.usage.input_tokens, output_tokens: message.usage.output_tokens },
  };
}

export interface Span {
  text: string;
  /** Set when this span replaced ASR text. */
  original?: string;
  reason?: string;
  uncertain?: boolean;
  /** Index into `corrections` of the correction behind this span. */
  fix?: number;
  /** A correction the user reverted: `text` is the ASR text again, `corrected` what Claude suggested. */
  reverted?: boolean;
  corrected?: string;
}

export interface DisplayLine extends TranscriptLine {
  spans: Span[];
}

/**
 * Applies corrections and uncertainty marks to the transcript. A correction whose `original` is not
 * literally in its line is skipped, so a confused answer can never garble the transcript.
 */
export function applyMinutes(
  lines: TranscriptLine[],
  minutes: (HaikuMinutes & { rejected?: number[] }) | null,
): DisplayLine[] {
  const out: DisplayLine[] = lines.map((l) => ({ ...l, spans: [{ text: l.text }] }));
  if (!minutes) return out;

  const mark = (lineNo: number, find: string, replace: (found: string) => Span) => {
    const line = out[lineNo - 1];
    if (!line || !find) return false;
    for (let i = 0; i < line.spans.length; i++) {
      const span = line.spans[i];
      // Only plain ASR text can be marked; never mark inside an earlier replacement.
      if (span.original !== undefined || span.uncertain || span.reverted) continue;
      const at = span.text.indexOf(find);
      if (at < 0) continue;
      const parts: Span[] = [];
      if (at > 0) parts.push({ text: span.text.slice(0, at) });
      parts.push(replace(find));
      const rest = span.text.slice(at + find.length);
      if (rest) parts.push({ text: rest });
      line.spans.splice(i, 1, ...parts);
      return true;
    }
    return false;
  };

  const rejected = new Set(minutes.rejected ?? []);
  minutes.corrections.forEach((c, fix) => {
    if (c.original === c.corrected) return;
    mark(c.line, c.original, () =>
      rejected.has(fix)
        ? { text: c.original, reverted: true, corrected: c.corrected, reason: c.reason, fix }
        : { text: c.corrected, original: c.original, reason: c.reason, fix },
    );
  });
  for (const u of minutes.uncertain) {
    mark(u.line, u.text, (found) => ({ text: found, uncertain: true }));
  }
  return out;
}

/** Corrections that were actually applied (the rest did not match the transcript). */
export function appliedCorrections(lines: DisplayLine[]): Span[] {
  return lines.flatMap((l) => l.spans.filter((s) => s.original !== undefined));
}

/** Groups display lines into paragraphs starting at the 1-based line numbers Haiku chose. */
export function toParagraphs(lines: DisplayLine[], starts: number[]): DisplayLine[][] {
  const cuts = [...new Set([0, ...starts.map((s) => s - 1).filter((s) => s > 0 && s < lines.length)])].sort(
    (a, b) => a - b,
  );
  return cuts.map((begin, i) => lines.slice(begin, cuts[i + 1] ?? lines.length));
}

export function mmss(seconds: number): string {
  const s = Math.max(0, Math.floor(seconds));
  const pad = (n: number) => String(n).padStart(2, "0");
  return s >= 3600 ? `${Math.floor(s / 3600)}:${pad(Math.floor(s / 60) % 60)}:${pad(s % 60)}` : `${pad(Math.floor(s / 60))}:${pad(s % 60)}`;
}

/** Copied text reads clean: corrections show their corrected words; uncertain spans keep a 〔?〕. */
const spanText = (s: Span) => (s.uncertain ? `${s.text}〔?〕` : s.text);

/** Joins short ASR lines into readable sentences, as spike/haiku.py does. */
export function paragraphText(lines: DisplayLine[]): string {
  let joined = "";
  for (const l of lines) {
    const text = l.spans.map(spanText).join("").trim();
    if (joined && !/[。，、？！?!,.：:；;）)]$/.test(joined)) joined += "，";
    joined += text;
  }
  return /[。？！?!]$/.test(joined) ? joined : `${joined}。`;
}

export interface MeetingInfo {
  title: string;
  created_at: string;
  duration_s: number;
}

function header(meeting: MeetingInfo, minutes: MinutesRecord | null, suffix = ""): string[] {
  const date = new Date(meeting.created_at).toLocaleString("zh-TW", { hour12: false });
  return [`# ${minutes?.summary.title || meeting.title}${suffix}`, "", `- 日期：${date}`, `- 長度：${mmss(meeting.duration_s)}`];
}

/** 「複製會議記錄」: the minutes only. */
export function minutesMarkdown(meeting: MeetingInfo, minutes: MinutesRecord): string {
  const s = minutes.summary;
  const out = header(meeting, minutes);
  if (s.participants.length) out.push(`- 與會者：${s.participants.join("、")}`);
  out.push("", "## 摘要", "", ...s.overview.map((o) => `- ${o}`), "", "## 討論重點", "");
  for (const t of s.topics) out.push(`### ${t.title}`, "", ...t.points.map((p) => `- ${p}`), "");
  out.push("## 決議", "", ...(s.decisions.length ? s.decisions.map((d) => `- ${d}`) : ["（無）"]), "");
  out.push("## 待辦事項", "");
  if (s.action_items.length) {
    out.push("| 事項 | 負責人 | 期限 |", "|---|---|---|", ...s.action_items.map((a) => `| ${a.task} | ${a.owner} | ${a.due} |`));
  } else {
    out.push("（無）");
  }
  return out.join("\n") + "\n";
}

/** 「複製逐字稿」: the transcript with Claude's corrections applied, in paragraphs, labelled by speaker. */
export function transcriptMarkdown(meeting: MeetingInfo, lines: TranscriptLine[], minutes: MinutesRecord | null): string {
  const display = applyMinutes(lines, minutes);
  const out = [...header(meeting, minutes, "（逐字稿）"), ""];
  const paragraphs = minutes ? toParagraphs(display, minutes.paragraphs) : display.map((l) => [l]);
  for (const p of paragraphs) {
    for (const run of speakerRuns(p)) {
      const who = run[0].speaker ? `**${speakerLabel(run[0].speaker)}**：` : "";
      out.push(`[${mmss(run[0].start)}] ${who}${paragraphText(run)}`, "");
    }
  }
  return out.join("\n");
}
