using Newtonsoft.Json;

namespace game.Entity.Game
{
   

    public class ThaiPlatformSeamlessPlaceBetRequestModel : BaseThaiPlatformSeamlessRequestModel
    {

        public string id { get; set; }
        public decimal amount { get; set; }
        public string username { get; set; }
        public string gamecode { get; set; }
        public string roundid { get; set; }

    }


    public class ThaiPlatformSeamlessPlaceBetResponseModel : BaseThaiPlatformSeamlessApiResponseModel
    {

        public decimal Balance { get; set; }
    }



}
