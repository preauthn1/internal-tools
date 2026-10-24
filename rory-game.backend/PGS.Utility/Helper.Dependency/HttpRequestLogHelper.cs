using Microsoft.AspNetCore.Http;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using example.invalid;
using target.Models;
using System.Linq;
using target.Utility.Extensions.Web;

namespace target.Utility.Helper.Dependency
{
    public class HttpRequestLogHelper
    {

        public static void SetRequestInfo(HttpRequest request , ApiRequestLogEntity model)
        {
            try
            {
                model.RequestUrl = $"{request.Scheme}://{request.Host.Value}{request.Path}{request.QueryString}";
                model.RelativeUrl = request.Path;
                model.HttpMethod = request.Method;
                model.Doamin = request.Host.Value;
                model.Headers = request.GetUsefulHeaders().ConvertToJsonString();
                model.IdentityKey = GetDeviceIdentityOrMemo(request);

                var apiConfig = ApiConfigDefineHelper.GetApiConfigSettingByCurrentRequest();

                //若request開關有設關閉/啟用
                if (apiConfig.IsRequestBodyLogEnable)
                {
                    model.RequestBody = request.GetRequestBodyString();
                }
                else
                {
                    model.RequestBody = "**Disable**";
                }

                //若開關有設關閉/啟用
                if (!apiConfig.IsHeaderLogEnable)
                {
                    model.Headers = "**Disable**";
                }

                //若response開關有設關閉/啟用
                if (!apiConfig.IsResponseLogEnable)
                {
                    model.ResponseText = "**Disable**";
                }
                
                string jwt = request.GetHeaderValueByKey(SystemConstant.HeaderKey_AccessToken);
                model.UserIdentity = GetUserIdentity(request);

                model.DebugInfo2 = request.GetHeaderValueByKey("DebugInfo2");
               
            }
            catch(Exception ex)
            {
                string debugMsg = ex.Message;
            }
            
            return;
        }



        private static string GetUserIdentity(HttpRequest request)
        {            
            string jwt = request.GetHeaderValueByKey(SystemConstant.HeaderKey_AccessToken);

            string userID = MyJwtHelper.GetUserIdentityTextFromJwt(jwt);

            if (!string.IsNullOrWhiteSpace(userID))
                return userID;

            bool isTryGetPlayerID = ApiConfigDefineHelper.GetApiConfigSettingByCurrentRequest().IsLoggingPlayerID;

            if (!isTryGetPlayerID)
                return string.Empty;

            if(request.Method == "GET")
            {
                return request.Query["PlayerID"];
            }

            string body = request.GetRequestBodyString();

            if (string.IsNullOrWhiteSpace(body))
                return string.Empty;

            if (body == "***MASK***")
                return string.Empty;

            return TryGetPlayerID(body);          
        }

        private static string GetDeviceIdentityOrMemo(HttpRequest request)
        {
            //有些API會用此欄位來記錄額外資訊
            string idenityMemo = request.GetHeaderValueByKey(SystemConstant.HeaderKey_IdentityMemo);

            if (!string.IsNullOrWhiteSpace(idenityMemo))
                return idenityMemo;

            var deviceID = HeaderValueHelper.GetHDCode();
            return deviceID;
        }


        private static string TryGetPlayerID(string jsonBody)
        {
            try
            {
                var parseModel = jsonBody.TryJsonConvertToModel<Dictionary<string, object>>();

                if (parseModel.ContainsKey("PlayerID"))
                    return parseModel["PlayerID"].ToString();

                if (parseModel.ContainsKey("playerID"))
                    return parseModel["playerID"].ToString();

                if (parseModel.ContainsKey("playerid"))
                    return parseModel["playerid"].ToString();
            }
            catch(Exception ex)
            {

            }

            return string.Empty;
        }


        public static void SetResponseInfo(HttpResponse response, ApiRequestLogEntity model)
        {
            try
            {
                model.ResponseStatusCode = response.StatusCode;
                model.ResponseLength = response.Body.Length;
             
                //We need to read the response stream from the beginning...
                response.Body.Seek(0, SeekOrigin.Begin);

                //...and copy it into a string
                string text = new StreamReader(response.Body).ReadToEndAsync().Result;

                //We need to reset the reader for the response so that the client can read it.
                response.Body.Seek(0, SeekOrigin.Begin);
                model.ResponseText = text;

            }
            catch(Exception ex)
            {
                string t = ex.Message;
            }


            //若response開關有設關閉/啟用
            if (!ApiConfigDefineHelper.GetApiConfigSettingByCurrentRequest().IsResponseLogEnable)
            {
                model.ResponseText = "**Disable**";
            }

            return;
        }


    }
}
