using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;

namespace MSC.Models.Enums
{
  
    /// <summary>
    /// 某些資料要管理分類 簡單用個enums去分類
    /// </summary>
    public enum DataCategoryEnum
    {

        /// <summary>
        /// 一般
        /// </summary>
        Common = 1000,

        /// <summary>
        /// 用戶端設置
        /// </summary>
        ClientSideSetting = 1300,


        /// <summary>
        /// 安全防護設置
        /// </summary>
        Security = 3000,

        /// <summary>
        /// 財務設置
        /// </summary>
        FinanceSetting = 1500,


        /// <summary>
        /// 域名配置
        /// </summary>
        DoaminSetting = 6000,


        /// <summary>
        /// 其他配置
        /// </summary>
        Others = 9500,

    }
}
