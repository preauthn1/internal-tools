<script lang="ts" setup>
import type { IntentExampleApi } from '#/api/intent-examples';
import type { ViewingApi } from '#/api/viewing';

import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';
import { preferences } from '@vben/preferences';

import {
  Alert,
  Button,
  Card,
  Collapse,
  CollapsePanel,
  message,
  Popconfirm,
  Spin,
  Tag,
  Tooltip,
} from 'ant-design-vue';

import { getIntentsApi } from '#/api/intent-examples';
import {
  endHandoffApi,
  getConversationApi,
  translateLogApi,
} from '#/api/viewing';
import { $t } from '#/locales';

defineOptions({ name: 'ConversationDetail' });

const route = useRoute();
const router = useRouter();
const platform = String(route.params.platform);
const userID = String(route.params.userID);

const loading = ref(true);
const conversation = ref<ViewingApi.Conversation>();
const intents = ref<IntentExampleApi.Intent[]>([]);
// 翻譯結果只放在畫面上（不存起來）：訊息編號 → 譯文
const translations = ref<Record<number, string>>({});
const translating = ref<Record<number, boolean>>({});

function formatTime(s: string) {
  return new Date(s).toLocaleString();
}

/** 意圖的白話名稱：依目前介面語言，沒有就用繁中 */
function intentLabel(code: string) {
  const labels = intents.value.find((i) => i.code === code)?.labels;
  return labels?.[example.invalid.locale] || labels?.['zh-TW'] || code;
}

const customActions: Record<string, string> = {
  fixed_text: 'page.customIntents.fixedText',
  handoff: 'page.customIntents.handoff',
  knowledge: 'page.customIntents.knowledge',
};

const faqChoices = new Set(['none', 'not_a_question', 'price_calc', 'unclear']);

/** 驗證、翻譯擋下的原因：result 是「原因」或「原因:附帶的值」（number_not_in_quote:15） */
function reason(result: string) {
  const i = result.indexOf(':');
  const [key, arg] = i === -1 ? [result, ''] : [result.slice(0, i), result.slice(i + 1)];
  return $t(`page.conversations.reasons.${key}`, { arg });
}

/** 一步 AI 判斷的白話結果 */
function traceResult(t: ViewingApi.Trace) {
  if (t.detail?.error === 'ai_failed' || t.result === 'ai_failed') {
    return $t('page.conversations.reasons.ai_failed');
  }
  switch (t.step) {
    case 'answer': {
      return $t(`page.conversations.answerable.${t.result}`);
    }
    case 'ask': {
      // 挑不到知識時確認是不是在問店家：寒暄、看不懂或問商家
      return $t(`page.conversations.kinds.${t.result}`);
    }
    case 'chat': {
      // 寒暄回的固定回覆：打招呼、道謝、道別或應答（#54）
      return $t(`page.conversations.chatKinds.${t.result}`);
    }
    case 'custom': {
      return customActions[t.result] ? $t(customActions[t.result]!) : t.result;
    }
    case 'intent': {
      return intentLabel(t.result);
    }
    case 'reask': {
      return $t(`page.conversations.fields.${t.result}`);
    }
    case 'translate':
    case 'verify': {
      return t.result === 'ok' ? $t('page.conversations.passed') : reason(t.result);
    }
    case 'faq': {
      // 挑到的 FAQ 代號，或算價、看不懂、不是在問、都不能回答
      return faqChoices.has(t.result)
        ? $t(`page.conversations.faqChoices.${t.result}`)
        : t.result;
    }
    default: {
      // pick：挑到的知識段落（當時知識清單的第幾段，K1+K2）
      return t.result === 'none' ? $t('page.conversations.none') : t.result;
    }
  }
}

/** 被擋下或失敗的步驟用紅色標示 */
function blocked(t: ViewingApi.Trace) {
  if (t.detail?.error || t.result === 'ai_failed') {
    return true;
  }
  if (t.step === 'verify' || t.step === 'translate') {
    return t.result !== 'ok';
  }
  return t.step === 'answer' && t.result !== 'yes';
}

/** 判斷的細節（顧客問的事、根據、AI 寫的回覆…）；kind 是寒暄（chat）或問商家（question） */
function detailEntries(t: ViewingApi.Trace) {
  return Object.entries(t.detail ?? {})
    .filter(([k, v]) => k !== 'error' && v)
    .map(([k, v]): [string, string] => [
      k,
      k === 'kind' ? $t(`page.conversations.kinds.${v}`) : v,
    ]);
}

async function translate(m: ViewingApi.Message) {
  translating.value[m.id] = true;
  try {
    const { translation } = await translateLogApi(m.id);
    translations.value[m.id] = translation;
  } finally {
    translating.value[m.id] = false;
  }
}

/**
 * 結束轉人工（同員工群組的按鈕），之後重新載入對話（多一筆「結束轉人工」）；
 * 失敗時也重新載入（員工已經在群組結束、或自動結束了，拿掉按鈕）
 */
const ending = ref(false);
async function endHandoff() {
  ending.value = true;
  try {
    await endHandoffApi(platform, userID);
    message.success($t('page.conversations.handoffEnded'));
  } finally {
    ending.value = false;
    conversation.value = await getConversationApi(platform, userID);
  }
}

onMounted(async () => {
  try {
    const [c, i] = await Promise.all([
      getConversationApi(platform, userID),
      getIntentsApi(),
    ]);
    conversation.value = c;
    intents.value = i;
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <Page :title="userID">
    <template #description>
      <Button class="px-0" type="link" @click="router.push('/conversations')">
        ← {{ $t('page.conversations.back') }}
      </Button>
      <span class="ml-2 text-gray-500">{{ platform }}</span>
    </template>
    <Alert
      :message="$t('page.conversations.retention')"
      class="mb-4"
      show-icon
      type="info"
    />
    <Alert
      v-if="conversation?.human"
      :message="$t('page.conversations.handingOff')"
      class="mb-4"
      show-icon
      type="warning"
    >
      <template #action>
        <Popconfirm
          :title="$t('page.conversations.confirmEndHandoff')"
          @confirm="endHandoff"
        >
          <Button :loading="ending" size="small" type="primary">
            {{ $t('page.conversations.endHandoff') }}
          </Button>
        </Popconfirm>
      </template>
    </Alert>
    <Spin :spinning="loading">
      <Card v-if="conversation" size="small">
        <div
          v-for="m in conversation.messages"
          :key="m.id"
          :class="m.direction === 'customer' ? 'items-start' : 'items-end'"
          class="mb-3 flex flex-col"
        >
          <div class="mb-1 text-xs text-gray-500">
            {{ $t(`page.pendingQuestions.directions.${m.direction}`) }}・{{
              formatTime(m.time)
            }}
            <Tooltip
              v-if="m.send_failed"
              :title="$t('page.conversations.sendFailedHint')"
            >
              <Tag class="ml-1 mr-0" color="error">
                {{ $t('page.conversations.sendFailed') }}
              </Tag>
            </Tooltip>
          </div>
          <div
            :class="{
              'bg-gray-100 dark:bg-gray-800': m.direction === 'customer',
              'bg-blue-50 dark:bg-blue-950': m.direction === 'bot',
              'bg-green-50 dark:bg-green-950': m.direction === 'staff',
            }"
            class="max-w-[90%] rounded-lg px-3 py-2 md:max-w-[70%]"
          >
            <div v-if="m.text" class="whitespace-pre-wrap break-words">
              {{ m.text }}
            </div>
            <div v-else class="text-gray-400">
              {{ $t('page.conversations.silent') }}
            </div>
            <div
              v-if="translations[m.id]"
              class="mt-1 whitespace-pre-wrap break-words border-t pt-1 text-gray-600 dark:text-gray-300"
            >
              {{ translations[m.id] }}
            </div>
            <Button
              v-if="m.translatable && !translations[m.id]"
              :loading="translating[m.id]"
              class="px-0"
              size="small"
              type="link"
              @click="translate(m)"
            >
              {{ $t('page.conversations.translate') }}
            </Button>
            <Collapse
              v-if="m.traces.length > 0"
              :bordered="false"
              class="mt-1 bg-transparent"
              ghost
              size="small"
            >
              <CollapsePanel
                key="traces"
                :header="$t('page.conversations.traces')"
              >
                <div v-for="(t, i) in m.traces" :key="i" class="mb-1">
                  <Tag :color="blocked(t) ? 'red' : undefined">
                    {{ $t(`page.conversations.steps.${t.step}`) }}
                  </Tag>
                  <span :class="blocked(t) ? 'text-red-500' : ''">
                    {{ traceResult(t) }}
                  </span>
                  <div
                    v-for="[k, v] in detailEntries(t)"
                    :key="k"
                    class="ml-2 break-words text-xs text-gray-500"
                  >
                    {{ $t(`page.conversations.details.${k}`) }}：{{ v }}
                  </div>
                </div>
              </CollapsePanel>
            </Collapse>
          </div>
        </div>
      </Card>
    </Spin>
  </Page>
</template>
