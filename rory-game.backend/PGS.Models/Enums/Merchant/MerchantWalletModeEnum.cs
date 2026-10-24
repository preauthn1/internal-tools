using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    
    public enum MerchantWalletModeEnum
    {
        [Description(i18nKeys.Enums_MerchantWalletModeEnum_TransferMode)]
        TransferMode = 1,

        [Description(i18nKeys.Enums_MerchantWalletModeEnum_SingleWalletMode)]
        SingleWalletMode = 2,       

    }
}
