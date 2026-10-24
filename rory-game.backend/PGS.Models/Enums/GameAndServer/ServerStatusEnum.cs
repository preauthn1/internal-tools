using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    public enum ServerStatusEnum
    {
        [Description(i18nKeys.Enums_ServerStatusEnum_Disable)]
        Disable = 0,

        [Description(i18nKeys.Enums_ServerStatusEnum_Enable)]
        Enable = 1,        

        [Description(i18nKeys.Enums_ServerStatusEnum_ServerDown)]
        ServerDown = 9,


        Delete = 999,
    }
}
