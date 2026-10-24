using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class CardConsumeRecordViewModel
    {
        public long ID { get; set; }
        public string CardNo { get; set; }
        public decimal Point { get; set; }
        public string OrderNo { get; set; }
        public string MemberCode { get; set; }
        public int MerchantID { get; set; }
        public string MerchantName { get; set; }
        public string MerchantRequestID { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }
        public string IPAddress { get; set; }
        public string DeviceID { get; set; }
        public bool IsNotifyCompleted { get; set; }


    }



}
