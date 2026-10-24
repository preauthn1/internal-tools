import { requestClient } from '#/api/request';

export namespace CustomIntentApi {
  /** 動作：回固定文字、轉人工（text 是提示）、只根據指定知識回答 */
  export type Action = 'fixed_text' | 'handoff' | 'knowledge';

  /** 自訂意圖；code 由伺服器產生（新的留空），畫面上不顯示 */
  export interface CustomIntent {
    action: Action;
    code: string;
    description: string;
    examples: string[];
    knowledge: string[];
    name: string;
    text: string;
  }

  export interface Body {
    custom_intents: CustomIntent[];
    /** 知識段落的標題（只根據指定知識回答時選段落用） */
    knowledge_titles: string[];
  }

  /** 自我測試的一句：expected 是標的意圖，got 是分類結果（AI 失敗時是空字串） */
  export interface SelfTestItem {
    expected: string;
    got: string;
    state?: string;
    text: string;
  }

  export interface SelfTestReport {
    failed: SelfTestItem[];
    conflicts: SelfTestItem[];
    total: number;
  }

  /** 分類自我測試的狀態；結果只放伺服器記憶體，重開後是 idle */
  export interface SelfTest {
    done: number;
    report?: SelfTestReport;
    status: 'done' | 'failed' | 'idle' | 'running';
    total: number;
  }
}

/**
 * 讀自訂意圖和知識段落標題
 */
export async function getCustomIntentsApi() {
  return requestClient.get<CustomIntentApi.Body>('/custom-intents');
}

/**
 * 整份存自訂意圖（沒通過檢查回 400，例句和題庫重複回 409）；存好後伺服器在背景跑分類自我測試
 */
export async function putCustomIntentsApi(
  list: CustomIntentApi.CustomIntent[],
) {
  return requestClient.put<CustomIntentApi.Body>('/custom-intents', {
    custom_intents: list,
  });
}

/**
 * 分類自我測試的狀態
 */
export async function getSelfTestApi() {
  return requestClient.get<CustomIntentApi.SelfTest>(
    '/custom-intents/self-test',
  );
}

/**
 * 手動重跑分類自我測試
 */
export async function runSelfTestApi() {
  return requestClient.post<CustomIntentApi.SelfTest>(
    '/custom-intents/self-test',
    {},
  );
}
