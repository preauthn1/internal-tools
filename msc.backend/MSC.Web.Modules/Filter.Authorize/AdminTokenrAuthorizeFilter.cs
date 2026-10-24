using Microsoft.AspNetCore.Mvc.Filters;
using MSC.BL.Interface;
using MSC.Models.Models;
using MSC.Utility.Helper;
using MSC.Utility.Helper.HttpContext;
using MSC.Utility.IOC;
using System;
using System.Collections.Generic;

namespace MSC.Web.Modules.Filter
{

    public class AdminSysTokenrAuthorizeFilter : Attribute, IAuthorizationFilter
    {
        public void OnAuthorization(AuthorizationFilterContext context)
        {
            var jwtToken = [REDACTED];

            var service = IOCHelper.Resolve<IAdminSysTokenService>();

            var resultTokenModel = service.GetAndVerifyTokenIsValid(jwtToken);
        }

    }


}
