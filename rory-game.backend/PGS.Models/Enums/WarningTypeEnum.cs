using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
  
    public enum WarningTypeEnum
    {
        [Description("無鎖寫分")]
        WriteGameResultWithoutGameRoom = 1000,

        [Description("錢包相關")]
        Wallet_Warning = 2000,
    }
}
