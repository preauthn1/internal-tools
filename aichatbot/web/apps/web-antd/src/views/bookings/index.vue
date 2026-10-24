<script lang="ts" setup>
import type { Dayjs } from 'dayjs';

import type { SettingsApi } from '#/api';
import type { ViewingApi } from '#/api/viewing';

import { computed, onMounted, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';
import { usePreferences } from '@vben/preferences';

import {
  Alert,
  Card,
  DatePicker,
  Empty,
  List,
  ListItem,
  Segmented,
  Spin,
  Table,
  Tag,
} from 'ant-design-vue';

import { getBookingSettingApi } from '#/api';
import { getBookingsApi } from '#/api/viewing';
import { $t } from '#/locales';

defineOptions({ name: 'Bookings' });

const { isMobile } = usePreferences();

const status = ref<'' | ViewingApi.BookingStatus>('');
const range = ref<[Dayjs, Dayjs]>();
const list = ref<ViewingApi.Booking[]>([]);
const loading = ref(false);
/** 預約欄位（#67）：人數關掉時不顯示人數，自訂欄位各一欄 */
const setting = ref<SettingsApi.BookingSetting>({
  enabled: true,
  fields: [],
  people_enabled: true,
});

const statusOptions = computed(() => [
  { label: $t('page.bookings.all'), value: '' },
  { label: $t('page.bookings.statuses.pending'), value: 'pending' },
  { label: $t('page.bookings.statuses.confirmed'), value: 'confirmed' },
  { label: $t('page.bookings.statuses.rejected'), value: 'rejected' },
]);

// List、Table 的插槽拿到的資料沒有型別，模板裡用 as 轉成 BookingStatus
const statusColors: Record<ViewingApi.BookingStatus, string> = {
  confirmed: 'green',
  pending: 'orange',
  rejected: 'red',
};

const columns = computed(() => [
  { key: 'when', title: $t('page.bookings.when'), width: 160 },
  ...(setting.value.people_enabled
    ? [{ dataIndex: 'people', key: 'people', title: $t('page.bookings.people'), width: 80 }]
    : []),
  ...setting.value.fields.map((f) => ({ key: `extra:${f.key}`, title: f.name })),
  { dataIndex: 'name', key: 'name', title: $t('page.bookings.name') },
  { dataIndex: 'phone', key: 'phone', title: $t('page.bookings.phone') },
  { key: 'status', title: $t('page.bookings.status'), width: 160 },
  { key: 'created_at', title: $t('page.bookings.createdAt'), width: 180 },
]);

function formatTime(s: string) {
  return new Date(s).toLocaleString();
}

async function load() {
  loading.value = true;
  try {
    list.value = await getBookingsApi({
      from: range.value?.[0].format('YYYY-MM-DD'),
      status: status.value || undefined,
      to: range.value?.[1].format('YYYY-MM-DD'),
    });
  } finally {
    loading.value = false;
  }
}

watch([status, range], load);
onMounted(async () => {
  await load();
  setting.value = await getBookingSettingApi();
});
</script>

<template>
  <Page
    :description="$t('page.bookings.description')"
    :title="$t('page.bookings.title')"
  >
    <Card>
      <div class="mb-4 flex flex-wrap items-center gap-3">
        <Segmented v-model:value="status" :options="statusOptions" />
        <DatePicker.RangePicker
          v-model:value="range"
          :placeholder="[$t('page.bookings.from'), $t('page.bookings.to')]"
          :class="isMobile ? 'w-full' : ''"
        />
      </div>
      <Alert
        :message="$t('page.bookings.handleInGroup')"
        class="mb-4"
        show-icon
        type="info"
      />

      <!-- 手機：卡片清單（表格太寬） -->
      <Spin v-if="isMobile" :spinning="loading">
        <Empty v-if="list.length === 0" />
        <List v-else :data-source="list" item-layout="vertical" size="small">
          <template #renderItem="{ item }">
            <ListItem>
              <div class="flex items-center justify-between">
                <span class="font-bold">{{ item.date }} {{ item.time }}</span>
                <Tag :color="statusColors[item.status as ViewingApi.BookingStatus]">                  {{ $t(`page.bookings.statuses.${item.status}`) }}
                </Tag>
              </div>
              <div>
                {{ item.name }}
                <template v-if="setting.people_enabled">
                  ・{{ item.people }} {{ $t('page.bookings.peopleUnit') }}
                </template>
              </div>
              <template v-for="f in setting.fields" :key="f.key">
                <div v-if="item.extra?.[f.key]">
                  {{ f.name }}：{{ item.extra[f.key] }}
                </div>
              </template>
              <div>
                <a :href="`tel:${item.phone}`">{{ item.phone }}</a>
              </div>
              <div class="text-gray-500 text-xs">
                {{ $t('page.bookings.createdAt') }}
                {{ formatTime(item.created_at) }}
                <template v-if="item.handled_by">
                  ・{{ $t('page.bookings.handledBy', { name: item.handled_by }) }}
                </template>
              </div>
            </ListItem>
          </template>
        </List>
      </Spin>

      <Table
        v-else
        :columns="columns"
        :data-source="list"
        :loading="loading"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'when'">
            {{ record.date }} {{ record.time }}
          </template>
          <template v-else-if="column.key === 'status'">
            <Tag :color="statusColors[record.status as ViewingApi.BookingStatus]">              {{ $t(`page.bookings.statuses.${record.status}`) }}
            </Tag>
            <div v-if="record.handled_by" class="text-gray-500 text-xs">
              {{ $t('page.bookings.handledBy', { name: record.handled_by }) }}
            </div>
          </template>
          <template v-else-if="column.key === 'created_at'">
            {{ formatTime(record.created_at) }}
          </template>
          <template v-else-if="String(column.key).startsWith('extra:')">
            {{ record.extra?.[String(column.key).slice('extra:'.length)] }}
          </template>
        </template>
      </Table>
    </Card>
  </Page>
</template>
