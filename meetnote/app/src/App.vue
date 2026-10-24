<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { open } from "@tauri-apps/plugin-dialog";
import {
  api,
  onMeetingUpdated,
  onModelLoad,
  onRecordingChanged,
  onRecordingCovered,
  onRecordingDraft,
  onRecordingError,
  onRecordingLevel,
  onTrayError,
  type Meeting,
  type ModelLoad,
  type ModelStatus,
  type RecordingLevel,
  type RecordingState,
} from "./lib/api";
import { summarize } from "./lib/summarizer";
import { mmss } from "./lib/minutes";
import { applyCovered, applyDraftEvent, emptyDraft, type DraftState } from "./lib/draft";
import MeetingList from "./components/MeetingList.vue";
import MeetingView from "./components/MeetingView.vue";
import SettingsView from "./components/SettingsView.vue";
import SearchResults from "./components/SearchResults.vue";
import SplashScreen from "./components/SplashScreen.vue";

const meetings = ref<Meeting[]>([]);
const selectedId = ref<number | null>(null);
const view = ref<"meeting" | "settings">("meeting");
/** Transcription progress (0–1) by meeting id, from "meeting-updated" events. */
const progress = ref<Record<number, number>>({});
const models = ref<ModelStatus | null>(null);
const hasKey = ref(true);
const error = ref("");
/** Bumped to make the open meeting reload after a background change. */
const revision = ref(0);
const recording = ref<RecordingState | null>(null);
const level = ref<RecordingLevel | null>(null);
const stopping = ref(false);
/** SpeechAnalyzer's grey text for the meeting being recorded. */
const draft = ref<{ id: number; state: DraftState } | null>(null);
const query = ref("");
/** null until the first answer, so the main window never flashes before the loading screen. */
const modelLoad = ref<ModelLoad | null>(null);

/** Recording and importing both need the models in memory. */
const canTranscribe = computed(() => !!models.value?.ready && modelLoad.value?.state === "ready");
const cannotTranscribeReason = computed(() =>
  !models.value?.ready
    ? "請先到設定下載辨識模型"
    : modelLoad.value?.state === "failed"
      ? "辨識模型載入失敗，請重新開啟 App"
      : "",
);

function onLoad(load: ModelLoad) {
  modelLoad.value = load;
  if (load.state === "failed") error.value = `辨識模型載入失敗：${load.error}`;
}
/** Transcript line to show when a meeting opens from search. */
const focusLine = ref<number | null>(null);

async function refresh() {
  meetings.value = await api.listMeetings();
  revision.value++;
}

/** Hands meetings whose transcript is ready to Claude, including ones left over from a previous run. */
function summarizePending() {
  for (const m of meetings.value) {
    if (m.status === "summarizing") summarize(m.id).then(refresh);
  }
}

async function checkSetup() {
  models.value = await api.modelStatus();
  hasKey.value = !!(await api.apiKeyHint());
}

async function importAudio() {
  error.value = "";
  const path = await open({
    multiple: false,
    directory: false,
    filters: [{ name: "音檔", extensions: ["m4a", "mp3", "wav", "aac", "flac", "ogg", "caf", "mp4", "mov"] }],
  });
  if (typeof path !== "string") return;
  try {
    const id = await api.importAudio(path);
    await refresh();
    selectedId.value = id;
    view.value = "meeting";
  } catch (e) {
    error.value = String(e);
  }
}

async function syncRecording() {
  recording.value = await api.recordingState();
  if (!recording.value) {
    level.value = null;
    draft.value = null;
  }
}

async function startRecording() {
  error.value = "";
  try {
    const id = await api.startRecording();
    await syncRecording();
    await refresh();
    selectedId.value = id;
    view.value = "meeting";
  } catch (e) {
    error.value = String(e);
  }
}

async function stopRecording() {
  if (stopping.value) return;
  stopping.value = true;
  try {
    await api.stopRecording();
  } finally {
    stopping.value = false;
    await syncRecording();
  }
}

async function togglePause() {
  if (!recording.value) return;
  await api.setRecordingPaused(!recording.value.paused);
  await syncRecording();
}

function select(id: number, line: number | null = null) {
  selectedId.value = id;
  focusLine.value = line;
  view.value = "meeting";
}

async function onDeleted() {
  selectedId.value = null;
  await refresh();
}

async function onSettingsSaved() {
  await checkSetup();
}

const unlisteners: (() => void)[] = [];
onMounted(async () => {
  // Listen before asking, so a load that finishes in between is not missed.
  unlisteners.push(await onModelLoad(onLoad));
  const load = await api.modelLoadState();
  if (modelLoad.value === null) onLoad(load);
  unlisteners.push(
    await onRecordingLevel((e) => (level.value = e)),
    await onRecordingDraft((e) => {
      const current = draft.value?.id === e.id ? draft.value.state : emptyDraft();
      draft.value = { id: e.id, state: applyDraftEvent(current, e) };
    }),
    await onRecordingCovered(({ id, until }) => {
      if (draft.value?.id === id) draft.value = { id, state: applyCovered(draft.value.state, until) };
    }),
    await onRecordingChanged(async () => {
      await syncRecording();
      await refresh();
      if (recording.value) {
        selectedId.value = recording.value.meeting_id;
        view.value = "meeting";
      }
    }),
    await onRecordingError(async ({ message }) => {
      error.value = message;
      if (recording.value) await stopRecording();
    }),
    await onTrayError((message) => (error.value = message)),
  );
  unlisteners.push(await onMeetingUpdated(async ({ id, progress: p }) => {
    if (p !== null) {
      progress.value = { ...progress.value, [id]: p };
      return;
    }
    const { [id]: _, ...rest } = progress.value;
    progress.value = rest;
    await refresh();
    summarizePending();
  }));
  await checkSetup();
  await syncRecording();
  await refresh();
  summarizePending();
  if (!models.value?.ready || !hasKey.value) view.value = "settings";
});
onUnmounted(() => unlisteners.forEach((u) => u()));
</script>

<template>
  <Transition name="splash">
    <SplashScreen v-if="modelLoad === null || modelLoad.state === 'loading'" />
  </Transition>
  <div class="app">
    <aside class="sidebar">
      <div class="sidebar-head">
        <h1>會議筆記</h1>
        <button :disabled="!canTranscribe" :title="cannotTranscribeReason" @click="importAudio">匯入</button>
      </div>
      <div class="record-bar">
        <button v-if="!recording" class="primary record" :disabled="!canTranscribe" :title="cannotTranscribeReason" @click="startRecording">
          <span class="rec-dot" />開始錄音
        </button>
        <div v-else class="live-box">
          <button class="live-status" title="打開這場錄音" @click="select(recording.meeting_id)">
            <span class="rec-dot" :class="{ paused: recording.paused }" />
            {{ recording.paused ? "已暫停" : "錄音中" }}
            <span class="live-time">{{ mmss(level?.elapsed ?? 0) }}</span>
          </button>
          <div class="live-actions">
            <button :disabled="stopping" @click="togglePause">{{ recording.paused ? "繼續" : "暫停" }}</button>
            <button class="stop" :disabled="stopping" @click="stopRecording"><span class="stop-icon" />{{ stopping ? "收尾中…" : "停止" }}</button>
          </div>
        </div>
        <p v-if="!recording && cannotTranscribeReason" class="muted reason">{{ cannotTranscribeReason }}</p>
      </div>
      <p v-if="error" class="error">{{ error }}</p>
      <div class="search">
        <input v-model="query" type="search" placeholder="搜尋逐字稿和會議記錄" @keydown.esc="query = ''" />
      </div>
      <SearchResults v-if="query.trim()" :query="query" :meetings="meetings" @open="select" />
      <MeetingList v-else :meetings="meetings" :selected-id="view === 'meeting' ? selectedId : null" :progress="progress" @select="select" />
      <button class="settings-link" :class="{ active: view === 'settings' }" @click="view = 'settings'">
        設定
        <span v-if="!models?.ready || !hasKey" class="dot" title="需要設定" />
      </button>
    </aside>

    <main class="content">
      <SettingsView v-if="view === 'settings'" :models="models" @saved="onSettingsSaved" />
      <MeetingView
        v-else-if="selectedId !== null"
        :id="selectedId"
        :revision="revision"
        :progress="progress[selectedId] ?? null"
        :recording="recording?.meeting_id === selectedId ? recording : null"
        :level="level?.id === selectedId ? level : null"
        :focus-line="focusLine"
        :draft="draft?.id === selectedId ? draft.state : null"
        @changed="refresh"
        @deleted="onDeleted"
        @stop="stopRecording"
        @toggle-pause="togglePause"
      />
      <div v-else class="empty">
        <p>按「開始錄音」錄下會議，或「匯入」現有的錄音檔。</p>
        <p class="muted">支援 m4a、mp3、wav 等格式。錄音會在本機轉成逐字稿，只有逐字稿文字會送到 Claude 產生會議記錄。</p>
      </div>
    </main>
  </div>
</template>

<style>
:root {
  --bg: #ffffff;
  --bg-sidebar: #f5f5f4;
  --bg-hover: #ebebea;
  --bg-selected: #e2e8f0;
  --text: #1c1917;
  --muted: #78716c;
  --border: #e7e5e4;
  --accent: #2563eb;
  --accent-text: #ffffff;
  --danger: #dc2626;
  --mark-fix: #dcfce7;
  --mark-unsure: #fef3c7;
  --radius: 8px;
  font-family: -apple-system, BlinkMacSystemFont, "PingFang TC", "Microsoft JhengHei", "Noto Sans TC", sans-serif;
  font-size: 15px;
  line-height: 1.6;
  color: var(--text);
  background: var(--bg);
  -webkit-font-smoothing: antialiased;
}

@media (prefers-color-scheme: dark) {
  :root {
    --bg: #1c1c1e;
    --bg-sidebar: #232325;
    --bg-hover: #2c2c2e;
    --bg-selected: #33373f;
    --text: #f5f5f4;
    --muted: #a8a29e;
    --border: #3a3a3c;
    --accent: #3b82f6;
    --mark-fix: #14532d;
    --mark-unsure: #713f12;
  }
}

* {
  box-sizing: border-box;
}

/* A desktop window, not a web page: the window itself never scrolls, and no area rubber-bands
   past its ends (macOS WebKit bounces by default). */
html,
body {
  height: 100%;
  overflow: hidden;
  overscroll-behavior: none;
}

body {
  margin: 0;
  background: var(--bg);
}

* {
  overscroll-behavior: none;
}

/* Printing (匯出 PDF): let the open meeting flow across pages; everything else is hidden.
   The meeting view shows only its print layout. */
@page {
  margin: 16mm 15mm;
}

@media print {
  html,
  body,
  #app {
    height: auto !important;
    overflow: visible !important;
    background: #fff !important;
  }

  .sidebar,
  .splash {
    display: none !important;
  }

  .app {
    display: block !important;
    height: auto !important;
  }

  .content {
    overflow: visible !important;
  }
}

button {
  font: inherit;
  color: inherit;
  background: var(--bg-hover);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 6px 12px;
  cursor: pointer;
}

button:hover:not(:disabled) {
  filter: brightness(0.96);
}

button:disabled {
  opacity: 0.5;
  cursor: default;
}

button.primary {
  background: var(--accent);
  color: var(--accent-text);
  border-color: var(--accent);
}

button.danger {
  color: var(--danger);
}

input,
textarea {
  font: inherit;
  color: inherit;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 6px 10px;
}

.muted {
  color: var(--muted);
}

.error {
  color: var(--danger);
  font-size: 13px;
}
</style>

<style scoped>
.splash-leave-active {
  transition: opacity 0.35s ease;
}

.splash-leave-to {
  opacity: 0;
}

.app {
  display: grid;
  grid-template-columns: 280px 1fr;
  height: 100vh;
}

.sidebar {
  display: flex;
  flex-direction: column;
  background: var(--bg-sidebar);
  border-right: 1px solid var(--border);
  min-height: 0;
}

.sidebar-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
}

.sidebar-head h1 {
  font-size: 17px;
  margin: 0;
}

.sidebar .error {
  margin: 0 16px 8px;
}

.record-bar {
  padding: 0 16px 12px;
}

.reason {
  font-size: 12px;
  margin: 6px 0 0;
  text-align: center;
}

.search {
  padding: 0 16px 8px;
}

.search input {
  width: 100%;
  font-size: 13px;
  padding: 5px 10px;
}

.record {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px;
}

.live-box {
  border: 1px solid var(--danger);
  border-radius: var(--radius);
  background: var(--bg);
  padding: 8px;
}

.live-status {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 8px;
  border: none;
  background: none;
  padding: 2px 4px 8px;
  color: var(--danger);
  font-weight: 600;
}

.live-time {
  margin-left: auto;
  font-variant-numeric: tabular-nums;
  color: var(--text);
  font-weight: 400;
}

.live-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.live-actions .stop {
  background: var(--danger);
  border-color: var(--danger);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.stop-icon {
  width: 9px;
  height: 9px;
  background: currentColor;
  border-radius: 1px;
}

.rec-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: currentColor;
}

.record.primary .rec-dot {
  background: #ef4444;
}

.live-status .rec-dot:not(.paused) {
  animation: blink 1.2s ease-in-out infinite;
}

@keyframes blink {
  50% {
    opacity: 0.3;
  }
}

.settings-link {
  margin: 8px;
  text-align: left;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  gap: 6px;
}

.settings-link.active {
  background: var(--bg-selected);
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--danger);
}

.content {
  overflow: auto;
  min-width: 0;
}

.empty {
  max-width: 460px;
  margin: 20vh auto 0;
  padding: 0 24px;
  text-align: center;
}
</style>
