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

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class UserController : BaseController
    {       
        public IUserService UserService { get; set; }

        public IPermissionService PermissionService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_User_List)]
        [Route("users")]
        [HttpGet]
        public PagedList<UserViewModel> GetUsersByPagingFilter([FromQuery] QueryUserApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = UserService.GetByPagingFilter(userID, filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_User_List)]
        [Route("users")]
        [HttpPost]
        public void CreateUser([FromBody] UserCreateRequestModel filter)
        {
            int userID = GetCurrentUserID();
            UserService.Create(userID, filter);
            return;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_User_List)]
        [Route("users")]
        [HttpPut]
        public void UpdateUser([FromBody] UserCreateUpdateModel filter)
        {
            int userID = GetCurrentUserID();
            UserService.Update(userID, filter);
            return;
        }

        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_User_List)]
        [Route("users/password")]
        [HttpPut]
        public void ChangePassword([FromBody] UserChangePasswordModel model)
        {
            int currentUserID = GetCurrentUserID();
            UserService.ChangePassword(currentUserID, model);
            return;
        }


        [Route("users/self-password")]
        [HttpPut]
        public void ChangeSelfPassword([FromBody] UserChangeSelfPasswordModel model)
        {
            int userID = GetCurrentUserID();
            UserService.ChangeSelfPassword(userID, model);
            return;
        }


        /// <summary>
        /// 內部小工具 將所有權限塞給 admin role
        /// </summary>
        [Route("users/role/auto-assign-permissions")]
        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_InternalUtility)]
        [HttpPost]
        public void AssignAllPermissions()
        {
            int userID = GetCurrentUserID();
            int targetRoleID = 1;
            PermissionService.AssignAllPermissions(targetRoleID);
            return;
        }
    }
}
