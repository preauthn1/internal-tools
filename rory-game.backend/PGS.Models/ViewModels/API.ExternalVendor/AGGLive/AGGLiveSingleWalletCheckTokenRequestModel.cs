using Newtonsoft.Json;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class AGGLiveSingleWalletCheckTokenRequestModel
    {
        public string type { get; set; }
        public DateTime dtsent { get; set; }
        public AGGLiveBaseGame basegame { get; set; }
        public AGGLiveAuth auth { get; set; }
        public AGGLivePlayerDetailRequest Playerdetailrequest { get; set; }
    }


    public class AGGLiveSingleWalletAPICheckTokenResponseDataModel
    {
        [JsonProperty(PropertyName = "playerdetailresponse")]
        public AGGLiveStatusOnlyResponseModel Playerdetailresponse { get; set; }
        public string accountid { get; set; }
        public string accountname { get; set; }
        public decimal balance { get; set; }
        public string currencycode { get; set; }
    }



    public class AGGLivePlayerDetailRequest
    {
        public string token { get; set; }
        public string gamelaunch { get; set; }
    }



    public class AGGLiveStatusResponseModel
    {
        public bool success { get; set; }
        public bool autherror { get; set; }
        public bool nofunds { get; set; }
        public bool successdebit { get; set; }
        public bool successcredit { get; set; }
        public bool refundstatus { get; set; }
        public string message { get; set; }
    }




    public class AGGLiveStatusOnlyResponseModel
    {
        [JsonProperty(PropertyName = "status")]
        public AGGLiveStatusResponseModel status { get; set; }

        public decimal balance { get; set; }

        public string currencycode { get; set; }

    }



}
