using target.Models.CultureText;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{

    public enum WalletTransferTypeEnum
    {
        [Description(i18nKeys.Enums_WalletTransferTypeEnum_Transfer_In)]
        Transfer_In = 1,

        [Description(i18nKeys.Enums_WalletTransferTypeEnum_Transfer_Out)]
        Transfer_Out = 2,

    }
}
