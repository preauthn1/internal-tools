using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    [IsNotAutoKeyTable]
    public class PlayerTotalSummaryEntity
    {
        [Key]
        public int PlayerID { get; set; }
        public decimal SumTax { get; set; }
        public int SumLoginTimes { get; set; }
        public decimal SumFlowAmount { get; set; }

    }
}
