using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models.ViewModels;
using MSC.Web.AdminSys.Models;

namespace MSC.Web.AdminSys.Controllers
{
    public class LoginController : Controller
    {
       
        public IActionResult Index()
        {
            return View();
        }

    }
}
