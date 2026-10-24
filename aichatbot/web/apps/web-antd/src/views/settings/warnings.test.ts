import { describe, expect, it } from 'vitest';

import { splitWarnings, warningText } from './warnings';

/** 假的翻譯函式：回傳 key 和參數，看得出用了哪段文字、帶了什麼值 */
const t = (key: string, params?: Record<string, unknown>) =>
  `${key} ${JSON.stringify(params)}`;

describe('長度警告（#65）', () => {
  it('依種類顯示說明，帶上字數、上限和是哪幾段', () => {
    expect(
      warningText(
        {
          chars: 6606,
          kind: 'knowledge_too_long',
          limit: 6000,
          targets: ['菜單', '歷史'],
        },
        t,
      ),
    ).toBe(
      'page.settings.warnings.knowledge_too_long {"chars":6606,"count":0,"limit":6000,"targets":"菜單、歷史","top":0}',
    );
  });

  it('GET 的回應拆成表單內容和警告：PUT 不能帶 warnings（伺服器不收不認得的欄位）', () => {
    const [form, warnings] = splitWarnings({
      faq: [],
      warnings: [{ count: 51, kind: 'faq_filtered', top: 20 }],
    });
    expect(form).toEqual({ faq: [] });
    expect(warnings).toEqual([{ count: 51, kind: 'faq_filtered', top: 20 }]);
  });
});
