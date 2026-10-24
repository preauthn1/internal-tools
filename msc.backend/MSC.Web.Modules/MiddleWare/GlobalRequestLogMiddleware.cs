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
using Microsoft.AspNetCore.Routing;
using MSC.Utility.Helper;

namespace MSC.Web.Modules.Filter
{

    /// <summary>
    /// 全域api log日誌
    /// </summary>
    public class GlobalRequestLogMiddleware
    {
        private readonly RequestDelegate _next;

        public GlobalRequestLogMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        public async Task Invoke(HttpContext context)
        {
            var url = context.Request.Path.ToString();

            //如果是swaggerAPI 不使用
            if (url.Contains("swagger"))
                return;

            DateTime startTime = DateTimeHelper.GetCurrentTime();

            //記錄log資訊
            Stream originalBodyStream = context.Response.Body;
            using (var responseBody = new MemoryStream())
            {
                context.Response.Body = responseBody;
                await _next(context);
                
                bool isStatusCodeNeedWriteLog = CheckIsWriteLogByStatusCode(context.Response.StatusCode);

                DateTime endTime = DateTimeHelper.GetCurrentTime();
                int executedMilliseconds = (int)(endTime - startTime).TotalMilliseconds;

                //寫log
                WriteRequestAndResponseLog(context , isStatusCodeNeedWriteLog , executedMilliseconds);

                //Copy the contents of the new memory stream (which contains the response) to the original stream, which is then returned to the client.
                await responseBody.CopyToAsync(originalBodyStream);
            }

            return;
        }



        private bool CheckIsWriteLogByStatusCode(int statusCode)
        {
            return statusCode < 400 || statusCode == 404;
        }


        private void WriteRequestAndResponseLog(HttpContext context , bool isSuccessResponse , int executedMilliseconds)
        {
            try
            {
                //紀錄客製的log
                var logModel = new ApiRequestLogEntity();

                HttpRequestLogHelper.SetRequestInfo(context.Request, logModel);

                //沒執行這段 response會掛掉...! (重要)
                HttpRequestLogHelper.SetResponseInfo(context.Response, logModel);
                
                if(isSuccessResponse)
                {
                    logModel.DebugInfo3 = executedMilliseconds.ToString();
                    WriteApiRequestLogToTable(logModel);
                }
            }
            catch (Exception ex)
            {
                string msg = ex.Message;
            }
        }


        private static void WriteApiRequestLogToTable(ApiRequestLogEntity logModel)
        {
            bool isLogEnable = ApiConfigDefineHelper.GetApiConfigSettingByCurrentRequest().IsLogApiEnable;

            if (!isLogEnable)
                return;

            if (logModel.HttpMethod == "OPTIONS") //不記錄此 浪費資源
                return;

            string tableName = SettingHelper.ApiRequestLogTableName;
            var service = IOCHelper.Resolve<IApiRequestLogService>();
            service.CreateTryCatch(tableName, logModel);
        }

    }

}
