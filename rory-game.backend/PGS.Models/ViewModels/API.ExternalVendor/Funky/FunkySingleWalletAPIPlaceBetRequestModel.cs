namespace target.Models.ViewModels
{
    public class FunkySingleWalletAPIPlaceBetRequestModel
    {

        public FunkySingleWalletAPIPlaceBetDetailRequestModel Bet { get; set; }       

        public string SessionId { get; set; }

        public string PlayerIp { get; set; }

    }


    public class FunkySingleWalletAPIPlaceBetDetailRequestModel
    {
        public string GameCode { get; set; }

        public string RefNo { get; set; }

        public decimal Stake { get; set; }
    }


    public class FunkySingleWalletAPIPlaceBetResultViewModel
    {       
        public decimal Balance { get; set; }

    }





}
