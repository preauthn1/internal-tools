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
using Microsoft.AspNetCore.Routing;
using target.Utility.Helper;
using target.Utility.Extensions.Web;
using Microsoft.AspNetCore.Mvc.Controllers;
using System.Collections.Generic;

namespace target.Web.Modules.Filter
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

            //如果是swaggerAPI 強制不處理
            if (url.Contains("swagger"))
                return;

            DateTime startTime = DateTimeHelper.GetCurrentTime();

            //記錄log資訊的
            Stream originalBodyStream = context.Response.Body;
            using (var responseBody = new MemoryStream())
            {
                context.Response.Body = responseBody;
                await _next(context);
                
                bool isSuccessResponse = CheckIsSuccessResponse(context.Response.StatusCode);

                DateTime endTime = DateTimeHelper.GetCurrentTime();
                int executedMilliseconds = (int)(endTime - startTime).TotalMilliseconds;

                //寫log
                WriteRequestAndResponseLog(context , isSuccessResponse , executedMilliseconds);

                //Copy the contents of the new memory stream (which contains the response) to the original stream, which is then returned to the client.
                await responseBody.CopyToAsync(originalBodyStream);
            }

            return;
        }

        private void WriteRequestAndResponseLog(HttpContext context , bool isSuccessResponse , int executedMilliseconds)
        {
            try
            {
                //紀錄客製的log
                var logModel = new ApiRequestLogEntity();

                //沒執行這段 response會掛掉...! (重要)
                HttpRequestLogHelper.SetResponseInfo(context.Response, logModel);
                
                if(isSuccessResponse)
                {
                    HttpRequestLogHelper.SetRequestInfo(context.Request, logModel);
                    logModel.DebugInfo3 = executedMilliseconds.ToString();

                    WriteApiRequestLogToTable(logModel);
                }
            }
            catch (Exception ex)
            {
                string msg = ex.Message;
            }
        }


        private void WriteApiRequestLogToTable(ApiRequestLogEntity logModel)
        {
            bool isLogEnable = ApiConfigDefineHelper.GetApiConfigSettingByCurrentRequest().IsLogApiEnable;

            if (!isLogEnable)
                return;

            bool is404NotFound = logModel.ResponseStatusCode == 404;
            string tableName = DetermineTableName(is404NotFound);
            
            if (is404NotFound)
            {
                List<string> skipKeyword = new List<string>
                {
                    "acme-challenge" ,
                    ".env" ,
                    "favicon.ico" ,
                    "wp-content" ,
                    "robots.txt" ,
                    "owa/"
                };

                //好像有一些 爬蟲/惡意讀取的請求 這些404也不記錄了
                foreach(string word in skipKeyword)
                {
                    if(logModel.RelativeUrl.Contains(word))
                    {
                        return;
                    }
                }
            }              

            var service = IOCHelper.Resolve<IApiRequestLogService>();

            bool isSyncLocalFileMode = SettingHelper.ApiRequestLogBackgroundSyncEnable;

            //根據設置採用直接寫db 或 先寫到local
            if(isSyncLocalFileMode)
            {
                service.WriteToLocalFile(logModel);
            }
            else
            {
                service.CreateTryCatch(tableName, logModel);
            }
        }


        private string DetermineTableName(bool is404NotFound)
        {
            if (is404NotFound)
                return SettingHelper.ApiRequestErrorLogTableName;

            if (SettingHelper.ApiRequestLogDailyTableEnable)
            {
                string tableName = AppHelper.ConvertToDateSequenceTableName(SettingHelper.ApiRequestLogTableName);
                return tableName;
            }
                
            return SettingHelper.ApiRequestLogTableName;
        }


        private bool CheckIsSuccessResponse(int statusCode)
        {
            return statusCode < 400 || statusCode == 404;
        }


    }

}
