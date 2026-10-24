using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    
    public enum PlayerStatusEnum
    {
        [Description(i18nKeys.Enums_PlayerStatusEnum_Disable)]
        Disable = 0,

        [Description(i18nKeys.Enums_PlayerStatusEnum_Enable)]
        Enable = 1,

        [Description(i18nKeys.Enums_PlayerStatusEnum_PermanentDelete)]
        PermanentDelete = 999,
    }
}
