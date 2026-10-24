<script setup lang="ts">
import { computed, nextTick, ref, watch } from "vue";
import { convertFileSrc } from "@tauri-apps/api/core";
import { ask, save } from "@tauri-apps/plugin-dialog";
import { api, termsOf, type MeetingDetail, type RecordingLevel, type RecordingState, type UsageTotal } from "../lib/api";
import {
  appliedCorrections,
  applyMinutes,
  estimateCostUsd,
  formatTokens,
  formatUsd,
  mmss,
  minutesMarkdown,
  transcriptMarkdown,
  toParagraphs,
  speakerRuns,
  speakerLabel,
  type MinutesRecord,
  type Span,
} from "../lib/minutes";
import { regenerate, isSummarizing } from "../lib/summarizer";
import MinutesEditor from "./MinutesEditor.vue";
import { pendingDraft, type DraftState } from "../lib/draft";

const props = defineProps<{
  id: number;
  revision: number;
  progress: number | null;
  /** Set while this meeting is being recorded. */
  recording: RecordingState | null;
  level: RecordingLevel | null;
  /** A transcript line to scroll to, e.g. from search. */
  focusLine: number | null;
  /** Grey instant text while recording (macOS). */
  draft: DraftState | null;
}>();
const emit = defineEmits<{ changed: []; deleted: []; stop: []; togglePause: [] }>();

const detail = ref<MeetingDetail | null>(null);
const tab = ref<"minutes" | "transcript">("minutes");
const editingTitle = ref(false);
const titleDraft = ref("");
const editingMinutes = ref(false);
const editingTranscript = ref(false);
const notice = ref("");

/** Every Claude call made for this meeting, including regenerations and failures. */
const usage = ref<UsageTotal | null>(null);

async function load() {
  detail.value = await api.getMeeting(props.id);
  usage.value = await api.meetingUsage(props.id);
}
watch(() => [props.id, props.revision], load, { immediate: true });
watch(
  () => props.id,
  () => {
    tab.value = "minutes";
    notice.value = "";
    stopping.value = false;
    editingMinutes.value = editingTranscript.value = false;
    labelledByHand.value = false;
    popover.value = null;
    playingIdx.value = null;
  },
);

const meeting = computed(() => detail.value?.meeting);
const minutes = computed(() => meeting.value?.minutes ?? null);
const display = computed(() => (detail.value ? applyMinutes(detail.value.lines, minutes.value) : []));
const paragraphs = computed(() =>
  // Claude's paragraphs (or one per line without minutes), split where the speaker changes.
  (minutes.value ? toParagraphs(display.value, minutes.value.paragraphs) : display.value.map((l) => [l])).flatMap((p) => speakerRuns(p, true)),
);
const fixes = computed(() => appliedCorrections(display.value).length);
const title = computed(() => meeting.value?.minutes?.summary.title || meeting.value?.title || "");
const busy = computed(() => ["recording", "transcribing", "summarizing"].includes(meeting.value?.status ?? ""));
const isRecording = computed(() => meeting.value?.status === "recording");
const stopping = ref(false);



function flash(text: string) {
  notice.value = text;
  setTimeout(() => notice.value === text && (notice.value = ""), 2500);
}

// ---- recording -------------------------------------------------------------------------------

const draftLines = computed(() => (isRecording.value && props.draft ? pendingDraft(props.draft) : []));

// While recording, follow the newest text.
const liveEnd = ref<HTMLElement | null>(null);
watch(
  () => [detail.value?.lines.length, draftLines.value.map((l) => l.text).join("")],
  async () => {
    if (!isRecording.value) return;
    await nextTick();
    liveEnd.value?.scrollIntoView({ behavior: "smooth", block: "end" });
  },
);
watch(isRecording, (now) => {
  if (now) tab.value = "transcript";
  else stopping.value = false;
});

function stop() {
  stopping.value = true;
  emit("stop");
}

// ---- playback --------------------------------------------------------------------------------

const player = ref<HTMLAudioElement | null>(null);
const audioSrc = computed(() => (meeting.value && !isRecording.value ? convertFileSrc(meeting.value.audio_path) : ""));
const audioError = ref(false);
watch(audioSrc, () => (audioError.value = false));
const playingIdx = ref<number | null>(null);

function playFrom(seconds: number) {
  const audio = player.value;
  if (!audio) return;
  // A little lead-in, so the first word is not clipped.
  audio.currentTime = Math.max(0, seconds - 0.5);
  void audio.play();
}

function onTimeUpdate() {
  const t = player.value?.currentTime ?? 0;
  const lines = detail.value?.lines ?? [];
  let current: number | null = null;
  for (const l of lines) {
    if (l.start <= t + 0.2) current = l.idx;
    else break;
  }
  playingIdx.value = player.value?.paused ? null : current;
}

// ---- search focus ----------------------------------------------------------------------------

const focused = ref<number | null>(null);
async function focusOn(idx: number | null) {
  if (idx === null || !detail.value) return;
  tab.value = "transcript";
  focused.value = idx;
  await nextTick();
  document.querySelector(`[data-line="${idx}"]`)?.scrollIntoView({ behavior: "smooth", block: "center" });
  setTimeout(() => focused.value === idx && (focused.value = null), 2500);
}
watch(() => [props.focusLine, detail.value?.meeting.id], () => focusOn(props.focusLine));

// ---- corrections -----------------------------------------------------------------------------

const popover = ref<{ span: Span; x: number; y: number } | null>(null);

function openPopover(span: Span, e: MouseEvent) {
  const rect = (e.target as HTMLElement).getBoundingClientRect();
  popover.value = { span, x: rect.left, y: rect.bottom + 6 };
}

async function saveMinutes(next: MinutesRecord) {
  await api.saveMinutes(props.id, next);
  await load();
  emit("changed");
}

async function toggleCorrection(span: Span) {
  if (!minutes.value || span.fix === undefined) return;
  const rejected = new Set(minutes.value.rejected ?? []);
  if (rejected.has(span.fix)) rejected.delete(span.fix);
  else rejected.add(span.fix);
  popover.value = null;
  await saveMinutes({ ...minutes.value, rejected: [...rejected].sort((a, b) => a - b) });
}

async function addToGlossary(word: string) {
  const words = await api.getGlossary();
  if (!words.includes(word)) await api.setGlossary([...words, word]);
  popover.value = null;
  flash(`已把「${word}」加入詞彙表`);
}

// ---- editing ---------------------------------------------------------------------------------

async function saveSummary(summary: MinutesRecord["summary"]) {
  if (!minutes.value) return;
  await saveMinutes({ ...minutes.value, summary, edited_at: new Date().toISOString() });
  editingMinutes.value = false;
  flash("已儲存會議記錄");
}

async function saveLine(idx: number, text: string) {
  const line = detail.value?.lines.find((l) => l.idx === idx);
  if (!line || line.text === text.trim()) return;
  await api.setLineText(props.id, idx, text);
  await load();
}

// Names offered when labelling speakers: the ones already used here, then the minutes' participants.
const speakerNames = computed(() => {
  const used = (detail.value?.lines ?? []).map((l) => l.speaker).filter((s): s is string => !!s && !(s === "me" || s === "others"));
  return [...new Set([...used, ...(minutes.value?.summary.participants ?? [])])];
});

const labelledByHand = ref(false);

async function saveSpeaker(idx: number, name: string) {
  const line = detail.value?.lines.find((l) => l.idx === idx);
  if (!line) return;
  // Confirming a voice-matched name by retyping it also makes it a hand label.
  const shown = line.speaker ? speakerLabel(line.speaker) : "";
  if (shown === name.trim() && line.speaker_score == null) return;
  await api.setLineSpeaker(props.id, idx, name);
  line.speaker = name.trim() || null;
  line.speaker_score = null;
  labelledByHand.value = true;
}

// ---- speakers by voice -----------------------------------------------------------------------

const identifying = ref(false);

/** Labels the remaining lines by voice, from everyone labelled by hand so far. */
async function identify() {
  if (identifying.value) return;
  identifying.value = true;
  try {
    const n = await api.identifySpeakers(props.id);
    await load();
    flash(n ? `已依聲音標出 ${n} 句的講者` : "沒有標出任何講者。先在「編輯」幫每個人標幾句，之後會依聲音認出其他句");
  } catch (e) {
    flash(`無法辨認講者：${e}`);
  } finally {
    identifying.value = false;
  }
}

function toggleTranscriptEditing() {
  editingTranscript.value = !editingTranscript.value;
  // Done labelling: let the new labels name the rest of the meeting.
  if (!editingTranscript.value && labelledByHand.value) {
    labelledByHand.value = false;
    void identify();
  }
}

const voiceTitle = (score: number | null | undefined) =>
  score != null ? `依聲音判斷（相似度 ${score.toFixed(2)}），不對的話在「編輯」改掉` : undefined;

/** Enter in a speaker field saves it and moves to the next line's, for labelling while listening. */
function nextSpeakerField(e: KeyboardEvent) {
  const input = e.target as HTMLInputElement;
  input.blur();
  const row = input.closest(".edit-line")?.nextElementSibling;
  row?.querySelector<HTMLInputElement>("input.who")?.focus();
}

function startRename() {
  titleDraft.value = title.value;
  editingTitle.value = true;
}

async function finishRename() {
  editingTitle.value = false;
  const t = titleDraft.value.trim();
  if (!t || t === title.value) return;
  // The heading shows the minutes title when there is one, so that is what gets renamed.
  if (minutes.value) await saveMinutes({ ...minutes.value, summary: { ...minutes.value.summary, title: t } });
  else await api.renameMeeting(props.id, t);
  emit("changed");
}

// ---- meeting terms ---------------------------------------------------------------------------

const terms = computed(() => (meeting.value ? termsOf(meeting.value) : []));
const termDraft = ref("");

async function addTerms() {
  const added = termDraft.value.split(/[、,，\n]+/).map((t) => t.trim()).filter(Boolean);
  termDraft.value = "";
  if (!added.length) return;
  await api.setMeetingTerms(props.id, [...new Set([...terms.value, ...added])]);
  await load();
}

async function removeTerm(term: string) {
  await api.setMeetingTerms(props.id, terms.value.filter((t) => t !== term));
  await load();
}

// ---- actions ---------------------------------------------------------------------------------

async function copyMinutes() {
  if (!meeting.value || !minutes.value) return;
  await navigator.clipboard.writeText(minutesMarkdown(meeting.value, minutes.value));
  flash("已複製會議記錄");
}

async function copyTranscript() {
  if (!meeting.value || !detail.value) return;
  await navigator.clipboard.writeText(transcriptMarkdown(meeting.value, detail.value.lines, minutes.value));
  flash("已複製逐字稿");
}

async function exportPdf() {
  const path = await save({ defaultPath: `${title.value}.pdf`, filters: [{ name: "PDF", extensions: ["pdf"] }] });
  if (!path) return;
  try {
    await api.savePdf(path);
    flash("已存成 PDF");
  } catch (e) {
    flash(String(e));
  }
}

async function onRegenerate() {
  if (minutes.value?.edited_at || minutes.value?.rejected?.length) {
    const ok = await ask("重新產生會覆蓋你對會議記錄做的修改，以及還原過的修正。要繼續嗎？", {
      title: "重新產生會議記錄",
      kind: "warning",
      okLabel: "重新產生",
      cancelLabel: "取消",
    });
    if (!ok) return;
  }
  editingMinutes.value = false;
  const run = regenerate(props.id);
  emit("changed");
  await run;
  emit("changed");
}

async function remove() {
  const ok = await ask(`要刪除「${title.value}」嗎？錄音、逐字稿和會議記錄都會一起刪除，無法復原。`, {
    title: "刪除會議",
    kind: "warning",
    okLabel: "刪除",
    cancelLabel: "取消",
  });
  if (!ok) return;
  await api.deleteMeeting(props.id);
  emit("deleted");
}

const lastCost = computed(() =>
  minutes.value ? estimateCostUsd(minutes.value.model, minutes.value.usage.input_tokens, minutes.value.usage.output_tokens) : null,
);

const kindLabel: Record<string, string> = { meeting: "會議", talk: "演講", memo: "口述筆記", other: "其他" };
</script>

<template>
  <article v-if="meeting" class="meeting" @click="popover = null">
    <header>
      <input v-if="editingTitle" v-model="titleDraft" class="title-input" autofocus @blur="finishRename" @keydown.enter="finishRename" />
      <h2 v-else title="點兩下改標題" @dblclick="startRename">{{ title }}</h2>
      <div class="meta muted">
        {{ new Date(meeting.created_at).toLocaleString("zh-TW", { hour12: false }) }}
        <template v-if="meeting.duration_s">・{{ mmss(meeting.duration_s) }}</template>
        <template v-if="minutes">・{{ kindLabel[minutes.summary.kind] }}</template>
      </div>
      <div class="actions">
        <button :disabled="busy || !minutes" :title="minutes ? '複製成 Markdown，可以直接貼到 Slack、工單或 wiki' : '還沒有會議記錄'" @click="copyMinutes">複製會議記錄</button>
        <button :disabled="busy || !detail?.lines.length" title="複製成 Markdown，已套用修正、標出講者" @click="copyTranscript">複製逐字稿</button>
        <button :disabled="busy || !minutes" :title="minutes ? '把會議記錄存成 PDF' : '還沒有會議記錄'" @click="exportPdf">匯出 PDF</button>
        <button :disabled="busy || isSummarizing(id) || !detail?.lines.length" @click="onRegenerate">重新產生會議記錄</button>
        <button class="danger" :disabled="meeting.status === 'transcribing' || isRecording" @click="remove">刪除</button>
        <span v-if="notice" class="notice">{{ notice }}</span>
      </div>
      <div class="terms">
        <span class="muted">本場詞彙</span>
        <span v-for="t in terms" :key="t" class="chip">{{ t }}<button title="移除" @click="removeTerm(t)">✕</button></span>
        <input
          v-model="termDraft"
          :placeholder="terms.length ? '再加…' : '這場會議會提到的人名、專案名稱…'"
          @keydown.enter="addTerms"
          @blur="addTerms"
        />
      </div>
    </header>

    <div v-if="isRecording" class="recorder">
      <div class="clock">{{ mmss(level?.elapsed ?? 0) }}</div>
      <div class="meter" :class="{ paused: recording?.paused }">
        <div class="meter-fill" :style="{ width: `${Math.min(100, Math.sqrt(level?.level ?? 0) * 160)}%` }" />
      </div>
      <button :disabled="!recording || stopping" @click="emit('togglePause')">{{ recording?.paused ? "繼續" : "暫停" }}</button>
      <button class="primary stop" :disabled="!recording || stopping" @click="stop">{{ stopping ? "收尾中…" : "停止並產生會議記錄" }}</button>
    </div>
    <p v-if="isRecording" class="muted hint">
      {{ recording?.paused ? "已暫停，暫停期間的聲音不會錄進去。" : "灰字是即時草稿，每講完一段，會換成更準確的逐字稿。補上本場詞彙，下一段就會用到。" }}
    </p>
    <div v-else-if="meeting.status === 'transcribing'" class="banner">
      <span>在本機辨識逐字稿中{{ progress !== null ? `（${Math.round(progress * 100)}%）` : "…" }}</span>
      <progress :value="progress ?? undefined" max="1" />
    </div>
    <div v-else-if="meeting.status === 'summarizing'" class="banner">Claude 正在產生會議記錄…</div>
    <div v-else-if="meeting.status === 'failed'" class="banner failed">{{ meeting.error }}</div>

    <div v-if="audioSrc" class="player">
      <audio ref="player" :src="audioSrc" controls preload="metadata" @timeupdate="onTimeUpdate" @pause="onTimeUpdate" @error="audioError = true" />
      <p v-if="audioError" class="muted small">找不到或無法播放這場會議的錄音檔。</p>
    </div>

    <nav v-if="detail?.lines.length && !isRecording" class="tabs">
      <button :class="{ active: tab === 'minutes' }" @click="tab = 'minutes'">會議記錄</button>
      <button :class="{ active: tab === 'transcript' }" @click="tab = 'transcript'">
        逐字稿<span v-if="fixes" class="muted">（修正 {{ fixes }} 處）</span>
      </button>
      <span class="spacer" />
      <button v-if="tab === 'minutes' && minutes && !editingMinutes" class="edit" :disabled="busy" @click="editingMinutes = true">編輯</button>
      <button
        v-if="tab === 'transcript' && !editingTranscript"
        class="edit"
        :disabled="busy || identifying"
        title="依你標過講者的句子，認出其他句是誰講的"
        @click="identify"
      >
        {{ identifying ? "辨認講者中…" : "依聲音標講者" }}
      </button>
      <button v-if="tab === 'transcript'" class="edit" :class="{ on: editingTranscript }" :disabled="busy || identifying" @click="toggleTranscriptEditing">
        {{ editingTranscript ? "完成" : "編輯" }}
      </button>
    </nav>

    <section v-if="tab === 'minutes' && detail?.lines.length && !isRecording" class="minutes">
      <p v-if="!minutes" class="muted">還沒有會議記錄。逐字稿已經可以在「逐字稿」分頁查看。</p>
      <MinutesEditor v-else-if="editingMinutes" :minutes="minutes" @save="saveSummary" @cancel="editingMinutes = false" />
      <template v-else>
        <p v-if="minutes.summary.participants.length" class="muted">與會者：{{ minutes.summary.participants.join("、") }}</p>
        <h3>摘要</h3>
        <ul>
          <li v-for="(o, i) in minutes.summary.overview" :key="i">{{ o }}</li>
        </ul>
        <h3>討論重點</h3>
        <div v-for="(t, i) in minutes.summary.topics" :key="i" class="topic">
          <h4>{{ t.title }}</h4>
          <ul>
            <li v-for="(p, j) in t.points" :key="j">{{ p }}</li>
          </ul>
        </div>
        <h3>決議</h3>
        <ul v-if="minutes.summary.decisions.length">
          <li v-for="(d, i) in minutes.summary.decisions" :key="i">{{ d }}</li>
        </ul>
        <p v-else class="muted">（無）</p>
        <h3>待辦事項</h3>
        <table v-if="minutes.summary.action_items.length">
          <thead>
            <tr><th>事項</th><th>負責人</th><th>期限</th></tr>
          </thead>
          <tbody>
            <tr v-for="(a, i) in minutes.summary.action_items" :key="i">
              <td>{{ a.task }}</td>
              <td :class="{ muted: a.owner === '未指定' }">{{ a.owner }}</td>
              <td :class="{ muted: a.due === '未指定' }">{{ a.due }}</td>
            </tr>
          </tbody>
        </table>
        <p v-else class="muted">（無）</p>
        <p class="footnote muted">
          由 {{ minutes.model }} 產生・{{ new Date(minutes.generated_at).toLocaleString("zh-TW", { hour12: false }) }}
          <template v-if="minutes.edited_at">・已手動編輯</template>
          <br />
          輸入 {{ formatTokens(minutes.usage.input_tokens) }} / 輸出 {{ formatTokens(minutes.usage.output_tokens) }} tokens
          <template v-if="lastCost !== null">・約 {{ formatUsd(lastCost) }}</template>
          <template v-if="usage && usage.calls > 1">・這場會議共呼叫 {{ usage.calls }} 次，累計約 {{ formatUsd(usage.cost_usd) }}</template>
        </p>
      </template>
    </section>

    <section v-else-if="tab === 'transcript' || isRecording" class="transcript">
      <template v-if="editingTranscript && !isRecording">
        <p class="legend muted">
          直接修改每一句，離開輸入框就會儲存。左邊的欄位可以標講者：點時間播放那一句，填名字後按 Enter 會跳到下一句。每個人標幾句，按「完成」後會依聲音認出其他句；標著「依聲音」的是自動判斷的，改掉或重新輸入就變成你標的。
        </p>
        <datalist id="speaker-names">
          <option v-for="n in speakerNames" :key="n" :value="n" />
        </datalist>
        <div v-for="line in detail?.lines" :key="line.idx" class="edit-line" :class="{ playing: playingIdx === line.idx }" :data-line="line.idx">
          <button class="time" @click="playFrom(line.start)">{{ mmss(line.start) }}</button>
          <input
            class="who"
            :class="{ auto: line.speaker_score != null }"
            :title="voiceTitle(line.speaker_score)"
            list="speaker-names"
            placeholder="講者"
            :value="line.speaker ? speakerLabel(line.speaker) : ''"
            @change="saveSpeaker(line.idx, ($event.target as HTMLInputElement).value)"
            @keydown.enter.prevent="nextSpeakerField"
          />
          <span class="by-voice" :title="voiceTitle(line.speaker_score)">{{ line.speaker_score != null ? "依聲音" : "" }}</span>
          <textarea :value="line.text" rows="1" @change="saveLine(line.idx, ($event.target as HTMLTextAreaElement).value)" />
        </div>
      </template>
      <template v-else>
        <p v-if="minutes && !isRecording" class="legend muted">
          點時間可以從那裡播放。講者標籤：<span class="speaker">實心</span>是你標的，<span class="speaker auto">≈ 虛線</span>是依聲音判斷的。<span class="fix">綠底</span>是 Claude 依上下文修正的詞，點一下可以還原或加入詞彙表；<span class="unsure">黃底</span>是無法確定或疑似辨識幻覺的片段。
        </p>
        <p v-else-if="!isRecording" class="legend muted">
          點時間可以從那裡播放。講者標籤：<span class="speaker">實心</span>是你標的，<span class="speaker auto">≈ 虛線</span>是依聲音判斷的。
        </p>
        <p v-for="(p, i) in paragraphs" :key="i">
          <button class="time" :disabled="!audioSrc" @click="playFrom(p[0].start)">{{ mmss(p[0].start) }}</button>
          <span
            v-if="p[0].speaker"
            class="speaker"
            :class="[p[0].speaker, { auto: p[0].speaker_score != null }]"
            :title="voiceTitle(p[0].speaker_score)"
            ><template v-if="p[0].speaker_score != null">≈ </template>{{ speakerLabel(p[0].speaker) }}</span
          >
          <template v-for="(line, j) in p" :key="line.idx">
            <template v-if="j > 0"> </template>
            <span :data-line="line.idx" :class="{ playing: playingIdx === line.idx, focused: focused === line.idx }">
              <template v-for="(s, k) in line.spans" :key="k">
                <mark v-if="s.original !== undefined" class="fix" @click.stop="openPopover(s, $event)">{{ s.text }}</mark>
                <mark v-else-if="s.reverted" class="reverted" @click.stop="openPopover(s, $event)">{{ s.text }}</mark>
                <mark v-else-if="s.uncertain" class="unsure" title="無法確定或疑似辨識幻覺">{{ s.text }}</mark>
                <template v-else>{{ s.text }}</template>
              </template>
            </span>
          </template>
        </p>
        <p v-if="draftLines.length" class="draft" title="即時草稿，Breeze 辨識完後會換成正式文字">
          <span class="time">{{ mmss(draftLines[0].start) }}</span>
          <template v-for="(l, i) in draftLines" :key="i">{{ i ? " " : "" }}{{ l.text }}</template>
        </p>
        <p v-if="isRecording && !detail?.lines.length && !draftLines.length" class="muted">等待第一段文字…</p>
        <div ref="liveEnd" />
      </template>
    </section>

    <!-- Only shown while printing (see @media print in App.vue). -->
    <section v-if="minutes" class="print-only" aria-hidden="true">
      <h1>{{ title }}</h1>
      <p class="print-meta">
        {{ new Date(meeting.created_at).toLocaleString("zh-TW", { hour12: false }) }}
        <template v-if="meeting.duration_s">・{{ mmss(meeting.duration_s) }}</template>
        <template v-if="minutes?.summary.participants.length">・與會者：{{ minutes.summary.participants.join("、") }}</template>
      </p>
      <h2>摘要</h2>
      <ul>
        <li v-for="(o, i) in minutes.summary.overview" :key="i">{{ o }}</li>
      </ul>
      <h2>討論重點</h2>
      <template v-for="(t, i) in minutes.summary.topics" :key="i">
        <h3>{{ t.title }}</h3>
        <ul>
          <li v-for="(p, j) in t.points" :key="j">{{ p }}</li>
        </ul>
      </template>
      <h2>決議</h2>
      <ul v-if="minutes.summary.decisions.length">
        <li v-for="(d, i) in minutes.summary.decisions" :key="i">{{ d }}</li>
      </ul>
      <p v-else>（無）</p>
      <h2>待辦事項</h2>
      <table v-if="minutes.summary.action_items.length">
        <thead>
          <tr><th>事項</th><th>負責人</th><th>期限</th></tr>
        </thead>
        <tbody>
          <tr v-for="(a, i) in minutes.summary.action_items" :key="i">
            <td>{{ a.task }}</td>
            <td>{{ a.owner }}</td>
            <td>{{ a.due }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else>（無）</p>
      <p class="print-footer">由會議筆記產生・會議記錄由 {{ minutes.model }} 整理</p>
    </section>

    <div v-if="popover" class="popover" :style="{ left: `${popover.x}px`, top: `${popover.y}px` }" @click.stop>
      <template v-if="popover.span.reverted">
        <p>Claude 建議改成「<strong>{{ popover.span.corrected }}</strong>」（已還原）</p>
        <p v-if="popover.span.reason" class="muted small">{{ popover.span.reason }}</p>
        <div class="popover-actions">
          <button class="primary" @click="toggleCorrection(popover.span)">套用修正</button>
        </div>
      </template>
      <template v-else>
        <p>原文：「{{ popover.span.original }}」</p>
        <p v-if="popover.span.reason" class="muted small">{{ popover.span.reason }}</p>
        <div class="popover-actions">
          <button @click="toggleCorrection(popover.span)">還原成原文</button>
          <button @click="addToGlossary(popover.span.text)">把「{{ popover.span.text }}」加入詞彙表</button>
        </div>
      </template>
    </div>
  </article>
</template>

<style scoped>
.meeting {
  max-width: 820px;
  margin: 0 auto;
  padding: 24px 32px 48px;
}

h2 {
  font-size: 22px;
  margin: 0;
  cursor: text;
}

.title-input {
  width: 100%;
  font-size: 22px;
  font-weight: 600;
}

.meta {
  font-size: 13px;
  margin: 4px 0 12px;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.notice {
  font-size: 13px;
  color: var(--accent);
}

.terms {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  font-size: 13px;
}

.terms input {
  flex: 1;
  min-width: 180px;
  font-size: 13px;
  padding: 3px 8px;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 1px 4px 1px 8px;
  border-radius: 999px;
  background: var(--bg-selected);
}

.chip button {
  border: none;
  background: none;
  padding: 0 4px;
  font-size: 11px;
  color: var(--muted);
}

.recorder {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
  padding: 12px 16px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--bg-sidebar);
}

.clock {
  font-size: 22px;
  font-variant-numeric: tabular-nums;
  min-width: 72px;
}

.meter {
  flex: 1;
  height: 8px;
  border-radius: 4px;
  background: var(--bg-hover);
  overflow: hidden;
}

.meter-fill {
  height: 100%;
  background: #22c55e;
  transition: width 0.1s linear;
}

.meter.paused .meter-fill {
  background: var(--muted);
}

.stop {
  white-space: nowrap;
}

.hint {
  font-size: 13px;
  margin: 6px 0 0;
}

.banner {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 16px 0 0;
  padding: 10px 14px;
  border-radius: var(--radius);
  background: var(--bg-sidebar);
  border: 1px solid var(--border);
  font-size: 14px;
}

.banner progress {
  flex: 1;
}

.banner.failed {
  color: var(--danger);
}

.player {
  margin-top: 16px;
  position: sticky;
  top: 0;
  z-index: 2;
  background: var(--bg);
  padding: 6px 0;
}

.player audio {
  width: 100%;
  height: 36px;
}

.tabs {
  display: flex;
  gap: 4px;
  margin: 12px 0 8px;
  border-bottom: 1px solid var(--border);
}

.tabs button {
  border: none;
  border-bottom: 2px solid transparent;
  border-radius: 0;
  background: none;
  padding: 8px 12px;
}

.tabs button.active {
  border-bottom-color: var(--accent);
  font-weight: 600;
}

.tabs .spacer {
  flex: 1;
}

.tabs .edit {
  color: var(--accent);
}

.tabs .edit.on {
  font-weight: 600;
}

h3 {
  font-size: 16px;
  margin: 24px 0 8px;
}

h4 {
  font-size: 15px;
  margin: 12px 0 4px;
}

ul {
  margin: 0;
  padding-left: 22px;
}

table {
  border-collapse: collapse;
  width: 100%;
  font-size: 14px;
}

th,
td {
  text-align: left;
  padding: 6px 8px;
  border-bottom: 1px solid var(--border);
  vertical-align: top;
}

th {
  font-weight: 600;
  color: var(--muted);
}

.footnote {
  margin-top: 32px;
  font-size: 12px;
}

.transcript p {
  margin: 0 0 14px;
}

.time {
  display: inline-block;
  min-width: 52px;
  margin-right: 4px;
  padding: 0 4px;
  border: none;
  background: none;
  font-size: 12px;
  color: var(--accent);
  font-variant-numeric: tabular-nums;
  text-align: left;
}

.time:disabled {
  color: var(--muted);
  opacity: 1;
}

.speaker {
  display: inline-block;
  margin-right: 6px;
  padding: 0 6px;
  border-radius: 999px;
  font-size: 12px;
  background: var(--bg-selected);
}

.example.invalid {
  background: var(--mark-fix);
}

/* Matched by voice rather than labelled by hand. */
.speaker.auto {
  background: none;
  border: 1px dashed var(--accent);
  color: var(--accent);
}

.edit-line .by-voice {
  flex: none;
  width: 40px;
  padding-top: 7px;
  font-size: 11px;
  color: var(--accent);
}

.edit-line .who.auto {
  color: var(--muted);
  font-style: italic;
}

.draft {
  color: var(--muted);
}

.draft .time {
  color: var(--muted);
}

.playing {
  background: var(--bg-selected);
  border-radius: 3px;
}

.focused {
  outline: 2px solid var(--accent);
  border-radius: 3px;
}

mark {
  color: inherit;
  border-radius: 3px;
  padding: 0 2px;
}

mark.fix,
mark.reverted {
  cursor: pointer;
}

.fix,
mark.fix {
  background: var(--mark-fix);
}

mark.reverted {
  background: none;
  text-decoration: underline dashed var(--muted);
}

.unsure,
mark.unsure {
  background: var(--mark-unsure);
}

.legend {
  font-size: 13px;
}

.legend span {
  padding: 0 3px;
  border-radius: 3px;
  color: var(--text);
}

.edit-line {
  display: flex;
  align-items: flex-start;
  gap: 4px;
  margin-bottom: 6px;
}

.edit-line .time {
  padding-top: 7px;
}

.edit-line .who {
  width: 84px;
  flex: none;
  font-size: 13px;
  padding: 5px 6px;
}

.edit-line.playing {
  background: var(--bg-selected);
}

.edit-line textarea {
  flex: 1;
  resize: vertical;
  font-family: inherit;
  field-sizing: content;
}

.popover {
  position: fixed;
  z-index: 10;
  max-width: 340px;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--bg);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
  font-size: 14px;
}

.popover p {
  margin: 0 0 6px;
}

.popover-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.small {
  font-size: 12px;
}

.print-only {
  display: none;
}

@media print {
  /* Only the print layout goes on paper. */
  .meeting > *:not(.print-only) {
    display: none !important;
  }

  .meeting {
    max-width: none;
    padding: 0;
  }

  .print-only {
    display: block;
    color: #000;
    font-size: 11pt;
    line-height: 1.6;
  }

  .print-only h1 {
    font-size: 18pt;
    margin: 0 0 4pt;
  }

  .print-only h2 {
    font-size: 13pt;
    margin: 16pt 0 6pt;
    padding-bottom: 2pt;
    border-bottom: 1px solid #999;
    break-after: avoid;
  }

  .print-only h3 {
    font-size: 11.5pt;
    margin: 8pt 0 2pt;
    break-after: avoid;
  }

  .print-meta {
    color: #555;
    font-size: 10pt;
    margin: 0;
  }

  .print-only table {
    width: 100%;
    border-collapse: collapse;
    font-size: 10pt;
  }

  .print-only th,
  .print-only td {
    border: 1px solid #bbb;
    padding: 3pt 6pt;
    text-align: left;
    vertical-align: top;
  }

  .print-only tr {
    break-inside: avoid;
  }

  .print-footer {
    margin-top: 24pt;
    color: #888;
    font-size: 9pt;
  }
}
</style>
