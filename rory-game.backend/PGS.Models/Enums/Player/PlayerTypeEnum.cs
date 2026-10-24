using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    
    public enum PlayerTypeEnum
    {
        [Description(i18nKeys.Enums_PlayerTypeEnum_Normal)]
        Normal = 1,

        [Description(i18nKeys.Enums_PlayerTypeEnum_Robot)]
        Robot = 9,
    }
}
