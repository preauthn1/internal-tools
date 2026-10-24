using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.Models.Models;
using MSC.Web.AdminSys.Models;

namespace MSC.Web.AdminSys.Controllers
{
    public class HomeController : Controller
    {
        private readonly ILogger<HomeController> _logger;

        public HomeController(ILogger<HomeController> logger)
        {
            _logger = logger;
        }


        [AdminTokenrAuthorizeFilter]
        public IActionResult Index()
        {
            return View();
        }


        [AdminTokenrAuthorizeFilter]
        [Route("error")]
        public IActionResult Error(string msg = "" , int code = 0)
        {
            ViewBag.Msg = msg;
            ViewBag.Code = code;

            return View();
        }


        [AdminTokenrAuthorizeFilter]
        public IActionResult NoPermission()
        {
            return View();
        }

        public IActionResult NotImplement()
        {
            throw new MSCException(ErrorCode.Authorize_AccessToken_Expired);
            return View();
        }


    }
}
