using target.Models.CultureText;
using target.Models.Enums;
using target.Utility.Helper;
using System;
using System.ComponentModel;

namespace target.Models.ViewModels.GameVendor
{


    public class ExternalVendorLaunchUrlRequestModel
    {

        public int MerchantID { get; set; }
               
        public string ThirdPartyGameCode { get; set; }

        public int PlayerID { get; set; }

        public string ThirdPartyLaunchSessionToken { get; set; }

        public string PlayerName { get; set; }

        public string Currency { get; set; }

        public LanguageEnum Language { get; set; }

        public string BackUrl { get; set; }

        /// <summary>
        /// 是否為試玩模式
        /// </summary>
        public bool IsTrialMode { get; set; }


        /// <summary>
        /// 建議把這組ID存進requestLog 以利debug (會帶入 PlayerLoginLog 的ID)
        /// </summary>
        public string DebugReferenceUUID { get; set; }


    }



    public class ExternalVendorLaunchUrlViewModel
    {
     
        public string Url { get; set; }


    }


}
