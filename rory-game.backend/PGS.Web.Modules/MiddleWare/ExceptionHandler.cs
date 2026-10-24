using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using target.Utility.Extensions;
using target.Models;
using target.Models.Models;
using target.Utility.Helper;
using System;
using System.Linq;
using target.Utility;
using target.Utility.IOC;
using target.BL.Interface;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Logging;
using example.invalid;
using target.Models.Table;
using target.Utility.Helper.Dependency;
using target.Utility.Utility.ThirdParty;
using Microsoft.AspNetCore.Diagnostics;
using System.Data.SqlClient;
using System.Collections.Generic;
using target.Utility.Extensions.Web;

namespace target.Web.Modules.Filter
{

    public static class ExceptionHandler
    {

        //特定環境、特定的API SQLException 錯誤訊息不slack通知 (因為很吵)
        public static Dictionary<string, List<string>> SkipSlackNotifyApis = new Dictionary<string, List<string>>
        {
            { 
                SystemConstant.Env_Name_Dev , new List<string>
                {
                    "game/slot/bonus",
                    "game/slot/prob",
                    "game/slot/user-assigned-mode",
                    "system/server/re-register",
                    "system/server/robots"
                }
            },

            {
                SystemConstant.Env_Name_QA , new List<string>
                {
                    "game/slot/bonus",
                }
            },

            {
                SystemConstant.Env_Name_Production , new List<string>
                {
                    "game/slot/bonus",
                }
            },
        };

        public static void UseGlobalExceptionHandler(this IApplicationBuilder appBuilder)
        {
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


                        bool isSkipWriteStackTraace = false;
                        int statusCode = 400;
                        var customizeException = ex as TerproException;

                        var res = new ApiResponseErrorModel
                        {
                            Code = ApiErrorCode.InternalServerError.ToInteger(),
                            Message = ex.GetMessageAndErrorLine(),
                            Data = null
                        };

                        if (customizeException != null)
                        {
                            res.Code = customizeException.ErrroCode.ToInteger();
                            res.Data = customizeException.ErrorDetail;
                            res.Message = MultiLanguageTextHelper.GetErrorCodeMsg(customizeException.ErrroCode, customizeException.ErrorDetail);

                            //某些errorCode太基本 不需要紀錄stackTrace (浪費寫入空間)
                            if (res.Code >= 401001 && res.Code <= 401099)
                            {
                                isSkipWriteStackTraace = true;
                            }

                            if(customizeException.ErrorDetail != null && customizeException.ErrorDetail.IsAlwaysNotify)
                            {
                                string errorMsg = customizeException.Message;
                                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                                {
                                    Channel = GetNotifySlackChannel(),
                                    Message = $" `{context.Request.Path}` {errorMsg}",
                                });
                            }

                            logModel.DebugInfo2 = customizeException.ErrorDetail?.LogExtraInfo;

                        }
                        else
                        {
                            int rndText = RandomHelper.GetRandIntegerNumbers(1000, 9999);

                            //如果是production環境 錯誤訊息不顯示太直白
                            if (SettingHelper.IsProduction)
                            {
                                logModel.DebugInfo3 = rndText.ToString() + " " + res.Message;
                                res.Message = MultiLanguageTextHelper.DefaultErrorMsg(rndText.ToString());
                            }

                            SlackNotifyWhenCSharpException(ex, context.Request.Path);
                        }



                        logModel.StackTrace = isSkipWriteStackTraace ? "" : ex.StackTrace;
                        logModel.ResponseStatusCode = statusCode;
                        SetErrorInfoHelper.SetDebugInfoBeforeWriteLog(logModel, res);


                        string tableName = SettingHelper.ApiRequestErrorLogTableName;
                        var service = IOCHelper.Resolve<IApiRequestLogService>();
                        service.CreateTryCatch(tableName, logModel);


                        bool isRedirect = ApiConfigDefineHelper.GetApiConfigSettingByCurrentRequest().IsRedirectToErrorPage;

                        if (!isRedirect)
                        {
                            context.Response.StatusCode = statusCode;
                            context.Response.ContentType = "application/json";
                            await context.Response.WriteAsync(res.ConvertToJsonString());
                        }
                        else
                        {
                            context.Response.Redirect($"/error?code={res.Code}&msg={res.Message.UrlEncode()}");
                        }
                    }
                });
            });
        }


        private static void SlackNotifyWhenCSharpException(Exception ex, string url)
        {
            if (!SettingHelper.SlackNotify.IsNotifyEnable_When_SystemException)
                return;

            //定義特定的exception不發slack通知 因為很吵
            if (IsSkipSlackNotifyExceptionType(ex, url))
                return;
          
            MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
            {
                Channel = GetNotifySlackChannel(),
                Message = $" `{url}` ," + ex.GetMessageAndErrorLine(),
            });            
        }


        public static bool IsSkipSlackNotifyExceptionType(Exception ex, string url)
        {
            string currentEnv = SettingHelper.Env;

            if (!SkipSlackNotifyApis.ContainsKey(currentEnv))
                return false;

            var defineApis = SkipSlackNotifyApis[currentEnv];

            //是否為定義要跳過的API
            if (!defineApis.Any(keyword=> url.Contains(keyword)))
                return false;

            bool isSQLException = IsPgsDefinedSqlException(ex);

            if (isSQLException)
                return true;

            return false;
        }


        private static bool IsPgsDefinedSqlException(Exception ex)
        {
            var sqlEx = ex as SqlException;

            if (sqlEx == null)
                return false;

            //鍵值重複的錯誤碼
            bool isDataUniqueException = sqlEx.Number == 2601 || sqlEx.Number == 2627;

            bool isTimeoutException = sqlEx.Number == -2;

            //只有這幾種的sql錯誤會跳過
            if (isDataUniqueException || isTimeoutException)
                return true;

            return false;
        }


        public static SlackChannelEnum GetNotifySlackChannel()
        {
            string channelName = SettingHelper.SlackNotify.NotifyChannel_When_SystemException;

            if (string.IsNullOrWhiteSpace(channelName))
                return SlackChannelEnum.n_error_pgs_common;

            return GetValidChannelName(channelName);
        }


        private static SlackChannelEnum GetValidChannelName(string channelName)
        {
            return ExtensionOfEnum.GetEnumByTextValue<SlackChannelEnum>(channelName, SlackChannelEnum.n_error_game_api);
        }

    }



}
