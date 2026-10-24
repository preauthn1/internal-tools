import type { SupportedLanguagesType } from '@vben/locales';

import { baseRequestClient, requestClient } from '#/api/request';

export namespace AuthApi {
  /** 登入參數 */
  export interface LoginParams {
    password?: string;
    username?: string;
  }

  /** 登入結果：session 放在 HttpOnly cookie，回應只有帳號 */
  export interface LoginResult {
    username: string;
  }

  /** 目前使用者：language 是商家慣用語言對應的介面語言 */
  export interface MeResult {
    language: SupportedLanguagesType;
    username: string;
  }
}

/**
 * 登入
 */
export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', {
    password: [REDACTED],
    username: data.username,
  });
}

/**
 * 登出（寫入請求一律要 application/json，所以帶空物件）
 */
export async function logoutApi() {
  return baseRequestClient.post('/auth/logout', {});
}

/**
 * 目前使用者
 */
export async function getMeApi() {
  return requestClient.get<AuthApi.MeResult>('/auth/me');
}
