using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;
using target.Utility.Utility;
using target.Utility;

namespace target.BL.Services
{

    public class ExternalApiRequestLogRefService : BaseService, IExternalApiRequestLogRefService
    {
        public ExternalApiRequestLogRepository ExternalApiRequestLogRepository { get; set; }

        public SingleWalletMerchantApiRequestLogRepository SingleWalletMerchantApiRequestLogRepository { get; set; }


        public void Create(ApiRequestLogTemporaryInfoModel model)
        {
            try
            {
                string relativeUrl = model.RequestParameters.ApiUrl == APIRequestHelper.No_Url_Keyword ? "" : model.RequestParameters.ApiUrl;

                var logInfo = model.RequestParameters?.AdvancedSetting;

                var insertModel = new ExternalApiRequestLogEntity
                {
                    CreateTime = DateTimeHelper.GetCurrentTime(),
                    ResponseText = model.Response.ResponseContent,
                    ResponseTime = model.Response.ResponseTime,
                    ResponseStatusCode = model.Response.ResponseStatusCode,
                    RequestUrl = (model.RequestParameters.BaseUrl + "/" + relativeUrl).ReplaceDuplicateUrlChar(),
                    RequestContent = model.RequestContent,
                    RequestHeader = ToKeyValueDicJsonString(model.RequestParameters),
                    Memo = string.Empty,
                    WebSiteName = SettingHelper.WebSiteName,
                    Category = logInfo?.Log_ApiCategory,
                    SubCategory = logInfo?.Log_SubCategory,
                    RequestID = logInfo?.Log_RequestUUID,
                    DebugInfo = logInfo.Log_DebugInfo1,
                    DebugInfo2 = logInfo.Log_DebugInfo2,
                    DebugInfo3 = logInfo.Log_DebugInfo3,
                };

                ExternalApiRequestLogRepository.Create(insertModel);
            }
            catch(Exception ex)
            {
                string qqq = "";
            }

        }


        public void CreateSingleWalletApiLog(ApiRequestLogTemporaryInfoModel model)
        {
            try
            {
                string relativeUrl = model.RequestParameters.ApiUrl == APIRequestHelper.No_Url_Keyword ? "" : model.RequestParameters.ApiUrl;

                var logInfo = model.RequestParameters?.AdvancedSetting;

                var insertModel = new SingleWalletMerchantApiRequestLogEntity
                {
                    CreateTime = DateTimeHelper.GetCurrentTime(),
                    ResponseText = model.Response.ResponseContent,
                    ResponseTime = model.Response.ResponseTime,
                    ResponseStatusCode = model.Response.ResponseStatusCode,
                    RequestUrl = (model.RequestParameters.BaseUrl + "/" + relativeUrl).ReplaceDuplicateUrlChar(),
                    RequestContent = model.RequestContent,
                    RequestHeader = ToKeyValueDicJsonString(model.RequestParameters),
                    Memo = string.Empty,
                    WebSiteName = SettingHelper.WebSiteName,
                    Category = logInfo?.Log_ApiCategory,
                    SubCategory = logInfo?.Log_SubCategory,
                    RequestID = logInfo?.Log_RequestUUID,
                    DebugInfo = logInfo.Log_DebugInfo1,
                    DebugInfo2 = logInfo.Log_DebugInfo2,
                    DebugInfo3 = logInfo.Log_DebugInfo3,
                };

                SingleWalletMerchantApiRequestLogRepository.Create(insertModel);
            }
            catch (Exception ex)
            {
                string qqq = "";
            }

        }



        private string ToKeyValueDicJsonString(RequestParameters request)
        {
            if (request.Headers.IsEmptyOrNull())
                return null;

            try
            {
                var dic = request.Headers?.ToDictionary(x => x.Key, z => z.Value);
                dic.Add("User-Agent", request.UserAgent);
                dic.Add("Content-Type", request.ContentType);

                return dic.ConvertToJsonString();
            }
            catch
            {

            }

            return null;
        }


    }
}
