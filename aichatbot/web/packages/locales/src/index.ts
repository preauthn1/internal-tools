import {
  i18n,
  loadLocaleMessages,
  loadLocalesMap,
  loadLocalesMapFromDir,
  setupI18n,
} from './i18n';

const $t = i18n.global.t;
const $te = i18n.global.te;

/**
 * aichatbot：使用者在語言切換選單點了某個語言（就算跟目前的語言一樣）時，在 window 上發出的事件，
 * detail 是選的語言。app 用它分辨手動選的語言和自動決定的語言
 */
const LOCALE_SELECTED_EVENT = 'vben:locale-selected';

export {
  $t,
  $te,
  i18n,
  LOCALE_SELECTED_EVENT,
  loadLocaleMessages,
  loadLocalesMap,
  loadLocalesMapFromDir,
  setupI18n,
};
export {
  type ImportLocaleFn,
  type LocaleSetupOptions,
  type SupportedLanguagesType,
} from './typing';
export type { CompileError } from '@intlify/core-base';

export { useI18n } from 'vue-i18n';

export type { Locale } from 'vue-i18n';
