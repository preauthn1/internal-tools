using Newtonsoft.Json;

namespace game.Entity.Game
{


    public class ThaiPlatformSeamlessBonusRequestModel : BaseThaiPlatformSeamlessRequestModel
    {
        public string id { get; set; }
        public decimal amount { get; set; }
        public string username { get; set; }
        public string gamecode { get; set; }
        public string roundid { get; set; }
        public string description { get; set; }
        public string type { get; set; }
    }

    public class ThaiPlatformSeamlessBonusResponseModel : BaseThaiPlatformSeamlessApiResponseModel
    {

        public decimal Balance { get; set; }
    }





}
