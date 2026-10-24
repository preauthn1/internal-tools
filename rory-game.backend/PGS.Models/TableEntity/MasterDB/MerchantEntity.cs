using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class MerchantEntity
    {
        [Key]
        public int MerchantID { get; set; }
        public int Type { get; set; }
        public string Name { get; set; }
        public string Code { get; set; }
        public string Memo { get; set; }
        public string Currency { get; set; }
        public string ChargeCurrency { get; set; }
        public string DefaultLanguage { get; set; }
        public int Status { get; set; }
        public string ApiSecretKey { get; set; }
        public string CreateBy { get; set; }
        public DateTime CreateTime { get; set; }
        public string UpdateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public int WalletMode { get; set; }
        public double TimeZone { get; set; }

    }

    public class MerchantBriefEntity
    {
        [Key]
        public int Id { get; set; }        

        public int Status { get; set; }

        public string Currency { get; set; }

        public string DefaultLanguage { get; set; }

        public string ApiSecretKey { get; set; }


    }
}
