using MSC.Models.Attributes;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;


namespace MSC.Models.Enums
{
    
    public enum MerchantCustomizedSettingKeyEnum
    {

        /// <summary>
        /// 卡號輸入頁 域名網址
        /// </summary>
        [DataCategory(DataCategoryEnum.DoaminSetting)]
        [ReferenceType(typeof(string))]
        Basic_CardInputPage_DomainUrl = 1001,


        /// <summary>
        /// 允許的金額設置
        /// </summary>
        [DataCategory(DataCategoryEnum.ClientSideSetting)]
        [ReferenceType(typeof(string))] 
        Basic_AvailableAmounts = 1011,


        /// <summary>
        /// 強制語系
        /// </summary>
        [DataCategory(DataCategoryEnum.ClientSideSetting)]
        [ReferenceType(typeof(string))]
        ForceChangeLang = 1061,


        /// <summary>
        /// 是否顯示按鈕
        /// </summary>
        [DataCategory(DataCategoryEnum.ClientSideSetting)]
        [ReferenceType(typeof(bool))]
        IsBackButtonVisible = 1091,


        /// <summary>
        /// 使用的樣式模板
        /// </summary>
        [DataCategory(DataCategoryEnum.ClientSideSetting)]
        [ReferenceType(typeof(string))]
        Style_BrandingTypeLayout = 2211,


        /// <summary>
        /// 每個帳號 每小時 最多請求幾次單
        /// </summary>
        [DataCategory(DataCategoryEnum.ClientSideSetting)]
        [ReferenceType(typeof(string))] 
        MaxTimesOfRequestLinkPerMember = 5001,



        /// <summary>
        /// 預設每筆交易抽的手續費比例 (e.g.  0.15)
        /// </summary>
        [DataCategory(DataCategoryEnum.FinanceSetting)]
        [ReferenceType(typeof(string))]
        [IsDisabledEdit]
        DefaultTransactionFeeRate = 4101,

    }
}
