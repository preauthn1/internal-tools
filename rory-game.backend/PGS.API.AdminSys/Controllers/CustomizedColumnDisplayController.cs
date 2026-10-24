using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Web.Modules.Filter;
using target.Models.Paging;
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class CustomizedColumnDisplayController : BaseController
    {
        public ICustomizedColumnDisplayService CustomizedColumnDisplayService { get; set; }



        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagementRecord)]
        [Route("customized-columns/dynamicQuery/setting/{querySqlID}")]
        [HttpGet]
        public DynamicQueryCommandViewModel GetDynamicQueryCommandSetting(int querySqlID)
        {
            int userID = GetCurrentUserID();
            var query = CustomizedColumnDisplayService.GetDynamicQueryCommandSetting(querySqlID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagementRecord)]
        [Route("customized-columns")]
        [HttpGet]
        public CustomizedColumnDynamicQueryResultViewModel GetDynamicQueryResult([FromQuery] CustomizedColumnDisplayRequestModel filter)
        {
            int userID = GetCurrentUserID();
            var query = CustomizedColumnDisplayService.GetDynamicQueryResult(userID, filter);
            return query;
        }

    }
}
