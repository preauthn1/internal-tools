<script lang="ts" setup>
import type { SettingsApi } from '#/api';

import { onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  Button,
  Card,
  Checkbox,
  Form,
  FormItem,
  Input,
  message,
  Spin,
  Switch,
  Textarea,
} from 'ant-design-vue';

import { getBookingSettingApi, putBookingSettingApi } from '#/api';
import { $t } from '#/locales';

import { move } from './list';

defineOptions({ name: 'SettingsBooking' });

/** 自訂欄位最多幾個（和伺服器的 merchant.MaxBookingFields 相同） */
const maxFields = 5;

const loading = ref(true);
const saving = ref(false);
const form = ref<SettingsApi.BookingSetting>({
  enabled: true,
  fields: [],
  people_enabled: true,
});

onMounted(async () => {
  try {
    form.value = await getBookingSettingApi();
  } finally {
    loading.value = false;
  }
});

function addField() {
  form.value.fields.push({
    description: '',
    key: '',
    name: '',
    required: true,
  });
}

/** 前端只做必填和代號格式的提示；其他檢查（代號重複、內建欄位）由伺服器回錯誤訊息，request.ts 會顯示 */
async function save() {
  saving.value = true;
  try {
    form.value = await putBookingSettingApi({ ...form.value });
    message.success($t('page.settings.saved'));
  } finally {
    saving.value = false;
  }
}

const required = [{ message: $t('page.settings.required'), required: true }];
const keyRules = [
  ...required,
  {
    message: $t('page.settings.booking.keyFormat'),
    pattern: /^[a-z][a-z0-9_]{0,19}$/,
  },
];
</script>

<template>
  <Page
    :title="$t('page.settings.booking.title')"
    :description="$t('page.settings.booking.help')"
  >
    <Spin :spinning="loading">
      <Form :model="form" layout="vertical" class="max-w-3xl" @finish="save">
        <FormItem :label="$t('page.settings.booking.enabled')">
          <Switch v-model:checked="form.enabled" />
          <template #extra>
            {{ $t('page.settings.booking.offReply') }}
            <RouterLink
              :to="{ path: '/settings/replies', query: { key: 'booking_off' } }"
            >
              {{ $t('page.settings.booking.offReplyLink') }}
            </RouterLink>
          </template>
        </FormItem>

        <h3 class="mb-1 mt-6 text-base font-semibold">
          {{ $t('page.settings.booking.fields') }}
        </h3>
        <p class="mb-3 text-sm text-foreground/60">
          {{ $t('page.settings.booking.fieldsHelp') }}
        </p>
        <FormItem :label="$t('page.settings.booking.people')">
          <Switch v-model:checked="form.people_enabled" />
          <template #extra>
            {{ $t('page.settings.booking.peopleHelp') }}
          </template>
        </FormItem>

        <p class="mb-1 font-medium">
          {{ $t('page.settings.booking.custom', { max: maxFields }) }}
        </p>
        <p class="mb-3 text-sm text-foreground/60">
          {{ $t('page.settings.booking.customHelp') }}
        </p>
        <Card
          v-for="(field, i) in form.fields"
          :key="i"
          size="small"
          class="mb-3"
        >
          <FormItem
            :label="$t('page.settings.booking.fieldName')"
            :name="['fields', i, 'name']"
            :rules="required"
          >
            <Input v-model:value="field.name" />
          </FormItem>
          <FormItem
            :label="$t('page.settings.booking.fieldKey')"
            :name="['fields', i, 'key']"
            :rules="keyRules"
            :extra="$t('page.settings.booking.fieldKeyHelp')"
          >
            <Input v-model:value="field.key" />
          </FormItem>
          <FormItem
            :label="$t('page.settings.booking.fieldDescription')"
            :name="['fields', i, 'description']"
            :rules="required"
            :extra="$t('page.settings.booking.fieldDescriptionHelp')"
          >
            <Textarea
              v-model:value="field.description"
              :auto-size="{ minRows: 2 }"
            />
          </FormItem>
          <FormItem :extra="$t('page.settings.booking.fieldRequiredHelp')">
            <Checkbox v-model:checked="field.required">
              {{ $t('page.settings.booking.fieldRequired') }}
            </Checkbox>
          </FormItem>
          <div class="flex flex-wrap gap-2">
            <Button
              size="small"
              :disabled="i === 0"
              @click="move(form.fields, i, -1)"
            >
              {{ $t('page.settings.moveUp') }}
            </Button>
            <Button
              size="small"
              :disabled="i === form.fields.length - 1"
              @click="move(form.fields, i, 1)"
            >
              {{ $t('page.settings.moveDown') }}
            </Button>
            <Button size="small" danger @click="form.fields.splice(i, 1)">
              {{ $t('page.settings.remove') }}
            </Button>
          </div>
        </Card>
        <Button
          block
          class="mb-6"
          :disabled="form.fields.length >= maxFields"
          @click="addField"
        >
          {{ $t('page.settings.booking.addField') }}
        </Button>

        <Button type="primary" html-type="submit" block :loading="saving">
          {{ $t('page.settings.save') }}
        </Button>
      </Form>
    </Spin>
  </Page>
</template>
