using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class ThirdPartySlot88MD5RecordEntity
    {
        [Key]
        public int ID { get; set; }
        public int PlayerID { get; set; }
        public int GameID { get; set; }
        public string MD5 { get; set; }
        public string OriginalJson { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
