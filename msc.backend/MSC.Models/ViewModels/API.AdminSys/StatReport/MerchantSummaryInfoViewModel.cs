using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class MerchantSummaryInfoViewModel
    {

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public string Currency { get; set; }

        public int Status { get; set; }

        public string StatusText { get; set; }

        public List<decimal> AvailableAmounts { get; set; }

        public string AssignDomain { get; set; }

        public int TodayDepositCardCount { get; set; }

        public decimal TodayDepositAmounts { get; set; }

        public decimal WalletBalance { get; set; }




    }


}
