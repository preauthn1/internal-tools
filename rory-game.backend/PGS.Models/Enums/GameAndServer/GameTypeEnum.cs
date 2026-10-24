using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
   
    public enum GameTypeEnum
    {
        [Description(i18nKeys.Enums_GameTypeEnum_None)]
        None = 0,

        Slot = 1000,
        
        BoardGame = 2000,

        Arcade = 3000,

        Arena = 4000,

        PachinSlotGame = 6000,

        Fishing = 7000,
    }
}
