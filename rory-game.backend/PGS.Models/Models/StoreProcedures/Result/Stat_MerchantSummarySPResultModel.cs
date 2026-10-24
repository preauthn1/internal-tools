using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models.StoreProcedures
{
    public class Stat_MerchantSummarySPResultModel
    {
        public int MerchantID { get; set; }
        public string MerchantName { get; set; }
		public int TotalLoginPlayer { get; set; }
		public decimal TransferAmount { get; set; }
		public decimal TransferInAmount { get; set; }
		public decimal TransferOutAmount { get; set; }
		public decimal WalletAmount { get; set; }


		public int TotalLoginCount { get; set; }
		public int NewLoginPlayer { get; set; }
		public int MaxLoginPlayer { get; set; }
		public int AvgLoginPlayer { get; set; }

		

		public int TotalRound { get; set; }
		public int TotalBetPlayer { get; set; }
		public decimal WinAmount { get; set; }
		public decimal LoseAmount { get; set; }
		public decimal MoneyChangeAmount { get; set; }

		public decimal BetAmount { get; set; }
		public decimal PayoutAmount { get; set; }
		public decimal ValidAmount { get; set; }
		public decimal WinLossAmount { get; set; }

	}
}
