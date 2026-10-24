using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.Attributes;
using target.Models.CultureText;


namespace target.Models.Enums
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
        /// 遊戲點數配置
        /// </summary>
        GameSetting = 5100,


        /// <summary>
        /// 轉換比例配置
        /// </summary>
        ConvertRate = 5120,


        /// <summary>
        /// 單錢包相關配置
        /// </summary>
        SingleWalletSetting = 5200,



        /// <summary>
        /// 域名配置
        /// </summary>
        DoaminSetting = 6000,

        /// <summary>
        /// 密鑰配置
        /// </summary>
        SecretKey = 6050,

        /// <summary>
        /// 內部測試相關設置
        /// </summary>
        InternalTest = 7800,


        /// <summary>
        /// 維護設置
        /// </summary>
        SystemMaintain = 8010,

        /// <summary>
        /// 特殊分類 禁止編輯
        /// </summary>
        Special_DisableEdit = 9999,

    }
}
