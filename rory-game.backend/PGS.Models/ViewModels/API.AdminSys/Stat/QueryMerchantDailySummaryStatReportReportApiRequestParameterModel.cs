using System;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class QueryMerchantDailySummaryStatReportReportApiRequestParameterModel
    {
        public int? MerchantID { get; set; }

        public double TimeZone { get; set; }

        public DateTime? Begin { get; set; }

        public DateTime? End { get; set; }


    }
    

}
