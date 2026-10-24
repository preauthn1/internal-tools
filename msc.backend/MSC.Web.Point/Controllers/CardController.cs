using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.Base;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Utility.Utility;
using MSC.Web.Point.Models;
using MSC.Utility.Extensions;
using MSC.Web.Point.Helper;

namespace MSC.Web.Point.Controllers
{
    public class CardController : BaseController
    {
        private readonly ILogger<HomeController> _logger;

        public ITokenCheckRefService TokenCheckRefService { get; set; }


        public CardController(ILogger<HomeController> logger)
        {
            _logger = logger;
        }



        public IActionResult Input(string token = "", string lang = "")
        {
            SetForceLanguage();

            var tokenParse = CheckTokenIsValid(token);

            //設置語系
            string targetLang = GetLanguageFromTokenOrQueryString(tokenParse.Lang, lang);
            HeaderValueHelper.SetLanguageHeaders(targetLang);


            var debug = HeaderValueHelper.GetCurrentLanguageCode();

            ViewBag.MSCRequestID = tokenParse.MSCRequestID;
            ViewBag.Token = [REDACTED];
            ViewBag.AllowAmount = tokenParse.AllowAmount;
            ViewBag.Device = GetDeviceId_Or_AddNew();
            ViewBag.ReturnUrl = tokenParse.ReturnUrl;
            ViewBag.ShoppingStore = SettingDataHelper.GetShoppingStoreByCache();
            ViewBag.ShoppingIframeUrl = SettingDataHelper.GetShoppingIframeUrl();
            ViewBag.PageLayoutSettings = SettingDataHelper.GetPageLayoutSettings(tokenParse.MerchantID);
            

            return View();
        }


        private LinkTokenMemberInfoModel CheckTokenIsValid(string token)
        {
            //自己Web這一層先檢查基本格式 (減少request存取浪費)
            var tokenParse = TokenCheckRefService.CheckIsTokenFormatCorrect(token);

            var data = new TokenCheckRequestModel
            {
                Token = [REDACTED]
            };

            //呼叫API檢查
            string baseurl = SettingHelper.WebApplicationSetting.F2E_Point_API_Doamin_Internal;
            var res = APIRequestHelper.RequestApiResult<ApiResponseModel<LinkTokenMemberInfoModel>>(new RequestParameters
            {
                BaseUrl = baseurl,
                ApiUrl = "token/valid-check",
                Method = HttpMethodType.HttpPost

            }, data);

            if (res.Exception != null)
            {
                if (res.IsTimeoutException)
                {
                    throw new MSCException(ErrorCode.CommonError_Timeout);
                }

                throw new MSCException(ErrorCode.ThirdPartyApiFail_TokenExpired);
            }

            if (res.ResultObject.Data == null)
            {
                throw new MSCException(ErrorCode.ParsingError_ParseTokenObjectFail);
            }

            return tokenParse;
        }



        private string GetLanguageFromTokenOrQueryString(LanguageEnum tokenLang, string queryLang)
        {
            //if(SettingHelper.WebApplicationSetting.ForceLanguage.IsNotNullNotWhiteSpace())
            //  return HeaderValueHelper.ConvertLangCodeToEnum(SettingHelper.WebApplicationSetting.ForceLanguage).ToString();

            if (!string.IsNullOrWhiteSpace(queryLang))
                return HeaderValueHelper.ConvertLangCodeToEnum(queryLang).ToString();

            return tokenLang.ToString();
        }


        private void SetForceLanguage()
        {
            if (SettingHelper.WebApplicationSetting.ForceLanguage.IsNotNullNotWhiteSpace())
            {
                HeaderValueHelper.SetLanguageHeaders(SettingHelper.WebApplicationSetting.ForceLanguage);
            }
        }


        private string GetDeviceId_Or_AddNew()
        {
            var deviceID = HttpContext.Request.Cookies["device"];

            if (!string.IsNullOrWhiteSpace(deviceID))
                return deviceID;

            string guid = Guid.NewGuid().ToString();
            HttpContext.Response.Cookies.Append("device", guid);
            return guid;
        }


    }
}
