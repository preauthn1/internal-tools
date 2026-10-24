using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using target.Utility.Extensions;
using target.Models;
using target.Models.Models;
using target.Utility.Helper;
using System;
using target.Utility;
using target.Utility.IOC;
using target.BL.Interface;
using target.Models.Enums;
using System.Collections.Generic;

namespace target.Web.Modules.Filter
{
   
    public class AdminPermissionAuthorizeFilter : Attribute , IAuthorizationFilter
    {
        /// <summary>
        /// 多個token有其中一個即可
        /// </summary>
        public List<AdminSysPermissionKey> RequiredTokens = new List<AdminSysPermissionKey>();


        public AdminPermissionAuthorizeFilter(AdminSysPermissionKey permissionKey)
        {
            RequiredTokens.Add(permissionKey);
        }

        public AdminPermissionAuthorizeFilter(AdminSysPermissionKey permissionKey , AdminSysPermissionKey permissionKey2)
        {
            RequiredTokens.Add(permissionKey);
            RequiredTokens.Add(permissionKey2);
        }

        public AdminPermissionAuthorizeFilter(AdminSysPermissionKey permissionKey, AdminSysPermissionKey permissionKey2 , AdminSysPermissionKey permissionKey3)
        {
            RequiredTokens.Add(permissionKey);
            RequiredTokens.Add(permissionKey2);
            RequiredTokens.Add(permissionKey3);
        }


        public void OnAuthorization(AuthorizationFilterContext context)
        {
            var adminUser = AdminSysAuthorizeHelper.GetValidUserSessionByJwt();

            var service = IOCHelper.Resolve<IPermissionService>();

            //透過資料庫(cache) 反查是否權限符合
            bool isExist = service.IsPermissionExist(adminUser.LoginUserID, RequiredTokens);

            if (!isExist)
            {
                throw new TerproException(ApiErrorCode.Forbidden_PagePermissionNotAllow);
            }
        }

    }

}
