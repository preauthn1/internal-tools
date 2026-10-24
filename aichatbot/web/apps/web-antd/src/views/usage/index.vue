<script lang="ts" setup>
import type { ViewingApi } from '#/api/viewing';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { usePreferences } from '@vben/preferences';

import {
  Card,
  Col,
  Empty,
  List,
  ListItem,
  Progress,
  Row,
  Spin,
  Statistic,
  Table,
} from 'ant-design-vue';

import { getUsageApi } from '#/api/viewing';
import { $t } from '#/locales';

defineOptions({ name: 'Usage' });

const { isMobile } = usePreferences();
const loading = ref(true);
const usage = ref<ViewingApi.Usage>();

/** 美元，小數點後 4 位（和模型設定頁的上限同一個單位） */
function usd(n: number) {
  return `$${n.toFixed(4)}`;
}

/** 本月費用占上限的比例（上限 0 是不限） */
const percent = computed(() => {
  const u = usage.value;
  if (!u || u.monthly_budget_usd <= 0) {
    return null;
  }
  return Math.min(100, Math.round((u.total_cost_usd / u.monthly_budget_usd) * 100));
});

const month = computed(() => {
  const u = usage.value;
  // from 是台北時間的月初（+08:00），直接取字串的年月，不受瀏覽器時區影響
  return u ? u.from.slice(0, 7) : '';
});

const columns = computed(() => [
  { dataIndex: 'model', key: 'model', title: $t('page.usage.model') },
  { dataIndex: 'calls', key: 'calls', title: $t('page.usage.calls') },
  { dataIndex: 'failed', key: 'failed', title: $t('page.usage.failed') },
  { key: 'tokens', title: $t('page.usage.tokens') },
  { key: 'cost', title: $t('page.usage.cost') },
]);

onMounted(async () => {
  try {
    usage.value = await getUsageApi();
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <Page
    :description="$t('page.usage.description')"
    :title="$t('page.usage.title')"
  >
    <Spin :spinning="loading">
      <template v-if="usage">
        <Row :gutter="[16, 16]" class="mb-4">
          <Col :md="8" :xs="24">
            <Card>
              <Statistic
                :title="$t('page.usage.totalCost', { month })"
                :value="usd(usage.total_cost_usd)"
              />
              <div class="mt-2 text-gray-500 text-xs">
                {{
                  usage.monthly_budget_usd > 0
                    ? $t('page.usage.budget', { budget: usd(usage.monthly_budget_usd) })
                    : $t('page.usage.noBudget')
                }}
              </div>
              <Progress v-if="percent !== null" :percent="percent" size="small" />
            </Card>
          </Col>
          <Col :md="8" :xs="24">
            <Card>
              <Statistic
                :title="$t('page.usage.totalCalls', { month })"
                :value="usage.total_calls"
              />
            </Card>
          </Col>
          <Col :md="8" :xs="24">
            <Card>
              <Statistic
                :title="$t('page.usage.currentModel')"
                :value="$t(`page.settings.model.${usage.model}`)"
                :value-style="{ fontSize: '16px' }"
              />
            </Card>
          </Col>
        </Row>

        <Card :title="$t('page.usage.byModel')">
          <p class="mb-3 text-gray-500 text-xs">{{ $t('page.usage.afmFree') }}</p>
          <Empty v-if="usage.by_model.length === 0" />
          <!-- 手機：清單（表格太寬） -->
          <List v-else-if="isMobile" :data-source="usage.by_model" size="small">
            <template #renderItem="{ item }">
              <ListItem>
                <div class="w-full">
                  <div class="flex justify-between font-bold">
                    <span class="break-all">{{ item.model }}</span>
                    <span>{{ usd(item.cost_usd) }}</span>
                  </div>
                  <div class="text-gray-500 text-xs">
                    {{ $t('page.usage.calls') }} {{ item.calls }}・{{
                      $t('page.usage.failed')
                    }}
                    {{ item.failed }}・{{ $t('page.usage.tokens') }}
                    {{ item.input_tokens }} / {{ item.output_tokens }}
                  </div>
                </div>
              </ListItem>
            </template>
          </List>
          <Table
            v-else
            :columns="columns"
            :data-source="usage.by_model"
            :pagination="false"
            row-key="model"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'tokens'">
                {{ record.input_tokens }} / {{ record.output_tokens }}
              </template>
              <template v-else-if="column.key === 'cost'">
                {{ usd(record.cost_usd) }}
              </template>
            </template>
          </Table>
        </Card>
      </template>
    </Spin>
  </Page>
</template>
