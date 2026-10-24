using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.Table
{

    public class PlayerEntity
    {
        [Key]
        public int PlayerID { get; set; }
        public int MerchantID { get; set; }
        public string PlayerName { get; set; }
        public string DisplayName { get; set; }
        public int PlayerType { get; set; }
        public string Currency { get; set; }

        public int Status { get; set; }


        public string LastIPAddress { get; set; }

        public string CreateIPAddress { get; set; }

        public DateTime CreateTime { get; set; }

        public DateTime UpdateTime { get; set; }

        public string LastLoginToken { get; set; }

        public DateTime? LastLoginTime { get; set; }

        public int FaceID { get; set; }

        public int Sex { get; set; }

    }


    public class PlayerBriefModel
    {
        [Key]
        public int PlayerID { get; set; }
        public int MerchantID { get; set; }
        public string PlayerName { get; set; }
        public string Currency { get; set; }
        public int Status { get; set; }

    }

}
