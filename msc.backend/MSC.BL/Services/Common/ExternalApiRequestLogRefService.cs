using MSC.DA.Repository;
using MSC.Models.Table;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Utility.Utility;
using MSC.Utility.Extensions;
using System;
using System.Collections.Generic;


namespace MSC.BL.Services
{

    public class ExternalApiRequestLogRefService : BaseService, IExternalApiRequestLogRefService
    {
        public ExternalApiRequestLogRepository ExternalApiRequestLogRepository { get; set; }


        public void Create(ApiRequestLogTemporaryInfoModel model)
        {
            string relativeUrl = model.RequestParameters.ApiUrl == APIRequestHelper.No_Url_Keyword ? "" : model.RequestParameters.ApiUrl;

            var insertModel = new ExternalApiRequestLogEntity
            {
                CreateTime = DateTimeHelper.GetCurrentTime(),
                ResponseText = model.Response.ResponseContent,
                ResponseTime = model.Response.ResponseTime,
                ResponseStatusCode = model.Response.ResponseStatusCode,
                ApiCategory = model.RequestParameters.AdvancedSetting?.Log_ApiCategory,
                Group = model.RequestParameters.AdvancedSetting?.Log_Group,
                RequestUUID = model.RequestParameters.AdvancedSetting?.Log_RequestUUID,
                RequestUrl = (model.RequestParameters.BaseUrl + relativeUrl).ReplaceDuplicateUrlChar(),
                RequestContent = model.RequestContent,
                Memo = string.Empty,
                WebSiteName = SettingHelper.WebSiteName
            };

            ExternalApiRequestLogRepository.Create(insertModel);
        }


    }
}
