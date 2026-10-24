<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { ask } from "@tauri-apps/plugin-dialog";
import {
  api,
  onModelDownload,
  onModelsChanged,
  type CaptureSettings,
  type InputDevice,
  type ModelDownload,
  type ModelStatus,
  type Person,
  type UsageSummary,
} from "../lib/api";
import { formatTokens, formatUsd } from "../lib/minutes";

defineProps<{ models: ModelStatus | null }>();
const emit = defineEmits<{ saved: [] }>();

const savedKey = ref<string | null>(null);
const keyDraft = ref("");
const message = ref("");

const capture = ref<CaptureSettings | null>(null);

async function setSystem(enabled: boolean) {
  await api.setCaptureSystem(enabled);
  capture.value = await api.getCaptureSettings();
}

const devices = ref<InputDevice[]>([]);
/** "" means the system default. */
const device = ref("");
const deviceError = ref("");


async function loadDevices() {
  deviceError.value = "";
  try {
    devices.value = await api.listInputDevices();
    device.value = (await api.getInputDevice()) ?? "";
  } catch (e) {
    deviceError.value = String(e);
  }
}

async function saveDevice() {
  await api.setInputDevice(device.value || null);
}

const glossary = ref("");
const glossarySaved = ref("");
const glossaryMessage = ref("");

async function saveGlossary() {
  const words = glossary.value.split("\n").map((w) => w.trim()).filter(Boolean);
  await api.setGlossary(words);
  glossary.value = glossarySaved.value = words.join("\n");
  glossaryMessage.value = `已儲存 ${words.length} 個詞`;
}

const people = ref<Person[]>([]);
const peopleMessage = ref("");

async function loadPeople() {
  people.value = await api.listPeople();
}

async function savePerson(p: Person) {
  await api.setPerson(p.name, p.role);
}

async function renamePerson(p: Person, to: string) {
  to = to.trim();
  if (!to || to === p.name) return;
  const merge = people.value.some((q) => q.name === to);
  if (merge && !(await ask(`已經有「${to}」了。要把「${p.name}」的標註併到「${to}」嗎？`, { title: "合併講者", okLabel: "合併", cancelLabel: "取消" }))) {
    await loadPeople();
    return;
  }
  try {
    await api.renamePerson(p.name, to);
    peopleMessage.value = merge ? `已把「${p.name}」併到「${to}」` : `已改名為「${to}」`;
  } catch (e) {
    peopleMessage.value = String(e);
  }
  await loadPeople();
}

async function deletePerson(p: Person) {
  const ok = await ask(`要刪除「${p.name}」嗎？所有會議裡標成「${p.name}」的句子都會變回未標記，聲紋也會一起刪除。`, {
    title: "刪除講者",
    kind: "warning",
    okLabel: "刪除",
    cancelLabel: "取消",
  });
  if (!ok) return;
  await api.deletePerson(p.name);
  peopleMessage.value = `已刪除「${p.name}」`;
  await loadPeople();
}

const download = ref<ModelDownload | null>(null);
const downloadError = ref("");

async function startDownload() {
  downloadError.value = "";
  try {
    await api.downloadModels();
  } catch (e) {
    downloadError.value = String(e);
  } finally {
    download.value = null;
  }
}

const usage = ref<UsageSummary | null>(null);
const when = (iso?: string) =>
  iso ? new Date(iso).toLocaleString("zh-TW", { month: "numeric", day: "numeric", hour: "2-digit", minute: "2-digit", hour12: false }) : "";

const gb = (bytes: number) => `${(bytes / 1e9).toFixed(2)} GB`;

const unlisteners: (() => void)[] = [];
onMounted(async () => {
  savedKey.value = await api.apiKeyHint();
  glossary.value = glossarySaved.value = (await api.getGlossary()).join("\n");
  usage.value = await api.usageSummary();
  capture.value = await api.getCaptureSettings();
  await loadPeople();
  unlisteners.push(
    await onModelDownload((p) => (download.value = p)),
    await onModelsChanged(() => emit("saved")),
  );
  await loadDevices();
});
onUnmounted(() => unlisteners.forEach((u) => u()));

async function saveKey(key: string | null) {
  message.value = "";
  try {
    await api.setApiKey(key);
    savedKey.value = await api.apiKeyHint();
    keyDraft.value = "";
    message.value = key ? "已儲存到鑰匙圈" : "已移除";
    emit("saved");
  } catch (e) {
    message.value = String(e);
  }
}
</script>

<template>
  <div class="settings">
    <h2>設定</h2>

    <section>
      <h3>Claude API key</h3>
      <p class="muted">
        用來產生會議記錄。只有逐字稿文字會送到 Claude，錄音不會離開這台電腦。key 存在系統鑰匙圈，不會寫進任何檔案。
      </p>
      <p v-if="savedKey">目前：<code>{{ savedKey }}</code></p>
      <p v-else class="error">尚未設定，匯入的錄音只會產生逐字稿。</p>
      <form class="row" @submit.prevent="saveKey(keyDraft)">
        <input v-model="keyDraft" type="password" placeholder="sk-ant-…" autocomplete="off" spellcheck="false" />
        <button class="primary" type="submit" :disabled="!keyDraft.trim()">{{ savedKey ? "更換" : "儲存" }}</button>
        <button v-if="savedKey" type="button" class="danger" @click="saveKey(null)">移除</button>
      </form>
      <p v-if="message" class="muted">{{ message }}</p>
    </section>

    <section v-if="usage">
      <h3>Claude 用量</h3>
      <p class="muted small">依公開價格估算（Claude Haiku 5.5：輸入 US$0.10、輸出 US$0.50 / 百萬 tokens），實際金額以 Anthropic 帳單為準。</p>
      <div class="totals">
        <div>
          <span class="muted">本月</span>
          <strong>{{ formatUsd(usage.month.cost_usd) }}</strong>
          <span class="muted small">{{ usage.month.calls }} 次・{{ formatTokens(usage.month.input_tokens + usage.month.output_tokens) }} tokens</span>
        </div>
        <div>
          <span class="muted">全部</span>
          <strong>{{ formatUsd(usage.all.cost_usd) }}</strong>
          <span class="muted small">{{ usage.all.calls }} 次・{{ formatTokens(usage.all.input_tokens + usage.all.output_tokens) }} tokens</span>
        </div>
      </div>
      <table v-if="usage.recent.length" class="calls">
        <thead>
          <tr><th>時間</th><th>輸入</th><th>輸出</th><th>費用</th><th>耗時</th><th></th></tr>
        </thead>
        <tbody>
          <tr v-for="c in usage.recent" :key="c.id">
            <td>{{ when(c.created_at) }}</td>
            <td>{{ formatTokens(c.input_tokens) }}</td>
            <td>{{ formatTokens(c.output_tokens) }}</td>
            <td>{{ formatUsd(c.cost_usd) }}</td>
            <td>{{ (c.duration_ms / 1000).toFixed(1) }} 秒</td>
            <td :class="c.ok ? 'muted' : 'error'" :title="c.error ?? ''">{{ c.ok ? "成功" : "失敗" }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="muted">還沒有呼叫紀錄。</p>
    </section>

    <section v-if="capture">
      <h3>錄音來源</h3>
      <label class="choice">
        <input type="radio" :checked="!capture.system" @change="setSystem(false)" />
        <span><strong>只有麥克風</strong><br /><span class="muted small">實體會議，或只需要錄自己的聲音</span></span>
      </label>
      <label class="choice" :class="{ disabled: !capture.system_available }">
        <input type="radio" :checked="capture.system" :disabled="!capture.system_available" @change="setSystem(true)" />
        <span>
          <strong>麥克風 + 電腦播放的聲音</strong><br />
          <span class="muted small">線上會議（Google Meet、Slack、Zoom）：同時錄下其他人的聲音，逐字稿會標出「我」和「對方」</span>
        </span>
      </label>
      <ul v-if="capture.system" class="muted small notes">
        <li>開會時請<strong>戴耳機</strong>，不然麥克風會再收到一次對方的聲音，對方的話會重複出現。</li>
        <li>會錄下這台 Mac 播放的所有聲音，開會時請關掉音樂和影片。</li>
        <li>第一次錄音時 macOS 會詢問「系統音訊錄製」權限，請允許。用 <code>dev.sh</code> 啟動的開發版錄不到，要用 <code>bundle-dev.sh</code> 打包的版本。</li>
      </ul>
      <p v-if="!capture.system_available" class="muted small">這個版本不支援錄製電腦播放的聲音。</p>
    </section>

    <section>
      <h3>麥克風</h3>
      <div class="row">
        <select v-model="device" @change="saveDevice">
          <option value="">系統預設{{ devices.find((d) => d.is_default) ? `（${devices.find((d) => d.is_default)!.name}）` : "" }}</option>
          <option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }}</option>
        </select>
        <button type="button" @click="loadDevices">重新整理</button>
      </div>
      <p v-if="device && !devices.some((d) => d.id === device)" class="error">選的麥克風目前沒有連接，錄音時會失敗。</p>
      <p v-if="deviceError" class="error">{{ deviceError }}</p>
      <p class="muted">開會時建議用外接的會議麥克風，收音品質會直接影響逐字稿的準確度。</p>
    </section>

    <section>
      <h3>詞彙表</h3>
      <p class="muted">
        常用的專有名詞、產品名稱、人名，一行一個。辨識時會優先辨識成這些詞，Claude 修正錯字時也會參考。每場會議另外要用的詞，可以在會議頁的「本場詞彙」補。
      </p>
      <textarea v-model="glossary" rows="10" spellcheck="false" placeholder="Hyperf&#10;cert-manager&#10;有效投注" />
      <div class="row end">
        <span class="muted">{{ glossaryMessage }}</span>
        <button class="primary" :disabled="glossary === glossarySaved" @click="saveGlossary">儲存詞彙表</button>
      </div>
      <p class="muted small">辨識模型的提示有長度上限，只會用到最前面約 100 個字的詞；Claude 會拿到完整清單。最重要的詞請放前面。</p>
    </section>

    <section>
      <h3>人員名單</h3>
      <p class="muted">
        在逐字稿標過講者的人會出現在這裡，聲紋就是用那些句子建的。名字和職位會交給 Claude，用來判斷待辦事項的負責人、對應「PM」這類說法；名字請填大家平常的叫法，例如「阿彬」。
      </p>
      <p v-if="!people.length" class="muted">還沒有人。打開一場會議的逐字稿，按「編輯」幫講話的人標幾句。</p>
      <table v-else class="people">
        <thead>
          <tr><th>名字</th><th>職位／角色</th><th>聲紋</th><th /></tr>
        </thead>
        <tbody>
          <tr v-for="p in people" :key="p.name">
            <td><input :value="p.name" @change="renamePerson(p, ($event.target as HTMLInputElement).value)" /></td>
            <td><input v-model="p.role" placeholder="PM、後端…" @change="savePerson(p)" /></td>
            <td class="muted small" :class="{ thin: p.hand_lines < 4 }" :title="p.hand_lines < 4 ? '標的句子太少，多標幾句會認得比較準' : ''">
              {{ p.hand_lines }} 句・{{ Math.round(p.hand_seconds) }} 秒
            </td>
            <td><button class="link danger" title="刪除這個人的標註和聲紋" @click="deletePerson(p)">刪除</button></td>
          </tr>
        </tbody>
      </table>
      <p v-if="peopleMessage" class="muted small">{{ peopleMessage }}</p>
      <p class="muted small">聲紋只存在這台電腦，不會上傳。名字改成另一個已有的名字，會把兩個人合併。</p>
    </section>

    <section>
      <h3>辨識模型</h3>
      <template v-if="models?.ready">
        <p>Breeze ASR 25 已就緒。</p>
        <p class="muted path">{{ models.dir }}</p>
      </template>
      <template v-else-if="models">
        <p>第一次使用需要下載辨識模型，下載完就不需要網路：</p>
        <ul>
          <li v-for="m in models.missing" :key="m.name">{{ m.label }}（{{ gb(m.size) }}）</li>
        </ul>
        <div v-if="download" class="download">
          <span>下載{{ download.label }}…</span>
          <progress :value="download.downloaded" :max="download.total" />
          <span class="muted">{{ gb(download.downloaded) }} / {{ gb(download.total) }}</span>
        </div>
        <button v-else class="primary" :disabled="models.downloading" @click="startDownload">
          {{ models.downloading ? "下載中…" : "下載模型" }}
        </button>
        <p v-if="downloadError" class="error">{{ downloadError }}</p>
        <p class="muted small path">存放位置：{{ models.dir }}</p>
      </template>
    </section>
  </div>
</template>

<style scoped>
.settings {
  max-width: 640px;
  margin: 0 auto;
  padding: 24px 32px 48px;
}

h2 {
  font-size: 22px;
  margin: 0 0 8px;
}

h3 {
  font-size: 16px;
  margin: 28px 0 6px;
}

.row {
  display: flex;
  gap: 8px;
}

.row input,
.row select {
  flex: 1;
}

select {
  font: inherit;
  color: inherit;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 6px 10px;
}

code {
  font-size: 13px;
}

textarea {
  width: 100%;
  resize: vertical;
  font-family: inherit;
}

.row.end {
  justify-content: flex-end;
  align-items: center;
  margin-top: 8px;
}

.small {
  font-size: 12px;
}

.totals {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin: 8px 0 12px;
}

.totals div {
  display: flex;
  flex-direction: column;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.totals strong {
  font-size: 20px;
  font-variant-numeric: tabular-nums;
}

.calls {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  font-variant-numeric: tabular-nums;
}

.calls th,
.calls td {
  text-align: left;
  padding: 4px 6px;
  border-bottom: 1px solid var(--border);
}

.calls th {
  color: var(--muted);
  font-weight: 500;
}

.choice {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  padding: 8px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  margin-bottom: 8px;
  cursor: pointer;
}

.choice input {
  margin-top: 4px;
}

.choice.disabled {
  opacity: 0.5;
  cursor: default;
}

.notes {
  margin: 4px 0 0;
  padding-left: 20px;
}

.people {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.people th {
  text-align: left;
  font-weight: 600;
  color: var(--muted);
  padding: 4px;
}

.people td {
  padding: 3px 4px;
}

.people input {
  width: 100%;
  font-size: 13px;
  padding: 3px 6px;
}

.people .thin {
  color: var(--danger);
}

.people .link {
  border: none;
  background: none;
  padding: 0 4px;
  font-size: 12px;
}

.download {
  display: flex;
  align-items: center;
  gap: 10px;
}

.download progress {
  flex: 1;
}

.path {
  word-break: break-all;
  font-size: 12px;
}
</style>
