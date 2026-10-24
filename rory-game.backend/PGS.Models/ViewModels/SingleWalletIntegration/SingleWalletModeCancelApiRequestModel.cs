namespace game.Entity.Game
{
    public class SingleWalletModeCancelApiRequestModel
    {
        public int GameID { get; set; }


        public int MerchantID { get; set; }


        public int PlayerID { get; set; }


        /// <summary>
        /// 傳給三方還是統一叫做 transactionNo 交易單號
        /// </summary>
        public string TransactionNo { get; set; }


        public string DebugProcessUUID { get; set; }


        /// <summary>
        /// 是否為 外部遊戲商的遊戲 (非PGS自製)
        /// </summary>
        public bool IsExternalGameVendor { get; set; }

    }

    public class SingleWalletModeCancelApiResponseModel
    {

        public string TransactionNo { get; set; }

    }


}
