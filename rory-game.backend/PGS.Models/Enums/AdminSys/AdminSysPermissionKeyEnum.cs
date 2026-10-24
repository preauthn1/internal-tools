using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{

    [Description("對應Permission資料庫的Token這欄位")]
    public enum AdminSysPermissionKey
    {
        Home,
        MembersManagement,
        MembersManagementList,

        /// <summary>
        /// 用戶列表(玩家Player)
        /// </summary>
        MembersManagementList_Search,
        MembersManagementList_Lock,


        MembersManagementDetail,

        /// <summary>
        /// 用戶查詢Detail (玩家Player)
        /// </summary>
        MembersManagementDetail_Search,
        MembersManagementDetail_Lock,

        /// <summary>
        /// 金幣紀錄(帳變紀錄) transaction
        /// </summary>
        MembersManagementMoneyRecord,
        MembersManagementBetRecord,


        /// <summary>
        /// 封號紀錄 (列表頁)
        /// </summary>
        MembersManagementLock,
        MembersManagementTransfer,
        GameManagement,

        /// <summary>
        ///  游戏历程查询 (perGame)
        /// </summary>
        GamesManagementRecord,
        GamesManagementServersHealthSummary,
        GamesManagement_BetRecordSearch,
        GamesManagement_BetRecordSearch_SearchList,
        GamesManagement_BetRecord_ManualRefund,


        GameSetting,
        GameSettingSwitch,
        GameSettingSwitch_Search,
        GameSettingSwitch_Update,
        Statistics,
        StatisticsMerchantSummary,
        StatisticsGame,
        OperatorManagement,
        OperatorManagementAccounts,
        OperatorManagementAccounts_Search,
        OperatorManagementAccounts_Create,
        OperatorManagementAccounts_Update,
        OperatorManagementAccounts_Delete,
        SystemSetting,
        SystemSettingOrganizations,
        SystemSettingOrganizations_Search,
        SystemSettingOrganizations_Create,
        SystemSettingOrganizations_Update,
        SystemSettingOrganizations_Delete,
        SystemSettingMerchant,
        SystemSettingMerchant_Search,
        SystemSettingMerchant_Create,
        SystemSettingMerchant_Update,
        SystemSettingMerchant_Delete,
        SystemSettingServers,
        SystemSettingServers_Search,
        SystemSettingServers_Create,
        SystemSettingServers_Update,
        SystemSettingServers_Delete,
        SystemSettingServers_Dispatch,
        SystemSettingServerRobotConfig,
        SystemSettingServerRobotConfig_Search,
        SystemSettingServerRobotConfig_Create,
        SystemSettingServerRobotConfig_Update,
        SystemSettingServerRobotConfig_Delete,
        SystemSettingRoles,
        SystemSettingRoles_Search,
        SystemSettingRoles_Create,
        SystemSettingRoles_Update,
        SystemSettingRoles_Delete,
        SystemSettingRoles_Rename,
    }
}
