using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class MerchantGameSettingEntity
    {
        [Key]
        public int Id { get; set; }
        public int MerchantID { get; set; }
        public int GameID { get; set; }
        public int Status { get; set; }
        public int Sort { get; set; }


        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }



    }
}
