using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{

    public enum SexEnum
    {
        [Description(i18nKeys.Enums_SexEnum_Male)]
        Male = 0,

        [Description(i18nKeys.Enums_SexEnum_Female)]
        Female = 1,      
    }
}
