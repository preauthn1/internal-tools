<script lang="ts" setup>
import type { IntentExampleApi } from '#/api/intent-examples';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { preferences } from '@vben/preferences';

import {
  Button,
  Card,
  Input,
  message,
  Popconfirm,
  Select,
  Space,
  Table,
  Tag,
} from 'ant-design-vue';

import {
  addIntentExampleApi,
  deleteIntentExampleApi,
  getIntentExamplesApi,
  getIntentsApi,
} from '#/api/intent-examples';
import { $t } from '#/locales';

defineOptions({ name: 'IntentExamples' });

const intents = ref<IntentExampleApi.Intent[]>([]);
const examples = ref<IntentExampleApi.Example[]>([]);
const loading = ref(false);
const saving = ref(false);
const text = ref('');
const intent = ref<string>();
// 篩選的意圖；空字串是全部
const filter = ref('');

/** 意圖的白話名稱：依目前介面語言，沒有就用繁中（不顯示代號） */
function intentLabel(code: string) {
  const labels = intents.value.find((i) => i.code === code)?.labels;
  return labels?.[example.invalid.locale] || labels?.['zh-TW'] || code;
}

const intentOptions = computed(() =>
  intents.value.map((i) => ({ label: intentLabel(i.code), value: i.code })),
);

const filterOptions = computed(() => [
  { label: $t('page.intentExamples.filterAll'), value: '' },
  ...intentOptions.value,
]);

const rows = computed(() =>
  filter.value
    ? examples.value.filter((e) => e.intent === filter.value)
    : examples.value,
);

const columns = computed(() => [
  { dataIndex: 'text', key: 'text', title: $t('page.intentExamples.text') },
  {
    key: 'intent',
    title: $t('page.intentExamples.intent'),
    width: 220,
  },
  {
    key: 'source',
    title: $t('page.intentExamples.source'),
    width: 120,
  },
  { key: 'action', width: 100 },
]);

async function load() {
  loading.value = true;
  try {
    [intents.value, examples.value] = await Promise.all([
      getIntentsApi(),
      getIntentExamplesApi(),
    ]);
  } finally {
    loading.value = false;
  }
}

// 錯誤訊息（和題庫重複、句子不正確）由 request 的攔截器顯示後端回的 error
async function add() {
  if (!text.value.trim() || !intent.value) {
    return;
  }
  saving.value = true;
  try {
    const added = await addIntentExampleApi(text.value, intent.value);
    examples.value = [...examples.value, added];
    text.value = '';
    message.success($t('page.intentExamples.added'));
  } finally {
    saving.value = false;
  }
}

async function remove(id: number) {
  await deleteIntentExampleApi(id);
  examples.value = examples.value.filter((e) => e.id !== id);
  message.success($t('page.intentExamples.deleted'));
}

onMounted(load);
</script>

<template>
  <Page
    :description="$t('page.intentExamples.description')"
    :title="$t('page.intentExamples.title')"
  >
    <Card class="mb-4">
      <Space wrap>
        <Input
          v-model:value="text"
          :maxlength="200"
          :placeholder="$t('page.intentExamples.placeholder')"
          style="width: 360px"
          @press-enter="add"
        />
        <Select
          v-model:value="intent"
          :options="intentOptions"
          :placeholder="$t('page.intentExamples.intentPlaceholder')"
          style="width: 220px"
        />
        <Button
          :disabled="!text.trim() || !intent"
          :loading="saving"
          type="primary"
          @click="add"
        >
          {{ $t('page.intentExamples.add') }}
        </Button>
      </Space>
    </Card>
    <Card>
      <Space class="mb-4">
        <span>{{ $t('page.intentExamples.filter') }}</span>
        <Select
          v-model:value="filter"
          :options="filterOptions"
          style="width: 220px"
        />
      </Space>
      <Table
        :columns="columns"
        :data-source="rows"
        :loading="loading"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'intent'">
            {{ intentLabel(record.intent) }}
          </template>
          <template v-else-if="column.key === 'source'">
            <Tag v-if="record.builtin">
              {{ $t('page.intentExamples.builtin') }}
            </Tag>
            <Tag v-else color="blue">
              {{ $t('page.intentExamples.custom') }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <Popconfirm
              v-if="!record.builtin"
              :title="$t('page.intentExamples.confirmDelete')"
              @confirm="remove(record.id)"
            >
              <Button danger size="small" type="link">
                {{ $t('page.intentExamples.delete') }}
              </Button>
            </Popconfirm>
          </template>
        </template>
      </Table>
    </Card>
  </Page>
</template>
