using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class MerchantViewModel
    {
        public int MerchantID { get; set; }

        public int Type { get; set; }

        public string MerchantTypeName { get; set; }

        public string Name { get; set; }

        public string Memo { get; set; }

        public string Currency { get; set; }

        public string Code { get; set; }

        public string ChargeCurrency { get; set; }

        public string DefaultLanguage { get; set; }

        public decimal? DefaultBalance { get; set; }

        public int Status { get; set; }

        public string ApiSecretKey { get; set; }

        public string ApiValidatePlayerSessionUrl { get; set; }

        public string ApiSingleWalletGetAmountUrl { get; set; }


        public int WalletMode { get; set; }

        public string WalletModeText { get; set; }

        public string CreateBy { get; set; }

        public DateTime CreateTime { get; set; }

        public string UpdateBy { get; set; }

        public DateTime UpdateTime { get; set; }


        public int NumberOfGames { get; set; }

        public int NumberOfServers { get; set; }


    }


    public class MerchantCopyInfoPreviewViewModel
    {
        public MerchantViewModel Info { get; set; }

        public bool IsEnableGameEmpty { get; set; }

        public bool IsEnablePGSGameOnlyButServerEmpty { get; set; }

        public bool IsCodeNotSetup { get; set; }

        public int EnablePgsGames { get; set; }

        public int EnableThirdPartyGames { get; set; }

        public int EnablePgsGameServers { get; set; }




        public decimal PGSConvertRate { get; set; }

    }

}
