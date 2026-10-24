import { requestClient } from '#/api/request';

export namespace IntentExampleApi {
  /** 意圖：labels 是三語白話名稱（key 是介面語言 zh-TW、en-US、ja-JP），畫面上不顯示 code */
  export interface Intent {
    code: string;
    labels: Record<string, string>;
  }

  /** 意圖案例庫的一句；builtin 是內建案例（只能看） */
  export interface Example {
    builtin: boolean;
    id: number;
    intent: string;
    text: string;
  }
}

/**
 * 意圖清單
 */
export async function getIntentsApi() {
  return requestClient.get<IntentExampleApi.Intent[]>('/intents');
}

/**
 * 列出意圖案例庫
 */
export async function getIntentExamplesApi() {
  return requestClient.get<IntentExampleApi.Example[]>('/intent-examples');
}

/**
 * 新增商家的案例（和題庫或已有的句子重複時回 409）
 */
export async function addIntentExampleApi(text: string, intent: string) {
  return requestClient.post<IntentExampleApi.Example>('/intent-examples', {
    intent,
    text,
  });
}

/**
 * 刪除商家的案例（寫入請求一律要 application/json，所以帶空物件）
 */
export async function deleteIntentExampleApi(id: number) {
  return requestClient.delete(`/intent-examples/${id}`, { data: {} });
}
