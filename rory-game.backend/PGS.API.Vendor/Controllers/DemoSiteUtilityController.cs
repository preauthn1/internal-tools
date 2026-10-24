using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.API.Vendor.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    [Route("v{version:int}")]
    public class DemoSiteUtilityController : BaseController
    {
        public IInternalUtilityService InternalUtilityService { get; set; }


        /// <summary>
        /// 這隻API 基本上是for demo站使用
        /// </summary>
        [Route("internal-api/merchant/{merchantID}/servers")]
        [HttpGet]
        public List<DemoGameCategoryViewModels> GetServers(int merchantID)
        {
            var query = InternalUtilityService.GetMerchantServers(merchantID);
            return query;
        }

    }
}
