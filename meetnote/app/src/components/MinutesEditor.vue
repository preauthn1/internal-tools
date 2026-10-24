<script setup lang="ts">
// Edits a copy of the summary; nothing is saved until 儲存.
import { ref } from "vue";
import type { MinutesRecord } from "../lib/minutes";

const props = defineProps<{ minutes: MinutesRecord }>();
const emit = defineEmits<{ save: [summary: MinutesRecord["summary"]]; cancel: [] }>();

// A plain deep copy: props are reactive proxies, which structuredClone rejects.
const draft = ref<MinutesRecord["summary"]>(JSON.parse(JSON.stringify(props.minutes.summary)));

function save() {
  const d = draft.value;
  const clean = (items: string[]) => items.map((s) => s.trim()).filter(Boolean);
  emit("save", {
    ...d,
    title: d.title.trim(),
    overview: clean(d.overview),
    participants: clean(d.participants),
    decisions: clean(d.decisions),
    topics: d.topics
      .map((t) => ({ title: t.title.trim(), points: clean(t.points) }))
      .filter((t) => t.title || t.points.length),
    action_items: d.action_items
      .map((a) => ({ task: a.task.trim(), owner: a.owner.trim() || "未指定", due: a.due.trim() || "未指定" }))
      .filter((a) => a.task),
  });
}

const participants = ref(draft.value.participants.join("、"));
function setParticipants() {
  draft.value.participants = participants.value.split(/[、,，\s]+/).filter(Boolean);
}
</script>

<template>
  <div class="editor">
    <label class="field">
      <span>標題</span>
      <input v-model="draft.title" />
    </label>
    <label class="field">
      <span>與會者</span>
      <input v-model="participants" placeholder="用頓號或逗號分隔" @change="setParticipants" />
    </label>

    <h3>摘要</h3>
    <div v-for="(_, i) in draft.overview" :key="i" class="item">
      <textarea v-model="draft.overview[i]" rows="2" />
      <button class="icon" title="刪除" @click="draft.overview.splice(i, 1)">✕</button>
    </div>
    <button class="add" @click="draft.overview.push('')">＋ 摘要</button>

    <h3>討論重點</h3>
    <div v-for="(t, i) in draft.topics" :key="i" class="topic">
      <div class="item">
        <input v-model="t.title" class="topic-title" placeholder="主題" />
        <button class="icon" title="刪除主題" @click="draft.topics.splice(i, 1)">✕</button>
      </div>
      <div v-for="(_, j) in t.points" :key="j" class="item indent">
        <input v-model="t.points[j]" />
        <button class="icon" title="刪除" @click="t.points.splice(j, 1)">✕</button>
      </div>
      <button class="add indent" @click="t.points.push('')">＋ 要點</button>
    </div>
    <button class="add" @click="draft.topics.push({ title: '', points: [''] })">＋ 主題</button>

    <h3>決議</h3>
    <div v-for="(_, i) in draft.decisions" :key="i" class="item">
      <input v-model="draft.decisions[i]" />
      <button class="icon" title="刪除" @click="draft.decisions.splice(i, 1)">✕</button>
    </div>
    <button class="add" @click="draft.decisions.push('')">＋ 決議</button>

    <h3>待辦事項</h3>
    <div class="actions-head muted"><span>事項</span><span>負責人</span><span>期限</span><span /></div>
    <div v-for="(a, i) in draft.action_items" :key="i" class="action-row">
      <input v-model="a.task" />
      <input v-model="a.owner" />
      <input v-model="a.due" />
      <button class="icon" title="刪除" @click="draft.action_items.splice(i, 1)">✕</button>
    </div>
    <button class="add" @click="draft.action_items.push({ task: '', owner: '', due: '' })">＋ 待辦</button>

    <div class="buttons">
      <button @click="emit('cancel')">取消</button>
      <button class="primary" @click="setParticipants(); save()">儲存</button>
    </div>
  </div>
</template>

<style scoped>
.field {
  display: grid;
  grid-template-columns: 64px 1fr;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

h3 {
  font-size: 16px;
  margin: 20px 0 8px;
}

.item {
  display: flex;
  gap: 6px;
  margin-bottom: 6px;
}

.item input,
.item textarea {
  flex: 1;
}

textarea {
  resize: vertical;
  font-family: inherit;
}

.indent {
  margin-left: 24px;
}

.topic {
  margin-bottom: 12px;
}

.topic-title {
  font-weight: 600;
}

.icon {
  border: none;
  background: none;
  color: var(--muted);
  padding: 4px 8px;
}

.add {
  border-style: dashed;
  background: none;
  font-size: 13px;
  padding: 3px 10px;
}

.actions-head,
.action-row {
  display: grid;
  grid-template-columns: 1fr 120px 120px 32px;
  gap: 6px;
  margin-bottom: 6px;
}

.actions-head {
  font-size: 13px;
}

.buttons {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 24px;
  position: sticky;
  bottom: 0;
  padding: 12px 0;
  background: var(--bg);
}
</style>
