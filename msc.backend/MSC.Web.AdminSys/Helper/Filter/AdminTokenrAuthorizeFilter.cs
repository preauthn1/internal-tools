using Microsoft.AspNetCore.Mvc.Filters;
using MSC.Models.Models;
using System;
using System.Linq;
using System.Collections.Generic;

namespace MSC.Web.AdminSys.Models
{

    public class AdminTokenrAuthorizeFilter : Attribute, IAuthorizationFilter
    {
        public void OnAuthorization(AuthorizationFilterContext context)
        {
            string url = context.HttpContext.Request.Path;

            bool isNotLogin = AuthTokenHelper.IsNotLogin();

            if (isNotLogin)
                throw new MSCException(ErrorCode.Authorize_PleaseLogin);

            return;
        }
    }


    public class AdminPagePermissionAuthorizeFilter : Attribute, IAuthorizationFilter
    {
        public void OnAuthorization(AuthorizationFilterContext context)
        {
            bool isNotLogin = AuthTokenHelper.IsNotLogin();

            if (isNotLogin)
                throw new MSCException(ErrorCode.Authorize_PleaseLogin);

            var allPermissions = PermissionHelper.GetMyPermissions();

            string url = context.HttpContext.Request.Path;

            bool isExist = allPermissions.SelectMany(z => z.Childs).Any(z => z.Url == url);

            if (!isExist)
                throw new MSCException(ErrorCode.Forbidden_PagePermissionNotAllow);

            return;
        }
    }

}
