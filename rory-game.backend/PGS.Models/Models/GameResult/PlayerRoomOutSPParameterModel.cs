using System;

namespace target.Models.ViewModels
{
    public class PlayerRoomOutSPParameterModel
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
        public long GameTimeSecond { get; set; }
        public DateTime LogoutTime { get; set; }
        public int LastWinLose { get; set; }       

        public decimal BetAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public decimal ValidAmount { get; set; }

    }


}
