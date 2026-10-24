using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using example.invalid.Http;
using System.Collections.Generic;
using System.Linq;
using Microsoft.AspNetCore.Mvc.Controllers;
using target.Utility.Extensions;
using target.Models;
using target.Utility.Helper;

namespace target.Web.Modules.Filter
{

    /// <summary>
    /// 部分API站台 要求任何請求都要帶某個特定header (證書)
    /// 因為 AuthorizationFilter 的順序更優先 所以寫成 AuthorizationFilter
    /// </summary>
    public class GlobalApiHeaderTemporaryInfoFilter : System.Attribute, IAuthorizationFilter
    {

        public void OnAuthorization(AuthorizationFilterContext context)
        {
            var actionDescriptor = context.ActionDescriptor as ControllerActionDescriptor;

            string uniqueRouteName = $"{actionDescriptor.ControllerName}_{actionDescriptor.ActionName}";

            //將此資訊存在httpHeader，寫log的時候，會用來判斷一些東西
            if (!context.HttpContext.Request.Headers.ContainsKey(SystemConstant.HeaderKey_Temperary_ApiUniqueRouteName))
            {
                context.HttpContext.Request.Headers.Add(SystemConstant.HeaderKey_Temperary_ApiUniqueRouteName, uniqueRouteName);
            }
        }

    }


}
