using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class SlotGameHistoryLogEntity
    {
        [Key]
        public long ID { get; set; }
        public int GameID { get; set; }
        public int ServerID { get; set; }
        public int PlayerID { get; set; }
        public decimal CurrentAmount { get; set; }
        public decimal BetAmount { get; set; }
        public decimal WinAmount { get; set; }
        public decimal FinalAmount { get; set; }
        public int GameType { get; set; }
        public int LogType { get; set; }
        public string GameData { get; set; }
        public DateTime GameTime { get; set; }
        public DateTime CreateTime { get; set; }

    }
}
