using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{




    /// <summary>
    /// 請求參數
    /// </summary>
    public class GamePlayRequestModel
    {

        /// <summary>
        /// e.g. QZNN 搶莊牛牛
        /// </summary>
        public string GameCode { get; set; }

    }


    /// <summary>
    /// 回應參數
    /// </summary>
    public class GamePlayViewModel
    {
        public decimal WalletCashAmount { get; set; }

        public int GameServerID { get; set; }

        public int GameID { get; set; }

        public decimal MinBetAmount { get; set; }

        public decimal MaxBetAmount { get; set; }

    }







}
