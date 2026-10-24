using MSC.Models.Enums;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.Enums
{
    public enum AdminPermissionTokenEnum
    {

        Report_Management,

        //卡號
        Card_Management,
        Card_Card_List,
        Card_Card_GenerateCardNo,


        //商戶管理
        Merchant_Management,
        System_Merchant_List,
        Merchant_MerchantSummary,
        Record_MerchantWalletLog_List,
        Merchant_Withdrawal_List,


        //歷史紀錄
        History_Management,
        Record_CardConsume_List,
        Record_CardConsumeMerchantNotifyQueue_List,
        Record_CardConsumeMerchantNotifyQueue_ViewExternalLogs,
        Record_CardConsumeMerchantNotifyQueue_ManualNotify,
        Record_MerchantLinkRequest_List,
        Record_UserOperateLogs_List,
        Record_RequestForbiddenLog_List,


        //報表
        Report_DailySummary_List,


        //系統設置
        System_Management,
        System_User_List,
        System_SystemSettingConfig_List,
        System_ExternalLink_List,
        System_InternalUtility

    }
}
