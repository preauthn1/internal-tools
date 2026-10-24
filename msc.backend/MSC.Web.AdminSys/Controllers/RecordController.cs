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
    public class RecordController : Controller
    {
       

        public IActionResult CardConsumeRecordList()
        {
            return View();
        }


        public IActionResult CardConsumeMerchantNotifyQueueList()
        {
            return View();
        }

        public IActionResult CardLinkRequestRecordList()
        {
            return View();
        }

        public IActionResult RequestForbiddenLogList()
        {
            return View();
        }

        public IActionResult UserOperateLogs()
        {
            return View();
        }

        public IActionResult MerchantWalletLogs()
        {
            return View();
        }

        
    }
}
