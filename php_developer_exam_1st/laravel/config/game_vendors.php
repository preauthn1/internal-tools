<?php

return [
    'target' => [
        'name' => 'PGSname',
        'api_url' => env('PGS_API_URL', 'https://example.invalid/v1'),
        'merchant_id' => env('GAME_VENDOR_PGS_MERCHANT_ID', ''),
        'merchant_secret' => env('GAME_VENDOR_PGS_MERCHANT_SECRET', ''),
        'api_timeout' => 30,
        'auth_method' => 'payload',
        'icon' => 'https://example.invalid/cdn-betmonsters-com-cdn-bucket/1673072398100.gif',
        'game_list_cache_ttl' => 86400 + 3600, // Cache game list for 1 hour
        'game_records_page_size' => 10000,
    ],
    //'bwg' => [
        //'name' => 'BWGname',
        //'api_url' => env('BWG_API_URL', ''),
        //'id' => env('GAME_VENDOR_BWG_MERCHANT_ID', ''),
        //'secret' => env('GAME_VENDOR_BWG_MERCHANT_SECRET', ''),
        //'api_timeout' => 30,
        //'auth_method' => 'http_basic',
        //'icon' => 'https://example.invalid/cdn-betmonsters-com-cdn-bucket/174306850781.png',
    //],
];
