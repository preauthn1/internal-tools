using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class MerchantGameSummaryDataEntity
    {
        public int MerchantID { get; set; }
        public int GameID { get; set; }
        public DateTime StatTime { get; set; }
        public int TotalRound { get; set; }
        public int BetPlayer { get; set; }
        public decimal WinAmount { get; set; }
        public decimal LoseAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal FlowAmount { get; set; }
        public decimal JackpotAmount { get; set; }
        public decimal BetAmount { get; set; }
        public decimal ValidAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public decimal WinLossAmount { get; set; }
        public DateTime UpdateTime { get; set; }

    }
}
