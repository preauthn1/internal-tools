import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    name: 'CustomIntents',
    path: '/custom-intents',
    component: () => import('#/views/custom-intents/index.vue'),
    meta: {
      icon: 'lucide:shapes',
      title: $t('page.customIntents.title'),
    },
  },
];

export default routes;
