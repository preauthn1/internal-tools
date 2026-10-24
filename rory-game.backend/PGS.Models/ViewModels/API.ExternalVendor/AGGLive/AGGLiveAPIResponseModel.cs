using Newtonsoft.Json;

namespace target.Models.ViewModels
{
    public class AGGLiveAPIResponseModel
    {
        public AGGLiveAPIResponseModel()
        {
            
        }

        [JsonProperty("statusCode")]
        public string StatusCode { get; set; }

    }





}
