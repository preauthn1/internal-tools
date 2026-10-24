using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;

namespace target.Models.Enums
{
    
    public enum OrganizationTypeEnum
    {
        [Description("超级管理")]
        SuperAdmin = 1,

        [Description("一般维运平台")]
        Normal = 2,
    }
}
