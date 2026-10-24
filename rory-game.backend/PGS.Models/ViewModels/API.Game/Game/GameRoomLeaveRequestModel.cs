using System;

namespace target.Models.ViewModels
{
    /// <summary>
    /// 遊戲退房
    /// </summary>
    public class GameRoomLeaveRequestModel
    {

        public int PlayerID { get; set; }
        public int ServerID { get; set; }
        public int GameID { get; set; }
        public long TradeID { get; set; }
        public decimal UpdateAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal FlowAmount { get; set; }
        public int WinCount { get; set; }
        public int LoseCount { get; set; }
        public int DrawCount { get; set; }
        public string IPAddress { get; set; }
        public DateTime LoginTime { get; set; }
        public DateTime LogoutTime { get; set; }
        public int LastWinLose { get; set; }
        public decimal BetAmount { get; set; }
        public decimal PayoutAmount { get; set; }

        /// <summary>
        /// 不確定跟 FlowAmount 是不是一樣的東西
        /// </summary>
        public decimal ValidAmount { get; set; }

    }


    /// <summary>
    /// 遊戲退房
    /// </summary>
    public class GameRoomLeaveViewModel
    {


    }




}
