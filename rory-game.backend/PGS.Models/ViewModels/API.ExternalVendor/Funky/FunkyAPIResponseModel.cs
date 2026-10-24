using Newtonsoft.Json;

namespace target.Models.ViewModels
{
    public class FunkyAPIResponseModel<T>
    {

        public FunkyAPIResponseModel()
        {
            this.ErrorMessage = "No Error";
        }

        [JsonProperty("errorCode")]
        public int ErrorCode { get; set; }

        [JsonProperty("errorMessage")]
        public string ErrorMessage { get; set; }

        [JsonProperty("data")]
        public T Data { get; set; }
    }



}
