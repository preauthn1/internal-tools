export type Game = {
    provider: string;
    // 數字型 gameID（進入遊戲用）
    gameId: string;
    // gameCode（顯示 / 識別用；與 gameId 常不同）
    code: string;
    name: string;
    isPrePaidMode: boolean;
    isLaunchAvailable: boolean;
    image: string | null;
};

export type GameRecord = {
    id: number;
    provider: string;
    player_name: string;
    record_id: string;
    game_id: string | null;
    game_code: string | null;
    game_name: string | null;
    bet_amount: string;
    payout_amount: string;
    valid_amount: string;
    currency: string | null;
    bet_at: string | null;
};

export type WalletTransaction = {
    id: number;
    type: 'initial' | 'transfer_in' | 'transfer_out' | 'refund';
    type_label: string;
    provider: string | null;
    provider_name: string | null;
    amount: string;
    balance_after: string;
    created_at: string | null;
};

// Laravel 分頁器（length-aware）精簡型別
export type Paginated<T> = {
    data: T[];
    current_page: number;
    last_page: number;
    per_page: number;
    total: number;
    links: { url: string | null; label: string; active: boolean }[];
};
