namespace game.Entity.Game
{

    public class SingleWalletModePlaceBetParameter
    {
        public int PlayerID { get; set; }

        public string SessionToken { get; set; }

        public int GameID { get; set; }

        public decimal BetAmount { get; set; }

        public string BetRecordNoPGS { get; set; }

        public string DebugProcessUUID { get; set; }

    }


    public class SingleWalletModePlaceBetApiRequestModel
    {

        public string PlayerUniqueID { get; set; }

        public string SessionToken { get; set; }

        public int GameID { get; set; }

        public decimal BetAmount { get; set; }


        /// <summary>
        /// 傳給三方還是統一叫做 transactionNo 交易單號
        /// </summary>
        public string TransactionNo { get; set; }

    }


    public class SingleWalletModePlaceBetApiResponseModel
    {
        public string TransactionNo { get; set; }

        public decimal Balance { get; set; }

    }


}
