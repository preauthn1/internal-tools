using System;
using System.Collections.Generic;
using target.BL.Interface;
using target.BL.Services;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Utility.Utility;
using System.Web;
using target.Utility.IOC;
using target.Models.ViewModels.Base;

namespace target.BL.Services
{

    /// <summary>
    /// 跟供應商對接的API
    /// </summary>
    public class MerchantIntergratedExternalApiService : BaseVendorService, IMerchantIntergratedExternalApiService
    {
        
        /// <summary>
        /// 檢查 玩家Session是否合法
        /// </summary>
        public MerchantIntergrationApiResultModel CheckPlayerSessionIsValid(MerchantIntergrationApiParameterModel model)
        {
            var req = new RequestParameters
            {
                BaseUrl = model.ApiValidatePlayerSessionUrl,
                ApiUrl = APIRequestHelper.No_Url_Keyword,
                Method = HttpMethodType.HttpPost,
                ContentType = "application/x-www-form-urlencoded",
                LogAction = WriteRequestLog,
                AdvancedSetting = new AdvancedSetting
                {
                    Log_ApiCategory = "valid-session"
                }
            };

            //參考PG文件.pdf  (第19頁)  5.2.1 令牌验证 API 
            var body = new MerchantIntergrationApiRequestModel
            {
                PlayerName = model.PlayerName,
                PlayerValidateSession = model.PlayerValidateSession,
                MerchantID = model.MerchantID,
                MerchantSecretKey = model.MerchantSecretKey                
            };

            string formBody = ConvertToFormUrlFormat(body);
            var res = APIRequestHelper.RequestApiResult<MerchantIntergrationApiResultModel>(req, formBody);

            //檢查http錯誤、model解析是否成功
            if (res.Exception != null)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseError);

            if (res.ResultObject == null)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ParsingFail);
            
            MerchantIntergrationApiResultModel resultObject = res.ResultObject;

            //檢查回應欄位是否OK
            if(!resultObject.IsValidOK)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_IsNotValid);

            if(resultObject.PlayerName?.ToLower() != model.PlayerName.ToLower())
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_SignCheckNotMatch);

            return resultObject;
        }


        private string ConvertToFormUrlFormat(object baseReuqestModel)
        {
            // 請求FDI存款接口
            var properties = from p in baseReuqestModel.GetType().GetProperties()
                             where p.GetValue(baseReuqestModel, null) != null
                             select p.Name.ToCamelCase() + "=" + HttpUtility.UrlEncode(p.GetValue(baseReuqestModel, null).ToString());

            string formPostBody = String.Join("&", properties);
            return formPostBody.Trim();
        }
        

        private void WriteRequestLog(ApiRequestLogTemporaryInfoModel model)
        {
            var service = IOCHelper.Resolve<IExternalApiRequestLogRefService>();
            service.Create(model);
        }

    }
}
