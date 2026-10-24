using System;
using System.Collections.Generic;

namespace target.Models.ViewModels
{

	public class MerchantDailySummaryStatReportInfoViewModel
	{		
		public List<MerchantDailySummaryStatReportViewModel> Items { get; set; }

		public MerchantDailySummaryStatReportViewModel TotalSummary { get; set; }

	}


	public class MerchantDailySummaryStatReportViewModel
    {
		public string Date { get; set; }
		public int MerchantID { get; set; }
		public string MerchantName { get; set; }
		public int TotalLoginCount { get; set; }
		public decimal TransferInAmount { get; set; }
		public decimal TransferOutAmount { get; set; }
		public int TotalLoginPlayer { get; set; }
		public decimal TaxAmount { get; set; }
		public decimal BetAmount { get; set; }
		public decimal PayoutAmount { get; set; }
		public decimal ValidAmount { get; set; }
		public decimal MerchantWinLossAmount { get; set; }
		public decimal WinLossAmount { get; set; }
		public decimal RTP { get; set; }

	}



}
