using Newtonsoft.Json;

namespace target.Models.ViewModels
{
    public class AGGAPIResponseModel
    {
        public AGGAPIResponseModel()
        {
            
        }

        [JsonProperty("statusCode")]
        public string StatusCode { get; set; }

    }



}
