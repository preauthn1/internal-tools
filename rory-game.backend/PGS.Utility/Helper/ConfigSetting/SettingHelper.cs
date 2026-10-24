using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using target.Utility.Extensions;
using example.invalid;
using example.invalid;
using Microsoft.AspNetCore.Http;
using System.Reflection;
using System.Runtime.Versioning;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.FileProviders;
using target.Models;

namespace target.Utility
{
    public static class SettingHelper
    {


        public static bool IsDev
        {
            get
            {
                return Env == SystemConstant.Env_Name_Development || Env == SystemConstant.Env_Name_Dev;
            }
        }
        

        /// <summary>
        /// 會強轉小寫
        /// </summary>
        public static string Env
        {
            get
            {
                return (BaseSettingHelper.GetSettingValue("Env") ?? "").ToLower();
            }
        }

        /// <summary>
        /// 是否是測試環境正在運行
        /// </summary>
        public static bool IsMSTest
        {
            get
            {
                return WebSiteName == "target.MSTest.Local";
            }
        }


        public static bool IsProduction
        {
            get
            {
                return Env == SystemConstant.Env_Name_Production;
            }
        }


        public static string WebSiteName
        {
            get
            {
                return BaseSettingHelper.GetSettingValue("WebSiteName");
            }
        }

        /// <summary>
        /// 是否使用 redis/遠端 型的cache (而非memory cache)
        /// </summary>
        public static bool IsRemoteCache
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<bool>("IsRemoteCache" , false);
            }
        }


        public static string IsEnableCache
        {
            get
            {
                return BaseSettingHelper.GetSettingValue("IsEnableCache");
            }
        }

        public static int MaxServerPlayerCount
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<int>("MaxServerPlayerCount", 300);
            }
        }

        public static string DefaultLang
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("DefaultLang", "zh-cn");
            }
        }


        public static string ApiRequestLogTableName
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("ApiRequestLogTableName", "[dbo].[ApiRequestLog]");
            }
        }


        public static string ApiRequestErrorLogTableName
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("ApiRequestErrorLogTableName", "[dbo].[ApiRequestLogError]");
            }
        }


        /// <summary>
        /// 是否採用每個日期分成一張table的模式
        /// </summary>
        public static bool ApiRequestLogDailyTableEnable
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<bool>("ApiRequestLogDailyTableEnable", false);
            }
        }


        /// <summary>
        /// 是否採用hostService 定時同步localFile至儲存空間
        /// </summary>
        public static bool ApiRequestLogBackgroundSyncEnable
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<bool>("ApiRequestLogBackgroundSyncEnable", false);
            }
        }

        public static bool IsKestrelWebHost
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<bool>("IsKestrelWebHost", false);
            }
        }


        /// <summary>
        /// 是否採用hostService 定時同步localFile至儲存空間
        /// </summary>
        public static bool GetConfig_ApiRequestLogDailyScheculeSyncEnable(IConfiguration config)
        {
            var settingValue = config.GetSection("ApiRequestLogDailyScheculeSyncEnable")?.Value;
            return BaseSettingHelper.ConvertSettingValueToValueType<bool>(settingValue, false);
        }


        public static string WebHostListenPort
        {
            get
            {
                return BaseSettingHelper.GetSettingValue("WebHostListenPort");
            }
        }



        public static SlackNotifySetting SlackNotify => new SlackNotifySetting();

        public static GameAPISetting GameApiSetting => new GameAPISetting();

        public static AdminSysAPISetting AdminSysAPISetting => new AdminSysAPISetting();

        public static ThirdPartySolutionSetting ThirdParty => new ThirdPartySolutionSetting();

        public static PublicApiSetting PublicApiSetting => new PublicApiSetting();

        public static FileApiSetting FileApiSetting => new FileApiSetting();

    }

    public class FileApiSetting
    {
        public string MinioApiEndPoint
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("MinioApiEndPoint", "");
            }
        }

        public string MinioApiBucket
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("MinioApiBucket", "");
            }
        }

        public string MinioApiAccessKey
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("MinioApiAccessKey", "");
            }
        }


        public string MinioApiSecretKey
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("MinioApiSecretKey", "");
            }
        }

        public string MinioImageUrl
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("MinioImageUrl", "");
            }
        }

        public bool IsMinioUseSSL
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<bool>("IsMinioUseSSL", false);
            }
        }
    }


    public class PublicApiSetting
    {
        public int Slot88MD5ExpiredMinutes
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<int>("Slot88MD5ExpiredMinutes", 10080);
            }
        }
    }


    public class ThirdPartySolutionSetting
    {

        private static string Prefix = "ThirdParty:";




        public string RedisConn
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>(Prefix + "RedisConn", "");
            }
        }

        /// <summary>
        /// 某些環境會共用redis，然後會用 DB編號區別開來
        /// </summary>
        public int RedisDBNumber
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<int>(Prefix + "RedisDBNumber", 0);
            }
        }

    }


    public class AdminSysAPISetting
    {
        private static string Prefix = "AdminSys:";

        /// <summary>
        /// 對外開放可連線的url (需請 SRE 開通防火牆)
        /// </summary>
        public string PublicUrlDomain
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>(Prefix + "PublicUrlDomain", "");
            }
        }

    }






    public class GameAPISetting
    {

        /// <summary>
        /// Game API 請求都要帶這個
        /// </summary>
        public string ApiCertificateValue
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("ApiCertificateValue", "851b-222e5f24024b");
            }
        }

    }




    public class SlackNotifySetting
    {
        private static string Prefix = "SlackNotify:";


        /// <summary>
        /// 當跳出一般C# exception的時候發通知
        /// </summary>
        public bool IsNotifyEnable_When_SystemException
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<bool>(Prefix + "IsNotifyEnable_When_SystemException", true);
            }
        }


        public string NotifyChannel_When_SystemException
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>(Prefix + "NotifyChannel_When_SystemException", "n_error_game_api");
            }
        }

    }





}
