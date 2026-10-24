using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class CardConsumeMerchantNotifyQueueEntity
    {
        [Key]
        public long Id { get; set; }
        public int MerchantID { get; set; }
        public string MerchantRequestID { get; set; }
        public int NotifyType { get; set; }
        public long? ReferenceID { get; set; }
        public string ReferenceOrderNo { get; set; }
        public string MemberCode { get; set; }
        public string CardNo { get; set; }
        public int Status { get; set; }
        public string DebugInfo { get; set; }
        public string ExtraData { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime UpdateTime { get; set; }
        public int RetryTimes { get; set; }
        public string LastTimeResponse { get; set; }
        public string LastTimeRequestUUID { get; set; }
        public string NotifyUrl { get; set; }


    }
}
