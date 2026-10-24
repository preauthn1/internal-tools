import { requestClient } from '#/api/request';

// 只能看的頁面：預約、對話紀錄、AI 費用

export namespace ViewingApi {
  /** 待確認、已確認、已拒絕（確認、拒絕在員工群組操作） */
  export type BookingStatus = 'confirmed' | 'pending' | 'rejected';

  export interface Booking {
    created_at: string;
    date: string;
    /** 自訂欄位的值（代號 → 值），沒有是 null */
    extra: null | Record<string, string>;
    handled_at: null | string;
    handled_by: string;
    id: number;
    name: string;
    people: string;
    phone: string;
    platform: string;
    status: BookingStatus;
    summary: string;
    time: string;
    user_id: string;
  }

  /** 篩選：status 不填是全部；from、to 是預約日期（YYYY-MM-DD，含頭尾） */
  export interface BookingFilter {
    from?: string;
    status?: BookingStatus;
    to?: string;
  }

  /** 對話紀錄的一位顧客 */
  export interface Customer {
    last_text: string;
    last_time: string;
    messages: number;
    platform: string;
    user_id: string;
  }

  /** 機器人回覆的一步 AI 判斷 */
  export interface Trace {
    detail?: Record<string, string>;
    result: string;
    step: string;
  }

  export interface Message {
    direction: 'bot' | 'customer' | 'staff';
    id: number;
    reply_key: string;
    /** 機器人的回覆沒有送到顧客 */
    send_failed: boolean;
    text: string;
    time: string;
    traces: Trace[];
    /** 不是商家慣用語言，可以翻譯 */
    translatable: boolean;
  }

  export interface Conversation {
    /** 顧客在轉人工中 */
    human: boolean;
    lang: string;
    messages: Message[];
    platform: string;
    preferred_language: string;
    user_id: string;
  }

  export interface ModelUsage {
    calls: number;
    cost_usd: number;
    failed: number;
    input_tokens: number;
    model: string;
    output_tokens: number;
  }

  /** 本月（台北時間 [from, to)）的 AI 用量 */
  export interface Usage {
    by_model: ModelUsage[];
    from: string;
    model: 'afm' | 'haiku';
    monthly_budget_usd: number;
    to: string;
    total_calls: number;
    total_cost_usd: number;
  }
}

/**
 * 預約列表（新申請的在前）
 */
export async function getBookingsApi(filter: ViewingApi.BookingFilter) {
  return requestClient.get<ViewingApi.Booking[]>('/bookings', {
    params: filter,
  });
}

/**
 * 有對話紀錄的顧客（最近有對話的在前）；q 比對顧客編號或訊息內容
 */
export async function getCustomersApi(q: string) {
  return requestClient.get<ViewingApi.Customer[]>('/conversations', {
    params: { q },
  });
}

/**
 * 一位顧客的對話紀錄（舊的在前）
 */
export async function getConversationApi(platform: string, userID: string) {
  return requestClient.get<ViewingApi.Conversation>(
    `/conversations/${encodeURIComponent(platform)}/${encodeURIComponent(userID)}`,
  );
}

/**
 * 把一則對話翻成商家慣用語言（每次都呼叫目前的模型，不存起來）
 */
export async function translateLogApi(id: number) {
  return requestClient.post<{ translation: string }>(
    `/conversations/logs/${id}/translate`,
    {},
  );
}

/**
 * 結束這位顧客的轉人工，恢復 AI 回覆（同員工群組的「結束轉人工」按鈕）
 */
export async function endHandoffApi(platform: string, userID: string) {
  return requestClient.post<{ ok: boolean }>(
    `/conversations/${encodeURIComponent(platform)}/${encodeURIComponent(userID)}/end-handoff`,
    {},
  );
}

/**
 * 本月的 AI 費用和呼叫次數
 */
export async function getUsageApi() {
  return requestClient.get<ViewingApi.Usage>('/usage');
}
