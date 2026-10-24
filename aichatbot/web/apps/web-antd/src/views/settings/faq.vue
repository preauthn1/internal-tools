<script lang="ts" setup>
import type { SettingsApi } from '#/api';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import {
  Alert,
  Button,
  Card,
  Form,
  FormItem,
  Input,
  message,
  Spin,
  Textarea,
} from 'ant-design-vue';

import { getFAQApi, putFAQApi } from '#/api';
import { $t } from '#/locales';

import { move } from './list';
import { warningText } from './warnings';

defineOptions({ name: 'SettingsFAQ' });

/** 畫面上的一題：其他問法用多行文字編輯，一行一個 */
interface Item {
  alt: string;
  answer: string;
  key: string;
  question: string;
}

const loading = ref(true);
const saving = ref(false);
const form = ref<{ items: Item[] }>({ items: [] });
/** 長度警告（#65）：存了之後跟著伺服器的回應更新 */
const warnings = ref<SettingsApi.Warning[]>([]);

function toItems(faq: SettingsApi.FAQ[]): Item[] {
  return faq.map((f) => ({
    alt: f.alt_questions.join('\n'),
    answer: f.answer,
    key: f.key,
    question: f.question,
  }));
}

onMounted(async () => {
  try {
    const r = await getFAQApi();
    form.value.items = toItems(r.faq);
    warnings.value = r.warnings;
  } finally {
    loading.value = false;
  }
});

/** 新的一題代號留空，由伺服器產生；商家看不到代號 */
function add() {
  form.value.items.push({ alt: '', answer: '', key: '', question: '' });
}

async function save() {
  saving.value = true;
  try {
    const r = await putFAQApi(
      form.value.items.map((it) => ({
        alt_questions: it.alt
          .split('\n')
          .map((s) => s.trim())
          .filter(Boolean),
        answer: it.answer,
        key: it.key,
        question: it.question,
      })),
    );
    form.value.items = toItems(r.faq);
    warnings.value = r.warnings;
    message.success($t('page.settings.saved'));
  } finally {
    saving.value = false;
  }
}

const required = [{ message: $t('page.settings.required'), required: true }];
</script>

<template>
  <Page
    :title="$t('page.settings.faq.title')"
    :description="$t('page.settings.faq.help')"
  >
    <Spin :spinning="loading">
      <Form :model="form" layout="vertical" class="max-w-3xl" @finish="save">
        <Alert
          v-for="(w, i) in warnings"
          :key="i"
          type="warning"
          show-icon
          class="mb-3"
          :message="warningText(w, $t)"
        />
        <p v-if="form.items.length === 0" class="mb-3 text-foreground/60">
          {{ $t('page.settings.faq.empty') }}
        </p>
        <Card
          v-for="(item, i) in form.items"
          :key="i"
          size="small"
          class="mb-3"
          :title="$t('page.settings.faq.item', { n: i + 1 })"
        >
          <FormItem
            :label="$t('page.settings.faq.question')"
            :name="['items', i, 'question']"
            :rules="required"
          >
            <Input v-model:value="item.question" />
          </FormItem>
          <FormItem :label="$t('page.settings.faq.altQuestions')">
            <Textarea v-model:value="item.alt" :auto-size="{ minRows: 1 }" />
          </FormItem>
          <FormItem
            :label="$t('page.settings.faq.answer')"
            :name="['items', i, 'answer']"
            :rules="required"
          >
            <Textarea v-model:value="item.answer" :auto-size="{ minRows: 2 }" />
          </FormItem>
          <div class="flex flex-wrap gap-2">
            <Button
              size="small"
              :disabled="i === 0"
              @click="move(form.items, i, -1)"
            >
              {{ $t('page.settings.moveUp') }}
            </Button>
            <Button
              size="small"
              :disabled="i === form.items.length - 1"
              @click="move(form.items, i, 1)"
            >
              {{ $t('page.settings.moveDown') }}
            </Button>
            <Button size="small" danger @click="form.items.splice(i, 1)">
              {{ $t('page.settings.remove') }}
            </Button>
          </div>
        </Card>
        <Button block class="mb-6" @click="add">
          {{ $t('page.settings.faq.add') }}
        </Button>

        <Button type="primary" html-type="submit" block :loading="saving">
          {{ $t('page.settings.save') }}
        </Button>
      </Form>
    </Spin>
  </Page>
</template>
