using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PlayerBriefInfoViewModel
    {

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public string DisplayName { get; set; }

        public int PlayerType { get; set; }

        public string Currency { get; set; }

        public int FaceID { get; set; }

        public int MerchantID { get; set; }

        /// <summary>
        /// Table其實沒這欄位 先寫死?
        /// </summary>
        public int Sex { get; set; }

    }



    public class PlayerWalletViewModel
    {

        /// <summary>
        /// 玩家錢包金額 (遊戲幣)
        /// </summary>
        public decimal Amount { get; set; }

        public int MerchantWalletMode { get; set; }

    }


}
