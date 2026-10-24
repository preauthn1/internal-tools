using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{
    public class MerchantCardStatisticReportEntity
    {
        [Key]
        public long Id { get; set; }
        public int MerchantID { get; set; }
        public DateTime StatisticTime { get; set; }
        public int TotalOrderCount { get; set; }
        public int TotalSuccessCount { get; set; }
        public int TotalMembers { get; set; }
        public decimal TotalAmount { get; set; }
        public DateTime UpdateTime { get; set; }


    }
}
