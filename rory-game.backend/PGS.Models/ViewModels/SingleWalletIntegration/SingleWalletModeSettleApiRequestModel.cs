using System;

namespace game.Entity.Game
{


    public class SingleWalletModeSettleApiParameter
    {
        public int MerchantID { get; set; }

        public int GameID { get; set; }

        public int PlayerID { get; set; }



        /// <summary>
        /// 若是玩三方遊戲商遊戲，請傳入PGS自訂的注單號
        /// </summary>
        public string BetRecordNoPGS { get; set; }

        public decimal BetAmount { get; set; }

        public decimal TaxAmount { get; set; }

        public decimal PayoutAmount { get; set; }

        public decimal ValidAmount { get; set; }

        public DateTime? SettlementTime { get; set; }

        /// <summary>
        /// debug用 (方便查詢 通知請求的API紀錄)
        /// </summary>
        public string DebugProcessUUID { get; set; }

    }


    public class SingleWalletModeSettleApiRequestModel
    {

        public int GameID { get; set; }

        public decimal BetAmount { get; set; }

        public decimal ValidAmount { get; set; }

        public decimal PayoutAmount { get; set; }

        public decimal TaxAmount { get; set; }

        /// <summary>
        /// 傳給三方還是統一叫做 transactionNo 交易單號
        /// </summary>
        public string TransactionNo { get; set; }

        public string SettlementTime { get; set; }

        public string PlayerUniqueID { get; set; }
    }



    public class SingleWalletModeSettleApiResponseModel
    {
        /// <summary>
        /// 對方需返回相同的單號、我方才會視為拿到成功回應
        /// </summary>
        public string TransactionNo { get; set; }

        
        public decimal Balance { get; set; }

    }



}
