using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Web.Modules.Filter;
using example.invalid;
using System.Web;
using MSC.Models.API.Middleware;
using Newtonsoft.Json.Linq;

namespace MSC.API.Job.Controllers
{

    public class AzureWebHookController : BaseController
    {       
        public IWebHookNotifyService WebHookNotifyService { get; set; }

        [ApiConfig(IsHeaderLogEnable = false)]
        [Route("webhook/azure-sql-notify")]
        [HttpPost]
        public void AzureNotifySQLWarning([FromBody]JObject model)
        {
            WebHookNotifyService.AzureNotifySQLWarning(model);
            return;
        }
    }
}
