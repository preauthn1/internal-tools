import { requestClient } from '#/api/request';

export namespace PendingQuestionApi {
  /** 待處理、已忽略、已處理（轉成 FAQ、知識、案例） */
  export type Status = 'ignored' | 'open' | 'resolved';

  /** 待補問題：text、lang 是第一次出現的原話和語言，translated 是翻成慣用語言的版本 */
  export interface Question {
    count: number;
    first_seen: string;
    id: number;
    lang: string;
    last_seen: string;
    resolution: '' | 'example' | 'faq' | 'knowledge';
    staff_answer: string;
    status: Status;
    text: string;
    translated: string;
  }

  /** 詳細頁的一則對話 */
  export interface LogLine {
    direction: 'bot' | 'customer' | 'staff';
    /** 機器人的回覆沒有送到顧客 */
    send_failed: boolean;
    text: string;
    time: string;
  }

  /** 每次出現：source 是 unanswerable、verify_failed、unclear */
  export interface Hit {
    conversation: LogLine[];
    lang: string;
    source: string;
    text: string;
    time: string;
    user_id: string;
  }

  export interface Detail extends Question {
    hits: Hit[];
  }
}

/**
 * 列出待補問題（次數多的在前）
 */
export async function getPendingQuestionsApi(
  status: PendingQuestionApi.Status,
) {
  return requestClient.get<PendingQuestionApi.Question[]>(
    '/pending-questions',
    { params: { status } },
  );
}

/**
 * 一筆待補問題的詳細資料（每次出現和前後的對話）
 */
export async function getPendingQuestionApi(id: number) {
  return requestClient.get<PendingQuestionApi.Detail>(
    `/pending-questions/${id}`,
  );
}

/**
 * 轉成 FAQ：在 FAQ 最後新增一題
 */
export async function pendingToFAQApi(
  id: number,
  data: { alt_questions: string[]; answer: string; question: string },
) {
  return requestClient.post<PendingQuestionApi.Detail>(
    `/pending-questions/${id}/faq`,
    data,
  );
}

/**
 * 加進知識段落：title 是現有段落就接在後面，否則新增一段
 */
export async function pendingToKnowledgeApi(
  id: number,
  data: { body: string; title: string },
) {
  return requestClient.post<PendingQuestionApi.Detail>(
    `/pending-questions/${id}/knowledge`,
    data,
  );
}

/**
 * 加成意圖案例（和題庫或已有的句子重複時回 409）
 */
export async function pendingToExampleApi(
  id: number,
  data: { intent: string; text: string },
) {
  return requestClient.post<PendingQuestionApi.Detail>(
    `/pending-questions/${id}/example`,
    data,
  );
}

/**
 * 忽略（之後相同的問題照樣計次，但不在待處理列表）
 */
export async function ignorePendingApi(id: number) {
  return requestClient.post<PendingQuestionApi.Detail>(
    `/pending-questions/${id}/ignore`,
    {},
  );
}
