// Grey instant text while recording: SpeechAnalyzer's draft for the part of the recording that
// Breeze has not transcribed yet.
import type { DraftEvent } from "./api";

export interface DraftLine {
  start: number;
  end: number;
  text: string;
}

export interface DraftState {
  finals: DraftLine[];
  volatile: DraftLine | null;
  /** Breeze's text covers the recording up to here. */
  covered: number;
}

export const emptyDraft = (): DraftState => ({ finals: [], volatile: null, covered: 0 });

export function applyDraftEvent(state: DraftState, e: DraftEvent): DraftState {
  const line = { start: e.start, end: e.end, text: e.text.trim() };
  if (e.kind === "final") return { ...state, finals: line.text ? [...state.finals, line] : state.finals, volatile: null };
  if (e.kind === "volatile") return { ...state, volatile: line.text ? line : null };
  return state;
}

export function applyCovered(state: DraftState, until: number): DraftState {
  // Keep a little overlap, so text at a piece boundary is never missing from both.
  return { ...state, covered: until, finals: state.finals.filter((l) => l.end > until + 0.3) };
}

/** The draft text still waiting for Breeze. */
export function pendingDraft(state: DraftState): DraftLine[] {
  const lines = state.finals.filter((l) => l.end > state.covered + 0.3);
  if (state.volatile && state.volatile.end > state.covered + 0.3) lines.push(state.volatile);
  return lines;
}
