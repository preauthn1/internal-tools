using Newtonsoft.Json;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;

namespace target.Models.ViewModels.Base
{
    public class MerchantDailySummaryExcelModel
    {


		[ExcelColumnSequence(100)]
		public string Date { get; set; }


		[ExcelColumnSequence(200)]
		public int MerchantID { get; set; }


		[ExcelColumnSequence(300)]
		public string MerchantName { get; set; }


		[ExcelColumnSequence(400)]
		public int TotalLoginCount { get; set; }



		[ExcelColumnSequence(500)]
		public decimal TransferInAmount { get; set; }


		[ExcelColumnSequence(600)]
		public decimal TransferOutAmount { get; set; }



		[ExcelColumnSequence(700)]
		public int TotalLoginPlayer { get; set; }


		[ExcelColumnSequence(800)]
		public decimal TaxAmount { get; set; }


		[ExcelColumnSequence(900)]
		public decimal BetAmount { get; set; }


		[ExcelColumnSequence(1000)]
		public decimal PayoutAmount { get; set; }


		[ExcelColumnSequence(1100)]
		public decimal ValidAmount { get; set; }


		[ExcelColumnSequence(1200)]
		public decimal MerchantWinLossAmount { get; set; }


		[ExcelColumnSequence(1300)]
		public decimal WinLossAmount { get; set; }


		[ExcelColumnSequence(1400)]
		public decimal RTP { get; set; }

	}
}
