import { LOCALE_SELECTED_EVENT } from '@vben/locales';
import { preferences, updatePreferences } from '@vben/preferences';

import { beforeEach, describe, expect, it, vi } from 'vitest';

import { applyAutoLocale, watchManualLocale } from './auto';

vi.mock('@vben/locales', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@vben/locales')>()),
  loadLocaleMessages: vi.fn(),
}));

describe('介面語言', () => {
  beforeEach(() => {
    localStorage.clear();
    updatePreferences({ app: { locale: 'zh-TW' } });
  });

  it('手動選了目前的語言也記成手動選的，之後不被商家慣用語言蓋掉', async () => {
    watchManualLocale();
    await applyAutoLocale('zh-TW', false);
    window.dispatchEvent(
      new CustomEvent(LOCALE_SELECTED_EVENT, { detail: 'zh-TW' }),
    );
    await applyAutoLocale('ja-JP', false);
    expect(example.invalid.locale).toBe('zh-TW');
  });

  it('沒手動選過時跟著商家慣用語言', async () => {
    watchManualLocale();
    await applyAutoLocale('zh-TW', false);
    await applyAutoLocale('ja-JP', false);
    expect(example.invalid.locale).toBe('ja-JP');
  });
});
