import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    name: 'PendingQuestions',
    path: '/pending-questions',
    component: () => import('#/views/pending-questions/index.vue'),
    meta: {
      icon: 'lucide:message-circle-question',
      title: $t('page.pendingQuestions.title'),
    },
  },
];

export default routes;
