export type Locale = 'en-US' | 'ja-JP' | 'zh-TW';

export const messages: Record<Locale, Record<string, string>> = {
  'en-US': {
    cancel: 'Cancel',
    collapse: 'Collapse',
    confirm: 'Confirm',
    expand: 'Expand',
    prompt: 'Prompt',
    reset: 'Reset',
    submit: 'Submit',
  },
  'ja-JP': {
    cancel: 'キャンセル',
    collapse: '折りたたむ',
    confirm: '確認',
    expand: '展開',
    prompt: 'お知らせ',
    reset: 'リセット',
    submit: '送信',
  },
  'zh-TW': {
    cancel: '取消',
    collapse: '收合',
    confirm: '確認',
    expand: '展開',
    prompt: '提示',
    reset: '重設',
    submit: '送出',
  },
};

export const getMessages = (locale: Locale) => messages[locale];
