using Newtonsoft.Json;
using target.Models.Paging;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 請求參數
    /// </summary>
    public class GameRecordRequestModel : PaginationParameter
    {
        public string MerchantSecretKey { get; set; }

        public int MerchantID { get; set; }

        public string PlayerName { get; set; }

        public DateTime StartTime { get; set; }

        public DateTime EndTime { get; set; }
    }



    /// <summary>
    /// 請求參數
    /// </summary>
    public class VendorReportSummaryRequestModel 
    {
        public string MerchantSecretKey { get; set; }

        public int MerchantID { get; set; }

        public DateTime StartTime { get; set; }

        public DateTime EndTime { get; set; }
    }


    /// <summary>
    /// 注單流水資料
    /// </summary>
    public class GameRecordViewModel
    {
        public string PlayerName { get; set; }

        public int GameID { get; set; }
        
        public decimal BeforeAmount { get; set; }

        public decimal WinLossAmount { get; set; }

        /// <summary>
        /// 局號
        /// </summary>
        public string RoundID { get; set; }

        /// <summary>
        /// 投注識別ID
        /// </summary>
        public string BetID { get; set; }

        public DateTime BetTime { get; set; }

        public DateTime RoundStartTime { get; set; }

        public DateTime RoundEndTime { get; set; }

        [JsonIgnore]
        public long PerGameLogID { get; set; }


        /// <summary>
        /// 投注金額
        /// </summary>
        public decimal BetAmount { get; set; }

        /// <summary>
        /// 派彩金額
        /// </summary>
        public decimal PayoutAmount { get; set; }


        /// <summary>
        /// 稅收金額
        /// </summary>
        public decimal TaxAmount { get; set; }

        /// <summary>
        /// 有效流水
        /// </summary>
        public decimal ValidAmount { get; set; }


    }


}
