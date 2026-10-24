namespace target.Models.ViewModels
{
    public class SingleWalletSettleBetRequestParameter : SingleWalletBasePlayerInfoModel
    {

        public string GameCode { get; set; }

        public string OrderNo { get; set; }

        public decimal BetAmount { get; set; }

        public string PlayerIPAddress { get; set; }
    }





}
