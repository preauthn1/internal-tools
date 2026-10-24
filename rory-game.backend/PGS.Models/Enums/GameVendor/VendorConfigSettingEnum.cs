using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.Attributes;
using target.Models.CultureText;


namespace target.Models.Enums
{
    
    public enum VendorConfigSettingEnum
    {

        /// <summary>
        /// 主要對接的API域名
        /// </summary>
        [ReferenceType(typeof(string))]
        [ExtraBooleanValue("允許編輯", false)]
        [DataCategoryAttribute(DataCategoryEnum.DoaminSetting)] 
        ApiUrl = 1001,


        /// <summary>
        /// (如果有分開) 投注紀錄/統計數據 的API域名
        /// </summary>
        [ReferenceType(typeof(string))]
        [ExtraBooleanValue("允許編輯", false)]
        [DataCategoryAttribute(DataCategoryEnum.DoaminSetting)] 
        BetReportApiUrl = 1002,


        /// <summary>
        /// 遊戲商分配給我們的 對接廠商識別ID
        /// </summary>
        [ReferenceType(typeof(string))]
        [ExtraBooleanValue("允許編輯", false)]
        [DataCategoryAttribute(DataCategoryEnum.SecretKey)]
        MainAssignVendorID = 1011,


        /// <summary>
        /// 遊戲商分配給我們的 對接API Token
        /// </summary>
        [ReferenceType(typeof(string))]
        [ExtraBooleanValue("允許編輯", false)]
        [DataCategoryAttribute(DataCategoryEnum.SecretKey)]
        MainSecretKeyOrApiToken = [REDACTED], 

    }
}
