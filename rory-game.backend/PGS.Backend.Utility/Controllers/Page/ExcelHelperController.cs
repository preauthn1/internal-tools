using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using target.Backend.Utility.Models;

namespace target.Backend.Utility.Controllers
{
    public class ExcelHelperController : Controller
    {
        private readonly ILogger<HomeController> _logger;


        public IActionResult ToNestedJson()
        {
            return View();
        }

    }
}
