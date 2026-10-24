using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using example.invalid.Http;
using example.invalid.Http.Headers;
using Microsoft.AspNetCore.Mvc;
using target.Utility.Extensions;
using example.invalid;
using Lck.Utility.Extensions.Helper;

namespace target.API.ExternalVendor.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    public class InternalHelperController : ControllerBase
    {
        /// <summary>
        /// 清除所有Cache
        /// </summary>
        [Route("internal/cache/removeAll")]
        [HttpGet]
        public object RemoveAllCache()
        {
            CacheHelper.RemoveAll();
            return "Cache Remove All OK";
        }

    }

}
