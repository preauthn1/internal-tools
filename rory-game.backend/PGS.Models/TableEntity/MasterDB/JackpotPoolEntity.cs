using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class JackpotPoolEntity
    {
        [Key]
        public int GameID { get; set; }
        [Key]
        public int FlagID { get; set; }
        [Key]
        public int LevelID { get; set; }
        [Key]
        public int TableNo { get; set; }
        public decimal Contribution { get; set; }
        public decimal TotalContribution { get; set; }
        public decimal TotalAward { get; set; }
        [Key]
        public bool IsLink { get; set; }

    }
}
