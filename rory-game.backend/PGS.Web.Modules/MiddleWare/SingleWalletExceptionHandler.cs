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
using target.Models.Enums;
using target.BL.Services;
using target.Utility.Extensions.Web;
using target.Models.ViewModels;

namespace target.Web.Modules.Filter
{

    /// <summary>
    /// 單一錢包API 專用的exception處理
    /// </summary>
    public static class SingleWalletExceptionHandler
    {

        public static void UseSingleWalletVendorExceptionHandler(this IApplicationBuilder appBuilder)
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

                            if (customizeException.ErrorDetail != null && customizeException.ErrorDetail.IsAlwaysNotify)
                            {
                                string errorMsg = customizeException.Message;
                                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                                {
                                    Channel = GetNotifySlackChannel(),
                                    Message = $" `{context.Request.Path}` {errorMsg}",
                                });
                            }
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

                        //會根據APIConfig定義的 遊戲商不同 轉成各自遊戲商所需要的格式
                        SingleWalletExceptionHandleModel customizedResponse = ConvertApiResponseModelByConfigDefine(res);


                        SetErrorInfoHelper.SetDebugInfoBeforeWriteLog(logModel, res);


                        statusCode = customizedResponse.StatusCode;
                        logModel.StackTrace = isSkipWriteStackTraace ? "" : ex.StackTrace;
                        logModel.ResponseStatusCode = statusCode;
                        logModel.ResponseText = customizedResponse.GameVendorCustomizedResponse == null ? "" : customizedResponse.GameVendorCustomizedResponse.ConvertToJsonString();


                        //寫log
                        string tableName = SettingHelper.ApiRequestErrorLogTableName;
                        var service = IOCHelper.Resolve<IApiRequestLogService>();
                        service.CreateTryCatch(tableName, logModel);


                        context.Response.StatusCode = statusCode;
                        context.Response.ContentType = "application/json";

                        await context.Response.WriteAsync(customizedResponse.GameVendorCustomizedResponse.ConvertToJsonString());

                    }
                });
            });
        }


        private static SingleWalletExceptionHandleModel ConvertApiResponseModelByConfigDefine(ApiResponseErrorModel model)
        {
            var apiCategory = ApiConfigDefineHelper.GetApiConfigSettingByCurrentRequest().ApiCategory;

            if (apiCategory == ApiConfigCategoryEnum.DebugApi ||
                apiCategory == ApiConfigCategoryEnum.Default)
            {
                return new SingleWalletExceptionHandleModel
                {
                    StatusCode = 400,
                    GameVendorCustomizedResponse = model
                };
            }

            if (apiCategory == ApiConfigCategoryEnum.Default)
                throw new TerproException(ApiErrorCode.NotImplement);

            GameVendorTypeEnum convertToVendorID = (GameVendorTypeEnum)apiCategory.ToInteger();

            if (!convertToVendorID.IsDefinedValue())
                throw new TerproException(ApiErrorCode.Forbidden_NotValidDefineValue);

            var vendorService = ServiceInstanceHelper.InstanceVendorService(convertToVendorID);
            return vendorService.FormatParsing_ParseExceptionModelToVendorApiResponse(model);
        }


        private static void SlackNotifyWhenCSharpException(Exception ex, string url)
        {
            if (!SettingHelper.SlackNotify.IsNotifyEnable_When_SystemException)
                return;

            //如果是C#原生的exception
            SlackChannelEnum enumParse = GetNotifySlackChannel();
            MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
            {
                Channel = enumParse,
                Message = $" `{url}` ," + ex.GetMessageAndErrorLine(),
            });
        }

        private static SlackChannelEnum GetNotifySlackChannel()
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
