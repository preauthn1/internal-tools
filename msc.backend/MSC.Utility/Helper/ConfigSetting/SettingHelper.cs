using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using MSC.Utility.Extensions;
using example.invalid;
using example.invalid;
using Microsoft.AspNetCore.Http;
using System.Reflection;
using System.Runtime.Versioning;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.FileProviders;
using MSC.Utility.Helper.HttpContext;
using MSC.Models.Enums;
using Lck.Utility.Extensions.Helper;

namespace MSC.Utility
{
    public static class SettingHelper
    {

        public static bool IsLocalHost
        {
            get
            {
                string url = MyHttpContext.Current.Request.Host.ToString();
                return url.Contains("localhost");
            }
        }


        public static bool IsDev
        {
            get
            {
                return Env == "Development" || Env == "Dev";
            }
        }

        public static int ApplicationListenPort
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<int>("ApplicationListenPort", 11441);
            }
        }


        public static string Env
        {
            get
            {
                return BaseSettingHelper.GetSettingValue("Env");
            }
        }

        /// <summary>
        /// 是否是測試環境正在運行
        /// </summary>
        public static bool IsMSTest
        {
            get
            {
                return WebSiteName == "MSC.MSTest.Local";
            }
        }

        public static bool IsUAT
        {
            get
            {
                return (Env ?? "").ToLower() == "uat";
            }
        }


        public static bool IsProduction
        {
            get
            {
                return (Env ?? "").ToLower() == "production";
            }
        }

        public static string WebSiteName
        {
            get
            {
                return BaseSettingHelper.GetSettingValue("WebSiteName");
            }
        }

        public static string IsEnableCache
        {
            get
            {
                return BaseSettingHelper.GetSettingValue("IsEnableCache");
            }
        }


        public static string DefaultLang
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("DefaultLang", "zh-CN");
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


        public static string WebHostListenPort
        {
            get
            {
                return BaseSettingHelper.GetSettingValue("WebHostListenPort");
            }
        }



        public static SlackNotifySetting SlackNotify => new SlackNotifySetting();

        public static ClientWebSetting WebApplicationSetting => new ClientWebSetting();

        public static AdminSysWebsiteSetting AdminSysWebsiteSetting => new AdminSysWebsiteSetting();

        public static ThirdPartySolutionSetting ThirdParty => new ThirdPartySolutionSetting();

        public static CacheSetting CacheSetting => new CacheSetting();

    }


    public class CacheSetting
    {
        private static string Prefix = "Cache:";

        public CacheTypeEnum CacheType
        {
            get
            {
                var setting = BaseSettingHelper.GetSettingValueOrDefaultValue<int>(Prefix + "CacheType", CacheTypeEnum.MemoryCache.ToInteger());
                return (CacheTypeEnum)setting;
            }
        }

        public string RedisConn
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>(Prefix + "RedisConn", "");
            }
        }


        public int RedisDBNumber
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<int>(Prefix + "RedisDBNumber", 1);
            }
        }
    }


    public class AdminSysWebsiteSetting
    {
        public string GoogleAuth_SystemName_EN
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("GoogleAuth_SystemName_EN", "");
            }
        }

        public string GoogleAuth_Saltkey
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("GoogleAuth_Saltkey", "");
            }
        }


        public string HeaderEnvText
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("HeaderEnvText", "");
            }
        }

        public string HeaderEnvColor
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("HeaderEnvColor", "");
            }
        }

        public int WebSiteBrandingType
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<int>("WebSiteBrandingType", WebSiteBrandingTypeEnum.MSCard.ToInteger());
            }
        }

        public string F2E_AdminSys_API_Doamin
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("F2E_AdminSys_API_Doamin", "");
            }
        }
    }


    public class ClientWebSetting
    {        

        public string ForceLanguage
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("ForceLanguage", "");
            }
        }


        public int BrandingType
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<int>("WebSiteBrandingType", WebSiteBrandingTypeEnum.MSCard.ToInteger());
            }
        }

        public string F2E_AdminSys_API_Doamin
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("F2E_AdminSys_API_Doamin", "");
            }
        }

        public string F2E_Point_API_Doamin
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("F2E_Point_API_Doamin", "");
            }
        }

        public string F2E_Point_API_Doamin_Internal
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("F2E_Point_API_Doamin_Internal", "");
            }
        }



        public string F2E_Vendor_API_Doamin
        {
            get
            {
                return BaseSettingHelper.GetSettingValueOrDefaultValue<string>("F2E_Vendor_API_Doamin", "");
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



    public class ThirdPartySolutionSetting
    {

        private static string Prefix = "ThirdParty:";

    }
}
