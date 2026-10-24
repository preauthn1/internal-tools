using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{
    public class MerchantWithdrawalEntity
    {
        [Key]
        public long Id { get; set; }
        public int Type { get; set; }
        public int MerchantID { get; set; }
        public decimal Amount { get; set; }
        public int Status { get; set; }
        public string AccountMemo { get; set; }
        public string Memo { get; set; }

        public string OrderNo { get; set; }


        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }


    }
}
