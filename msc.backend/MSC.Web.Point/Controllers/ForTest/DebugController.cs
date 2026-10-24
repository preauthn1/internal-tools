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
using MSC.Utility;

namespace target.Backend.Utility.Controllers
{
    public class DebugController : Controller
    {


        [HttpGet]
        public ActionResult YYYY(string filterClassName = "")
        {

            string now = DateTime.Now.ToString("yyyyMMddHHmmss");

            
            var myKeyValues = new
            {
                qqq = "2321",
                qqqqq = "2312",
                text_is_va = "132123"
            };

            var json = myKeyValues.ConvertToJsonString();

            var final = $" $t =  {json}";


            return Content(final);

        }


        [HttpGet]
        public ActionResult TestException()
        {
            throw new MSCException(ErrorCode.BusinessLogicDeny_CurrencyNotMatch);
        }


    }
}
