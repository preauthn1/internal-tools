using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    /// <summary>
    /// 幣別
    /// </summary>
    [Description("資料庫是存字串_編號本身沒意義")]
    public enum CurrencyEnum
    {

        /// <summary>
        /// 泰銖
        /// </summary>
        [Description("THB")]
        THB,

        /// <summary>
        /// 馬來幣
        /// </summary>
        [Description("MYR")]
        MYR,

        /// <summary>
        /// 印尼盾
        /// </summary>
        [Description("IDR")]
        IDR,


        /// <summary>
        /// 越南盾
        /// </summary>
        [Description("VND")]
        VND,


        /// <summary>
        /// 新加坡幣
        /// </summary>
        [Description("SGD")]
        SGD,



        /// <summary>
        /// 人民幣
        /// </summary>
        [Description("CNY")]
        CNY,

        /// <summary>
        /// 台幣
        /// </summary>
        [Description("TWD")]
        TWD,

        /// <summary>
        /// 美元
        /// </summary>
        [Description("USD")]
        USD,

      

        /// <summary>
        /// 澳幣
        /// </summary>
        [Description("AUD")]
        AUD,


        /// <summary>
        /// 印度盧比
        /// </summary>
        [Description("INR")]
        INR,

  
        /// <summary>
        /// BETA幣
        /// </summary>
        [Description("BETA")]
        BETA,

        /// <summary>
        /// 加拿大
        /// </summary>
        [Description("CAD")]
        CAD,

        /// <summary>
        /// 捷克幣
        /// </summary>
        [Description("CZK")]
        CZK,

        /// <summary>
        /// 丹麥
        /// </summary>
        [Description("DKK")]
        DKK,

        /// <summary>
        /// 歐元
        /// </summary>
        [Description("EUR")]
        EUR,

        /// <summary>
        /// 英鎊
        /// </summary>
        [Description("GBP")]
        GBP,

        /// <summary>
        /// 港元
        /// </summary>
        [Description("HKD")]
        HKD,

        /// <summary>
        /// 日幣
        /// </summary>
        [Description("JPY")]
        JPY,

        /// <summary>
        /// 韓元
        /// </summary>
        [Description("KRW")]
        KRW,

        /// <summary>
        /// 墨西哥
        /// </summary>
        [Description("MXN")]
        MXN,

        /// <summary>
        /// 緬甸元 (對方還有支援另一種 1:1000 MYK)
        /// </summary>
        [Description("MMK")]
        MMK,

        /// <summary>
        /// 挪威
        /// </summary>
        [Description("NOK")]
        NOK,

        /// <summary>
        /// 紐西蘭
        /// </summary>
        [Description("NZD")]
        NZD,

        /// <summary>
        /// 菲律賓
        /// </summary>
        [Description("PHP")]
        PHP,

        /// <summary>
        /// 波蘭
        /// </summary>
        [Description("PLZ")]
        PLZ,

        /// <summary>
        /// 瑞典克朗
        /// </summary>
        [Description("SEK")]
        SEK,

     
        /// <summary>
        /// 點數制
        /// </summary>
        [Description("MSC_POINT")]
        MSC_POINT,
    }
}
