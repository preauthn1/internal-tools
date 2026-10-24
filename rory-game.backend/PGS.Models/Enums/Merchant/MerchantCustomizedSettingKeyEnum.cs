using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.Attributes;
using target.Models.CultureText;


namespace target.Models.Enums
{
    
    public enum MerchantCustomizedSettingKeyEnum
    {

        /// <summary>
        /// 建立帳號，預設錢包金額
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DisableEdit]
        [SettingDefaultValue("0")]
        [DataCategoryAttribute(DataCategoryEnum.Special_DisableEdit)]
        DefaultWalletBalance = 5001,


        /// <summary>
        /// 是否為內部測試用商戶 (會開放額外功能)
        /// </summary>
        [ReferenceType(typeof(NormalYesOrNoEnum))]
        [DisableEdit]
        [SettingDefaultValue("0")]
        [DataCategoryAttribute(DataCategoryEnum.Special_DisableEdit)]
        IsInternalTestMerchant = 5025,



        /// <summary>
        /// 當遊戲載入時，是否使用靜態圖版本 (默認不使用: 則採用天馬icon動畫模式)
        /// </summary>
        [ReferenceType(typeof(MerchantLoadingStyleModeEnum))]
        [SettingDefaultValue((int)MerchantLoadingStyleModeEnum.Default)]
        [DataCategoryAttribute(DataCategoryEnum.Common)]
        Style_LoadingStyleMode = 6001,



        /// <summary>
        /// 是否啟用client端log開關
        /// </summary>
        [ReferenceType(typeof(MerchantLoadingStyleModeEnum))]
        [DataCategoryAttribute(DataCategoryEnum.Common)]
        [SettingDefaultValue((int)NormalActiveStatus.Disable)]
        System_ClientLogEnable = 9001,



        /// <summary>
        /// 單錢包模式，平台商的API域名
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.SingleWalletSetting)]
        [SettingDefaultValue("")]
        Integration_SingleWalletApiDomain = 7001,



        /// <summary>
        /// 單錢包模式，PGS送出的token
        /// </summary>
        [ReferenceType(typeof(string))]
        [DataCategoryAttribute(DataCategoryEnum.SingleWalletSetting)]
        [SettingDefaultValue("")]
        Integration_SingleWalletApiPgsValidToken = [REDACTED],



        /// <summary>
        /// 是否三方遊戲 永遠使用試玩模式
        /// </summary>
        [ReferenceType(typeof(NormalYesOrNoEnum))]
        [DataCategoryAttribute(DataCategoryEnum.GameSetting)]
        [SettingDefaultValue((int)NormalYesOrNoEnum.No)]
        System_IsAlwaysTrialMode = 9005,



        /// <summary>
        /// 是否 針對Slots遊戲 不擋卡房
        /// </summary>
        [ReferenceType(typeof(NormalYesOrNoEnum))]
        [DataCategoryAttribute(DataCategoryEnum.GameSetting)]
        [SettingDefaultValue((int)NormalYesOrNoEnum.Yes)]
        GameSetting_IsSlotsGameRoomLock = 9021,



        /// <summary>
        /// PGS遊戲商的 轉換比例
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DataCategoryAttribute(DataCategoryEnum.ConvertRate)]
        [SettingDefaultValue("1")]
        GameVendor_ConvertRate_PGS = 9000000,



        /// <summary>
        /// AGG遊戲商的 轉換比例
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DataCategoryAttribute(DataCategoryEnum.ConvertRate)]
        [SettingDefaultValue("1")]
        GameVendor_ConvertRate_AGG = 9290000,



        /// <summary>
        /// CrazyFish 遊戲商的 轉換比例
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DataCategoryAttribute(DataCategoryEnum.ConvertRate)]
        [SettingDefaultValue("1")]
        GameVendor_ConvertRate_CrazyFishing = 9310000,


        /// <summary>
        /// AGGLive 遊戲商的 轉換比例
        /// </summary>
        [ReferenceType(typeof(decimal))]
        [DataCategoryAttribute(DataCategoryEnum.ConvertRate)]
        [SettingDefaultValue("1")]
        GameVendor_ConvertRate_AGGLive = 9330000,


    }
}
