using System;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    /// <summary>
    /// 每局遊戲紀錄
    /// </summary>
    public class PerGameLogRequestModel
    {
        public PerGameLogItemRequestModel PerGameLog { get; set; }


        /// <summary>
        /// 有v1/v2版本的API 流程不太一樣
        /// </summary>
        [SwaggerIgnore]
        public bool IsPlaceBetProcess { get; set; }



        /// <summary>
        /// 如果未來能做到記錄每筆 投注的話
        /// </summary>
        public List<EachBetRecordViewModel> EachBetRecords { get; set; }

    }


    /// <summary>
    /// 每局遊戲紀錄
    /// </summary>
    public class PerGameLogItemRequestModel
    {
        public int ServerID { get; set; }



        /// <summary>
        /// 填寫純文字的遊戲歷程記錄
        /// </summary>
        public string GameData { get; set; }


        /// <summary>
        /// 這樣的格式 10001:2000:-10:20180401,10002:1000:-10:20180402,10003:3000:40:20180405,10003:3000:-20:20180409
        /// 分別是 {playerID} : {currentAmount} : {updateAmount} : {tradeID}
        /// </summary>
        public string UserGameDetail { get; set; }

        /// <summary>
        /// 2020/7/8 調整成物件型結構
        /// </summary>
        public List<PerGameDetailLogInfoViewModel> UserGameDetailLogs { get; set; }


        public DateTime GameBeginTime { get; set; }
        public DateTime GameEndTime { get; set; }
        public decimal TaxAmount { get; set; }
        public int GameID { get; set; }
        public int? WinPrizeID { get; set; }

    }


    public class PerGameDetailLogInfoViewModel
    {

        public int PlayerID { get; set; }

        public long TradeID { get; set; }

        public decimal CurrentAmount { get; set; }

        public decimal UpdateAmount { get; set; }

        public decimal BetAmount { get; set; }

        public decimal ValidAmount { get; set; }

        public decimal PayoutAmount { get; set; }

        public decimal TaxAmount { get; set; }

        /// <summary>
        /// 投注單號
        /// </summary>
        public string TransactionNo { get; set; }

        /// <summary>
        /// 投注單號 (多筆)
        /// </summary>
        public List<string> TransactionNoList { get; set; }
    }




    /// <summary>
    /// 每局遊戲紀錄
    /// </summary>
    public class EachBetRecordViewModel
    {
        public int PlayerID { get; set; }

        public decimal Amount { get; set; }

    }





    /// <summary>
    /// 每局遊戲紀錄
    /// </summary>
    public class PerGameLogViewModel
    {

        public long PerGameLogID { get; set; }

    }


}
