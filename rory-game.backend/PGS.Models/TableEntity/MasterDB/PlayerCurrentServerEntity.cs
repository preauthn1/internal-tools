using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    [IsNotAutoKeyTable]
    public class PlayerCurrentServerEntity
    {
        [Key]
        public int PlayerID { get; set; }
        public int ServerID { get; set; }
        public int HDType { get; set; }
        public string Token { get; set; }
        public long LastTradeID { get; set; }
        public DateTime CreateTime { get; set; }

    }
}
