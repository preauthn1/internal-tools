using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class MerchantEntity
    {
        [Key]
        public int MerchantID { get; set; }
        public string DisplayName { get; set; }
        public int Status { get; set; }
        public string Memo { get; set; }
        public string ApiSecretKey { get; set; }
        public string Currency { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }


    }
}
