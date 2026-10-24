// Turns a transcribed meeting into minutes. Rust marks a meeting "summarizing" once its transcript is
// saved; this module picks it up, calls Claude, and stores the result. On failure the meeting is
// marked failed with the reason, and the transcript stays usable.
import Anthropic from "@anthropic-ai/sdk";
import { api, termsOf } from "./api";
import { estimateCostUsd, generateMinutes, MinutesError, MODEL } from "./minutes";

const inFlight = new Set<number>();

function describe(e: unknown): string {
  if (e instanceof Anthropic.AuthenticationError) return "API key 無效，請到設定重新輸入";
  if (e instanceof Anthropic.RateLimitError) return "Claude API 請求太頻繁，請稍後按「重新產生」";
  if (e instanceof Anthropic.APIConnectionError) return "無法連線到 Claude API，請確認網路後按「重新產生」";
  if (e instanceof Anthropic.APIError) return `Claude API 錯誤（${e.status ?? "?"}）：${e.message}`;
  return e instanceof Error ? e.message : String(e);
}

/** Generates minutes for one meeting unless that is already running. Resolves when done. */
export async function summarize(id: number): Promise<void> {
  if (inFlight.has(id)) return;
  inFlight.add(id);
  try {
    const key = await api.getApiKey();
    if (!key) {
      await api.markFailed(id, "尚未設定 Claude API key，請到設定輸入後按「重新產生」");
      return;
    }
    try {
      // Names matched by voice let Claude see who said what; without them the minutes still work.
      await api.identifySpeakers(id);
    } catch (e) {
      console.error("無法辨認講者", e);
    }
    const { meeting, lines } = await api.getMeeting(id);
    if (lines.length === 0) {
      await api.markFailed(id, "錄音裡沒有辨識到任何語音");
      return;
    }
    const glossary = [...new Set([...termsOf(meeting), ...(await api.getGlossary())])];
    const started = performance.now();
    try {
      // Only the people who speak in this meeting; the rest would just cost tokens.
      const speaking = new Set(lines.map((l) => l.speaker));
      const people = (await api.listPeople()).filter((p) => speaking.has(p.name));
      const minutes = await generateMinutes(key, lines, glossary, people);
      await logCall(id, started, minutes.usage, null);
      await api.saveMinutes(id, minutes);
    } catch (e) {
      // A rejected request costs nothing; an answer that could not be used was still billed.
      await logCall(id, started, e instanceof MinutesError ? e.usage : null, describe(e));
      throw e;
    }
  } catch (e) {
    await api.markFailed(id, `會議記錄產生失敗：${describe(e)}`);
  } finally {
    inFlight.delete(id);
  }
}

/** Records every Claude call with its tokens and estimated cost, successful or not. */
async function logCall(
  meetingId: number,
  started: number,
  usage: { input_tokens: number; output_tokens: number } | null,
  error: string | null,
) {
  const input = usage?.input_tokens ?? 0;
  const output = usage?.output_tokens ?? 0;
  try {
    await api.logApiCall({
      meeting_id: meetingId,
      model: MODEL,
      input_tokens: input,
      output_tokens: output,
      cost_usd: estimateCostUsd(MODEL, input, output) ?? 0,
      duration_ms: Math.round(performance.now() - started),
      ok: error === null,
      error,
    });
  } catch (e) {
    // Losing a usage record must never cost the user their minutes.
    console.error("無法記錄 API 用量", e);
  }
}

/** "重新產生": also used to retry after a failure. */
export async function regenerate(id: number): Promise<void> {
  await api.markSummarizing(id);
  await summarize(id);
}

export const isSummarizing = (id: number) => inFlight.has(id);
