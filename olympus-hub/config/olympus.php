<?php

return [
    /*
    |--------------------------------------------------------------------------
    | 神聖指令 (Sacred Commands)
    |--------------------------------------------------------------------------
    */
    'commands' => [
        'game' => [
            '遊戲核心 (Game Core)' => [
                'game:sync-config'      => '同步遊戲配置 (Sync Game Config)',
                'game:maintenance-mode' => '切換維護模式 (Toggle Maintenance)',
            ],
        ],
        'finance' => [
            '金流商 (Cash Flow)' => [
                'cash-flow:sync' => '同步金流配置 (Sync Cash Flows)',
            ],
            '代理管理 (Agent)' => [
                'agent:add-agent-cash-flow' => '新增代理金流 (Add Agent Cash Flow)',
            ],
        ],
    ],

    /*
    |--------------------------------------------------------------------------
    | Consul 神聖鏈結設定 (Consul Tether Settings)
    |--------------------------------------------------------------------------
    */
    'consul' => [
        'base_url'    => env('CONSUL_HTTP_ADDR', 'http://consul:8500'),
        'token'       => env('CONSUL_HTTP_TOKEN'),
        'realms_path' => base_path('realms'),
    ],
];
