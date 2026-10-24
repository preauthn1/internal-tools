using MSC.Models.Attributes;
using MSC.Models.CultureText;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;


namespace MSC.Models.Enums
{
    /// <summary>
    /// 如果狀態只有 開啟/關閉/永久隱藏 這幾種
    /// 請統一使用這個enum
    /// </summary>
    public enum NormalActiveStatus
    {
        [ReferenceText(i18nKeys.Enums_NormalActiveStatus_Disable)]
        Disable = 0,

        [ReferenceText(i18nKeys.Enums_NormalActiveStatus_Enable)]
        Enable = 1,

        [ReferenceText(i18nKeys.Enums_NormalActiveStatus_PermanentDelete)]
        PermanentDelete = 999,

    }
}
