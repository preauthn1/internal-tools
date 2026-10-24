<?php

namespace App\Enums;

/**
 * 錢包帳本異動類型。
 */
enum WalletTransactionType: string
{
    case Initial = 'initial';           // 註冊初始金
    case TransferIn = 'transfer_in';    // 轉入遊戲商（平台餘額減少）
    case TransferOut = 'transfer_out';  // 自遊戲商轉回（平台餘額增加）
    case Refund = 'refund';             // 轉入失敗退回

    public function label(): string
    {
        return match ($this) {
            self::Initial => '初始金',
            self::TransferIn => '轉入遊戲',
            self::TransferOut => '從遊戲轉回',
            self::Refund => '轉入失敗退回',
        };
    }
}
