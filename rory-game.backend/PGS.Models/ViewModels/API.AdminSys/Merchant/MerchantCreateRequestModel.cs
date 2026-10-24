using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class MerchantCreateRequestModel
    {

        public string Name { get; set; }
        public string Memo { get; set; }
        public string Currency { get; set; }
        public string DefaultLanguage { get; set; }
        public int Status { get; set; }
        public string ApiSecretKey { get; set; }
        public string ApiValidatePlayerSessionUrl { get; set; }
        public string ApiSingleWalletGetAmountUrl { get; set; }
        public int WalletMode { get; set; }
        public string Code { get; set; }
        public int Type { get; set; }

    }


    public class MerchantUpdateRequestModel
    {
        public int MerchantID { get; set; }
        public string Name { get; set; }
        public string Memo { get; set; }
        public string Currency { get; set; }
        public string DefaultLanguage { get; set; }
        public int Status { get; set; }
        public string ApiSecretKey { get; set; }
        public string ApiValidatePlayerSessionUrl { get; set; }
        public string ApiSingleWalletGetAmountUrl { get; set; }
        public int WalletMode { get; set; }
        public string Code { get; set; }
        public int Type { get; set; }

    }



}
