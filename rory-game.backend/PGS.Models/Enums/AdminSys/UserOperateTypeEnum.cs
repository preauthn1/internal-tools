using target.Models.CultureText;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;

namespace target.Models.Enums
{

    /// <summary>
    /// UserOperateRecord 這張表在用的
    /// </summary>
    public enum UserOperateTypeEnum
    {       
        [Description(i18nKeys.Enums_UserOperateTypeEnum_User_UserLogin)]
        User_UserLogin = 7777,

        [Description(i18nKeys.Enums_UserOperateTypeEnum_Player_LockStatus)]
        Player_LockStatus = 1001,

        [Description(i18nKeys.Enums_UserOperateTypeEnum_Player_UnlockStatus)]
        Player_UnlockStatus = 1002,

        [Description(i18nKeys.Enums_UserOperateTypeEnum_Player_UnlockStatus)]
        Player_RoomLockRemove = 1011,




        /// <summary>
        /// 手動關閉注單 並退還金幣
        /// </summary>
        SingleWallet_CancelRefundWalletManual = 2051,

        /// <summary>
        /// 手動結算注單
        /// </summary>
        SingleWallet_ManualSettleBetRecord = 2052,



        ManualDepositWalletAmount = 3101,


    }
}
