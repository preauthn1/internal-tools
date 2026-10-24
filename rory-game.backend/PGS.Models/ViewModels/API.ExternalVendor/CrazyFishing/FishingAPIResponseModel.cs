using Newtonsoft.Json;

namespace target.Models.ViewModels
{
    public class FishingAPIResponseModel
    {
        public FishingAPIResponseModel()
        {
            
        }

        [JsonProperty("code")]
        public int Code { get; set; }


        [JsonProperty("msg")]
        public string Msg { get; set; }        

    }



}
