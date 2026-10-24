using target.Models.Enums;
using System;
using System.Collections.Generic;

namespace target.Models.ViewModels.GameVendor
{
    public class ExternalVendorGetReportRequestParameter
    {
       
        /// <summary>
        /// 這邊可能不是指 查單一MerchantID 而是指要用哪一組merchant的設置去查資料 (可能有多個merchant共用同一組設置)
        /// </summary>
        public int MerchantID { get; set; }


        /// <summary>
        /// 請統一傳UTC過來，在根據遊戲商不同 自行轉成對方接受的時區
        /// 若查詢 2020-01-01 當日的報表，StartTimeUTC 會帶入 2020-01-01 00:00:00
        /// </summary>
        public DateTime StartTimeUTC { get; set; }


        /// <summary>
        /// 若查詢 2020-01-01 當日的報表，EndTimeUTC 會帶入 2020-01-02 00:00:00
        /// 請注意各遊戲商API接收的endTime格式要再自行做轉換 (e.g. 起始結束 要同一日)
        /// </summary>
        public DateTime EndTimeUTC { get; set; }



        /// <summary>
        /// 配合遊戲商時區的 起始時間
        /// </summary>
        public DateTime StartTimeTargetTimeZone { get; set; }


        /// <summary>
        /// 配合遊戲商時區的 結束時間
        /// </summary>
        public DateTime EndTimeTargetTimeZone { get; set; }



        public string DebugRequestUUID { get; set; }

    }



    public class ExternalVendorGetReportResultModel
    {

        public string ResponseText { get; set; }


        public BetRecordSummaryReportModel DbSummary { get; set; }

        public BetRecordSummaryReportModel ApiSummary { get; set; }

    }


    public class BetRecordSummaryReportModel
    {

        public int TotalBetCount { get; set; }
        public decimal TotalBetAmount { get; set; }
        public decimal TotalValidBetAmount { get; set; }
        public decimal TotalPayoutAmount { get; set; }
        public decimal TotalWinLoss { get; set; }


        /// <summary>
        /// 定義哪幾個欄位要做對帳檢查
        /// </summary>
        public List<GameVendorReportApiAvailableColumnsEnum> AvailableColumns { get; set; }

    }







}
