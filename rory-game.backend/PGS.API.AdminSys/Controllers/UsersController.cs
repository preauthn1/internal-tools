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
    public class UsersController : BaseController
    {
        public IUsersService UsersService { get; set; }

        public IUserInfoService UserInfoService { get; set; }

        public IViewPermissionCheckService ViewPermissionCheckService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.OperatorManagementAccounts_Search)]
        [Route("users")]
        [HttpGet]
        public PagedList<UsersViewModel> GetUsersPagingByFilter([FromQuery]QueryUserApiRequestParameterModel filter)
        {
            var query = UsersService.GetUsersPagingByFilter(filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.OperatorManagementAccounts_Search)]
        [Route("users/{userID}")]
        [HttpGet]
        public UsersUpdatePageViewModel GetUserByID(int userID)
        {
            int currentUserID = GetCurrentUserID();
            ViewPermissionCheckService.CheckViewUserPermission(currentUserID , userID);
            var query = UsersService.GetUserByID(userID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.OperatorManagementAccounts_Create)]
        [ApiConfig(IsRequestBodyLogEnable = false)]
        [Route("users")]
        [HttpPost]
        public void Create(UserCreateModel model)
        {
            UsersService.Create(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.OperatorManagementAccounts_Update)]
        [ApiConfig(IsRequestBodyLogEnable = false)]
        [Route("users")]
        [HttpPut]
        public void Update(UserUpdateModel model)
        {
            ViewPermissionCheckService.CheckViewUserPermission(GetCurrentUserID(), model?.UserID ?? 0);
            UsersService.Update(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.OperatorManagementAccounts_Delete)]
        [Route("users/{userID}")]
        [HttpDelete]
        public void Delete(int userID)
        {
            ViewPermissionCheckService.CheckViewUserPermission(GetCurrentUserID(), userID);
            UsersService.Delete(userID);
        }

    }
}
