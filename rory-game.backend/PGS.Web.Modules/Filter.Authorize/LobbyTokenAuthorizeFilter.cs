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
   
    public class LobbyTokenAuthorizeFilter : Attribute , IAuthorizationFilter
    {
        public void OnAuthorization(AuthorizationFilterContext context)
        {
            string tokenRawText = context.HttpContext.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_AccessToken);

            var service = IOCHelper.Resolve<IPlayerJwtCheckService>();

            //反查DB 檢查token是否過期
            service.CheckIsLobbyTokenExist(tokenRawText);
        }

    }

}
