using Newtonsoft.Json;

namespace game.Entity.Game
{


    public class ThaiPlatformSeamlessGetBalanaceRequestModel : BaseThaiPlatformSeamlessRequestModel
    {

        public string username { get; set; }
    }

    public class ThaiPlatformSeamlessGetBalanaceResponseModel : BaseThaiPlatformSeamlessApiResponseModel
    {

        public decimal Balance { get; set; }
    }





}
