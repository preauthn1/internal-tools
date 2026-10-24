<script lang="ts" setup>
import { Page } from '@vben/common-ui';
import { useUserStore } from '@vben/stores';

import { $t } from '#/locales';

defineOptions({ name: 'Dashboard' });

const userStore = useUserStore();
</script>

<template>
  <Page
    :description="$t('page.dashboard.comingSoon')"
    :title="
      $t('page.dashboard.welcome', { name: userStore.userInfo?.realName ?? '' })
    "
  />
</template>
