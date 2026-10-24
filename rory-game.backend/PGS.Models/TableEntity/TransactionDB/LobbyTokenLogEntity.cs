using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    public class LobbyTokenLogEntity
    {
        [Key]
        public long Id { get; set; }
        public string IdentityHash { get; set; }
        public Guid? IdentityGuid { get; set; }
        public string RawJwtToken { get; set; }
        public int MerchantID { get; set; }
        public int PlayerID { get; set; }
        public DateTime ExpiredTime { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
