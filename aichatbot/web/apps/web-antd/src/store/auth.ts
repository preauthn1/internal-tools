import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { defineStore } from 'pinia';

import { getMeApi, loginApi, logoutApi } from '#/api';
import { applyAutoLocale } from '#/locales';

/**
 * 真正的登入憑證是 HttpOnly 的 session cookie（前端讀不到）。
 * Vben 的路由守衛看 accessToken 判斷有沒有登入，所以登入後放一個不是機密的佔位值；
 * session 失效時 API 回 401，request.ts 會清掉它、回登入頁。
 */
const SESSION_PLACEHOLDER = 'cookie-session';

/** 只有一種管理員角色 */
const ADMIN_ROLE = 'admin';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();

  const loginLoading = ref(false);

  /**
   * 登入：成功後取得目前使用者（含介面語言），再回原本要去的頁面
   * @param params 登入表單資料
   */
  async function authLogin(
    params: Recordable<any>,
    onSuccess?: () => Promise<void> | void,
  ) {
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      await loginApi(params);
      accessStore.setAccessToken(SESSION_PLACEHOLDER);
      userInfo = await fetchUserInfo();
      accessStore.setAccessCodes([]);

      if (accessStore.loginExpired) {
        accessStore.setLoginExpired(false);
      } else {
        onSuccess
          ? await onSuccess?.()
          : await router.push(
              userInfo.homePath || example.invalid.defaultHomePath,
            );
      }
    } finally {
      loginLoading.value = false;
    }

    return {
      userInfo,
    };
  }

  async function logout(redirect: boolean = true) {
    try {
      await logoutApi();
    } catch {
      // session 已經失效也沒關係
    }
    resetAllStores();
    accessStore.setLoginExpired(false);

    // 回登入頁，帶上目前的路由位址
    await router.replace({
      path: LOGIN_PATH,
      query: redirect
        ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
          }
        : {},
    });
  }

  /** 取得目前使用者；介面語言跟著商家慣用語言（管理員手動選過就不改） */
  async function fetchUserInfo() {
    const me = await getMeApi();
    const userInfo: UserInfo = {
      avatar: '',
      desc: '',
      homePath: example.invalid.defaultHomePath,
      realName: me.username,
      roles: [ADMIN_ROLE],
      token: '',
      userId: me.username,
      username: me.username,
    };
    userStore.setUserInfo(userInfo);
    await applyAutoLocale(me.language);
    return userInfo;
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    fetchUserInfo,
    loginLoading,
    logout,
  };
});
