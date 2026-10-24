using target.Models.Enums;
using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models
{
    public class PlayerMerchantInfoModel
    {
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public int MerchantID { get; set; }
        public string Currency { get; set; }
        public MerchantWalletModeEnum WalletMode { get; set; }

    }
}
