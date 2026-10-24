using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class PerGameLogEntity
    {
        [Key]
        public long PerGameLogID { get; set; }
        public int GameID { get; set; }
        public int ServerID { get; set; }
        public string GameData { get; set; }
        public decimal TaxAmount { get; set; }
        public DateTime BeginTime { get; set; }
        public DateTime EndTime { get; set; }
        public DateTime CreateTime { get; set; }

    }
}
