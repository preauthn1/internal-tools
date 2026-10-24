using Microsoft.AspNetCore.Http;
using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Text;
using MSC.Utility.Extensions;
using example.invalid;
using MSC.Models;
using System.Linq;

namespace MSC.Utility.Helper.Dependency
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
                model.IdentityKey = HeaderValueHelper.GetDeviceID();

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

                //若response開關有設關閉/啟用
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
                model.UserIdentity = MyJwtHelper.GetUserIdentityTextFromJwt(jwt);
            }
            catch(Exception ex)
            {
                string debugMsg = ex.Message;
            }
            

            return;
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
