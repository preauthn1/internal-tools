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
using Microsoft.AspNetCore.Mvc.Controllers;

namespace target.Web.Modules.Filter
{
   
    /// <summary>
    /// 部分API站台 要求任何請求都要帶某個特定header (證書)
    /// </summary>
    public class ApiCertificateAuthorizeFilter : Attribute , IAuthorizationFilter
    {
        
        public void OnAuthorization(AuthorizationFilterContext context)
        {
            var apiUrl = context.HttpContext.Request.Path.ToString() ?? string.Empty;

            if (apiUrl.Contains("swagger"))
                return;

            var controller = context.ActionDescriptor as ControllerActionDescriptor;

            if(controller != null)
            {
                var controllerAttrs = controller.ControllerTypeInfo.GetCustomAttributes(false);

                foreach(var item in controllerAttrs)
                {
                    //如果controller有標記attribute 不檢查 ApiCertificate
                    var myAttr = item as DoNotCheckCertificate;

                    if (myAttr != null)
                    {
                        return;
                    }
                }
            }

            var headerValue = HeaderValueHelper.GetApiCertificate();

            if (string.IsNullOrWhiteSpace(headerValue))
                throw new TerproException(ApiErrorCode.Forbidden_RequiredApiCertificate);

            var configSetting = SettingHelper.GameApiSetting.ApiCertificateValue;

            if(configSetting != headerValue)
                throw new TerproException(ApiErrorCode.Forbidden_RequiredApiCertificate);
        }

    }

}
