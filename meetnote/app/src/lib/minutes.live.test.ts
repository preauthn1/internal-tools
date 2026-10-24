// Calls the real API on the spike's script-01 transcript. Skipped unless ANTHROPIC_API_KEY is set:
//   ANTHROPIC_API_KEY=… npx vitest run src/lib/minutes.live.test.ts
import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";
import { appliedCorrections, applyMinutes, generateMinutes, minutesMarkdown, transcriptMarkdown } from "./minutes";

const TRANSCRIPT = new URL("../../../spike/out/segtest/rust-m4a.txt", import.meta.url);

describe.skipIf(!process.env.ANTHROPIC_API_KEY)("generateMinutes (live)", () => {
  it("returns minutes and corrections for script-01", { timeout: 120_000 }, async () => {
    const lines = readFileSync(TRANSCRIPT, "utf-8")
      .split("\n")
      .filter(Boolean)
      .map((l, idx) => {
        const [, m, s, text] = l.match(/^\[(\d+):(\d+)\] (.*)$/)!;
        return { idx, start: +m * 60 + +s, end: +m * 60 + +s, text };
      });

    const started = Date.now();
    const minutes = await generateMinutes(process.env.ANTHROPIC_API_KEY!, lines, ["Hyperf", "cert-manager", "Redis"]);
    const display = applyMinutes(lines, minutes);
    console.log(
      `${((Date.now() - started) / 1000).toFixed(1)}s, ${minutes.usage.input_tokens}/${minutes.usage.output_tokens} tokens,`,
      `applied ${appliedCorrections(display).length}/${minutes.corrections.length} corrections`,
    );
    const info = { title: "script-01", created_at: new Date().toISOString(), duration_s: 136 };
    console.log(minutesMarkdown(info, minutes));
    console.log(transcriptMarkdown(info, lines, minutes));

    expect(minutes.summary.action_items.some((a) => a.owner.includes("小明"))).toBe(true);
    expect(appliedCorrections(display).length).toBeGreaterThan(0);
  });
});
