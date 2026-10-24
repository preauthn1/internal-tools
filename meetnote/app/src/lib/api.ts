// Typed wrappers for the Rust commands in src-tauri/src/lib.rs.
import { invoke } from "@tauri-apps/api/core";
import { listen, type UnlistenFn } from "@tauri-apps/api/event";
import type { MinutesRecord, TranscriptLine } from "./minutes";

export type Status = "recording" | "transcribing" | "summarizing" | "done" | "failed";

export interface Meeting {
  id: number;
  title: string;
  created_at: string;
  duration_s: number;
  audio_path: string;
  status: Status;
  minutes: MinutesRecord | null;
  error: string | null;
  /** Extra glossary words for this meeting, one per line. */
  terms: string;
}

export interface MeetingDetail {
  meeting: Meeting;
  lines: TranscriptLine[];
}

export interface ModelStatus {
  dir: string;
  ready: boolean;
  missing: { name: string; label: string; size: number }[];
  downloading: boolean;
}

/** Whether the recognition models are in memory. */
export type ModelLoad =
  | { state: "idle" } // model files missing
  | { state: "loading" }
  | { state: "ready" }
  | { state: "failed"; error: string };

export interface ModelDownload {
  label: string;
  downloaded: number;
  total: number;
}

export interface ApiCall {
  id?: number;
  meeting_id: number | null;
  created_at?: string;
  model: string;
  input_tokens: number;
  output_tokens: number;
  cost_usd: number;
  duration_ms: number;
  ok: boolean;
  error: string | null;
}

export interface UsageTotal {
  calls: number;
  input_tokens: number;
  output_tokens: number;
  cost_usd: number;
}

/** Someone in the 人員名單: labelled by hand in some transcript, or given details. */
export interface Person {
  name: string;
  role: string;
  /** Lines labelled by hand; the voiceprint is built from these. */
  hand_lines: number;
  hand_seconds: number;
  meetings: number;
}

export interface UsageSummary {
  month: UsageTotal;
  all: UsageTotal;
  recent: ApiCall[];
}

export interface SearchHit {
  meeting_id: number;
  /** The title or the minutes matched. */
  in_meeting: boolean;
  lines: TranscriptLine[];
}

export interface InputDevice {
  id: string;
  name: string;
  is_default: boolean;
}

export interface CaptureSettings {
  /** Also record what the Mac plays (the other side of an online meeting). */
  system: boolean;
  /** Whether this platform and build can capture system audio at all. */
  system_available: boolean;
}

export interface RecordingState {
  meeting_id: number;
  paused: boolean;
}

export interface RecordingLevel {
  id: number;
  /** Peak RMS 0–1 over the last ~100 ms. */
  level: number;
  /** Seconds recorded (pauses excluded). */
  elapsed: number;
}

export interface MeetingUpdated {
  id: number;
  /** 0–1 while transcribing. */
  progress: number | null;
}

export const api = {
  listMeetings: () => invoke<Meeting[]>("list_meetings"),
  getMeeting: (id: number) => invoke<MeetingDetail>("get_meeting", { id }),
  importAudio: (path: string) => invoke<number>("import_audio", { path }),
  saveMinutes: (id: number, minutes: MinutesRecord) => invoke<void>("save_minutes", { id, minutes }),
  markFailed: (id: number, error: string) => invoke<void>("mark_failed", { id, error }),
  markSummarizing: (id: number) => invoke<void>("mark_summarizing", { id }),
  renameMeeting: (id: number, title: string) => invoke<void>("rename_meeting", { id, title }),
  deleteMeeting: (id: number) => invoke<void>("delete_meeting", { id }),
  modelStatus: () => invoke<ModelStatus>("model_status"),
  getGlossary: () => invoke<string[]>("get_glossary"),
  /** Saves the page's print layout (the open meeting's minutes) as a PDF at `path`. */
  savePdf: (path: string) => invoke<void>("save_pdf", { path }),
  /** Masked key for display, or null when none is set. Never touches the keychain. */
  apiKeyHint: () => invoke<string | null>("api_key_hint"),
  /** The real key; may make macOS ask for keychain access (at most once per launch). */
  getApiKey: () => invoke<string | null>("get_api_key"),
  setApiKey: [REDACTED] string | null) => invoke<void>("set_api_key", { key }),
  listInputDevices: () => invoke<InputDevice[]>("list_input_devices"),
  getInputDevice: () => invoke<string | null>("get_input_device"),
  setInputDevice: (deviceId: string | null) => invoke<void>("set_input_device", { deviceId }),
  getCaptureSettings: () => invoke<CaptureSettings>("get_capture_settings"),
  setCaptureSystem: (enabled: boolean) => invoke<void>("set_capture_system", { enabled }),
  recordingState: () => invoke<RecordingState | null>("recording_state"),
  startRecording: () => invoke<number>("start_recording"),
  setRecordingPaused: (paused: boolean) => invoke<void>("set_recording_paused", { paused }),
  stopRecording: () => invoke<number>("stop_recording"),
  downloadModels: () => invoke<void>("download_models"),
  modelLoadState: () => invoke<ModelLoad>("model_load_state"),
  logApiCall: (call: ApiCall) => invoke<void>("log_api_call", { call }),
  usageSummary: () => invoke<UsageSummary>("usage_summary"),
  meetingUsage: (id: number) => invoke<UsageTotal>("meeting_usage", { id }),
  search: (query: string) => invoke<SearchHit[]>("search", { query }),
  setLineText: (id: number, idx: number, text: string) => invoke<void>("set_line_text", { id, idx, text }),
  /** An empty name clears the label. */
  setLineSpeaker: (id: number, idx: number, speaker: string) => invoke<void>("set_line_speaker", { id, idx, speaker }),
  listPeople: () => invoke<Person[]>("list_people"),
  setPerson: (name: string, role: string) => invoke<void>("set_person", { name, role }),
  /** Renaming to an existing name merges the two people. */
  renamePerson: (from: string, to: string) => invoke<void>("rename_person", { from, to }),
  /** Removes their labels everywhere, and with them the voiceprint. */
  deletePerson: (name: string) => invoke<void>("delete_person", { name }),
  /** Labels the other lines by voice, from everyone labelled by hand; returns how many got a name. */
  identifySpeakers: (id: number) => invoke<number>("identify_speakers", { id }),
  setMeetingTerms: (id: number, terms: string[]) => invoke<void>("set_meeting_terms", { id, terms }),
  setGlossary: (words: string[]) => invoke<void>("set_glossary", { words }),
};

export const termsOf = (m: Meeting) =>
  m.terms
    .split("\n")
    .map((t) => t.trim())
    .filter(Boolean);

export const onMeetingUpdated = (handler: (e: MeetingUpdated) => void): Promise<UnlistenFn> =>
  listen<MeetingUpdated>("meeting-updated", (event) => handler(event.payload));

export const onRecordingLevel = (handler: (e: RecordingLevel) => void): Promise<UnlistenFn> =>
  listen<RecordingLevel>("recording-level", (event) => handler(event.payload));

export const onRecordingError = (handler: (e: { id: number; message: string }) => void): Promise<UnlistenFn> =>
  listen<{ id: number; message: string }>("recording-error", (event) => handler(event.payload));

/** Recording started, stopped or paused, from the app or the menu bar. */
export const onRecordingChanged = (handler: () => void): Promise<UnlistenFn> => listen("recording-changed", handler);

/** A menu bar action failed; the message is meant for the user. */
export const onTrayError = (handler: (message: string) => void): Promise<UnlistenFn> =>
  listen<string>("tray-error", (event) => handler(event.payload));

export const onModelDownload = (handler: (e: ModelDownload) => void): Promise<UnlistenFn> =>
  listen<ModelDownload>("model-download", (event) => handler(event.payload));

/** A model download finished or failed. */
export const onModelsChanged = (handler: () => void): Promise<UnlistenFn> => listen("models-changed", handler);

export const onModelLoad = (handler: (e: ModelLoad) => void): Promise<UnlistenFn> =>
  listen<ModelLoad>("model-load", (event) => handler(event.payload));

/** A line from the macOS SpeechAnalyzer draft; times are seconds into the recording. */
export interface DraftEvent {
  id: number;
  kind: "ready" | "volatile" | "final" | "error";
  start: number;
  end: number;
  text: string;
  message: string;
}

export const onRecordingDraft = (handler: (e: DraftEvent) => void): Promise<UnlistenFn> =>
  listen<DraftEvent>("recording-draft", (event) => handler(event.payload));

/** Breeze's transcript now covers the recording up to `until` seconds. */
export const onRecordingCovered = (handler: (e: { id: number; until: number }) => void): Promise<UnlistenFn> =>
  listen<{ id: number; until: number }>("recording-covered", (event) => handler(event.payload));
