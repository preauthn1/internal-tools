import { requestClient } from '#/api/request';

/** 商家設定（Go 的 internal/admin/settings.go）：每個區塊分開讀、分開存，存了馬上生效 */
export namespace SettingsApi {
  /** 知識段落 */
  export interface Section {
    body: string;
    title: string;
  }

  /** 商家資訊：名稱、介紹、聯絡方式必填，加上知識段落 */
  export interface Merchant {
    contact: string;
    intro: string;
    knowledge: Section[];
    name: string;
  }

  /**
   * 長度警告（#65）：內容太長時本機模型會答不出來，不擋儲存。只在 GET、PUT 的回應裡，送出時不能帶。
   * kind：faq_filtered（FAQ 有 count 題，先挑最像的 top 題）、faq_too_long、knowledge_too_long（合計 chars 字，上限 limit）、
   * knowledge_section_long（targets 這幾段單段超過 limit）
   */
  export interface Warning {
    chars?: number;
    count?: number;
    kind: string;
    limit?: number;
    targets?: string[];
    top?: number;
  }

  /** GET、PUT 商家資訊的回應：另外附上知識段落的長度警告 */
  export interface MerchantResponse extends Merchant {
    warnings: Warning[];
  }

  /** 一題 FAQ；key 是程式用的代號，新增時送空字串由伺服器產生 */
  export interface FAQ {
    alt_questions: string[];
    answer: string;
    key: string;
    question: string;
  }

  /** GET、PUT FAQ 的回應：另外附上長度警告 */
  export interface FAQResponse {
    faq: FAQ[];
    warnings: Warning[];
  }

  export interface Language {
    code: string;
    enabled: boolean;
  }

  export interface LanguageSettings {
    languages: Language[];
    preferred_language: string;
  }

  /** GET 另外附上可勾選（supported）和量過品質（measured）的語言 */
  export interface Languages extends LanguageSettings {
    measured: string[];
    supported: string[];
  }

  /** 一則固定回覆：labels 是後台介面語言 → 白話名稱，texts 是顧客語言 → 句子（機器人隨機挑一句） */
  export interface FixedReply {
    key: string;
    labels: Record<string, string>;
    texts: Record<string, string[]>;
  }

  /** 可以改的固定回覆；languages 是啟用的顧客語言 */
  export interface FixedReplies {
    languages: string[];
    preferred_language: string;
    replies: FixedReply[];
  }

  /** 預約設定；關閉時回的文字是固定回覆「預約關閉」（booking_off） */
  export interface BookingSetting {
    enabled: boolean;
    /** 自訂欄位（最多 5 個），照順序追問必填的 */
    fields: BookingField[];
    /** 問不問人數 */
    people_enabled: boolean;
  }

  /** 預約的自訂欄位：key 是代號（小寫英文開頭的英數底線），description 給 AI 抽取用 */
  export interface BookingField {
    description: string;
    key: string;
    name: string;
    required: boolean;
  }

  /** 員工服務時段：weekday 0 是星期日；start、end 是 HH:MM（end 不含） */
  export interface ServiceHours {
    end: string;
    start: string;
    weekday: number;
  }

  /** 轉人工設定（存的欄位）；enabled 是轉人工總開關，關閉時回固定回覆「轉人工關閉」（handoff_off） */
  export interface HandoffSettingInput {
    auto_end_minutes: number;
    enabled: boolean;
    on_no_answer: boolean;
    service_hours: ServiceHours[];
    streak_limit: number;
  }

  /** GET 另外附上只讀的綁定碼和有沒有綁好員工群組 */
  export interface HandoffSetting extends HandoffSettingInput {
    bind_code: string;
    bound: boolean;
  }

  /** 模型設定：haiku_available 是伺服器有沒有設定 API key（後台看不到 key 本身） */
  export interface ModelSettings {
    haiku_available: boolean;
    model: 'afm' | 'haiku';
    monthly_budget_usd: number;
  }
}

export async function getMerchantApi() {
  return requestClient.get<SettingsApi.MerchantResponse>('/merchant');
}

export async function putMerchantApi(data: SettingsApi.Merchant) {
  return requestClient.put<SettingsApi.MerchantResponse>('/merchant', data);
}

export async function getFAQApi() {
  return requestClient.get<SettingsApi.FAQResponse>('/faq');
}

export async function putFAQApi(faq: SettingsApi.FAQ[]) {
  return requestClient.put<SettingsApi.FAQResponse>('/faq', { faq });
}

export async function getLanguagesApi() {
  return requestClient.get<SettingsApi.Languages>('/languages');
}

export async function putLanguagesApi(data: SettingsApi.LanguageSettings) {
  return requestClient.put<SettingsApi.Languages>('/languages', data);
}

export async function getFixedRepliesApi() {
  return requestClient.get<SettingsApi.FixedReplies>('/fixed-replies');
}

/** 存一則固定回覆；改了慣用語言版本時伺服器會用 AI 翻成其他語言，本機模型可能要幾十秒，逾時放寬到 3 分鐘 */
export async function putFixedReplyApi(
  key: string,
  texts: Record<string, string[]>,
) {
  return requestClient.put<SettingsApi.FixedReplies>(
    `/fixed-replies/${encodeURIComponent(key)}`,
    { texts },
    { timeout: 180_000 },
  );
}

export async function getBookingSettingApi() {
  return requestClient.get<SettingsApi.BookingSetting>('/booking-setting');
}

export async function putBookingSettingApi(data: SettingsApi.BookingSetting) {
  return requestClient.put<SettingsApi.BookingSetting>('/booking-setting', data);
}

export async function getHandoffSettingApi() {
  return requestClient.get<SettingsApi.HandoffSetting>('/handoff-setting');
}

export async function putHandoffSettingApi(
  data: SettingsApi.HandoffSettingInput,
) {
  return requestClient.put<SettingsApi.HandoffSetting>(
    '/handoff-setting',
    data,
  );
}

/** 重新產生綁定碼：舊的碼就不能再用 */
export async function regenerateBindCodeApi() {
  return requestClient.post<SettingsApi.HandoffSetting>(
    '/handoff-setting/bind-code',
    {},
  );
}

export async function getModelApi() {
  return requestClient.get<SettingsApi.ModelSettings>('/model');
}

export async function putModelApi(data: {
  model: SettingsApi.ModelSettings['model'];
  monthly_budget_usd: number;
}) {
  return requestClient.put<SettingsApi.ModelSettings>('/model', data);
}
