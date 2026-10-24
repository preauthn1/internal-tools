using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class MerchantWithdrawalViewModel
    {
        public long Id { get; set; }
        public int Type { get; set; }
        public int MerchantID { get; set; }
        public string MerchantName { get; set; }
        public decimal Amount { get; set; }
        public int Status { get; set; }
        public string StatusText { get; set; }
        public string AccountMemo { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }
    }

    public class MerchantWithdrawalCreateRequestModel
    {
        public int MerchantID { get; set; }
        public decimal Amount { get; set; }
        public string AccountMemo { get; set; }
        public string Memo { get; set; }
    }



}
