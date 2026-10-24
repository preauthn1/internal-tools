export type User = {
    id: number;
    name: string;
    account: string;
    email: string;
    avatar?: string;
    email_verified_at: string | null;
    created_at: string;
    updated_at: string;
    [key: string]: unknown;
};

export type Auth = {
    user: User | null;
    // 平台錢包餘額（未登入為 null）
    balance: string | null;
};
