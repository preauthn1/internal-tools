/**
 * 介面語言的決定順序：管理員手動選的 > 商家慣用語言（登入後 /api/auth/me 回傳）> 瀏覽器語言 > 繁中。
 * Vben 的偏好設定本身會把語言存在 localStorage，分不出是手動選的還是自動決定的，
 * 所以手動選過另外記一個 key；之後就不再自動改。
 */
import type { SupportedLanguagesType } from '@vben/locales';

import { watch } from 'vue';

import { LOCALE_SELECTED_EVENT, loadLocaleMessages } from '@vben/locales';
import { preferences, updatePreferences } from '@vben/preferences';

const MANUAL_KEY = 'aichatbot-locale-manual';

/** 最後一次自動決定的語言；語言變成別的就是手動選的 */
let autoLocale: null | SupportedLanguagesType = null;

function isManual() {
  return localStorage.getItem(MANUAL_KEY) === '1';
}

/** 瀏覽器語言對到介面語言：中文一律繁中，日文、英文各自對應，其他用繁中 */
export function browserLocale(
  languages: readonly string[] = navigator.languages,
): SupportedLanguagesType {
  for (const lang of languages) {
    const code = lang.toLowerCase();
    if (code.startsWith('zh')) return 'zh-TW';
    if (code.startsWith('ja')) return 'ja-JP';
    if (code.startsWith('en')) return 'en-US';
  }
  return 'zh-TW';
}

/**
 * 沒手動選過時改成 lang。i18n 還沒裝好（啟動時）loadMessages 傳 false，
 * 之後 setupI18n 會用偏好設定裡的語言載入
 */
export async function applyAutoLocale(
  lang: SupportedLanguagesType,
  loadMessages = true,
) {
  if (isManual()) return;
  autoLocale = lang;
  if (example.invalid.locale === lang) return;
  updatePreferences({ app: { locale: lang } });
  if (loadMessages) await loadLocaleMessages(lang);
}

function markManual() {
  localStorage.setItem(MANUAL_KEY, '1');
}

/**
 * 管理員在語言切換選單點了語言（就算跟目前的語言一樣）：手動選的，記下來。
 * 偏好設定面板改語言不會發事件，所以也看語言有沒有變成不是自動決定的值
 */
export function watchManualLocale() {
  window.addEventListener(LOCALE_SELECTED_EVENT, markManual);
  watch(
    () => example.invalid.locale,
    (locale) => {
      if (autoLocale !== null && locale !== autoLocale) {
        markManual();
      }
    },
  );
}
