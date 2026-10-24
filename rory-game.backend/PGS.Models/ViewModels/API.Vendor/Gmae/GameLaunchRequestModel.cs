using Newtonsoft.Json;
using target.Models.Attributes;
using target.Models.Enums;
using target.Models.Table;
using target.Models.ViewModels.Base;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{


    /// <summary>
    /// 這個class有被兩個class繼承
    /// </summary>
    public class BaseGameLaunchRequestModel
    {
        public int MerchantID { get; set; }


        public string MerchantSecretKey { get; set; }



        public string Lang { get; set; }


        public int GameID { get; set; }

             

        /// <summary>
        /// 選填參數，可讓商戶選擇要進哪種等級的廳
        /// </summary>
        [SwaggerIgnore]
        public int? RoomLevelFlagID { get; set; }

    }


    /// <summary>
    /// 請求參數
    /// </summary>
    public class GameLaunchRequestModel : BaseGameLaunchRequestModel
    {
    
        /// <summary>
        /// 非正式使用(內部測試特殊場景)
        /// </summary>
        [SwaggerIgnore]
        public int? SecretParameter_ServerID { get; set; }



        public string PlayerName { get; set; }


        /// <summary>
        /// [單錢包模式] 要帶入此參數，跟商戶的API溝通需要用到此Token
        /// </summary>
        public string SessionToken { get; set; }


        /// <summary>
        /// 如果是Lobby模式，singleWalletToken 會允許使用同一組sessionToken
        /// </summary>
        public bool IsAllowDuplicateSessionToken { get; set; }

    }



    /// <summary>
    /// 启动游戏請求參數 (VendorAPI_转账模式)
    /// </summary>
    public class GameLaunchVendorApiTransferModeRequestModel : BaseGameLaunchRequestModel
    {


        public string PlayerName { get; set; }


        /// <summary>
        /// 非正式使用(內部測試特殊場景)
        /// </summary>
        [SwaggerIgnore]
        public int? SecretParameter_ServerID { get; set; }

    }




    /// <summary>
    /// 启动游戏請求參數 (VendorAPI_单钱包模式)
    /// </summary>
    public class GameLaunchVendorApiSingleWalletModeRequestModel : BaseGameLaunchRequestModel
    {

        /// <summary>
        /// [單錢包模式] 要帶入此參數 (用戶唯一識別碼)
        /// </summary>
        public string PlayerUniqueID { get; set; }



        /// <summary>
        /// [單錢包模式] 要帶入此參數，跟商戶的API溝通需要用到此Token
        /// </summary>
        public string SessionToken { get; set; }



    }





    /// <summary>
    /// 回應參數
    /// </summary>
    public class GameLaunchViewModel
    {
        public bool IsMultipleServers { get; set; }

        public bool IsLobbyLink { get; set; }

        public PlayerTokenJwtModel PlayerJwtModel { get; set; }

        public GameLaunchRequiredInfoViewModel Info { get; set; }

        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public ExternalVendorLaunchInfoViewModel ExternalVendorInfo { get; set; }

    }

    public class ExternalVendorLaunchInfoViewModel
    {
        public string ExternalVendorUrl { get; set; }
    }


    /// <summary>
    /// 這個model包的東西  client & GameServer 會用到
    /// </summary>
    public class GameLaunchRequiredInfoViewModel
    {
        public GameLaunchRequiredInfoViewModel()
        {
            this.ThirdPartyInfo = new ThirdPartyExtraInfoModel();
        }


        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public string Lang { get; set; }

        public int ServerID { get; set; }

        public int GameID { get; set; }

        public string GameCode { get; set; }

        [ReferenceType(typeof(MerchantLoadingStyleModeEnum))]
        public int StyleMode { get; set; }

        public string ServerConnection { get; set; }

        public string Access_Token { get; set; }


        /// <summary>
        /// Client端的log開關  0: 關閉，  1: 啟用， 未來也許有第三種? 到時候在定義enum
        /// </summary>
        public int Log { get; set; }

        public ThirdPartyExtraInfoModel ThirdPartyInfo { get; set; }
    }



    public class GameLinkGenerateResultViewModel
    {

        public string Link { get; set; }
    }




    public class ThirdPartyExtraInfoModel
    {
        public int GameVendorID { get; set; }


        public string MD5Text { get; set; }
    }



}
