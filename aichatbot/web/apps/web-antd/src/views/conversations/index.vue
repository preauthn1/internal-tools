<script lang="ts" setup>
import type { ViewingApi } from '#/api/viewing';

import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import { Alert, Card, Empty, InputSearch, List, ListItem, Spin, Tag } from 'ant-design-vue';

import { getCustomersApi } from '#/api/viewing';
import { $t } from '#/locales';

defineOptions({ name: 'Conversations' });

const router = useRouter();
const query = ref('');
const list = ref<ViewingApi.Customer[]>([]);
const loading = ref(false);

function formatTime(s: string) {
  return new Date(s).toLocaleString();
}

async function load() {
  loading.value = true;
  try {
    list.value = await getCustomersApi(query.value.trim());
  } finally {
    loading.value = false;
  }
}

function open(c: ViewingApi.Customer) {
  router.push({
    name: 'ConversationDetail',
    params: { platform: c.platform, userID: c.user_id },
  });
}

onMounted(load);
</script>

<template>
  <Page
    :description="$t('page.conversations.description')"
    :title="$t('page.conversations.title')"
  >
    <Alert
      :message="$t('page.conversations.retention')"
      class="mb-4"
      show-icon
      type="info"
    />
    <Card>
      <InputSearch
        v-model:value="query"
        :placeholder="$t('page.conversations.search')"
        allow-clear
        class="mb-4"
        enter-button
        @search="load"
      />
      <Spin :spinning="loading">
        <Empty v-if="list.length === 0" />
        <!-- 清單在手機和桌面都能用：點一位顧客進對話內容頁 -->
        <List v-else :data-source="list" size="small">
          <template #renderItem="{ item }">
            <ListItem class="cursor-pointer" @click="open(item)">
              <div class="min-w-0 flex-1">
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <span class="font-bold break-all">
                    <Tag>{{ item.platform }}</Tag>{{ item.user_id }}
                  </span>
                  <span class="text-gray-500 text-xs">
                    {{ formatTime(item.last_time) }}
                  </span>
                </div>
                <div class="truncate text-gray-500">{{ item.last_text }}</div>
                <div class="text-gray-400 text-xs">
                  {{ $t('page.conversations.messages', { n: item.messages }) }}
                </div>
              </div>
            </ListItem>
          </template>
        </List>
      </Spin>
    </Card>
  </Page>
</template>
