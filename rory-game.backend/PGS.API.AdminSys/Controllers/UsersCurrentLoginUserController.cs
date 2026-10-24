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

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class UsersCurrentLoginUserController : BaseController
    {
        public IUsersService UsersService { get; set; }

        public IUserInfoService UserInfoService { get; set; }

        public IViewPermissionCheckService ViewPermissionCheckService { get; set; }



        [Route("users/profile")]
        [HttpGet]
        public UserInfoViewModel GetUserProfile()
        {
            int userID = GetCurrentUserID();
            var query = UserInfoService.GetUserProfile(userID);
            return query;
        }


        [Route("users/language")]
        [HttpPut]
        public void UpdateUserLanguage([FromBody]UpdateCurrentUserLangdRequestModel model)
        {
            int userID = GetCurrentUserID();
            UserInfoService.UpdateUserLanguage(userID , model?.Lang);
        }


        [Route("users/password")]
        [HttpPut]
        public void UpdateUserPassword([FromBody]UpdateCurrentUserPasswordRequestModel model)
        {
            int userID = GetCurrentUserID();
            UserInfoService.UpdateUserPassword(userID, model);
        }

    }
}
