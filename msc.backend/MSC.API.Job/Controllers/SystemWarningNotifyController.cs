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

namespace MSC.API.Job.Controllers
{

    public class SystemWarningNotifyController : BaseController
    {       
        public IJobRoutineTaskService JobRoutineTaskService { get; set; }

        [ApiConfig(IsHeaderLogEnable = false)]
        [Route("routine/system-warning")]
        [HttpGet]
        public void CheckSystemWarningAndSendNotify()
        {
            JobRoutineTaskService.CheckSystemWarningAndSendNotify();
            return;
        }
    }
}
