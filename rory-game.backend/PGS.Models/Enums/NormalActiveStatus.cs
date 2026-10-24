using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    /// <summary>
    /// 如果狀態只有 開啟/關閉/永久隱藏 這幾種
    /// 請統一使用這個enum
    /// </summary>
    public enum NormalActiveStatus
    {
        [Description(i18nKeys.Enums_NormalActiveStatus_Disable)]
        Disable = 0,

        [Description(i18nKeys.Enums_NormalActiveStatus_Enable)]
        Enable = 1,

        [Description(i18nKeys.Enums_NormalActiveStatus_PermanentDelete)]
        PermanentDelete = 999,
    }  
    

    /// <summary>
    /// 如果狀態只有 開啟/關閉/永久隱藏 這幾種
    /// 請統一使用這個enum
    /// </summary>
    public enum NormalYesOrNoEnum
    {
        [Description(i18nKeys.Enums_NormalActiveStatus_Enable)]
        Yes = 1,

        [Description(i18nKeys.Enums_NormalActiveStatus_Disable)]
        No = 0,      
    }
}
