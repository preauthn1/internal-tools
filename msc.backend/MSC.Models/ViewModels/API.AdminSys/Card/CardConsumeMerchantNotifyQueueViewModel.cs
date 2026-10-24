using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class CardConsumeMerchantNotifyQueueViewModel
    {

        public int MerchantID { get; set; }
        public string MerchantName { get; set; }
        public string MerchantRequestID { get; set; }
        public int NotifyType { get; set; }
        public long? ReferenceID { get; set; }
        public string ReferenceOrderNo { get; set; }
        public string MemberCode { get; set; }
        public string CardNo { get; set; }
        public int Status { get; set; }
        public string StatusText { get; set; }

        public string DebugInfo { get; set; }
        public string ExtraData { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime UpdateTime { get; set; }
        public int RetryTimes { get; set; }
        public string LastTimeResponse { get; set; }
        public string LastTimeRequestUUID { get; set; }

    }



}
