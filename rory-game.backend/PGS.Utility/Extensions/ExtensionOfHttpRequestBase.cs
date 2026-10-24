using Microsoft.AspNetCore.Http;
using target.Models;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using example.invalid;
using System.Linq;
using System.Linq.Expressions;
using System.Reflection;
using System.Text;
using System.Web;

namespace target.Utility.Extensions
{
    public static class ExtensionOfHttpRequestBase
    {

        /// <summary>
        /// 有些太廢話、對debug無幫助的header就不特別紀錄了 避免佔空間
        /// </summary>
        private static List<string> skipHeaderKeys = new List<string>
        {
            "Accept",
            "Accept-Encoding",
            "Accept-Language",
            "Connection",
            "Cache-Control",
            "Content-Length",
            "X-Forwarded-Host",
            "X-Forwarded-Proto",
            "X-Forwarded-Port",
            "X-Scheme",
            "X-Different-Name",
            "X-Using-Nginx-Controller",
            "X-Request-ID",
            "Host",
            SystemConstant.HeaderKey_TemperaryBody,
            SystemConstant.HeaderKey_Temperary_ApiUniqueRouteName,
            SystemConstant.HeaderKey_Temperary_IsLogApiEnable,
            SystemConstant.HeaderKey_Temperary_IsRequestBodyLogEnable,
            SystemConstant.HeaderKey_Temperary_IsResponseLogEnable,
        };


        private static List<string> passwordApis = new List<string>
        {
            "login",
            "logon",
            "signin",
            "password",
            "resetp",
        };

        public static string GetQueryStringValueFromUrlFormat(this string url , string parameterName)
        {
            if (string.IsNullOrWhiteSpace(url))
                return string.Empty;

            if(url.Contains("?"))
            {
                url = url.Split('?')[1];
            }
        
            return HttpUtility.ParseQueryString(url).Get(parameterName);
        }

        
        public static bool IsFormUrlEncoded(this HttpRequest request)
        {
            return request.ContentType == "application/x-www-form-urlencoded";
        }


        public static bool IsPasswordApi(this HttpRequest request)
        {
            var relativeUrl = request.Path.ToString();

            if (passwordApis.Any(keyword => relativeUrl.Contains(keyword)))
                return true;

            return false;
        }

        public static Dictionary<string, string> GetUsefulHeaders(this HttpRequest request)
        {
            Dictionary<string, string> dic = new Dictionary<string, string>();

            foreach (string key in request.Headers.Keys)
            {
                if (dic.ContainsKey(key))
                    continue;

                if (skipHeaderKeys.Contains(key))
                    continue;

                dic.Add(key, request.Headers[key]);
            }

            return dic;
        }
        

        public static string GetHeaderValueByKey(this HttpRequest request , string headerKey)
        {
            string headerValue = request.Headers[headerKey];
            return headerValue;
        }


        public static Dictionary<string, string> GetAllHeaders(this HttpRequest request)
        {
            Dictionary<string, string> dic = new Dictionary<string, string>();

            foreach(string key in request.Headers.Keys)
            {
                if (dic.ContainsKey(key))
                    continue;

                dic.Add(key, request.Headers[key]);
            }

            return dic;
        }


        public static string GetRequestBodyString(this HttpRequest request , bool passwordProtected = true)
        {
            try
            {
                if (request.Method == "GET")
                    return string.Empty;

                if (passwordProtected && IsPasswordApi(request))
                    return "***MASK***";

                //如果是formurlEncode
                if(IsFormUrlEncoded(request))
                {
                    string body = string.Empty;

                    if(request.Form.Keys == null ||
                       request.Form.Keys.Count == 0)
                    {
                        return string.Empty;
                    }

                    foreach(var key in request.Form.Keys)
                    {
                        string keyValue = key.ToString();
                        body += $"&{keyValue}={request.Form[keyValue]}";
                    }

                    if (body.FirstOrDefault() == '&')
                        body = body.Substring(1, body.Length - 1);

                    return body;
                }
                else
                {
                    var queryTemp = request.Headers[SystemConstant.HeaderKey_TemperaryBody];

                    //如果有放在暫存中
                    if (!string.IsNullOrWhiteSpace(queryTemp))
                        return queryTemp;

                    //注意* .net core 第二次讀取會有問題。找不到好方法處理之

                    var body = request.Body;
                    
                    //This line allows us to set the reader for the request back at the beginning of its stream.
                    request.EnableBuffering();

                    //We now need to read the request stream.  First, we create a new byte[] with the same length as the request stream...
                    var buffer = new byte[Convert.ToInt32(request.ContentLength)];

                    //...Then we copy the entire request stream into the new buffer.
                    request.Body.ReadAsync(buffer, 0, buffer.Length);

                    //We convert the byte[] into a string using UTF8 encoding...
                    var bodyAsText = Encoding.UTF8.GetString(buffer);
                                        
                    //..and finally, assign the read body back to the request body, which is allowed because of EnableRewind()
                    request.Body = body;


                    return bodyAsText;
                }
            }
            catch(Exception ex)
            {
                string msg = ex.Message;
            }

            return string.Empty;
        }

        
    }

}
