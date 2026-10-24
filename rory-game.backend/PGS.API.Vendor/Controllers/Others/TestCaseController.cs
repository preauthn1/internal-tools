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

namespace target.API.Vendor.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    public class TestCaseController : BaseController
    {
      
        [Route("testcase/TestSetRedisValue")]
        [HttpGet]
        public object TestSetRedisValue(string cacheKey = "Test")
        {
            //DbEntity_VendorSessionToken_5736cdc6-1f44-4e5e-bb72-2029e6c11040

            string cacheKeySample = CacheKey.DbEntity_VendorSessionToken("5736cdc6-1f44-4e5e-bb72-2029e6c11040");
            CacheHelper.SetCacheValue(cacheKey, "XXXXXXXXXX", new TimeSpan(12, 0, 0));
            return "";
        }

    }
}
