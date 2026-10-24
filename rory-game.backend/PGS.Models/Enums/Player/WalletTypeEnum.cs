using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    
    public enum WalletTypeEnum
    {
        [Description(i18nKeys.Enums_WalletTypeEnum_GameWallet)]
        GameWallet = 1,

        [Description(i18nKeys.Enums_WalletTypeEnum_PlayFoFirstTry)]
        PlayFoFirstTry = 99

    }
}
