using System;

namespace target.Models.ViewModels
{
    public class UnionBetRecordQueryModel
    {
        /// <summary>
        /// 對應 BetRecordTypeEnum (僅PGS自製遊戲會有此)
        /// </summary>
        public int Type { get; set; }
        public string BetRecordNo { get; set; }
        public string BetRecordNoPGS { get; set; }     
        public int MerchantID { get; set; }
        public int GameVendor { get; set; }
        public int GameID { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public decimal BetAmount { get; set; }
        public decimal WinLossAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public int Status { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime? SettlementTime { get; set; }

    }


}
