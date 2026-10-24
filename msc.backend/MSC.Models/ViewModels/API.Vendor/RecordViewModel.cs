using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class MerchantRequestRecordViewModel
    {

        public int MerchantID { get; set; }
        public string MerchantRequestID { get; set; }
        public string MemberCode { get; set; }
        public string ReturnUrl { get; set; }
        public string NotifyUrl { get; set; }
        public DateTime CreateTime { get; set; }
        public string IPAddress { get; set; }
        public decimal MerchantOrderAmount { get; set; }

        public CardInputInfoViewModel CardInputInfo { get; set; }
        public SuccessNotifyInfoViewModel NotifyInfo { get; set; }
    }



    public class CardInputInfoViewModel
    {
        public string CardNo { get; set; }

        public decimal Amount { get; set; }

        public string IPAddress { get; set; }

        public bool IsSuccess { get; set; }

        public DateTime? CreateTime { get; set; }
    }



    public class SuccessNotifyInfoViewModel
    {
        public bool IsSuccess { get; set; }

        public int TotalRequestTimes { get; set; }

        public string LastTimeResponseText { get; set; }

        public DateTime? LastNotifyTime { get; set; }
    }



}
