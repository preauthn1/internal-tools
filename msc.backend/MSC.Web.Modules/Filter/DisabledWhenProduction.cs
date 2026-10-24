using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using MSC.Models.ViewModels.Base;
using MSC.Utility.Extensions;
using MSC.Models;
using MSC.Models.Models;
using MSC.Utility.Helper;
using System;
using MSC.Utility;
using MSC.Utility.IOC;
using MSC.BL.Interface;

namespace MSC.Web.Modules.Filter
{
   
    public class DisabledWhenProduction : Attribute , IAuthorizationFilter
    {
        public void OnAuthorization(AuthorizationFilterContext context)
        {
            if (SettingHelper.IsProduction)
                throw new MSCException(ErrorCode.Forbidden_DoNotAllowWhenProduction);

            return;
        }

    }

}
