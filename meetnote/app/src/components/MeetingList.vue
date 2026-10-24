<script setup lang="ts">
import type { Meeting } from "../lib/api";
import { mmss } from "../lib/minutes";

defineProps<{ meetings: Meeting[]; selectedId: number | null; progress: Record<number, number> }>();
defineEmits<{ select: [id: number] }>();

function when(iso: string) {
  return new Date(iso).toLocaleString("zh-TW", { month: "numeric", day: "numeric", hour: "2-digit", minute: "2-digit", hour12: false });
}

function statusLabel(m: Meeting, p: number | undefined) {
  switch (m.status) {
    case "recording":
      return "錄音中";
    case "transcribing":
      return p === undefined ? "準備辨識…" : `辨識中 ${Math.round(p * 100)}%`;
    case "summarizing":
      return "產生會議記錄中…";
    case "failed":
      return m.minutes ? "" : "需要處理";
    default:
      return "";
  }
}
</script>

<template>
  <ul class="list">
    <li v-if="meetings.length === 0" class="muted placeholder">還沒有會議</li>
    <li
      v-for="m in meetings"
      :key="m.id"
      :class="{ selected: m.id === selectedId }"
      @click="$emit('select', m.id)"
    >
      <div class="title">{{ m.minutes?.summary.title || m.title }}</div>
      <div class="meta">
        <span>{{ when(m.created_at) }}</span>
        <span v-if="m.duration_s">・{{ mmss(m.duration_s) }}</span>
        <span v-if="statusLabel(m, progress[m.id])" class="status" :class="m.status">{{ statusLabel(m, progress[m.id]) }}</span>
      </div>
    </li>
  </ul>
</template>

<style scoped>
.list {
  list-style: none;
  margin: 0;
  padding: 0 8px;
  flex: 1;
  overflow-y: auto;
}

li {
  padding: 8px 10px;
  border-radius: var(--radius);
  cursor: pointer;
}

li:hover {
  background: var(--bg-hover);
}

li.selected {
  background: var(--bg-selected);
}

li.placeholder {
  cursor: default;
  background: none;
  font-size: 13px;
}

.title {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta {
  font-size: 12px;
  color: var(--muted);
  display: flex;
  flex-wrap: wrap;
  gap: 2px 0;
}

.status {
  margin-left: auto;
  color: var(--accent);
}

.status.failed {
  color: var(--danger);
}
</style>
