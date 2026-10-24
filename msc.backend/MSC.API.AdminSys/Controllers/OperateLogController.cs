using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.Enums;
using MSC.Models.ViewModels;
using MSC.Web.Modules.Filter;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class OperateLogController : BaseController
    {

        public IUserOperateLogService UserOperateLogService { get; set; }


        [Route("operate-logs")]
        [HttpPost]
        public void WriteOperateLog([FromBody] UserOperateLogCreateRequestModel model)
        {
            int userID = GetCurrentUserID();
            UserOperateLogService.WriteOperateLog(userID , model);
            return;
        }

    }
}
