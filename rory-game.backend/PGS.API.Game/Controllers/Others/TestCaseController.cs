using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.Web.Modules.Filter;
using target.Utility.Helper.HttpContext;
using Lck.Utility.Extensions.Helper;
using target.Utility;

namespace target.API.Game.Controllers
{
    [DoNotCheckCertificate]
    [ApiExplorerSettings(IgnoreApi = true)]
    public class TestCaseController : BaseController
    {

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


        [Route("testcase/cache")]
        [HttpGet]
        public object testCache()
        {
            CacheHelper.SetCacheValue("API_Game", "API_Game", null);
            return CacheHelper.GetCacheValue<string>("API_Game");
        }



        [Route("appSetting/test")]
        [HttpGet]
        public object TestAppSetting()
        {
            return new
            {
                ApiRequestLogTableName = SettingHelper.ApiRequestLogTableName,
                WebSiteName = SettingHelper.WebSiteName,
                Env = SettingHelper.Env
            };
        }




        [Route("ex")]
        [HttpGet]
        public object TestExceptionn222222222()
        {
            throw new System.Exception("XXXXXXXXXXXX");
            return "";
        }


        [Route("CreateUsers")]
        [HttpPost]
        public PostMemberRequestModel CreateUsers([FromBody]PostMemberRequestModel model)
        {
            return model;
        }

        [Route("CreateUsersXXXXX")]
        [HttpPost]
        public string CreateUsersXXXXX(DateTime startTime)
        {
            return "XXXXXXXXX";
        }


        [Route("testcaseXXXXXXXXXX")]
        [HttpGet]
        public string testcaseXXXXXXXXXX()
        {
            return "XXXXXXXXX";
        }


        [Route("testcase/TestStringResult2")]
        [HttpGet]
        public string TestStringResult()
        {
            return "XXXXXXXXX";
        }


        [HttpGet("testcase/TestVoidResult")]
        public void TestVoidResult()
        {
            return;
        }

        [Route("exception2")]
        [HttpPost, HttpGet]
        public void Exception()
        {
            throw new Exception("怪怪的");
        }


        [Route("exception2/yyyyyyyy")]
        [HttpPost, HttpGet]
        public string yyyyyyyyHHHHHHHHH()
        {
            return "yyyyyyyy";
        }


        [Route("exception2/{id}")]
        [HttpPost, HttpGet]
        public object GetByID(string id)
        {
            return "XXXXXXXX" + id;
        }

    }
}
