using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using example.invalid;
using System.Linq;
using System.Reflection;
using System.Runtime.Versioning;
using System.Threading;
using System.Threading.Tasks;
using Lck.Utility.Extensions.Helper;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using target.BL.Interface;
using target.Models.ViewModels;
using target.Models.ViewModels.Base;
using target.Utility;
using target.Utility.Helper;
using target.Utility.Helper.HttpContext;
using target.Utility.Extensions;
using target.Utility.Utility;
using target.Utility.Utility.ThirdParty;
using target.Web.Modules.Filter;
using target.Models;

namespace target.API.AdminSys.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    public class TestCaseController : BaseController
    {

        public IGameRecordService GameRecordService { get; set; }


        [Route("testCase/GenerateDetailInfoLink")]
        [HttpGet]
        public BetRecordDetailInfoLinkGenerateViewModel GenerateDetailInfoLink([FromQuery] BetRecordDetailInfoLinkGenerateRequestModel model)
        {
            var result = GameRecordService.GenerateDetailInfoLink(model);
            return result;
        }


        [Route("internal/cache/removeAll")]
        [HttpGet]
        public object RemoveAllCache()
        {
            CacheHelper.RemoveAll();
            return "OK";
        }

        [Route("testcase/headers")]
        [HttpGet]
        public object GetHeaders()
        {
            return MyHttpContext.Current.Request.Headers;
        }


        [Route("testcase/testVoid")]
        [HttpGet]
        public void testVoid()
        {

        }



        [Route("testcase/delay")]
        [HttpGet]
        public object DelayTest()
        {
            Thread.Sleep(3000);
            return "xxx";
        }


        [Route("testcase/001")]
        [HttpGet]
        public object Test001()
        {
            var qqq = CacheHelper.SetCacheValue("test_key_001", "xxxx", null);
            return CacheHelper.GetCacheValue<string>("test_key_001");
        }


        [Route("testcase/testDelayShutdown")]
        [HttpGet]
        public object testDelayShutdown(int delay = 10)
        {
            Task.Factory.StartNew(() =>
            {
                var startTime = DateTime.Now;

                Thread.Sleep(delay * 1000);

                string msg = $"testDelayShutdown_收到請求_{startTime.ToString("HH:mm:ss")}_Sleep結束_{ DateTime.Now.ToString("HH:mm:ss")}";

                MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
                {
                    Message = msg,
                    TelegramChannelEnum = TelegramChannelEnum.pgs_common_error
                });
            });

            return "xxxxxxx";
        }


        [Route("testcase/static_response")]
        [HttpGet]
        public object Test_StaticResponse()
        {            
            //無資料庫存取相關，檢查時間有點慢 是否是跟DB層有關
            return "hello";
        }


        [AdminTokenrAuthorizeFilter]
        [Route("testcase/static_response_loginToken")]
        [HttpGet]
        public object Test_StaticResponse_RequiredToken()
        {
            //要求要登入才能存取
            return "hello";
        }



        [Route("testcase/TestException")]
        [HttpGet]
        public object TestException()
        {
            throw new Exception("qqqqqqqq");
            var qqq1 = HeaderValueHelper.GetHDCode();
            return "";
        }



    }
}
