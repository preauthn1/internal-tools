using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.Base;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Utility.Utility;
using MSC.Web.Point.Models;
using MSC.Utility.Extensions;
using MSC.Web.Point.Helper;

namespace MSC.Web.Point.Controllers
{
    public class WidgetController : BaseController
    {
        private readonly ILogger<HomeController> _logger;
        

        public WidgetController()
        {

        }


        [Route("Shoppings")]
        public IActionResult ShoppingStore()
        {
            ViewBag.ShoppingStore = SettingDataHelper.GetShoppingStoreByCache();
            return View();
        }

       

    }
}
