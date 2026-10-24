using MSC.Models.Attributes;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;


namespace MSC.Models.Enums
{
    
    public enum SystemSettingConfigEnum
    {

        /// <summary>
        /// 系統是否開放使用中 (若 false = 維護中)
        /// </summary>
        [DataCategory(DataCategoryEnum.Common)]
        [ReferenceType(typeof(bool))]
        [SettingDefaultValue("true")]
        System_SystemEnableFlag = 2201,


        /// <summary>
        /// Point API 每小時超過多少次錯誤要發警告
        /// </summary>
        [DataCategory(DataCategoryEnum.Security)]
        [ReferenceType(typeof(string))]
        [SettingDefaultValue("true")]
        SystemWarning_PointApiSystemWarningConditionPerHour = 3901,


        /// <summary>
        /// 每組IP 每小時 最多請求幾次卡號輸入頁
        /// </summary>
        [DataCategory(DataCategoryEnum.Security)]
        [ReferenceType(typeof(string))]
        [SettingDefaultValue("50")]
        MaxTimesOfRequestCardInputPerIPAddress = 7051,



        /// <summary>
        /// 合作商家的iframe配置
        /// </summary>
        [DataCategory(DataCategoryEnum.Common)]
        [ReferenceType(typeof(string))]
        [SettingDefaultValue("")]
        ShoppingsIframeUrl = 5102,




    }
}
