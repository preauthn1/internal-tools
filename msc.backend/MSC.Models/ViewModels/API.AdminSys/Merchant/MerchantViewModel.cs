using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class MerchantViewModel
    {
        public int MerchantID { get; set; }
        public string DisplayName { get; set; }
        public int Status { get; set; }
        public string Memo { get; set; }
        public string ApiSecretKey { get; set; }
        public string Currency { get; set; }

        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }
        public string StatusText { get; set; }



        public decimal WalletAmount { get; set; }
        public decimal FeeRate { get; set; }
        public int MaxTimesOfRequestLinkPerMember { get; set; }
        public string AvailableAmounts { get; set; }
        public string CardInputPage_DomainUrl { get; set; }

    }




}
