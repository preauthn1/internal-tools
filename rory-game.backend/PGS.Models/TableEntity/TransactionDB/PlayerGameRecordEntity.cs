using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class PlayerGameRecordEntity
    {
        [Key]
        public long ID { get; set; }
        public int PlayerID { get; set; }

        public string PlayerName { get; set; }
        public int MerchantID { get; set; }
        public int GameID { get; set; }
        public int ServerID { get; set; }
        public long TradeID { get; set; }
        public decimal LoginAmount { get; set; }
        public decimal UpdateAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal FlowAmount { get; set; }
        public int WinCount { get; set; }
        public int LoseCount { get; set; }
        public int DrawCount { get; set; }
        public string IPAddress { get; set; }
        public long GameTimeSecond { get; set; }
        public DateTime LoginTime { get; set; }
        public DateTime? LogoutTime { get; set; }


    }
}
