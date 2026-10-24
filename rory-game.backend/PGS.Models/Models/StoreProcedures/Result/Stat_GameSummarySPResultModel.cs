using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models.StoreProcedures
{
    public class Stat_GameSummarySPResultModel
    {
        public int MerchantID { get; set; }
        public string MerchantName { get; set; }
        public int GameID { get; set; }
        public string GameName { get; set; }
        public int TotalRound { get; set; }
        public int OnlinePlayer { get; set; }
        public decimal UpdateAmount { get; set; }
        public decimal WinAmount { get; set; }
        public decimal LoseAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal FlowAmount { get; set; }


        public decimal BetAmount { get; set; }
        public decimal ValidAmount { get; set; }
        public decimal PayoutAmount { get; set; }


        public decimal JackpotPoolAmount { get; set; }

    }
}
