using MSC.Models.Enums;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.Models
{

    public class MerchantDailySummarySPRequestParameter
    {

        public int MerchantID { get; set; }

        /// <summary>
        /// 請傳入UTC時間
        /// </summary>
        public DateTime BeginTime { get; set; }

        public DateTime EndTime { get; set; }

        public double TimeZone { get; set; }

    }



    public class MerchantDailySummarySPResult
    {
  
        public int MerchantID { get; set; }
        public DateTime Date { get; set; }
        public int TotalCardCount { get; set; }
        public int TotalMemberCount { get; set; }
        public decimal TotalPoint { get; set; }

    }



    public class MerchantDailySummarySPResultV2
    {
        public int MerchantID { get; set; }
        public DateTime Date { get; set; }
        public int TotalOrderCount { get; set; }
        public int TotalSuccessCount { get; set; }
        public int TotalMembers { get; set; }
        public decimal TotalAmount { get; set; }
    }



}
