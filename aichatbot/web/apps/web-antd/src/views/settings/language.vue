<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';
import { preferences } from '@vben/preferences';

import {
  Alert,
  Button,
  Checkbox,
  Form,
  FormItem,
  message,
  Select,
  Spin,
  Tag,
} from 'ant-design-vue';

import { getLanguagesApi, getMeApi, putLanguagesApi } from '#/api';
import { $t, applyAutoLocale } from '#/locales';

defineOptions({ name: 'SettingsLanguage' });

const loading = ref(true);
const saving = ref(false);
const supported = ref<string[]>([]);
const measured = ref<string[]>([]);
/** 存過的語言順序（含沒勾的）；存檔時照這個順序，新勾的接在後面 */
const savedOrder = ref<string[]>([]);
const enabled = ref<Record<string, boolean>>({});
const preferred = ref('');

function apply(data: Awaited<ReturnType<typeof getLanguagesApi>>) {
  supported.value = data.supported;
  measured.value = data.measured;
  savedOrder.value = data.languages.map((l) => l.code);
  enabled.value = Object.fromEntries(
    data.languages.map((l) => [l.code, l.enabled]),
  );
  preferred.value = data.preferred_language;
}

onMounted(async () => {
  try {
    apply(await getLanguagesApi());
  } finally {
    loading.value = false;
  }
});

/** 語言名稱用目前介面語言顯示（瀏覽器內建的名稱表） */
const names = computed(
  () => new Intl.DisplayNames([example.invalid.locale], { type: 'language' }),
);
function nameOf(code: string) {
  return names.value.of(code) ?? code;
}

const checked = computed(() =>
  supported.value.filter((code) => enabled.value[code]),
);
const unmeasured = computed(() =>
  checked.value.filter((code) => !measured.value.includes(code)),
);
const preferredOptions = computed(() =>
  checked.value.map((code) => ({ label: nameOf(code), value: code })),
);

// 慣用語言被取消勾選時，改成第一個有勾的語言（慣用語言只能選有勾的）
watch(checked, (codes) => {
  if (!codes.includes(preferred.value)) preferred.value = codes[0] ?? '';
});

async function save() {
  saving.value = true;
  try {
    // 已經不在可勾選清單裡的語言不送（伺服器會擋），存了就清掉
    const order = [
      ...savedOrder.value.filter((code) => supported.value.includes(code)),
      ...checked.value.filter((code) => !savedOrder.value.includes(code)),
    ];
    apply(
      await putLanguagesApi({
        languages: order.map((code) => ({
          code,
          enabled: !!enabled.value[code],
        })),
        preferred_language: preferred.value,
      }),
    );
    message.success($t('page.settings.saved'));
    // 介面語言跟著新的慣用語言（管理員手動選過語言就不改）
    await applyAutoLocale((await getMeApi()).language);
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <Page :title="$t('page.settings.language.title')">
    <Spin :spinning="loading">
      <Form layout="vertical" class="max-w-3xl" @finish="save">
        <FormItem
          :label="$t('page.settings.language.customer')"
          :extra="$t('page.settings.language.customerHelp')"
        >
          <div class="flex flex-col gap-2">
            <Checkbox
              v-for="code in supported"
              :key="code"
              v-model:checked="enabled[code]"
            >
              {{ nameOf(code) }}
              <Tag v-if="!measured.includes(code)" color="warning" class="ml-1">
                {{ $t('page.settings.language.unmeasured') }}
              </Tag>
            </Checkbox>
          </div>
        </FormItem>
        <Alert
          v-if="unmeasured.length > 0"
          type="warning"
          show-icon
          class="mb-6"
          :message="
            $t('page.settings.language.unmeasuredWarning', {
              langs: unmeasured.map(nameOf).join('、'),
            })
          "
        />
        <FormItem
          :label="$t('page.settings.language.preferred')"
          :extra="$t('page.settings.language.preferredHelp')"
        >
          <Select v-model:value="preferred" :options="preferredOptions" />
        </FormItem>

        <Button type="primary" html-type="submit" block :loading="saving">
          {{ $t('page.settings.save') }}
        </Button>
      </Form>
    </Spin>
  </Page>
</template>
