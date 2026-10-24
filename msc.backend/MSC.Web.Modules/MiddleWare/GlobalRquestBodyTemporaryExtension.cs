using Microsoft.AspNetCore.Mvc.Filters;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using System;
using example.invalid;
using System.Text;
using MSC.Utility.Extensions;
using example.invalid;
using Newtonsoft.Json;
using MSC.Models.Table;
using MSC.Utility.Helper.Dependency;
using MSC.Utility.IOC;
using MSC.BL;
using MSC.Utility.Helper.HttpContext;
using MSC.Utility;
using MSC.BL.Interface;
using Microsoft.AspNetCore.Builder;
using MSC.Models;

namespace MSC.Web.Modules.Filter
{

    public static class GlobalRquestBodyTemporaryExtension
    {
        public static void SetRequestBodyTemporaryHttpHeader(this IApplicationBuilder appBuilder)
        {
            appBuilder.Use(async (context, next) =>
            {
                context.Request.EnableBuffering();
                var stream = context.Request.Body;

                using (var reader = new StreamReader(stream))
                {
                    if(context.Request.Method != "GET")
                    {
                        string requestBodyAsString = await reader.ReadToEndAsync();

                        //因為一直找不到很好的方法 又能記錄到requestBody + responseBody  
                        //只好寫一個很奇怪的東西 (將request header 暫存在body)
                        context.Request.Headers.Add(SystemConstant.HeaderKey_TemperaryBody, requestBodyAsString);
                    }

                    if (stream.CanSeek)
                        stream.Seek(0, SeekOrigin.Begin);

                    await next.Invoke();
                }
                
                return;
            });
        }


        private static void WriteRequestAndResponseLog(HttpContext context)
        {
            try
            {
                //紀錄客製的log
                var logModel = new ApiRequestLogEntity();

                HttpRequestLogHelper.SetRequestInfo(context.Request, logModel);

                //404的話不讀取
                if (context.Response.StatusCode != 404)
                {
                    HttpRequestLogHelper.SetResponseInfo(context.Response, logModel);
                }

                string tableName = SettingHelper.ApiRequestLogTableName;
                var service = IOCHelper.Resolve<IApiRequestLogService>();
                service.CreateTryCatch(tableName, logModel);

            }
            catch (Exception ex)
            {
                string msg = ex.Message;
            }
        }






    }

}
