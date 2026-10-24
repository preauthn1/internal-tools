using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Lck.Utility.Extensions.Helper;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;

namespace MSC.API.AdminSys.Controllers
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


        [Route("internal/cache-remove-ipForbidden")]
        [HttpGet, HttpPost]
        public object RemoveIPAddressForbiddenCache()
        {
            CacheHelper.RemoveKeywordCache(RequestForbiddenTypeEnum.Request_CardInput_CheckCardNumber_IPAddress.ToString());
            CacheHelper.RemoveKeywordCache(RequestForbiddenTypeEnum.Request_CardInput_CheckCardNumber_TokenHash.ToString());
            return "OK";
        }

    }
}
