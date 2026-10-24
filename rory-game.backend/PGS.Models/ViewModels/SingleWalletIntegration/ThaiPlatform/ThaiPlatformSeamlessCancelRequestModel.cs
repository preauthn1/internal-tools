using Newtonsoft.Json;

namespace game.Entity.Game
{
   

    public class ThaiPlatformSeamlessCancelRequestModel : BaseThaiPlatformSeamlessRequestModel
    {
        public string id { get; set; }
        public string username { get; set; }
        public string gamecode { get; set; }
        public string roundid { get; set; }
        public string betid { get; set; }
    }


    public class ThaiPlatformSeamlessCancelResponseModel : BaseThaiPlatformSeamlessApiResponseModel
    {

        public decimal Balance { get; set; }
    }



}
