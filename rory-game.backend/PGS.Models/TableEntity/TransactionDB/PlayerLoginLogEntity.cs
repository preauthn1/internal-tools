using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class PlayerLoginLogEntity
    {
        [Key]
        public long Id { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public int MerchantID { get; set; }
        public int GameID { get; set; }
        public int HDType { get; set; }
        public string HDCode { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateIPAddress { get; set; }
        public string Memo { get; set; }
        public string MerchantSession { get; set; }
        public string UserAgent { get; set; }
        public int Status { get; set; }

    }
}
