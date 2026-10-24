using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models.StoreProcedures
{
    public class Transaction_RecordSPResultModel
    {
        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int MerchantID { get; set; }
        public string MerchantName { get; set; }

        public int WalletMode { get; set; }
        public int TransactionType { get; set; }
        public decimal CurrentAmount { get; set; }
        public decimal AfterAmount { get; set; }
        public decimal UpdateAmount { get; set; }
        public string IPAddress { get; set; }
        public DateTime CreateTime { get; set; }
	}
}
