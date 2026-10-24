using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using example.invalid.Http;
using System.Collections.Generic;
using System.Linq;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Models.Models;
using target.Models;
using System;
using target.Utility.Helper.Dependency;

namespace target.Web.Modules.Filter
{

    public class SignKeyVerifyFilter : Attribute , IActionFilter
    {

        public void OnActionExecuting(ActionExecutingContext context)
        {
            var signKey = HeaderValueHelper.GetSignKey();

            if (string.IsNullOrWhiteSpace(signKey))
                throw new TerproException(ApiErrorCode.InvalidParameter_SignKeyError);

            //神奇萬用signKey
            if (signKey == SystemConstant.PublicAPI_SignKey_MagicPassword)
                return;

            var method = context.HttpContext.Request.Method;

            bool isGET = method == "GET";

            if (isGET)
            {              
                string expect_signKey = SignKeyRuleHelper.GenerateF2E_HttpGet_SignKey(context.HttpContext.Request);

                if (signKey != expect_signKey)
                {
                    throw new TerproException(ApiErrorCode.InvalidParameter_SignKeyError);
                }
            }
        }

        public void OnActionExecuted(ActionExecutedContext context)
        {

        }
    }
}
