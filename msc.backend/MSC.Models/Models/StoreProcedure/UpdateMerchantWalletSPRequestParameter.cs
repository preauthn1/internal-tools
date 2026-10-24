using MSC.Models.Enums;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.Models
{

    public class UpdateMerchantWalletSPRequestParameter
    {
        public int MerchantID { get; set; }

        public decimal SourceAmount { get; set; }

        public decimal FeeRate { get; set; }

        public TransactionTypeEnum TransactionType { get; set; }

        public string RefNo { get; set; }


    }
}
