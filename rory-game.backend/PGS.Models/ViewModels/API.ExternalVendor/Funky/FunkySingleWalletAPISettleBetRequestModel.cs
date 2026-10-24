namespace target.Models.ViewModels
{


    public class FunkySingleWalletAPISettleBetRequestModel
    {

        public string refNo { get; set; }


        public FunkySingleWalletAPISettleBetDetailRequestModel betResultReq { get; set; }

    }

    public class FunkySingleWalletAPISettleBetDetailRequestModel
    {

        public decimal winAmount { get; set; }

        public decimal stake { get; set; }

        public decimal effectiveStake { get; set; }

        public int playerId { get; set; }

        public string gameCode { get; set; }

    }


    public class FunkySingleWalletAPISettleBetResponseModel
    {
        public string RefNo { get; set; }

        public decimal Balance { get; set; }

        public string PlayerId { get; set; }

        public string StatementDate { get; set; }

    }


}
