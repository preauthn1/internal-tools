using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using MSC.GoogleApi;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using System.Reflection;
using MSC.Models.Enums;
using MSC.Models;
using MSC.Utility.Extensions;
using MSC.Models.Models;
using MSC.Models.CultureText;

namespace target.Backend.Utility.Controllers
{
    public class FakeReponseController : Controller
    {
        
        [HttpPost]
        public ActionResult FakeResponseOK()
        {
            return Content("OK , (this is fake API)");
        }


    }
}
