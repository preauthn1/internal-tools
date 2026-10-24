using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class VendorSessionTokenEntity
    {
        [Key]
        public long Id { get; set; }
        public string SessionToken { get; set; }
        public int VendorID { get; set; }

        public int GameID { get; set; }
        public int PlayerID { get; set; }
        public DateTime ExpiredTime { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime UpdateTime { get; set; }


    }
}
