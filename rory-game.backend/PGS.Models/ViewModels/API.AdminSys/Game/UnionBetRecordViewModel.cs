using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class UnionBetRecordViewModel
    {

        /// <summary>
        /// 對應 BetRecordTypeEnum
        /// </summary>
        public int Type { get; set; }
        public string TypeName { get; set; }

        public string BetRecordNo { get; set; }
        public string BetRecordNoPGS { get; set; }
        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public int GameVendor { get; set; }
        public string GameVendorName { get; set; }

        public int GameID { get; set; }
        public string GameName { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public decimal BetAmount { get; set; }
        public decimal WinLossAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public int Status { get; set; }
        public string StatusText { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime? SettlementTime { get; set; }

        /// <summary>
        /// 重發retry紀錄
        /// </summary>
        public long MerchantNotifyQueueID { get; set; }
    }


}
