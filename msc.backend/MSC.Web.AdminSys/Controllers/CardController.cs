using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.Web.AdminSys.Models;

namespace MSC.Web.AdminSys.Controllers
{
    [AdminPagePermissionAuthorizeFilter]
    public class CardController : Controller
    {        
        public CardController()
        {

        }

        public IActionResult List()
        {
            return View();
        }

        public IActionResult GenerateCardNo()
        {
            return View();
        }

    }
}
