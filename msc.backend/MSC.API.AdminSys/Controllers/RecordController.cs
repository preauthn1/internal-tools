using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using MSC.BL.Interface;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Web.Modules.Filter;
using MSC.Models.Paging;
using MSC.Models.Enums;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class RecordController : BaseController
    {
    
        public IExternalApiRequestLogService ExternalApiRequestLogService { get; set; }

        public IRequestForbiddenLogService RequestForbiddenLogService { get; set; }

        public IUserOperateLogService UserOperateLogService { get; set; }


    
        [Route("records/system-forbidden")]
        [HttpGet]
        public PagedList<RequestForbiddenLogViewModel> GetForbiddenLogsByPagingFilter([FromQuery]QueryRequestForbiddenLogApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = RequestForbiddenLogService.GetByPagingFilter(userID, filter);
            return query;
        }

        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Record_UserOperateLogs_List)]
        [Route("records/operation-logs")]
        [HttpGet]
        public PagedList<UserOperateLogViewModel> GetUserOperateLogs([FromQuery]QueryUserOperateLogApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = UserOperateLogService.GetByPagingFilter(userID, filter);
            return query;
        }


        [Route("records/external-request-logs")]
        [HttpGet]
        public PagedList<ExternalApiRequestLogViewModel> GetExternalApiLogByPagingFilter([FromQuery]QueryExternalApiRequestLogApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = ExternalApiRequestLogService.GetByPagingFilter(userID, filter);
            return query;
        }

    }
}
