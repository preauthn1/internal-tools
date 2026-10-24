using target.Models.Enums;
using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models
{
    public class SystemConstant
    {
        public const int JwtToken_Expired_Seconds = 86400;

        public const string JwtToken_SaltKey = "y7dw2t6tvc2ujel618";

        public const string Password_SaltKey = "8f3c0a484f9fc1";

        public const string ThirdPartyMD5_SaltKey = "4942335a69c2";

        public const string PublicAPI_Singkey_SaltKey = "2b3c7342aadf";

        public const string PublicAPI_Lobby_Singkey_SaltKey = "z85pc7w42a6dc";

        /// <summary>
        /// public API，萬用加密Key
        /// </summary>
        public const string PublicAPI_SignKey_MagicPassword = [REDACTED]";

        public const string HeaderKey_AccessToken = [REDACTED]";

        public const string HeaderKey_IdentityMemo = "identity_memo";

        public const string HeaderKey_TemperaryBody = "X-TempBody";

        public const string HeaderKey_Temperary_ApiUniqueRouteName = "X-Temp-ApiUniqueRouteName";

        public const string HeaderKey_Temperary_IsLogApiEnable = "X-Temp-IsLogApiEnable";

        public const string HeaderKey_Temperary_IsLoggingPlayerID = "X-Temp-LoggingPlayerID";

        public const string HeaderKey_Temperary_IsRequestBodyLogEnable = "X-Temp-IsRequestBodyLogEnable";

        public const string HeaderKey_Temperary_IsResponseLogEnable = "X-Temp-IsResponseLogEnable";

        public const string HeaderKey_Temperary_IsRedirectToErrorPage = "X-Temp-IsRedirectToErrorPage";

        public const string HeaderKey_Temperary_ApiCategory = "X-Temp-ApiCategory";

        public const string HeaderKey_ApiCertificate = "X-API-Certificate";

        public const string HeaderKey_SignKey = "X-SignKey";

        public const string HeaderKey_Lang = "locale";


        public const string Excel_File_ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";

        public const string HeaderKey_HDCode = "HDCode";



        /// <summary>
        /// 站台關閉的delay時間
        /// </summary>
        public const int ApplicationShutdownDelaySeconds = 30;

        public const string ApiUrlRemoveAllCache = "internal/cache/removeAll";

        public const string DefaultCurrency = "CNY";

        public const LanguageEnum DefaultLanguageEnum = LanguageEnum.zh_cn;

        public const string DefaultLang = "zh-cn";

        public const int JwtTokenVersionAdminSys = 2001;

        public const int JwtTokenVersionPlayer = 1001;

        public const string CompanyIPAddress = "0.0.0.0";

        public const int DefaultMerchantID = 175001;

        public const int SkipValidUserID = -5566;


        /// <summary>
        /// 玩家帳號字串長度
        /// </summary>
        public const int PlayerAccountMinLength = 4;

        public const int PlayerAccountMaxLength = 50;


        /// <summary>
        /// 玩家密碼字串長度
        /// </summary>
        public const int PlayerPasswordMinLength = 6;

        public const int PlayerPasswordMaxLength = 30;

        

        //各環境名稱定義 (判斷的時候請轉小寫)
        public const string Env_Name_Dev = "dev";
        public const string Env_Name_Development = "development";
        public const string Env_Name_QA = "qa";
        public const string Env_Name_Production = "production";


        public const int CrazyFishingDefaultGameID = 310001;

        public const string CrazyFishingDefaultGameCode = "happyfishing";


        //包裝成docker container要依照這些設定監聽port
        public const int Port_API_Game = 3100;
        public const int Port_API_Vendor = 3101;
        public const int Port_API_AdminSys = 3102;
        public const int Port_API_Public = 3109;


        public const string Public_API_Always_Session_OK = "Always_Session_OK";


        public const string AGGGameVendorDefaultSuccessResponse = "RS_OK";


        public readonly static List<string> RD_ComputerNames = new List<string>
        {
            "DESKTOP-EMCG6HU"
        };


    }
}
