using MSC.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{
    [IsNotAutoKeyTable]
    public class MerchantWalletEntity
    {
        public int MerchantID { get; set; }

        public int WalletType { get; set; }

        public decimal Amount { get; set; }

    }
}
