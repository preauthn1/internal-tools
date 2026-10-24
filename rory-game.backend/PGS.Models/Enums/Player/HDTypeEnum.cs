using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    
    public enum HDTypeEnum
    {
        [Description("PC")]
        PC = 1,

        [Description("安卓")]
        Android = 2,

        [Description("苹果")]
        IOS = 3

    }
}
