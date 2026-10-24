using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class CardLinkRequestRecordEntity
    {
        [Key]
        public long Id { get; set; }
        public string SystemRequestID { get; set; }
        public int MerchantID { get; set; }
        public string MerchantRequestID { get; set; }
        public string MemberCode { get; set; }
        public string ReturnUrl { get; set; }
        public string NotifyUrl { get; set; }
        public string IPAddress { get; set; }
        public decimal Amount { get; set; }
        public string TokenHashValue { get; set; }
        public bool IsDisabled { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
