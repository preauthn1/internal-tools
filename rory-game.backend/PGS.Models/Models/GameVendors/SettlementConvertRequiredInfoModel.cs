using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{


    public class SettlementConvertRequiredInfoModel
    {
        public string BetRecordNo { get; set; }

        /// <summary>
        /// 更新後金額 (法幣)
        /// </summary>
        public decimal AfterAmount { get; set; }

        /// <summary>
        /// 更新後金額 (遊戲幣)
        /// </summary>
        public decimal AfterGamePoiint { get; set; }

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public string Currency { get; set; }

    }


}
