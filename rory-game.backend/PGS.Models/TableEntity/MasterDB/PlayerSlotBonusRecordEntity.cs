using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class PlayerSlotBonusRecordEntity
    {
        [Key]
        public int ID { get; set; }
        public int PlayerID { get; set; }
        public int GameID { get; set; }
        public int FlagID { get; set; }
        public int OriginID { get; set; }
        public int TurnID { get; set; }
        public int Location { get; set; }
        public decimal BetAmount { get; set; }
        public int BetRate { get; set; }
        public string Data { get; set; }

        /// <summary>
        /// 此遊戲獲得的來源 Pattern
        /// </summary>
        public int PatternID { get; set; }
    }
}
