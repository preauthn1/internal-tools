using Microsoft.AspNetCore.Mvc.Filters;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using System;
using example.invalid;
using System.Text;
using target.Utility.Extensions;
using example.invalid;
using Newtonsoft.Json;
using target.Models.Table;
using target.Utility.Helper.Dependency;
using target.Utility.IOC;
using target.BL;
using target.Utility.Helper.HttpContext;
using target.Utility;
using target.BL.Interface;
using Microsoft.AspNetCore.Builder;
using target.Models;

namespace target.Web.Modules.Filter
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
                        //只好寫一個很奇怪的東西 (將request body 暫存在 httpHeader中)
                        context.Request.Headers.Add(SystemConstant.HeaderKey_TemperaryBody, requestBodyAsString);
                    }

                    if (stream.CanSeek)
                        stream.Seek(0, SeekOrigin.Begin);

                    await next.Invoke();
                }
                
                return;
            });
        }

    }

}
