using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class StatReportViewModel
    {
        public List<StatReportMerchantDailySummaryViewModel> Items { get; set; }

        public StatReportMerchantDailySummaryViewModel SummaryInfo { get; set; }
    }


    public class StatReportMerchantDailySummaryViewModel
    {
        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public string Currency { get; set; }


        public string Date { get; set; }

        /// <summary>
        /// 訂單完成率
        /// </summary>
        public decimal DepositCompletedRatio { get; set; }

        /// <summary>
        /// 總訂單數
        /// </summary>
        public int TotalOrderCount { get; set; }

        /// <summary>
        /// 總完成數
        /// </summary>
        public int TotalSuccessCardCount { get; set; }


        /// <summary>
        /// 總人數
        /// </summary>
        public int TotalMemberCount { get; set; }

        /// <summary>
        /// 總金額
        /// </summary>
        public decimal TotalPoint { get; set; }


    }


}
