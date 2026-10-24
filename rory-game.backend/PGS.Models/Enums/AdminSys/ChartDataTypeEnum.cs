using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{

    /// <summary>
    /// Chart Data 用的
    /// </summary>
    public enum ChartDataTypeEnum
    {
       
        [Description(i18nKeys.Enums_ChartDataTypeEnum_Game_TotalGame)]
        Summary_TotalLogin = 1,
     
        [Description(i18nKeys.Enums_ChartDataTypeEnum_Summary_AverageOnline)]
        Summary_AverageOnline = 2,

        [Description(i18nKeys.Enums_ChartDataTypeEnum_Summary_TotalGame)]
        Summary_TotalGame = 3,
     
        [Description(i18nKeys.Enums_ChartDataTypeEnum_Game_TotalGame)]
        Game_TotalGame = 101,
      
        [Description(i18nKeys.Enums_ChartDataTypeEnum_Game_OnlineUser)]
        Game_OnlineUser = 102,
    }
}
