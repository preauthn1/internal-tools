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
using target.DA.Repository;
using target.Utility.Helper;

namespace target.API.Public.Controllers
{
    [DoNotCheckCertificate]
    [ApiExplorerSettings(IgnoreApi = true)]
    public class TestCaseController : BaseController
    {


        /// <summary>
        /// 測試redis是否有寫入
        /// </summary>
        [Route("internal/cache/testCacheInsert")]
        [HttpGet]
        public object testCacheInsert()
        {
            string testKey = SettingHelper.WebSiteName + "XXXX_YYYY_TTTT";
            CacheHelper.SetCacheValue(testKey, "XXXXXXXX_" + DateTimeHelper.GetCurrentTimeTimestamp().ToString() , new TimeSpan(0, 5, 0));

            var myValue = CacheHelper.GetCacheValue<string>(testKey);
            return myValue;
        }


    }
}
