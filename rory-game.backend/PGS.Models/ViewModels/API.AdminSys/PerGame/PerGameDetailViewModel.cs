using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 回應參數
    /// </summary>
    public class PerGameDetailViewModel
    {
        public long PerGameDetailLogID { get; set; }

        public long PerGameLogID { get; set; }

        public int ServerID { get; set; }

        public string ServerName { get; set; }

        public int GameID { get; set; }

        public string GameName { get; set; }

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }


        public decimal CurrentAmount { get; set; }

        public decimal AfterAmount { get; set; }

        public decimal UpdateAmount { get; set; }

        public decimal BetAmount { get; set; }

        public decimal PayoutAmount { get; set; }

        public decimal ValidAmount { get; set; }

        public decimal JackpotAmount { get; set; }

        public decimal TaxAmount { get; set; }


        public DateTime GameTime { get; set; }

        public bool IsExternalVendorGame { get; set; }

        public string ThirdPartyBetRecordNo { get; set; }

        public string ThirdPartyGameCode { get; set; }

        public bool IsSlotGame { get; set; }

        public bool IsHideAfterAmountInfo { get; set; }

        public bool IsSupportExternalLink { get; set; }

        public string VendorName { get; set; }

        public string TransactionNo { get; set; }

        public decimal ConvertRate { get; set; }
    }



}
