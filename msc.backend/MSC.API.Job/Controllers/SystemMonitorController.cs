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

    public class SystemMonitorController : BaseController
    {       
        public IJobSystemMonitorService JobSystemMonitorService { get; set; }

        [ApiConfig(IsHeaderLogEnable = false)]
        [Route("routine/monitor/point-api-errors")]
        [HttpGet]
        public void DetectPointApiRequestUnusual()
        {
            JobSystemMonitorService.DetectPointApiRequestUnusual();
            return;
        }


        [ApiConfig(IsHeaderLogEnable = false)]
        [Route("routine/monitor/forbidden-logs")]
        [HttpGet]
        public void DetectForbiddenRequestRecord()
        {
            JobSystemMonitorService.DetectForbiddenRequestRecord();
            return;
        }

    }
}
