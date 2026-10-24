using System;

namespace target.Models.ViewModels.GameVendor
{
    public class BetRecordGeneralColumnModel
    {
        public long Id { get; set; }


        /// <summary>
        /// 三方的注單號
        /// </summary>
        public string BetRecordNo { get; set; }
        public string BetRecordNoPGS { get; set; }
        public string ThirdPartyGameCode { get; set; }
        public int MerchantID { get; set; }
        public int GameID { get; set; }
        public string GameName { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public string PlayerIPAddress { get; set; }
        public decimal BetAmount { get; set; }
        public decimal ValidAmount { get; set; }
        public decimal WinLossAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public int Status { get; set; }
        public string StatusText { get; set; }

        public DateTime CreateTime { get; set; }
        public DateTime UpdateTime { get; set; }
        public DateTime? SettlementTime { get; set; }
        public int? ManualSettleUserID { get; set; }
        public decimal ConvertRate { get; set; }

    }





}
