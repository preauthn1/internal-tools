using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.Attributes;
using target.Models.CultureText;


namespace target.Models.Enums
{
  
    public enum SystemConfigSettingEnum
    {

        /// <summary>
        /// 遊戲網址
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.Common)]
        GamePlayDomain = 1001,


        /// <summary>
        /// 遊戲Lobby網址
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.Common)]
        LobbyPageUrl = 1002,


        /// <summary>
        /// 遊戲Lobby 遊戲圖片網址來源
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.Common)]
        LobbyPageGameImageUrl = 1003,


        /// <summary>
        /// 後台網址 (產生detail info 外部連結需要)
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.Common)]
        AdminSysDomain_ExternalDetailInfoUrl = 2201,


        /// <summary>
        /// 內部測試帳號PlayerID (逗號隔開)
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.InternalTest)]
        InternalTestAccountPlayerIDs = 7801,


        /// <summary>
        /// 常用商戶IDs(影響下拉選單排序) (逗號隔開)
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.InternalTest)]
        InternalTest_Favorit_MerchantIDs = 7802,


        /// <summary>
        /// 模擬單錢包API，模擬錯誤機率(投注API)
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DataCategoryAttribute(DataCategoryEnum.InternalTest)]
        FakeSingleWalletApi_ErrorRate_PlaceBet = 7806,


        /// <summary>
        /// 模擬單錢包API，模擬錯誤機率(結算API)
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DataCategoryAttribute(DataCategoryEnum.InternalTest)]
        FakeSingleWalletApi_ErrorRate_Settlement = 7807,



        /// <summary>
        /// 遊戲商 [投注]，模擬成功_但是處理逾時
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DataCategoryAttribute(DataCategoryEnum.InternalTest)]
        GameVendor_PlaceBet_SuccessUpdateBut_DelayTimeout = 7809,


        /// <summary>
        /// 遊戲商 [結算]，模擬成功_但是處理逾時
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DataCategoryAttribute(DataCategoryEnum.InternalTest)]
        GameVendor_Settle_SuccessUpdateBut_DelayTimeout = 7810,



        /// <summary>
        /// 是否開啟全系統維護開關
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.SystemMaintain)]
        SystemMaintain_IsSwtichToMaintain = 8011,


        
        /// <summary>
        /// 維護時間文案設置 (e.g. 設置  12:00- 14:00)
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.SystemMaintain)]
        SystemMaintain_MaintainTime = 8012,
        

    }
}
