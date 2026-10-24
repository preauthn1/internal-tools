import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

// 只能看的頁面：預約、對話紀錄、AI 費用
const routes: RouteRecordRaw[] = [
  {
    name: 'Bookings',
    path: '/bookings',
    component: () => import('#/views/bookings/index.vue'),
    meta: {
      icon: 'lucide:calendar-check',
      title: $t('page.bookings.title'),
    },
  },
  {
    name: 'Conversations',
    path: '/conversations',
    component: () => import('#/views/conversations/index.vue'),
    meta: {
      icon: 'lucide:messages-square',
      title: $t('page.conversations.title'),
    },
  },
  {
    name: 'ConversationDetail',
    path: '/conversations/:platform/:userID',
    component: () => import('#/views/conversations/detail.vue'),
    meta: {
      activePath: '/conversations',
      hideInMenu: true,
      title: $t('page.conversations.detail'),
    },
  },
  {
    name: 'Usage',
    path: '/usage',
    component: () => import('#/views/usage/index.vue'),
    meta: {
      icon: 'lucide:receipt',
      title: $t('page.usage.title'),
    },
  },
];

export default routes;
