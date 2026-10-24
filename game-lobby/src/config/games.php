<?php

return [

    /*
    |--------------------------------------------------------------------------
    | 預設遊戲商
    |--------------------------------------------------------------------------
    | 對應 game_providers.key。遊戲商設定改由 game_providers 資料表存放
    | （見 GameProviderSeeder / GameProviderManager），不再放在此 config。
    */
    'default' => env('GAME_DEFAULT_PROVIDER', 'target'),

    /*
    |--------------------------------------------------------------------------
    | 註冊初始金額（平台幣）
    |--------------------------------------------------------------------------
    */
    'initial_balance' => env('GAME_INITIAL_BALANCE', 10000),

    'provider_defaults' => [
        'target' => [
            'merchant_id' => env('PGS_MERCHANT_ID'),
            'merchant_secret' => env('PGS_MERCHANT_SECRET'),
            'base_url' => env('PGS_API_URL', 'https://example.invalid/v1'),
            'timeout' => (int) env('PGS_TIMEOUT', 10),
            'locale' => env('PGS_LOCALE', 'zh-cn'),
            'player_prefix' => env('GAME_PLAYER_PREFIX', 'target'),
            'games_cache_ttl' => (int) env('PGS_GAMES_CACHE_TTL', 300),
            'invalid_state_codes' => array_map('intval', array_values(array_filter(array_map('trim', explode(',', (string) env('PGS_INVALID_STATE_CODES', '403601')))))),
            'records_window_minutes' => (int) env('PGS_RECORDS_WINDOW', 60),
            'records_lookback_minutes' => (int) env('PGS_RECORDS_LOOKBACK', 120),
            'records_page_size' => (int) env('PGS_RECORDS_PAGE_SIZE', 200),
            // 用 ?: 讓「變數存在但為空」也退回預設，避免空值被轉成 0 而形同關閉稽核延遲。
            'audit_delay_minutes' => (int) (env('PGS_AUDIT_DELAY_MINUTES') ?: 30),
            'audit_retry_minutes' => (int) env('PGS_AUDIT_RETRY_MINUTES', 5),
            'audit_review_after_minutes' => (int) env('PGS_AUDIT_REVIEW_AFTER_MINUTES', 120),
            'audit_max_attempts' => (int) env('PGS_AUDIT_MAX_ATTEMPTS', 24),
        ],
    ],

];

/*
|--------------------------------------------------------------------------
| 送遊戲商的 playerName 前綴
|--------------------------------------------------------------------------
| 綁定時以「{prefix}_{英數亂數}_{user_id}」產生遊戲商玩家識別（如 pgs_a1B2c3D4_5）。
| 前綴完全由各遊戲商自行持有（game_providers.settings.player_prefix，見
| GameProviderSeeder）；未設定前綴的遊戲商，PlayGameService 會退回以 provider key
| 當前綴，確保同一會員在不同遊戲商產生的玩家識別必然不同。中段為隨機英數（不外流平台
| account），user_id 保證唯一；產生後持久化於 game_provider_accounts.player_name 沿用。
*/
