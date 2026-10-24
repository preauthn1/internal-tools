using target.Models.Enums;
using System;
using System.Collections.Generic;

namespace target.Models.ViewModels.GameVendor
{
    public class ExternalVendorGetSingleBetRecordRequestParameter
    {       
        public int MerchantID { get; set; }

        public string GameCode { get; set; }

        /// <summary>
        /// 要帶三方的單號
        /// </summary>
        public string BetRecordNo { get; set; }

        public string DebugReferenceUUID { get; set; }

    }


}
