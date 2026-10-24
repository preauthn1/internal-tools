using MSC.Models.ViewModels.Base;
using MSC.Utility.Helper;
using MSC.Utility.Helper.HttpContext;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using Lck.Utility.Extensions.Helper;

namespace MSC.Web.AdminSys.Models
{
    public class AuthTokenHelper
    {
        public static string NOT_LOGIN_NAME = "Unknow";


        public static bool IsNotLogin()
        {
            var user = GetUserModelByCookieToken();
            return user.LoginName == NOT_LOGIN_NAME;
        }


        public static AdminUserTokenJwtModel GetUserModelByCookieToken()
        {
            string jwtToken = [REDACTED]"token"];

            var jwtModel = MyJwtHelper.GetAdminUserModelFromJwtString(jwtToken);

            if(jwtModel == null)
            {
                return new AdminUserTokenJwtModel
                {
                    LoginName = NOT_LOGIN_NAME,
                };
            }

            return jwtModel;
        }
    }


}
