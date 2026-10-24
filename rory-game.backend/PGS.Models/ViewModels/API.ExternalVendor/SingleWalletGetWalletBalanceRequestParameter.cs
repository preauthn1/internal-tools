using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{



    public class SingleWalletGetWalletBalanceRequestParameter : SingleWalletBasePlayerInfoModel
    {
    
    }



    public class SingleWalletGetWalletBalanceResultViewModel
    {

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        /// <summary>
        /// 錢包金額 (法幣)
        /// </summary>
        public decimal BalanceAmount { get; set; }

        /// <summary>
        /// 錢包金額 (遊戲幣)
        /// </summary>
        public decimal BalanceGamePoint { get; set; }

        public string Currency { get; set; }
    }



}
