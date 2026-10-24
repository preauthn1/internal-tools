<script lang="ts" setup>
import type { SettingsApi } from '#/api';

import { onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { Page } from '@vben/common-ui';

import { useClipboard } from '@vueuse/core';
import {
  Alert,
  Button,
  Card,
  Form,
  FormItem,
  InputNumber,
  message,
  Popconfirm,
  Select,
  Space,
  Spin,
  Switch,
  TimePicker,
  Typography,
} from 'ant-design-vue';

import {
  getHandoffSettingApi,
  putHandoffSettingApi,
  regenerateBindCodeApi,
} from '#/api';
import { $t } from '#/locales';

defineOptions({ name: 'SettingsHandoff' });

const loading = ref(true);
const saving = ref(false);
const regenerating = ref(false);
const form = ref<SettingsApi.HandoffSettingInput>({
  auto_end_minutes: 0,
  enabled: true,
  on_no_answer: true,
  service_hours: [],
  streak_limit: 3,
});
const bindCode = ref('');
const bound = ref(false);

function apply(data: SettingsApi.HandoffSetting) {
  const { bind_code, bound: isBound, ...rest } = data;
  form.value = rest;
  bindCode.value = bind_code;
  bound.value = isBound;
}

onMounted(async () => {
  try {
    apply(await getHandoffSettingApi());
  } finally {
    loading.value = false;
  }
});

/** 星期幾：0 是星期日，選單從星期一排到星期日 */
const weekdays = [1, 2, 3, 4, 5, 6, 0].map((value) => ({
  label: $t(`page.settings.handoff.weekdays.${value}`),
  value,
}));

function addHours() {
  form.value.service_hours.push({ end: '18:00', start: '09:00', weekday: 1 });
}

function removeHours(index: number) {
  form.value.service_hours.splice(index, 1);
}

/** 存檔；格式不對時伺服器回 400，錯誤訊息由 request 攔截器顯示 */
async function save() {
  saving.value = true;
  try {
    apply(await putHandoffSettingApi({ ...form.value }));
    message.success($t('page.settings.saved'));
  } finally {
    saving.value = false;
  }
}

async function regenerate() {
  regenerating.value = true;
  try {
    // 只換綁定碼，表單上還沒存的修改留著
    const data = await regenerateBindCodeApi();
    bindCode.value = data.bind_code;
    bound.value = data.bound;
    message.success($t('page.settings.handoff.regenerated'));
  } finally {
    regenerating.value = false;
  }
}

const { copy } = useClipboard({ legacy: true });
async function copyCode() {
  await copy(bindCode.value);
  message.success($t('page.settings.handoff.copied'));
}
</script>

<template>
  <Page
    :title="$t('page.settings.handoff.title')"
    :description="$t('page.settings.handoff.help')"
  >
    <Spin :spinning="loading">
      <div class="max-w-3xl">
        <Card
          size="small"
          class="mb-3"
          :title="$t('page.settings.handoff.bind')"
        >
          <Alert
            :type="bound ? 'success' : 'warning'"
            show-icon
            class="mb-3"
            :message="
              bound
                ? $t('page.settings.handoff.bound')
                : $t('page.settings.handoff.unbound')
            "
          />
          <Space wrap>
            <Typography.Text code class="text-lg">
              {{ bindCode }}
            </Typography.Text>
            <Button @click="copyCode">
              {{ $t('page.settings.handoff.copy') }}
            </Button>
            <Popconfirm
              :title="$t('page.settings.handoff.confirmRegenerate')"
              @confirm="regenerate"
            >
              <Button :loading="regenerating">
                {{ $t('page.settings.handoff.regenerate') }}
              </Button>
            </Popconfirm>
          </Space>
          <div class="mt-3 text-sm opacity-75">
            {{ $t('page.settings.handoff.bindHelp', { code: bindCode }) }}
          </div>
        </Card>

        <Form :model="form" layout="vertical" @finish="save">
          <FormItem :label="$t('page.settings.handoff.enabled')">
            <Switch v-model:checked="form.enabled" />
            <template #extra>
              {{ $t('page.settings.handoff.offReply') }}
              <RouterLink
                :to="{
                  path: '/settings/replies',
                  query: { key: 'handoff_off' },
                }"
              >
                {{ $t('page.settings.handoff.offReplyLink') }}
              </RouterLink>
            </template>
          </FormItem>
          <FormItem
            :label="$t('page.settings.handoff.onNoAnswer')"
            :extra="$t('page.settings.handoff.onNoAnswerHelp')"
          >
            <Switch v-model:checked="form.on_no_answer" />
          </FormItem>
          <FormItem
            :label="$t('page.settings.handoff.streakLimit')"
            :extra="$t('page.settings.handoff.streakLimitHelp')"
          >
            <InputNumber
              v-model:value="form.streak_limit"
              :min="1"
              :max="10"
              :precision="0"
              class="w-48"
            />
          </FormItem>
          <FormItem
            :label="$t('page.settings.handoff.autoEnd')"
            :extra="$t('page.settings.handoff.autoEndHelp')"
          >
            <InputNumber
              v-model:value="form.auto_end_minutes"
              :min="0"
              :max="1440"
              :precision="0"
              class="w-48"
            />
          </FormItem>
          <FormItem
            :label="$t('page.settings.handoff.serviceHours')"
            :extra="$t('page.settings.handoff.serviceHoursHelp')"
          >
            <div v-if="form.service_hours.length === 0" class="mb-2 opacity-75">
              {{ $t('page.settings.handoff.allDay') }}
            </div>
            <div
              v-for="(h, i) in form.service_hours"
              :key="i"
              class="mb-2 flex flex-wrap items-center gap-2"
            >
              <Select
                v-model:value="h.weekday"
                :options="weekdays"
                class="w-32"
              />
              <TimePicker
                v-model:value="h.start"
                format="HH:mm"
                value-format="HH:mm"
                :allow-clear="false"
              />
              <span>～</span>
              <TimePicker
                v-model:value="h.end"
                format="HH:mm"
                value-format="HH:mm"
                :allow-clear="false"
              />
              <Button danger @click="removeHours(i)">
                {{ $t('page.settings.remove') }}
              </Button>
            </div>
            <Button @click="addHours">
              {{ $t('page.settings.handoff.addHours') }}
            </Button>
          </FormItem>
          <FormItem :label="$t('page.settings.handoff.prompt')">
            {{ $t('page.settings.handoff.promptHelp') }}
            <RouterLink
              :to="{ path: '/settings/replies', query: { key: 'handoff' } }"
            >
              {{ $t('page.settings.handoff.promptLink') }}
            </RouterLink>
          </FormItem>

          <Button type="primary" html-type="submit" block :loading="saving">
            {{ $t('page.settings.save') }}
          </Button>
        </Form>
      </div>
    </Spin>
  </Page>
</template>
