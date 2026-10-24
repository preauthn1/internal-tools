using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class CardLinkRequestRecordViewModel
    {
        public string SystemRequestID { get; set; }
        public int MerchantID { get; set; }
        public string MerchantName { get; set; }
        public string MerchantRequestID { get; set; }
        public string MemberCode { get; set; }
        public string ReturnUrl { get; set; }
        public string NotifyUrl { get; set; }
        public decimal Amount { get; set; }
        public string IPAddress { get; set; }
        public string TokenHashValue { get; set; }

        /// <summary>
        /// 連結失效 代表有完成輸入卡號
        /// </summary>
        public bool IsDisabled { get; set; }

        /// <summary>
        /// 是否已完成商戶通知
        /// </summary>
        public bool IsNotifyCompleted { get; set; }

        public bool IsNotifyExist { get; set; }


        public DateTime CreateTime { get; set; }

    }



}
