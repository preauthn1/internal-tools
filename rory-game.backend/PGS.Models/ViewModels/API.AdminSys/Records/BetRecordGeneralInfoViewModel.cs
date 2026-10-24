using target.Models.Table;
using target.Models.ViewModels.GameVendor;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class BetRecordQueryRequestModel
    {
        public string TransactionNo { get; set; }
    }

    public class BetRecordManualRefundRequestModel
    {
        public string TransactionNo { get; set; }

        public string Memo { get; set; }
    }



    public class BetRecordGeneralInfoViewModel
    {

        public int GameVendorID { get; set; }


        public string GameVendorName { get; set; }



        public bool IsSupportExternalLink { get; set; }



        /// <summary>
        /// 所有注單共通型欄位
        /// </summary>
        public BetRecordGeneralColumnViewModel GeneralInfo { get; set; }


        public PerGameDetailInfo PerGameLogInfo { get; set; }
    }



    public class BetRecordGeneralColumnViewModel : BetRecordGeneralColumnModel
    {
        public string MerchantName { get; set; }

        public string ManualSettlementUserName { get; set; }

        public bool IsManualSettlement { get; set; }
    }



    public class PerGameDetailInfo
    {
        public long PerGameDetailLogID { get; set; }

        public long PerGameLogID { get; set; }

        public int ServerID { get; set; }

    }





}
