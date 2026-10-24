using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using MSC.Models.ViewModels.Base;
using MSC.Utility.Extensions;
using MSC.Models;
using MSC.Models.Models;
using MSC.Utility.Helper;
using System;
using MSC.Utility;
using MSC.Utility.IOC;
using MSC.BL.Interface;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Logging;
using example.invalid;
using MSC.Models.Table;
using MSC.Utility.Helper.Dependency;
using MSC.Utility.Utility.ThirdParty;
using Microsoft.AspNetCore.Diagnostics;

namespace MSC.Web.Modules.Filter
{

    public static class ExceptionHandler
    {
        public static void UseGlobalExceptionHandler(this IApplicationBuilder appBuilder)
        {
            //突然想不起來當初為什麼用這個 而不是用exceptionHanlder了 (目前這種方案 有造成一些優化機制 無法實作)
            //是因為middleware、filter的一些 優先順序問題嗎...?
            appBuilder.UseExceptionHandler(app =>
            {
                app.Run(async context =>
                {
                    var ex = context.Features.Get<IExceptionHandlerFeature>()?.Error;

                    if (ex != null)
                    {
                        //紀錄客製的log
                        var logModel = new ApiRequestLogEntity();
                        HttpRequestLogHelper.SetRequestInfo(context.Request, logModel);


                        int statusCode = 400;
                        var customizeException = ex as MSCException;

                        var res = new ApiResponseErrorModel
                        {
                            Code = 400001,
                            Message = ex.GetMessageAndErrorLine(),
                            Data = null
                        };

                        if (customizeException != null)
                        {
                            res.Code = customizeException.ErrroCode.ToInteger();
                            res.Data = customizeException.ErrorDetail;
                            res.Message = MultiLanguageTextHelper.GetErrorCodeMsg(customizeException.ErrroCode, customizeException.ErrorDetail);
                        }
                        else
                        {
                            int rndText = RandomHelper.GetRandIntegerNumbers(1000, 9999);

                            //如果是production環境 錯誤訊息不顯示太直白
                            if(SettingHelper.IsProduction)
                            {
                                logModel.DebugInfo3 = rndText.ToString() + " " + res.Message;
                                res.Message = MultiLanguageTextHelper.DefaultErrorMsg(rndText.ToString());
                            }

                            SlackNotifyWhenCSharpException(ex , context.Request.Path);
                        }


                        logModel.StackTrace = ex.StackTrace;
                        logModel.ResponseStatusCode = statusCode;
                        SetErrorInfoHelper.SetDebugInfoBeforeWriteLog(logModel, res);


                        string tableName = SettingHelper.ApiRequestErrorLogTableName;
                        var service = IOCHelper.Resolve<IApiRequestLogService>();
                        service.CreateTryCatch(tableName, logModel);

                        context.Response.StatusCode = statusCode;
                        context.Response.ContentType = "application/json";
                        await context.Response.WriteAsync(res.ConvertToJsonString());
                    }
                });
            });
        }


        private static void SlackNotifyWhenCSharpException(Exception ex , string url)
        {
            if (!SettingHelper.SlackNotify.IsNotifyEnable_When_SystemException)
                return;

            //如果是C#原生的exception，檢查是否有設slack通知
            string channelName = SettingHelper.SlackNotify.NotifyChannel_When_SystemException;

            if (!string.IsNullOrWhiteSpace(channelName))
            {
                SlackChannelEnum enumParse = GetValidChannelName(channelName);
                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                {
                    Channel = enumParse,
                    Message = $" `{url}` ," + ex.GetMessageAndErrorLine(),
                });
            }
        }


        private static SlackChannelEnum GetValidChannelName(string channelName)
        {
            return ExtensionOfEnum.GetEnumByTextValue<SlackChannelEnum>(channelName, SlackChannelEnum.notify_msc_errors);
        }

    }



}
