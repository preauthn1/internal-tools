import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    name: 'IntentExamples',
    path: '/intent-examples',
    component: () => import('#/views/intent-examples/index.vue'),
    meta: {
      icon: 'lucide:message-square-quote',
      title: $t('page.intentExamples.title'),
    },
  },
];

export default routes;
