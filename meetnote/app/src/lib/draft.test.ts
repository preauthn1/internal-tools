import { describe, expect, it } from "vitest";
import { applyCovered, applyDraftEvent, emptyDraft, pendingDraft } from "./draft";

const ev = (kind: "volatile" | "final", start: number, end: number, text: string) =>
  ({ id: 1, kind, start, end, text, message: "" }) as const;

describe("draft", () => {
  it("shows finals and the current volatile text until Breeze covers them", () => {
    let s = emptyDraft();
    s = applyDraftEvent(s, ev("volatile", 4, 6, "好，我們"));
    s = applyDraftEvent(s, ev("final", 4, 15, "好，我們開始今天的站會"));
    s = applyDraftEvent(s, ev("volatile", 15, 18, "先講"));
    expect(pendingDraft(s).map((l) => l.text)).toEqual(["好，我們開始今天的站會", "先講"]);

    s = applyCovered(s, 15.2);
    expect(pendingDraft(s).map((l) => l.text)).toEqual(["先講"]);
  });

  it("a final replaces the volatile text", () => {
    let s = applyDraftEvent(emptyDraft(), ev("volatile", 0, 2, "好"));
    s = applyDraftEvent(s, ev("final", 0, 3, "好的"));
    expect(s.volatile).toBeNull();
    expect(pendingDraft(s).map((l) => l.text)).toEqual(["好的"]);
  });
});
