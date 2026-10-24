namespace target.Models.ViewModels
{






    public class SingleWalletBetSettleRequestParameter : SingleWalletBasePlayerInfoModel
    {       
        /// <summary>
        /// 投注單號
        /// </summary>
        public string BetRecordNo { get; set; }


        /// <summary>
        /// 派彩金額 (更新錢包金額)
        /// </summary>
        public decimal PayoutAmount { get; set; }

        /// <summary>
        /// 三方廠商的遊戲代碼
        /// </summary>
        public string GameCode { get; set; }

    }


    public class SingleWalletVendorCustomizedResponseConvertModel 
    {
        /// <summary>
        /// 根據不同遊戲商 填入各種不同的自訂回應物件
        /// 請參考 interface FormatParsing_ConvertToSettleBetResponseResult 是如何實作的
        /// </summary>
        public dynamic GameVendorCustomizedResponse { get; set; }

    }




}
