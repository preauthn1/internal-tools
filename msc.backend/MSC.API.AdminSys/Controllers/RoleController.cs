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
    public class RoleController : BaseController
    {

        public IRoleService RoleService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_User_List)]
        [Route("role/{roleID}/permissions")]
        [HttpGet]
        public RolePermissionViewModel GetByID(int roleID)
        {
            int userID = GetCurrentUserID();
            var query = RoleService.GetRolePermissions(roleID);
            return query;
        }

    }
}
