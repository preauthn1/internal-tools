using Newtonsoft.Json;

namespace target.Models.ViewModels
{
    public class AGGLiveErrorResponseModel
    {
        public AGGLiveErrorResponseModel()
        {

        }

        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public AGGLiveStatusOnlyResponseModel playerdetailresponse { get; set; }

        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public AGGLiveStatusOnlyResponseModel fundtransferresponse { get; set; }

    }

}
