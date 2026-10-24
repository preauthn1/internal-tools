<script lang="ts" setup>
import type { IntentExampleApi } from '#/api/intent-examples';
import type { PendingQuestionApi } from '#/api/pending-questions';

import { computed, onMounted, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';
import { preferences } from '@vben/preferences';

import {
  AutoComplete,
  Button,
  Card,
  Descriptions,
  DescriptionsItem,
  Drawer,
  Form,
  FormItem,
  Input,
  message,
  Popconfirm,
  Segmented,
  Select,
  Table,
  TabPane,
  Tabs,
  Tag,
  Textarea,
} from 'ant-design-vue';

import { getMerchantApi } from '#/api/core';
import { getIntentsApi } from '#/api/intent-examples';
import {
  getPendingQuestionApi,
  getPendingQuestionsApi,
  ignorePendingApi,
  pendingToExampleApi,
  pendingToFAQApi,
  pendingToKnowledgeApi,
} from '#/api/pending-questions';
import { $t } from '#/locales';

defineOptions({ name: 'PendingQuestions' });

const status = ref<PendingQuestionApi.Status>('open');
const list = ref<PendingQuestionApi.Question[]>([]);
const loading = ref(false);
const detail = ref<PendingQuestionApi.Detail>();
const open = ref(false);
const saving = ref(false);
const intents = ref<IntentExampleApi.Intent[]>([]);
const sectionTitles = ref<string[]>([]);

// 四種動作的表單（打開詳細頁時預填：問題是顧客原話，答案是員工在群組的回覆）
const faq = ref({ alt: '', answer: '', question: '' });
const knowledge = ref({ body: '', title: '' });
const example = ref({ intent: undefined as string | undefined, text: '' });

const statusOptions = computed(() => [
  { label: $t('page.pendingQuestions.open'), value: 'open' },
  { label: $t('page.pendingQuestions.ignored'), value: 'ignored' },
  { label: $t('page.pendingQuestions.resolved'), value: 'resolved' },
]);

const columns = computed(() => [
  { dataIndex: 'count', key: 'count', title: $t('page.pendingQuestions.count'), width: 80 },
  { key: 'text', title: $t('page.pendingQuestions.text') },
  { key: 'last_seen', title: $t('page.pendingQuestions.lastSeen'), width: 180 },
  { key: 'action', width: 100 },
]);

/** 意圖的白話名稱：依目前介面語言，沒有就用繁中（不顯示代號） */
function intentLabel(code: string) {
  const labels = intents.value.find((i) => i.code === code)?.labels;
  return labels?.[example.invalid.locale] || labels?.['zh-TW'] || code;
}

const intentOptions = computed(() =>
  intents.value.map((i) => ({ label: intentLabel(i.code), value: i.code })),
);

const sectionOptions = computed(() =>
  sectionTitles.value.map((t) => ({ value: t })),
);

function formatTime(s: string) {
  return new Date(s).toLocaleString();
}

function sourceLabel(source: string) {
  return $t(`page.pendingQuestions.sources.${source}`);
}

async function load() {
  loading.value = true;
  try {
    list.value = await getPendingQuestionsApi(status.value);
  } finally {
    loading.value = false;
  }
}

async function show(id: number) {
  const d = await getPendingQuestionApi(id);
  detail.value = d;
  faq.value = { alt: '', answer: d.staff_answer, question: d.text };
  knowledge.value = { body: d.staff_answer, title: '' };
  example.value = { intent: undefined, text: d.text };
  open.value = true;
}

// 動作成功後關掉詳細頁、重新讀列表；錯誤訊息（驗證沒過、和題庫重複）由 request 的攔截器顯示
async function act(run: () => Promise<unknown>, done: string) {
  saving.value = true;
  try {
    await run();
    message.success(done);
    open.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

function toFAQ(id: number) {
  const alt = faq.value.alt
    .split('\n')
    .map((s) => s.trim())
    .filter(Boolean);
  return act(
    () =>
      pendingToFAQApi(id, {
        alt_questions: alt,
        answer: faq.value.answer,
        question: faq.value.question,
      }),
    $t('page.pendingQuestions.toFAQDone'),
  );
}

function toKnowledge(id: number) {
  return act(
    () => pendingToKnowledgeApi(id, knowledge.value),
    $t('page.pendingQuestions.toKnowledgeDone'),
  );
}

function toExample(id: number) {
  const { intent, text } = example.value;
  if (!intent) {
    return;
  }
  return act(
    () => pendingToExampleApi(id, { intent, text }),
    $t('page.pendingQuestions.toExampleDone'),
  );
}

function ignore(id: number) {
  return act(() => ignorePendingApi(id), $t('page.pendingQuestions.ignoreDone'));
}

watch(status, load);

onMounted(async () => {
  await load();
  const [i, m] = await Promise.all([getIntentsApi(), getMerchantApi()]);
  intents.value = i;
  sectionTitles.value = m.knowledge.map((s) => s.title);
});
</script>

<template>
  <Page
    :description="$t('page.pendingQuestions.description')"
    :title="$t('page.pendingQuestions.title')"
  >
    <Card>
      <Segmented
        v-model:value="status"
        :options="statusOptions"
        class="mb-4"
      />
      <Table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'text'">
            <div>{{ record.text }}</div>
            <div v-if="record.translated" class="text-gray-500">
              {{ record.translated }}
            </div>
            <Tag v-if="record.staff_answer" color="green">
              {{ $t('page.pendingQuestions.hasStaffAnswer') }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'last_seen'">
            {{ formatTime(record.last_seen) }}
          </template>
          <template v-else-if="column.key === 'action'">
            <Button size="small" type="link" @click="show(record.id)">
              {{ $t('page.pendingQuestions.view') }}
            </Button>
          </template>
        </template>
      </Table>
    </Card>

    <Drawer
      v-model:open="open"
      :title="$t('page.pendingQuestions.detail')"
      width="720"
    >
      <template v-if="detail">
        <Descriptions :column="1" bordered size="small">
          <DescriptionsItem :label="$t('page.pendingQuestions.text')">
            {{ detail.text }}
          </DescriptionsItem>
          <DescriptionsItem
            v-if="detail.translated"
            :label="$t('page.pendingQuestions.translated')"
          >
            {{ detail.translated }}
          </DescriptionsItem>
          <DescriptionsItem :label="$t('page.pendingQuestions.count')">
            {{ detail.count }}
          </DescriptionsItem>
          <DescriptionsItem
            v-if="detail.staff_answer"
            :label="$t('page.pendingQuestions.staffAnswer')"
          >
            <div class="whitespace-pre-wrap">{{ detail.staff_answer }}</div>
          </DescriptionsItem>
        </Descriptions>

        <Tabs v-if="detail.status !== 'resolved'" class="mt-4">
          <TabPane key="faq" :tab="$t('page.pendingQuestions.toFAQ')">
            <Form layout="vertical">
              <FormItem :label="$t('page.settings.faq.question')" required>
                <Input v-model:value="faq.question" />
              </FormItem>
              <FormItem :label="$t('page.settings.faq.altQuestions')">
                <Textarea v-model:value="faq.alt" :rows="2" />
              </FormItem>
              <FormItem :label="$t('page.settings.faq.answer')" required>
                <Textarea v-model:value="faq.answer" :rows="3" />
              </FormItem>
              <Button
                :disabled="!faq.question.trim() || !faq.answer.trim()"
                :loading="saving"
                type="primary"
                @click="toFAQ(detail.id)"
              >
                {{ $t('page.pendingQuestions.toFAQ') }}
              </Button>
            </Form>
          </TabPane>
          <TabPane
            key="knowledge"
            :tab="$t('page.pendingQuestions.toKnowledge')"
          >
            <Form layout="vertical">
              <FormItem
                :extra="$t('page.pendingQuestions.sectionHelp')"
                :label="$t('page.settings.merchant.sectionTitle')"
                required
              >
                <AutoComplete
                  v-model:value="knowledge.title"
                  :options="sectionOptions"
                />
              </FormItem>
              <FormItem :label="$t('page.settings.merchant.sectionBody')" required>
                <Textarea v-model:value="knowledge.body" :rows="3" />
              </FormItem>
              <Button
                :disabled="!knowledge.title.trim() || !knowledge.body.trim()"
                :loading="saving"
                type="primary"
                @click="toKnowledge(detail.id)"
              >
                {{ $t('page.pendingQuestions.toKnowledge') }}
              </Button>
            </Form>
          </TabPane>
          <TabPane key="example" :tab="$t('page.pendingQuestions.toExample')">
            <Form layout="vertical">
              <FormItem :label="$t('page.intentExamples.text')" required>
                <Input v-model:value="example.text" :maxlength="200" />
              </FormItem>
              <FormItem :label="$t('page.intentExamples.intent')" required>
                <Select
                  v-model:value="example.intent"
                  :options="intentOptions"
                  :placeholder="$t('page.intentExamples.intentPlaceholder')"
                />
              </FormItem>
              <Button
                :disabled="!example.text.trim() || !example.intent"
                :loading="saving"
                type="primary"
                @click="toExample(detail.id)"
              >
                {{ $t('page.pendingQuestions.toExample') }}
              </Button>
            </Form>
          </TabPane>
          <TabPane
            v-if="detail.status === 'open'"
            key="ignore"
            :tab="$t('page.pendingQuestions.ignore')"
          >
            <p class="mb-4">{{ $t('page.pendingQuestions.ignoreHelp') }}</p>
            <Popconfirm
              :title="$t('page.pendingQuestions.confirmIgnore')"
              @confirm="ignore(detail.id)"
            >
              <Button :loading="saving" danger>
                {{ $t('page.pendingQuestions.ignore') }}
              </Button>
            </Popconfirm>
          </TabPane>
        </Tabs>

        <h3 class="mb-2 mt-6 font-bold">
          {{ $t('page.pendingQuestions.hits') }}
        </h3>
        <Card
          v-for="(hit, i) in detail.hits"
          :key="i"
          class="mb-3"
          size="small"
        >
          <template #title>
            {{ formatTime(hit.time) }}・{{ hit.user_id }}
            <Tag class="ml-2">{{ sourceLabel(hit.source) }}</Tag>
          </template>
          <p v-if="hit.conversation.length === 0" class="text-gray-500">
            {{ $t('page.pendingQuestions.noConversation') }}
          </p>
          <div
            v-for="(line, j) in hit.conversation"
            :key="j"
            :class="line.direction === 'customer' ? '' : 'text-right'"
            class="mb-1"
          >
            <Tag>{{ $t(`page.pendingQuestions.directions.${line.direction}`) }}</Tag>
            <span class="whitespace-pre-wrap">{{ line.text }}</span>
          </div>
        </Card>
      </template>
    </Drawer>
  </Page>
</template>
