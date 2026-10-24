namespace target.Models.ViewModels
{



    public class SingleWalletBetDebitAmountRequestParameter : SingleWalletBasePlayerInfoModel
    {

        public string GameCode { get; set; }

        public string BetRecordNo { get; set; }

        public decimal BetAmount { get; set; }

        public string PlayerIPAddress { get; set; }
    }


    public class SingleWalletPlaceBetDebitAmountResponseParameter
    {
        
        public string Currency { get; set; }

        /// <summary>
        /// 更新後金額 (法幣)
        /// </summary>
        public decimal AfterAmount { get; set; }

        /// <summary>
        /// 更新後金額 (遊戲幣)
        /// </summary>
        public decimal AfterGamePoint { get; set; }

    }





}
