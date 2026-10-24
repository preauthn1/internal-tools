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
using target.Utility.Helper.HttpContext;
using target.Models;
using target.Models.Enums;

namespace target.Utility.Helper
{

    public static class HeaderValueHelper
    {       

        public static string GetLocaleHeader()
        {
            var value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_Lang);
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


        public static LanguageEnum ConvertLangCodeToEnum(string value)
        {
            if (string.IsNullOrWhiteSpace(value))
                return SystemConstant.DefaultLanguageEnum;

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

        public static string GetSignKey()
        {
            string value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_SignKey);
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


        public static string GetHDCode()
        {
            string value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_HDCode);
            return value;
        }


        public static string HeaderTemp_ApiUniqueRouteName()
        {
            string value = MyHttpContext.Current.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_Temperary_ApiUniqueRouteName);
            return value;
        }


        public static string GetHeaderValueByKey(string key)
        {
             return MyHttpContext.Current.Request.GetHeaderValueByKey(key);
        }


        public static void AppendDebugInfo2(string info)
        {
            string key = "DebugInfo2";

            if (MyHttpContext.Current.Request.Headers.ContainsKey(key))
            {
                var matchValue = MyHttpContext.Current.Request.Headers[key];
                MyHttpContext.Current.Request.Headers[key] = matchValue + info;
            }
            else
            {
                MyHttpContext.Current.Request.Headers.Add(key, info);
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
