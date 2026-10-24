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
using MSC.Utility.Helper;
using MSC.Utility.Helper.HttpContext;

namespace MSC.Utility
{

    public static class RuntimeHelper
    {

        public static string RandomID = Guid.NewGuid().ToString();


        public static string EnvironmentName { get; set; }


        public static string MyVersion { get; set; }



        /// <summary>
        /// **重要 請在startup的時候呼叫此
        /// </summary>
        /// <param name="envName"></param>
        public static void Initialize(string envName)
        {
            EnvironmentName = envName;
            MyVersion = GetMyVersion();
        }
        

        /// <summary>
        /// 取得當前環境 設定
        /// </summary>
        /// <returns></returns>
        public static string GetCurrentEnviromentSettingName()
        {
            return EnvironmentName;
        }



        public static string GetMyVersion()
        {
            try
            {
                string path = AppDomain.CurrentDomain.BaseDirectory + "AppData\\AppVersion.txt";
                example.invalid.StreamReader file = new example.invalid.StreamReader(path);
                var text = file.ReadToEnd();
                return text?.Replace("\r\n", "");
            }
            catch(Exception ex)
            {
                string exx = ex.Message;
            }

            return "";
        }

        public static string GetGitSHA1()
        {
            try
            {
                return Environment.GetEnvironmentVariable("Git_SHA1");
            }
            catch
            {

            }

            return "";
        }

        public static string Get_Gitlab_Pipeline_ID()
        {
            try
            {
                return Environment.GetEnvironmentVariable("Pipeline_ID");
            }
            catch
            {

            }

            return "";
        }


        /// <summary>
        /// 判斷當前環境是否為 .net core
        /// </summary>
        /// <returns></returns>
        public static bool IsNetCoreRuntimeFramework()
        {
            var ver = GetRuntimeFrameworkVersion();

            if (string.IsNullOrWhiteSpace(ver))
                return false;

            if (ver.Contains("NETCoreApp"))
                return true;

            return false;
        }


        /// <summary>
        /// 取得當前運行的frameword版本 (.net core、.net frameword)
        /// </summary>
        /// <returns></returns>
        public static string GetRuntimeFrameworkVersion()
        {
            var framework = Assembly
              .GetEntryAssembly()?
              .GetCustomAttribute<TargetFrameworkAttribute>()?
              .FrameworkName;

            return framework;
        }



        public static DateTime? GetBinFileLastModifiedTime()
        {

            try
            {
                var location = System.Reflection.Assembly.GetExecutingAssembly().Location;
                var time = File.GetLastWriteTime(location);
                return time;
            }
            catch
            {

            }

            return null;
        }


        /// <summary>
        /// production環境，在未帶入password的情況下，首頁可以顯示的資訊有限制
        /// </summary>
        public static object GetHideInformationForProduction()
        {
            return new
            {
                EnvironmentName = SettingHelper.Env,
                WebSiteName = SettingHelper.WebSiteName,
            };
        }



        /// <summary>
        /// 首頁可以顯示的資訊
        /// </summary>
        /// <returns></returns>
        public static object GetHomePageDebugInformation()
        {
          
            return new
            {
                MyVersion = GetMyVersion(),
                EnvironmentName = SettingHelper.Env,
                WebSiteName = SettingHelper.WebSiteName,
                IsEnableCache = SettingHelper.IsEnableCache,
                configName = BaseSettingHelper.GetJsonSettingFileNameByEnv(),
                SHA1 = GetGitSHA1(),
                Pipeline_ID = Get_Gitlab_Pipeline_ID(),
                YourIP = IPHelper.GetIP(),
                LastModifiedTimeTW = GetBinFileLastModifiedTimeStringFormat(),
                //Headers = HeaderValueHelper.GetHeaders(),
            };
        }

        public static string GetBinFileLastModifiedTimeStringFormat()
        {
            DateTime? date = GetBinFileLastModifiedTime();

            if (date == null)
                return string.Empty;

            return date.Value.AddHours(8).ToString("yyyy-MM-dd HH:mm:ss");
        }

    }

}
