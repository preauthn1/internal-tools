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
    public class LoginController : BaseController
    {
        public IUserLoginService UserLoginService { get; set; }


        [Route("auth/login")]
        [HttpPost]
        public UserLoginViewModel Login([FromBody]UserLoginRequestModel model)
        {
            var result = UserLoginService.Login(model);
            return result;
        }

        /// <summary>
        /// 強制登出所有帳號
        /// </summary>
        [Route("auth/force-logout-all")]
        [HttpPost]
        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_InternalUtility)]
        public void ForceLogoutAllUsers()
        {
            UserLoginService.ForceLogoutAllUsers();
            return;
        }


    }
}
