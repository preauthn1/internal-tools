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
using target.Utility.Extensions;
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{

    [AdminTokenrAuthorizeFilter]
    public class PermissionController : BaseController
    {
        public IPermissionService PermissionService { get; set; }


        [ApiConfig(IsResponseLogEnable = false)]
        [Route("users/permissions")]
        [HttpGet]
        public UserPermissionViewModel GetUserAvailablePermissions()
        {
            int userID = GetCurrentUserID();
            var query = PermissionService.GetAvailablePermissions(userID);
            return query;
        }


        [Route("permissions/nested")]
        [HttpGet]
        public UserPermissionViewModel GetPermissionsNestedFormat()
        {
            int userID = GetCurrentUserID();
            var query = PermissionService.GetPermissionsNestedFormat(userID);
            return query;
        }
        

        [Route("permissions")]
        [HttpGet]
        public List<PermissionViewModel> GetAvailablePermissions()
        {
            int userID = GetCurrentUserID();
            var query = PermissionService.GetPermissions(userID);
            return query;
        }

    }
}
