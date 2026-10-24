using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class GameWarningLogEntity
    {
        [Key]
        public long ID { get; set; }
        public int WarningType { get; set; }
        public int ServerID { get; set; }
        public int PlayerID { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }

    }
}
