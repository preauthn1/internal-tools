using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    /// <summary>
    /// 如果狀態只有 開啟/關閉/永久隱藏 這幾種
    /// </summary>
    public enum JobActionTypeEnum
    {
        [Description("為定義")]
        None = 0,

        [Description("寫遊戲每局紀錄")]
        WritePerGame = 0,
    }
}
