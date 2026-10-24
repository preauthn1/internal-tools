<script setup lang="ts">
import { ref, watch } from "vue";
import { api, type Meeting, type SearchHit } from "../lib/api";
import { mmss } from "../lib/minutes";

const props = defineProps<{ query: string; meetings: Meeting[] }>();
const emit = defineEmits<{ open: [id: number, line: number | null] }>();

const hits = ref<SearchHit[]>([]);
const searching = ref(false);
let timer: ReturnType<typeof setTimeout> | undefined;

watch(
  () => props.query,
  (q) => {
    clearTimeout(timer);
    searching.value = true;
    timer = setTimeout(async () => {
      const results = await api.search(q);
      // Ignore answers to a query the user has already typed past.
      if (q === props.query) {
        hits.value = results;
        searching.value = false;
      }
    }, 200);
  },
  { immediate: true },
);

const titleOf = (id: number) => {
  const m = props.meetings.find((m) => m.id === id);
  return m?.minutes?.summary.title || m?.title || "";
};

/** Splits text around case-insensitive matches of the query, for highlighting. */
function parts(text: string): { text: string; hit: boolean }[] {
  const q = props.query.trim().toLowerCase();
  if (!q) return [{ text, hit: false }];
  const out: { text: string; hit: boolean }[] = [];
  const lower = text.toLowerCase();
  let i = 0;
  for (let at = lower.indexOf(q); at >= 0; at = lower.indexOf(q, i)) {
    if (at > i) out.push({ text: text.slice(i, at), hit: false });
    out.push({ text: text.slice(at, at + q.length), hit: true });
    i = at + q.length;
  }
  if (i < text.length) out.push({ text: text.slice(i), hit: false });
  return out;
}
</script>

<template>
  <ul class="results">
    <li v-if="!searching && hits.length === 0" class="muted placeholder">找不到「{{ query.trim() }}」</li>
    <li v-for="h in hits" :key="h.meeting_id">
      <button class="meeting" @click="emit('open', h.meeting_id, null)">
        <template v-for="(p, i) in parts(titleOf(h.meeting_id))" :key="i"><mark v-if="p.hit">{{ p.text }}</mark><template v-else>{{ p.text }}</template></template>
      </button>
      <span v-if="h.in_meeting && !h.lines.length" class="muted small">會議記錄裡有提到</span>
      <button v-for="l in h.lines" :key="l.idx" class="line" @click="emit('open', h.meeting_id, l.idx)">
        <span class="time">{{ mmss(l.start) }}</span>
        <template v-for="(p, i) in parts(l.text)" :key="i"><mark v-if="p.hit">{{ p.text }}</mark><template v-else>{{ p.text }}</template></template>
      </button>
    </li>
  </ul>
</template>

<style scoped>
.results {
  list-style: none;
  margin: 0;
  padding: 0 8px;
  flex: 1;
  overflow-y: auto;
}

li {
  padding: 6px 4px 10px;
  border-bottom: 1px solid var(--border);
}

li.placeholder {
  border: none;
  font-size: 13px;
}

button {
  display: block;
  width: 100%;
  text-align: left;
  border: none;
  background: none;
  padding: 3px 6px;
}

button:hover {
  background: var(--bg-hover);
}

.meeting {
  font-weight: 600;
}

.line {
  font-size: 13px;
  color: var(--muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.time {
  margin-right: 6px;
  font-variant-numeric: tabular-nums;
}

mark {
  background: var(--mark-unsure);
  color: var(--text);
  border-radius: 2px;
}

.small {
  font-size: 12px;
  padding-left: 6px;
}
</style>
