<script lang="ts" setup>
import type { SettingsApi } from '#/api';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import {
  Alert,
  Button,
  Form,
  FormItem,
  InputNumber,
  message,
  Radio,
  RadioGroup,
  Spin,
} from 'ant-design-vue';

import { getModelApi, putModelApi } from '#/api';
import { $t } from '#/locales';

defineOptions({ name: 'SettingsModel' });

const loading = ref(true);
const saving = ref(false);
const haikuAvailable = ref(false);
const form = ref<{
  model: SettingsApi.ModelSettings['model'];
  monthly_budget_usd: number;
}>({ model: 'afm', monthly_budget_usd: 0 });

function apply(data: SettingsApi.ModelSettings) {
  haikuAvailable.value = data.haiku_available;
  form.value = { model: data.model, monthly_budget_usd: data.monthly_budget_usd };
}

onMounted(async () => {
  try {
    apply(await getModelApi());
  } finally {
    loading.value = false;
  }
});

async function save() {
  saving.value = true;
  try {
    apply(await putModelApi({ ...form.value }));
    message.success($t('page.settings.saved'));
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <Page :title="$t('page.settings.model.title')">
    <Spin :spinning="loading">
      <Form :model="form" layout="vertical" class="max-w-3xl" @finish="save">
        <FormItem :label="$t('page.settings.model.model')">
          <RadioGroup v-model:value="form.model">
            <div class="flex flex-col gap-2">
              <Radio value="afm">{{ $t('page.settings.model.afm') }}</Radio>
              <Radio value="haiku" :disabled="!haikuAvailable">
                {{ $t('page.settings.model.haiku') }}
              </Radio>
            </div>
          </RadioGroup>
        </FormItem>
        <Alert
          v-if="!haikuAvailable"
          type="info"
          show-icon
          class="mb-6"
          :message="$t('page.settings.model.noKey')"
        />
        <FormItem
          :label="$t('page.settings.model.budget')"
          :extra="$t('page.settings.model.budgetHelp')"
          name="monthly_budget_usd"
          :rules="[{ message: $t('page.settings.required'), required: true }]"
        >
          <InputNumber
            v-model:value="form.monthly_budget_usd"
            :min="0"
            :step="1"
            :precision="2"
            prefix="$"
            class="w-48"
          />
        </FormItem>

        <Button type="primary" html-type="submit" block :loading="saving">
          {{ $t('page.settings.save') }}
        </Button>
      </Form>
    </Spin>
  </Page>
</template>
