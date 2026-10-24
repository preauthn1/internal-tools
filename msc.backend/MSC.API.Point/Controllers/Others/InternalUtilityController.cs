using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Lck.Utility.Extensions.Helper;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility.Helper;

namespace MSC.API.Game.Controllers
{
    public class InternalUtilityController : BaseController
    {       
      
        [Route("internal/cache-flush-all")]
        [HttpGet, HttpPost]
        public object Cache_FlushAll()
        {
            CacheHelper.FlushAll();
            return "OK";
        }

        [Route("internal/test-cache")]
        [HttpGet, HttpPost]
        public object TestGetCacheValue()
        {
            string cacheKey = "Debug_001825";
            CacheHelper.SetCacheValue(cacheKey, "當前時間為:" + DateTimeHelper.GetCurrentTime(8) , 10);
            return CacheHelper.GetCacheValue<string>(cacheKey);
        }

    }
}
