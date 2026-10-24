using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Collections.Generic;
using MSC.Models;
using MSC.Utility.Helper;
using MSC.Utility.Helper.HttpContext;
using MSC.Models.Models;
using MSC.Utility.IOC;
using MSC.BL.Interface;
using MSC.Models.Enums;

namespace MSC.Web.Modules.Filter
{
   
    public class AdminPermissionAuthorizeFilter : Attribute , IAuthorizationFilter
    {
        /// <summary>
        /// 多個token有其中一個即可 
        /// </summary>
        public List<string> RequiredTokens = new List<string>();



        public AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum permissionKey)
        {
            RequiredTokens.Add(permissionKey.ToString());
        }

        public AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum permissionKey , AdminPermissionTokenEnum permissionKey2)
        {
            RequiredTokens.Add(permissionKey.ToString());
            RequiredTokens.Add(permissionKey2.ToString());
        }


        public AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum permissionKey, AdminPermissionTokenEnum permissionKey2 , AdminPermissionTokenEnum permissionKey3)
        {
            RequiredTokens.Add(permissionKey.ToString());
            RequiredTokens.Add(permissionKey2.ToString());
            RequiredTokens.Add(permissionKey3.ToString());
        }


        public void OnAuthorization(AuthorizationFilterContext context)
        {
            string jwtToken = [REDACTED];
            var jwtModel = MyJwtHelper.GetAdminUserModelFromJwtString(jwtToken);

            if (jwtModel == null)
            {
                throw new MSCException(ErrorCode.Authorize_AccessToken_Expired);
            }

            var service = IOCHelper.Resolve<IAdminSysTokenService>();
            service.CheckPermissionExist(jwtModel.LoginUserID, RequiredTokens);
        }

    }

}
