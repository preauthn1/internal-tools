using Newtonsoft.Json;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 不同遊戲商要response的格式可能不同
    /// </summary>
    public class SingleWalletExceptionHandleModel
    {        

        /// <summary>
        /// http status code (有的遊戲商可能要求回 200、400....)
        /// </summary>
        public int StatusCode { get; set; }

        /// <summary>
        /// 每家遊戲商各自的客製回覆
        /// </summary>
        public object GameVendorCustomizedResponse { get; set; }


        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public object PGSDebug { get; set; }

    }





}
