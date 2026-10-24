/**
 * 後台 API（Go 的 internal/admin）：登入狀態放在 HttpOnly 的 session cookie，前端不拿 token；
 * 回應直接是 JSON（不是 Vben 範例的 {code, data}），錯誤是 {"error": 訊息}。
 */
import type { RequestClientOptions } from '@vben/request';

import { useAppConfig } from '@vben/hooks';
import { preferences } from '@vben/preferences';
import { errorMessageResponseInterceptor, RequestClient } from '@vben/request';
import { useAccessStore } from '@vben/stores';

import { message } from 'ant-design-vue';

import { useAuthStore } from '#/store';

const { apiURL } = useAppConfig(import.meta.env, import.meta.env.PROD);

function createRequestClient(baseURL: string, options?: RequestClientOptions) {
  const client = new RequestClient({
    ...options,
    baseURL,
    withCredentials: true,
  });

  /**
   * session 失效（過期、登出、被清掉）：回登入頁
   */
  async function doReAuthenticate() {
    const accessStore = useAccessStore();
    const authStore = useAuthStore();
    accessStore.setAccessToken(null);
    if (
      example.invalid.loginExpiredMode === 'modal' &&
      accessStore.isAccessChecked
    ) {
      accessStore.setLoginExpired(true);
    } else {
      await authStore.logout();
    }
  }

  client.addRequestInterceptor({
    fulfilled: async (config) => {
      config.headers['Accept-Language'] = example.invalid.locale;
      return config;
    },
  });

  // 回應本身就是資料：直接回 body（RequestClient 預設回整個 axios response；非 2xx 已經被 axios 擋成錯誤）。
  // 個別請求設 responseReturn: 'raw'（例如下載檔案）時照舊回整個 response
  client.addResponseInterceptor({
    fulfilled: (response) =>
      response.config.responseReturn === 'raw' ? response : response.data,
  });

  // 401：登入本身失敗（帳號或密碼錯誤）照常顯示錯誤，其他的是 session 失效
  client.addResponseInterceptor({
    rejected: async (error) => {
      if (
        error?.response?.status === 401 &&
        !error.config?.url?.endsWith('/auth/login')
      ) {
        await doReAuthenticate();
      }
      throw error;
    },
  });

  client.addResponseInterceptor(
    errorMessageResponseInterceptor((msg: string, error) => {
      const responseData = error?.response?.data ?? {};
      message.error(responseData?.error || msg);
    }),
  );

  return client;
}

export const requestClient = createRequestClient(apiURL, {
  responseReturn: 'body',
});

/** 沒有攔截器（登出用：失敗也不跳錯誤訊息） */
export const baseRequestClient = new RequestClient({
  baseURL: apiURL,
  withCredentials: true,
});
