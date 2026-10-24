using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.API.Game.Controllers;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility.Helper;
using MSC.Web.Modules.Filter;

namespace MSC.API.Point.Controllers
{
    public class ConfigSettingController : BaseController
    {       
        public IConfigSettingService ConfigSettingService { get; set; }


        /// <summary>
        /// 取得合作商家
        /// </summary>
        [Route("settings/shopping-store")]
        [HttpGet]
        public List<ShoppingStoreViewModel> GetCardShoppingStores()
        {
            var result = ConfigSettingService.GetCardShoppingStores();
            return result;
        }

        /// <summary>
        /// 取得合作商家iframe
        /// </summary>
        [Route("settings/shopping-iframe-url")]
        [HttpGet]
        public string GetShoppingsIframeUrl()
        {
            var result = ConfigSettingService.GetShoppingsIframeUrl();
            return result;
        }


        /// <summary>
        /// 取得 輸入卡號頁的配置
        /// </summary>
        [Route("settings/card-page/{merchantID}")]
        [HttpGet]
        public CardInputPageSettingViewModel GetCardInputPageSetting(int merchantID)
        {
            var result = ConfigSettingService.GetCardInputPageSetting(merchantID);
            return result;
        }


    }
}
