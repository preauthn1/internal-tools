using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
  
    public enum TransactionTypeEnum
    {
        [Description(i18nKeys.Enums_TransactionTypeEnum_TransferWallet_To_GameWallet)]
        TransferWallet_To_GameWallet = 1,

        [Description(i18nKeys.Enums_TransactionTypeEnum_TransferWallet_To_PlatformWallet)]
        TransferWallet_To_PlatformWallet = 2,

        [Description(i18nKeys.Enums_TransactionTypeEnum_BetGame_Settlement)]
        BetGame_Settlement = 3,


        ManualDepositWalletAmount = 11,

        [Description(i18nKeys.Enums_TransactionTypeEnum_BetGame_Settlement)]
        SingleWallet_PlaceBetDebit = 21,

        [Description(i18nKeys.Enums_TransactionTypeEnum_BetGame_Settlement)]
        SingleWallet_BetSettlement = 22,


        /// <summary>
        /// 系統取消注單 自動退款
        /// </summary>
        [Description(i18nKeys.Enums_TransactionTypeEnum_BetGame_Settlement)]
        SingleWallet_CancelRefundWallet = 27,

        /// <summary>
        /// 手動人工退還金幣
        /// </summary>
        [Description(i18nKeys.Enums_TransactionTypeEnum_BetGame_Settlement)]
        SingleWallet_CancelRefundWalletManual = 39,
    }
}
