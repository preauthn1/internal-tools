using Newtonsoft.Json.Linq;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class FishingSingleWalletCheckAccountTokenRequestModel
    {

        public string Account { get; set; }


        public string Token { get; set; }

        /// <summary>
        /// 時間戳記 精確到秒
        /// </summary>
        public string Time { get; set; }

        /// <summary>
        /// MD5 加密簽
        /// </summary>
        public string Sign { get; set; }


    }


    public class FishingSingleWalletCheckAccountTokenResponseModel : FishingAPIResponseModel
    {

        public string Account { get; set; }


        public string Nick { get; set; }

        /// <summary>
        /// 頭像圖片網址
        /// </summary>
        public string Icon { get; set; }

        /// <summary>
        /// 平台幣
        /// </summary>
        public decimal Wealth { get; set; }

    }


}
