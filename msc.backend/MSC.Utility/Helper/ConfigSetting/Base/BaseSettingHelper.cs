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

namespace MSC.Utility
{
    public static class BaseSettingHelper
    {

        public static string GetJsonSettingFileNameByEnv()
        {
            string envName = RuntimeHelper.GetCurrentEnviromentSettingName();
            return $"appsettings.{envName}.json";
        }


        public static string GetSettingValue(string keyName)
        {
            var config = InstanceConfigProvider();

            //object結構的話 需代入類似這樣的路徑字串 ConnectionStrings:MyDB
            var settingValue = config.GetSection(keyName).Value;
            return settingValue;
        }


        public static T GetSettingValueOrDefaultValue<T>(string keyName , T defaultValue)
        {
            var config = InstanceConfigProvider();

            //object結構的話 需代入類似這樣的路徑字串 ConnectionStrings:MyDB
            var settingValue = config.GetSection(keyName).Value;

            if (string.IsNullOrWhiteSpace(settingValue))
                return defaultValue;

            return (T)Convert.ChangeType(settingValue, typeof(T)); ;
        }


        private static IConfigurationRoot InstanceConfigProvider()
        {
            var builder = new ConfigurationBuilder()
                  .SetFileProvider(new PhysicalFileProvider(Directory.GetCurrentDirectory()))
                   .AddJsonFile("appsettings.json")
                   .AddJsonFile(GetJsonSettingFileNameByEnv(), optional: true);

            var config = builder.Build();
            return config;
        }

    }

}
