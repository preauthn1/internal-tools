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
   
    public class PlayerTokenAuthorizeFilter : Attribute , IAuthorizationFilter
    {
        public void OnAuthorization(AuthorizationFilterContext context)
        {
            string jwtToken = [REDACTED];

            //jwt檢查包含 檢查偽造、過期、解析失敗
            if (string.IsNullOrWhiteSpace(jwtToken))
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            PlayerTokenJwtModel convertUserModel = MyJwtHelper.GetPlayerModelFromJwtString(jwtToken);

            if(convertUserModel == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);            

            var now = DateTimeHelper.GetCurrentTimeTimestamp();

            if (convertUserModel.TokenExpiredTime <= now)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);


            string jwtHash = jwtToken.MD5Hash();

            var service = IOCHelper.Resolve<IPlayerJwtCheckService>();

            //反查DB 檢查token是否過期
            var isExist = service.CheckIsTokenExist(jwtHash);

            if(!isExist)
                throw new TerproException(ApiErrorCode.Authorize_TokenHasChanged);
        }

    }

}
