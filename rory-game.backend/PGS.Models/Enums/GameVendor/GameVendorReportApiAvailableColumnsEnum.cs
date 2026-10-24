using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{

    /// <summary>
    /// 定義遊戲商 報表API，提供有提供哪些欄位、可作為對帳檢查判斷用
    /// </summary>
    public enum GameVendorReportApiAvailableColumnsEnum
    {
        TotalBetCount,
        TotalBetAmount,
        TotalValidBetAmount,
        TotalPayoutAmount,
        TotalWinLoss,
    }
}
