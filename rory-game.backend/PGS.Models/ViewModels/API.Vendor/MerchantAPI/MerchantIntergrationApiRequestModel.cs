using Newtonsoft.Json;

namespace target.Models.ViewModels
{
    public class MerchantIntergrationApiParameterModel : BaseVendorRequestParameter
    {

        /// <summary>
        /// 後半段url e.g. /target/VadlidateSession
        /// </summary>
        public string ApiValidatePlayerSessionUrl { get; set; }


        public string PlayerName { get; set; }

        public string PlayerValidateSession { get; set; }

    }


    public class MerchantIntergrationApiRequestModel
    {
        public string PlayerName { get; set; }

        public string PlayerValidateSession { get; set; }

        public int MerchantID { get; set; }

        public string MerchantSecretKey { get; set; }

    }



    public class MerchantIntergrationApiResultModel 
    {     

        public bool IsValidOK { get; set; }

        public string PlayerName { get; set; }

    }


}
