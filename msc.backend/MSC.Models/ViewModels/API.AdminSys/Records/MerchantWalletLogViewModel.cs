using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class MerchantWalletLogViewModel
    {
        public long Id { get; set; }
        public int TransactionType { get; set; }

        public string TransactionTypeText { get; set; }

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public decimal SourceAmount { get; set; }
        public decimal FeeRate { get; set; }
        public decimal Fee { get; set; }
        public decimal UpdateAmount { get; set; }
        public decimal BeforeAmount { get; set; }
        public decimal AfterAmount { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }


    }



}
