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

import { getMerchantApi, putMerchantApi } from '#/api';
import { $t } from '#/locales';

import { move } from './list';
import { splitWarnings, warningText } from './warnings';

defineOptions({ name: 'SettingsMerchant' });

const loading = ref(true);
const saving = ref(false);
const form = ref<SettingsApi.Merchant>({
  contact: '',
  intro: '',
  knowledge: [],
  name: '',
});
/** 知識段落的長度警告（#65）：和表單分開放，存的時候不送回去 */
const warnings = ref<SettingsApi.Warning[]>([]);

onMounted(async () => {
  try {
    [form.value, warnings.value] = splitWarnings(await getMerchantApi());
  } finally {
    loading.value = false;
  }
});

function addSection() {
  form.value.knowledge.push({ body: '', title: '' });
}

/** 前端只做必填提示；其他檢查（標題重複等）由伺服器回錯誤訊息，request.ts 會顯示 */
async function save() {
  saving.value = true;
  try {
    [form.value, warnings.value] = splitWarnings(
      await putMerchantApi(form.value),
    );
    message.success($t('page.settings.saved'));
  } finally {
    saving.value = false;
  }
}

const required = [{ message: $t('page.settings.required'), required: true }];
</script>

<template>
  <Page :title="$t('page.settings.merchant.title')">
    <Spin :spinning="loading">
      <Form :model="form" layout="vertical" class="max-w-3xl" @finish="save">
        <FormItem
          :label="$t('page.settings.merchant.name')"
          name="name"
          :rules="required"
        >
          <Input v-model:value="form.name" />
        </FormItem>
        <FormItem
          :label="$t('page.settings.merchant.intro')"
          name="intro"
          :rules="required"
        >
          <Textarea v-model:value="form.intro" :auto-size="{ minRows: 2 }" />
        </FormItem>
        <FormItem
          :label="$t('page.settings.merchant.contact')"
          name="contact"
          :rules="required"
        >
          <Textarea v-model:value="form.contact" :auto-size="{ minRows: 2 }" />
        </FormItem>

        <h3 class="mb-1 mt-6 text-base font-semibold">
          {{ $t('page.settings.merchant.knowledge') }}
        </h3>
        <p class="mb-3 text-sm text-foreground/60">
          {{ $t('page.settings.merchant.knowledgeHelp') }}
        </p>
        <Alert
          v-for="(w, i) in warnings"
          :key="i"
          type="warning"
          show-icon
          class="mb-3"
          :message="warningText(w, $t)"
        />
        <p v-if="form.knowledge.length === 0" class="mb-3 text-foreground/60">
          {{ $t('page.settings.merchant.empty') }}
        </p>
        <Card
          v-for="(section, i) in form.knowledge"
          :key="i"
          size="small"
          class="mb-3"
        >
          <FormItem
            :label="$t('page.settings.merchant.sectionTitle')"
            :name="['knowledge', i, 'title']"
            :rules="required"
          >
            <Input v-model:value="section.title" />
          </FormItem>
          <FormItem
            :label="$t('page.settings.merchant.sectionBody')"
            :name="['knowledge', i, 'body']"
            :rules="required"
          >
            <Textarea
              v-model:value="section.body"
              :auto-size="{ minRows: 3 }"
            />
          </FormItem>
          <div class="flex flex-wrap gap-2">
            <Button
              size="small"
              :disabled="i === 0"
              @click="move(form.knowledge, i, -1)"
            >
              {{ $t('page.settings.moveUp') }}
            </Button>
            <Button
              size="small"
              :disabled="i === form.knowledge.length - 1"
              @click="move(form.knowledge, i, 1)"
            >
              {{ $t('page.settings.moveDown') }}
            </Button>
            <Button size="small" danger @click="form.knowledge.splice(i, 1)">
              {{ $t('page.settings.remove') }}
            </Button>
          </div>
        </Card>
        <Button block class="mb-6" @click="addSection">
          {{ $t('page.settings.merchant.addSection') }}
        </Button>

        <Button type="primary" html-type="submit" block :loading="saving">
          {{ $t('page.settings.save') }}
        </Button>
      </Form>
    </Spin>
  </Page>
</template>
