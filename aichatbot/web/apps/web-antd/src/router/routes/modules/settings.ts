import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    name: 'Settings',
    path: '/settings',
    redirect: '/settings/merchant',
    meta: {
      icon: 'lucide:store',
      order: 0,
      title: $t('page.settings.title'),
    },
    children: [
      {
        name: 'SettingsMerchant',
        path: '/settings/merchant',
        component: () => import('#/views/settings/merchant.vue'),
        meta: { title: $t('page.settings.merchant.title') },
      },
      {
        name: 'SettingsFAQ',
        path: '/settings/faq',
        component: () => import('#/views/settings/faq.vue'),
        meta: { title: $t('page.settings.faq.title') },
      },
      {
        name: 'SettingsLanguage',
        path: '/settings/language',
        component: () => import('#/views/settings/language.vue'),
        meta: { title: $t('page.settings.language.title') },
      },
      {
        name: 'SettingsReplies',
        path: '/settings/replies',
        component: () => import('#/views/settings/replies.vue'),
        meta: { title: $t('page.settings.replies.title') },
      },
      {
        name: 'SettingsBooking',
        path: '/settings/booking',
        component: () => import('#/views/settings/booking.vue'),
        meta: { title: $t('page.settings.booking.title') },
      },
      {
        name: 'SettingsHandoff',
        path: '/settings/handoff',
        component: () => import('#/views/settings/handoff.vue'),
        meta: { title: $t('page.settings.handoff.title') },
      },
      {
        name: 'SettingsModel',
        path: '/settings/model',
        component: () => import('#/views/settings/model.vue'),
        meta: { title: $t('page.settings.model.title') },
      },
    ],
  },
];

export default routes;
