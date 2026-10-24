using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class PerGameDetailLogZeroEntity : IPlayerRowDataColumn
    {
        [Key]
        public long Id { get; set; }
        public long PerGameLogID { get; set; }
        public int ServerID { get; set; }
        public int GameID { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public long TradeID { get; set; }
        public int MerchantID { get; set; }
        public decimal CurrentAmount { get; set; }
        public decimal UpdateAmount { get; set; }
        public DateTime GameTime { get; set; }
        public DateTime CreateTime { get; set; }
        public decimal BetAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public decimal ValidAmount { get; set; }


    }
}
