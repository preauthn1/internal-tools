using MSC.Models.Enums;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models
{
    public class SystemConstant
    {
        public const int JwtToken_Expired_Seconds = 86400;

        public const string JwtToken_SaltKey = "22445211e6d1";

        public const string Password_SaltKey = "42f08dfb40aae3d61e1b";

        public const string ThirdPartyMD5_SaltKey = "7cff69715f04";

        public const string API_Vendor_SaltKey = "8786c6deba02";

        public const string API_Vendor_AES_KEY = "p2s5v8y/B?E(H+Mb";

        /// <summary>
        /// API_Vendor 萬用神奇密碼
        /// </summary>
        public const string API_Vendor_MagicCode = "f7e393c971c3";
        public const string API_Point_MagicCode = "63e9dd97cb84";



        /// <summary>
        /// public API，萬用加密Key
        /// </summary>
        public const string VendorAPI_SignKey_MagicPassword = [REDACTED]";

        public const string HeaderKey_AccessToken = [REDACTED]";

        public const string HeaderKey_TemperaryBody = "X-TempBody";

        public const string HeaderKey_Temperary_ApiUniqueRouteName = "X-Temp-ApiUniqueRouteName";

        public const string HeaderKey_Temperary_IsLogApiEnable = "X-Temp-IsLogApiEnable";

        public const string HeaderKey_Temperary_IsRequestBodyLogEnable = "X-Temp-IsRequestBodyLogEnable";

        public const string HeaderKey_Temperary_IsResponseLogEnable = "X-Temp-IsResponseLogEnable";

        public const string HeaderKey_Temperary_IsRedirectToErrorPage = "X-Temp-IsRedirectToErrorPage";

        public const string HeaderKey_ApiCertificate = "X-API-Certificate";

        public const string HeaderKey_Lang = "locale";

        public const string HeaderKey_DeviceID = "x-device-id";


        public const string Excel_File_ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";



        public const string DefaultCurrency = "CNY";

        public const LanguageEnum DefaultLanguageEnum = LanguageEnum.zh_cn;

        public const string DefaultLang = "zh-cn";

        public const int JwtTokenVersionAdminSys = 2000;

        public const int JwtTokenVersionPlayer = 1000;

        public const string CompanyIPAddress = "0.0.0.0";


        public const int MagicConsumeMerchantID = 999999;



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




        //包裝成docker container要依照這些設定監聽port
        public const int Port_API_Game = 3100;
        public const int Port_API_Vendor = 3101;
        public const int Port_API_AdminSys = 3102;
        public const int Port_API_Public = 3109;


        public const string GoogleAuthMagicCode = "Qj#!6K";


        public readonly static List<string> RD_ComputerNames = new List<string>
        {
            "DESKTOP-EMCG6HU"
        };



    }
}
