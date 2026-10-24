using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.Attributes;
using target.Models.CultureText;


namespace target.Models.Enums
{
    
    public enum MerchantTypeEnum
    {     
        None = 0,

        /// <summary>
        /// 平台對接 正式環境
        /// </summary>
        Plaform_ProductionEnvironment = 20,

        /// <summary>
        /// 自己人內部測試用
        /// </summary>
        InternalTest = 30,

        /// <summary>
        /// 平台對接 測試環境
        /// </summary>
        Plaform_DevelopEnvironment = 40,

    }
}
