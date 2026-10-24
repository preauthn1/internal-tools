<script lang="ts" setup>
import type { SettingsApi } from '#/api';

import { computed, nextTick, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { preferences } from '@vben/preferences';

import {
  Button,
  Card,
  Form,
  FormItem,
  message,
  Spin,
  Textarea,
} from 'ant-design-vue';

import { getFixedRepliesApi, putFixedReplyApi } from '#/api';
import { $t } from '#/locales';

defineOptions({ name: 'SettingsReplies' });

/** 畫面上的一則：每種語言用多行文字編輯，一行一句 */
interface Item {
  key: string;
  labels: Record<string, string>;
  texts: Record<string, string>;
}

const loading = ref(true);
const saving = ref('');
const preferred = ref('');
const languages = ref<string[]>([]);
const items = ref<Item[]>([]);

function toItem(r: SettingsApi.FixedReply): Item {
  return {
    key: r.key,
    labels: r.labels,
    texts: Object.fromEntries(
      Object.entries(r.texts).map(([lang, lines]) => [lang, lines.join('\n')]),
    ),
  };
}

/** 從預約設定、轉人工設定頁連過來時帶 ?key=，捲到那一則 */
const route = useRoute();

onMounted(async () => {
  try {
    const data = await getFixedRepliesApi();
    preferred.value = data.preferred_language;
    languages.value = data.languages;
    items.value = data.replies.map((r) => toItem(r));
  } finally {
    loading.value = false;
  }
  const key = route.query.key;
  if (typeof key === 'string') {
    await nextTick();
    document
      .querySelector(`#reply-${CSS.escape(key)}`)
      ?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
});

/** 語言名稱用目前介面語言顯示（瀏覽器內建的名稱表） */
const names = computed(
  () => new Intl.DisplayNames([example.invalid.locale], { type: 'language' }),
);
function nameOf(code: string) {
  return names.value.of(code) ?? code;
}

function labelOf(item: Item) {
  return item.labels[example.invalid.locale] || item.labels['zh-TW'] || item.key;
}

/** 慣用語言排第一，其他語言照設定的順序 */
const ordered = computed(() => [
  preferred.value,
  ...languages.value.filter((code) => code !== preferred.value),
]);

/** 存一則：慣用語言版本有改時，伺服器用 AI 翻成其他語言，存好後顯示翻好的內容 */
async function save(item: Item) {
  saving.value = item.key;
  try {
    const texts = Object.fromEntries(
      ordered.value.map((lang) => [
        lang,
        (item.texts[lang] ?? '')
          .split('\n')
          .map((s) => s.trim())
          .filter(Boolean),
      ]),
    );
    const data = await putFixedReplyApi(item.key, texts);
    // 只換掉這一則，其他則還沒存的修改留著
    const saved = data.replies.find((r) => r.key === item.key);
    if (saved) item.texts = toItem(saved).texts;
    message.success($t('page.settings.saved'));
  } finally {
    saving.value = '';
  }
}
</script>

<template>
  <Page
    :title="$t('page.settings.replies.title')"
    :description="$t('page.settings.replies.help')"
  >
    <Spin :spinning="loading">
      <div class="max-w-3xl">
        <Card
          v-for="item in items"
          :id="`reply-${item.key}`"
          :key="item.key"
          size="small"
          class="mb-3"
          :title="labelOf(item)"
        >
          <Spin
            :spinning="saving === item.key"
            :tip="$t('page.settings.replies.translating')"
          >
            <Form layout="vertical" @finish="save(item)">
              <FormItem
                v-for="lang in ordered"
                :key="lang"
                :label="
                  lang === preferred
                    ? $t('page.settings.replies.preferred', {
                        lang: nameOf(lang),
                      })
                    : nameOf(lang)
                "
                :extra="
                  lang === preferred
                    ? $t('page.settings.replies.preferredHelp')
                    : undefined
                "
              >
                <Textarea
                  v-model:value="item.texts[lang]"
                  :auto-size="{ minRows: 1 }"
                  :placeholder="
                    lang === preferred
                      ? undefined
                      : $t('page.settings.replies.emptyTranslation')
                  "
                />
              </FormItem>
              <Button
                type="primary"
                html-type="submit"
                block
                :loading="saving === item.key"
                :disabled="saving !== '' && saving !== item.key"
              >
                {{ $t('page.settings.save') }}
              </Button>
            </Form>
          </Spin>
        </Card>
      </div>
    </Spin>
  </Page>
</template>
