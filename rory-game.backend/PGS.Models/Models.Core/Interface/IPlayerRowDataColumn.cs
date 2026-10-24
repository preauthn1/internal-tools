using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models
{

    /// <summary>
    /// 基本上record型 資料表都有此欄位 (反正規化存法)
    /// </summary>
    public interface IPlayerRowDataColumn
    {
        int PlayerID { get; set; }

        int MerchantID { get; set; }

        string PlayerName { get; set; }
    }
}
