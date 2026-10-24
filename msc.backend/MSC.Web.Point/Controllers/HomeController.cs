using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Utility.Utility;
using MSC.Web.Point.Models;
using MSC.Models.ViewModels;
using MSC.Utility.Helper;
using MSC.Models.ViewModels.Base;
using MSC.Models;
using MSC.Utility;
using MSC.Utility.Utility.ThirdParty;
using MSC.Models.Models;

namespace MSC.Web.Point.Controllers
{
    public class HomeController : BaseController
    {

        public ILinkGenerateService LinkGenerateService { get; set; }


        public HomeController()
        {

        }


        public IActionResult Index()
        {
            return View();
        }


        public IActionResult Success()
        {
            return View();
        }



        [Route("error")]
        public IActionResult Error(string msg = "系統忙碌中" , int code = 500001)
        {
            ViewBag.Msg = msg; 
            ViewBag.Code = code;
            return View();
        }


        public IActionResult ToCardInputLink()
        {
            bool isLocalhost = Request.Host.Value.Contains("localhost") || Request.Host.Value.Contains("0.0.0.0:11449");
            string targetDomain = isLocalhost ? $"http://{Request.Host}" : "http://0.0.0.0:11449";

            var data = new LinkGenerateRequestModel
            {
                SignKey = SystemConstant.API_Vendor_MagicCode,
                RequestInfo = new LinkGenerateRequestInfoModel
                {
                    MerchantID = 175001,
                    RandomID = Guid.NewGuid().ToString(),
                    MemberCode = "marshall001",
                    Lang = "en-us",
                    RequestID = RandomHelper.GetOrderNo("TT"),
                    TS = DateTimeHelper.GetCurrentTimeTimestamp(),
                    NotifyUrl = $"{targetDomain}/FakeReponse/FakeResponseOK",
                    ReturnUrl = $"http://{Request.Host}",
                    Amount = 150,
                }
            };

            var requestParam = new RequestParameters
            {
                BaseUrl = SettingHelper.WebApplicationSetting.F2E_Vendor_API_Doamin,
                ApiUrl = "card/link-request",
                Method = HttpMethodType.HttpPost
            };

            var res = APIRequestHelper.RequestApiResult<ApiResponseModel<LinkGenerateViewModel>>(requestParam, data);


            if(res.Exception != null || res.ResultObject.Code >= 400000)
            {
                var msg = "(card/link-request) " +  res.ResponseText;
                throw new MSCException(ErrorCode.InternalServerError, msg);
            }
            else
            {
                //請注意這邊的導頁網址是資料庫設定的 MerchantCustomizedSetting
                string url = res.ResultObject.Data.Link;
                string replaceText = $"http://{Request.Host}";

                if (isLocalhost)
                {
                    url = url.Replace("http://0.0.0.0:11449", replaceText);
                    url = url.Replace("https://example.invalid" ,replaceText);
                }

                return Redirect(url);
            }
          

        }




    }
}
