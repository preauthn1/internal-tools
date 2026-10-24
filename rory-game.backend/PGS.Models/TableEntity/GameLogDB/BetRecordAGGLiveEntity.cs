using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    public class BetRecordAGGLiveEntity : IBetRecordEntity
    {
        [Key]
        public long Id { get; set; }
        public string BetRecordNo { get; set; }
        public string BetRecordNoPGS { get; set; }

        public string ThirdPartyGameCode { get; set; }
        public int MerchantID { get; set; }
        public int GameID { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public string PlayerIPAddress { get; set; }
        public decimal BetAmount { get; set; }
        public decimal ValidAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public decimal WinLossAmount { get; set; }


        /// <summary>
        /// AGG特有的欄位
        /// </summary>
        public decimal JackpotAmount { get; set; }

        public int Status { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime UpdateTime { get; set; }
        public DateTime? SettlementTime { get; set; }
        public string DebugInfo { get; set; }
        public int? ManualSettleUserID { get; set; }
        public decimal ConvertRate { get; set; }

    }
}
