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

namespace target.Web.Modules.Filter
{
   
    public class DisabledWhenProduction : Attribute , IAuthorizationFilter
    {
        public void OnAuthorization(AuthorizationFilterContext context)
        {
            if (SettingHelper.IsProduction)
                throw new TerproException(ApiErrorCode.Forbidden_DoNotAllowWhenProduction);

            return;
        }

    }

}
