using MSC.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{


    [IsNotAutoKeyTable]
    public class CardInfoEntity
    {
       
        public string CardNo { get; set; }
        public string CardPassword { get; set; }
        public decimal Point { get; set; }
        public string Memo { get; set; }
        public string Currency { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }
        public int? ConsumeMerchantID { get; set; }

    }
}
