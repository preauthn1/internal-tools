<script lang="ts" setup>
import type { CustomIntentApi } from '#/api/custom-intents';
import type { IntentExampleApi } from '#/api/intent-examples';

import { computed, onMounted, onUnmounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { preferences } from '@vben/preferences';

import {
  Alert,
  Button,
  Card,
  Form,
  FormItem,
  Input,
  message,
  Progress,
  Radio,
  RadioGroup,
  Select,
  Spin,
  Table,
  Textarea,
} from 'ant-design-vue';

import {
  getCustomIntentsApi,
  getSelfTestApi,
  putCustomIntentsApi,
  runSelfTestApi,
} from '#/api/custom-intents';
import { getIntentsApi } from '#/api/intent-examples';
import { $t } from '#/locales';

defineOptions({ name: 'CustomIntents' });

/** 最多幾個自訂意圖、每個至少幾句例句（和伺服器的檢查相同） */
const MAX = 10;
const MIN_EXAMPLES = 3;

/** 畫面上的一個自訂意圖：例句用多行文字編輯，一行一句 */
interface Item {
  action: CustomIntentApi.Action;
  code: string;
  description: string;
  examples: string;
  knowledge: string[];
  name: string;
  text: string;
}

const loading = ref(true);
const saving = ref(false);
const form = ref<{ items: Item[] }>({ items: [] });
const titles = ref<string[]>([]);
const intents = ref<IntentExampleApi.Intent[]>([]);
const selfTest = ref<CustomIntentApi.SelfTest>();
let timer: ReturnType<typeof setInterval> | undefined;

function toItems(list: CustomIntentApi.CustomIntent[]): Item[] {
  return list.map((ci) => ({ ...ci, examples: ci.examples.join('\n') }));
}

function lines(s: string) {
  return s
    .split('\n')
    .map((x) => x.trim())
    .filter(Boolean);
}

/** 意圖的白話名稱：依目前介面語言，沒有就用繁中（不顯示代號） */
function intentLabel(code: string) {
  const labels = intents.value.find((i) => i.code === code)?.labels;
  return labels?.[example.invalid.locale] || labels?.['zh-TW'] || code;
}

const titleOptions = computed(() =>
  titles.value.map((t) => ({ label: t, value: t })),
);

/** 執行中每 2 秒問一次進度，跑完就停 */
async function refreshSelfTest() {
  selfTest.value = await getSelfTestApi();
  if (selfTest.value.status === 'running') {
    timer ??= setInterval(refreshSelfTest, 2000);
  } else if (timer) {
    clearInterval(timer);
    timer = undefined;
    intents.value = await getIntentsApi();
  }
}

onMounted(async () => {
  try {
    const data = await getCustomIntentsApi();
    form.value.items = toItems(data.custom_intents);
    titles.value = data.knowledge_titles;
    intents.value = await getIntentsApi();
    await refreshSelfTest();
  } finally {
    loading.value = false;
  }
});

onUnmounted(() => clearInterval(timer));

/** 新的意圖代號留空，由伺服器產生 */
function add() {
  form.value.items.push({
    action: 'fixed_text',
    code: '',
    description: '',
    examples: '',
    knowledge: [],
    name: '',
    text: '',
  });
}

async function save() {
  saving.value = true;
  try {
    const data = await putCustomIntentsApi(
      form.value.items.map((it) => ({
        ...it,
        examples: lines(it.examples),
        // 只有「只根據指定知識回答」用得到段落
        knowledge: it.action === 'knowledge' ? it.knowledge : [],
      })),
    );
    form.value.items = toItems(data.custom_intents);
    intents.value = await getIntentsApi();
    message.success($t('page.settings.saved'));
    await refreshSelfTest();
  } finally {
    saving.value = false;
  }
}

async function rerun() {
  selfTest.value = await runSelfTestApi();
  await refreshSelfTest();
}

const required = [{ message: $t('page.settings.required'), required: true }];

const examplesRule = [
  {
    validator: (_: unknown, v: string) =>
      lines(v ?? '').length >= MIN_EXAMPLES
        ? Promise.resolve()
        : Promise.reject(
            new Error(
              $t('page.customIntents.examplesMin', { n: MIN_EXAMPLES }),
            ),
          ),
  },
];

/** 同一句可能出現在不同狀態（例如預約中問姓名時），key 要帶上狀態 */
function itemKey(item: CustomIntentApi.SelfTestItem) {
  return `${item.state ?? ''}|${item.text}`;
}

const percent = computed(() => {
  const s = selfTest.value;
  return s && s.total > 0 ? Math.floor((s.done / s.total) * 100) : 0;
});

const columns = computed(() => [
  { dataIndex: 'text', key: 'text', title: $t('page.customIntents.text') },
  {
    key: 'expected',
    title: $t('page.customIntents.expected'),
    width: 220,
  },
  { key: 'got', title: $t('page.customIntents.got'), width: 220 },
]);
</script>

<template>
  <Page
    :title="$t('page.customIntents.title')"
    :description="$t('page.customIntents.help')"
  >
    <Spin :spinning="loading">
      <div class="max-w-3xl">
        <Form :model="form" layout="vertical" @finish="save">
          <p v-if="form.items.length === 0" class="mb-3 text-foreground/60">
            {{ $t('page.customIntents.empty') }}
          </p>
          <Card
            v-for="(item, i) in form.items"
            :key="i"
            size="small"
            class="mb-3"
            :title="item.name || $t('page.customIntents.item', { n: i + 1 })"
          >
            <FormItem
              :label="$t('page.customIntents.name')"
              :name="['items', i, 'name']"
              :rules="required"
            >
              <Input v-model:value="item.name" />
            </FormItem>
            <FormItem
              :label="$t('page.customIntents.description')"
              :extra="$t('page.customIntents.descriptionHelp')"
              :name="['items', i, 'description']"
              :rules="required"
            >
              <Textarea
                v-model:value="item.description"
                :auto-size="{ minRows: 1 }"
              />
            </FormItem>
            <FormItem
              :label="$t('page.customIntents.examples')"
              :extra="
                $t('page.customIntents.examplesHelp', { n: MIN_EXAMPLES })
              "
              :name="['items', i, 'examples']"
              :rules="examplesRule"
            >
              <Textarea
                v-model:value="item.examples"
                :auto-size="{ minRows: 3 }"
              />
            </FormItem>
            <FormItem :label="$t('page.customIntents.action')">
              <RadioGroup v-model:value="item.action">
                <Radio value="fixed_text">
                  {{ $t('page.customIntents.fixedText') }}
                </Radio>
                <Radio value="handoff">
                  {{ $t('page.customIntents.handoff') }}
                </Radio>
                <Radio value="knowledge">
                  {{ $t('page.customIntents.knowledge') }}
                </Radio>
              </RadioGroup>
            </FormItem>
            <FormItem
              v-if="item.action === 'fixed_text'"
              :label="$t('page.customIntents.text')"
              :extra="$t('page.customIntents.textHelp')"
              :name="['items', i, 'text']"
              :rules="required"
            >
              <Textarea v-model:value="item.text" :auto-size="{ minRows: 2 }" />
            </FormItem>
            <FormItem
              v-else-if="item.action === 'handoff'"
              :label="$t('page.customIntents.prompt')"
              :extra="$t('page.customIntents.promptHelp')"
            >
              <Textarea v-model:value="item.text" :auto-size="{ minRows: 1 }" />
            </FormItem>
            <FormItem
              v-else
              :label="$t('page.customIntents.sections')"
              :name="['items', i, 'knowledge']"
              :rules="required"
            >
              <Select
                v-model:value="item.knowledge"
                mode="multiple"
                :options="titleOptions"
                :placeholder="$t('page.customIntents.sectionsPlaceholder')"
              />
            </FormItem>
            <Button size="small" danger @click="form.items.splice(i, 1)">
              {{ $t('page.settings.remove') }}
            </Button>
          </Card>
          <Button
            block
            class="mb-6"
            :disabled="form.items.length >= MAX"
            @click="add"
          >
            {{
              form.items.length >= MAX
                ? $t('page.customIntents.full', { n: MAX })
                : $t('page.customIntents.add')
            }}
          </Button>

          <Button type="primary" html-type="submit" block :loading="saving">
            {{ $t('page.settings.save') }}
          </Button>
        </Form>

        <Card class="mt-6" :title="$t('page.customIntents.selfTest.title')">
          <p class="mb-3 text-foreground/60">
            {{ $t('page.customIntents.selfTest.help') }}
          </p>
          <template v-if="selfTest?.status === 'running'">
            <p>{{ $t('page.customIntents.selfTest.running') }}</p>
            <Progress :percent="percent" />
          </template>
          <p v-else-if="selfTest?.status === 'idle'">
            {{ $t('page.customIntents.selfTest.idle') }}
          </p>
          <Alert
            v-else-if="selfTest?.status === 'failed'"
            type="error"
            :message="$t('page.customIntents.selfTest.failed')"
          />
          <template v-else-if="selfTest?.report">
            <Alert
              class="mb-3"
              :type="
                selfTest.report.conflicts.length > 0 ? 'warning' : 'success'
              "
              :message="
                $t('page.customIntents.selfTest.summary', {
                  total: selfTest.report.total,
                  conflicts: selfTest.report.conflicts.length,
                  failed: selfTest.report.failed.length,
                })
              "
            />
            <Table
              v-if="selfTest.report.conflicts.length > 0"
              :columns="columns"
              :data-source="selfTest.report.conflicts"
              :pagination="false"
              :row-key="itemKey"
              size="small"
              class="mb-3"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'expected'">
                  {{ intentLabel(record.expected) }}
                </template>
                <template v-else-if="column.key === 'got'">
                  {{ intentLabel(record.got) }}
                </template>
              </template>
            </Table>
            <template v-if="selfTest.report.failed.length > 0">
              <p class="mb-2">
                {{ $t('page.customIntents.selfTest.failedItems') }}
              </p>
              <ul class="mb-3 list-disc pl-6">
                <li v-for="f in selfTest.report.failed" :key="itemKey(f)">
                  {{ f.text }}
                </li>
              </ul>
            </template>
          </template>
          <Button :disabled="selfTest?.status === 'running'" @click="rerun">
            {{ $t('page.customIntents.selfTest.rerun') }}
          </Button>
        </Card>
      </div>
    </Spin>
  </Page>
</template>
