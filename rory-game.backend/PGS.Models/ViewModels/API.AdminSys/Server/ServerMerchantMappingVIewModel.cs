using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class ServerMerchantMappingViewModel
    {
        public List<ServerMerchantMappingSettingVIewModel> AvailableSettings { get; set; }

    }


    public class ReversServerMerchantRequestViewModel
    {
        public int MerchantID { get; set; }

        public int ServerID { get; set; }
    }


    public class ReversServerMerchantResultViewModel
    {
    
        public bool AfterValue { get; set; }

    }



    public class AllServerMerchantMappingViewModel
    {
        /// <summary>
        /// 有伺服器的幣別
        /// </summary>
        public List<string> AvailableCurrency { get; set; }

        public List<CurrencyCategoryItemModel> CurrencyGroups { get; set; }

    }


    public class CurrencyCategoryItemModel
    {

        public string Currency { get; set; }

        public List<BatchEditMappingMerchantItemViewModel> Merchants { get; set; }

        public List<BatchEditMappingGameItemViewModel> Games { get; set; }

        public List<BatchEditMappingServerItemViewModel> Servers { get; set; }
    }



    public class BatchEditMappingGameItemViewModel
    {
        public string GameName { get; set; }

        public int GameID { get; set; }

        public List<BatchEditMappingServerItemViewModel> Servers { get; set; }
    }


    public class BatchEditMappingServerItemViewModel
    {
        public string Currency { get; set; }

        public string GameName { get; set; }

        public int GameID { get; set; }


        public string ServerName { get; set; }

        public int ServerID { get; set; }

        public List<BatchEditMappingMerchantFlagItemViewModel> MerchantEnableFlags { get; set; }
    }



    public class BatchEditMappingMerchantItemViewModel
    {
        public string Currency { get; set; }

        public string MerchantName { get; set; }

        public int MerchantID { get; set; }
    }


    public class BatchEditMappingMerchantFlagItemViewModel
    {
        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public string MerchantCurrency { get; set; }

        public bool IsEnable { get; set; }
    }




    public class ServerMerchantMappingSettingVIewModel
    {

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public bool IsEnable { get; set; }
    }


}
