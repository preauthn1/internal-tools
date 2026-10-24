using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{
    public class ProviderCardGenearteRecordEntity
    {
        [Key]
        public long Id { get; set; }
        public string SystemOrderNo { get; set; }
        public int ProviderID { get; set; }
        public string TransactionNo { get; set; }
        public string Currency { get; set; }
        public decimal Point { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime UpdateTime { get; set; }


    }
}
