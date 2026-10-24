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
using MSC.Utility.Helper.HttpContext;
using MSC.Models;
using MSC.Models.Enums;

namespace MSC.Utility.Helper
{

    public static class HeaderValueHelper
    {

        public static string GetDeviceID()
        {
            var value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_DeviceID);
            return value;
        }

        public static string GetLocaleHeader()
        {
            var value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_Lang);
            return value;
        }

        public static string GetCurrentLanguageCode()
        {
            var value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_Lang);

            if (string.IsNullOrWhiteSpace(value))
                return SettingHelper.DefaultLang;

            return value;
        }


        public static LanguageEnum GetCurrentLanguage()
        {
            var value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_Lang);

            if (string.IsNullOrWhiteSpace(value))
                return ConvertLangCodeToEnum(SettingHelper.DefaultLang);

            var parsing = ConvertLangCodeToEnum(value);
            return parsing;
        }


        public static LanguageEnum ConvertLangCodeToEnum(string value , string defaultLang = SystemConstant.DefaultLang)
        {
            if (string.IsNullOrWhiteSpace(value))
                value = defaultLang;

             value = value?.ToLower();
            value = value?.Replace("-", "_");

            var parsing = ExtensionOfEnum.GetEnumByTextValue<LanguageEnum>(value , SystemConstant.DefaultLanguageEnum);
            return parsing;
        }


        public static string GetUserAgent()
        {
            var value = MyHttpContext.Current.Request.GetHeaderValueByKey("User-Agent");
            return value;
        }

        public static string GetApiCertificate()
        {
            string value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_ApiCertificate);
            return value;
        }
        

        public static string GetJwtToken()
        {
            string value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_AccessToken);
            return value;
        }


        public static string HeaderTemp_ApiUniqueRouteName()
        {
            string value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_Temperary_ApiUniqueRouteName);
            return value;
        }


        public static void SetLanguageHeaders(string value)
        {
            if(MyHttpContext.Current.Request.Headers.ContainsKey(SystemConstant.HeaderKey_Lang))
            {
                MyHttpContext.Current.Request.Headers[SystemConstant.HeaderKey_Lang] = value;
            }
            else
            {
                MyHttpContext.Current.Request.Headers.Add(SystemConstant.HeaderKey_Lang, value);
            }
        }


        public static object GetHeaders()
        {
            try
            {
                var values = MyHttpContext.Current.Request.Headers;
                return values;
            }
            catch
            {

            }

            return null;
        }

    }

}
