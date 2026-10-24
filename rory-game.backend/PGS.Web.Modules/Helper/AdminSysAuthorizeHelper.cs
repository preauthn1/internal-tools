using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;
using target.Utility.Extensions;
using target.Utility.Helper.Dependency;
using target.Models.Table;
using target.Utility;
using target.Utility.IOC;
using target.BL.Interface;
using target.Models.Models;
using target.Models;
using target.Utility.Helper;
using target.DA.Repository;
using target.Models.Enums;

namespace target.Web.Modules.Filter
{
    public class AdminSysAuthorizeHelper
    {

        public static AdminUserTokenJwtModel GetValidUserSessionByJwt()
        {
            string jwtToken = [REDACTED];

            //jwt檢查包含 檢查偽造、過期、解析失敗
            if (string.IsNullOrWhiteSpace(jwtToken))
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            AdminUserTokenJwtModel convertUserModel = MyJwtHelper.GetAdminUserModelFromJwtString(jwtToken);

            if (convertUserModel == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            var now = DateTimeHelper.GetCurrentTimeTimestamp();

            if (convertUserModel.TokenExpiredTime <= now)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            string jwtHash = jwtToken.MD5Hash();

            var rep = new UserRepository();

            var queryUser = rep.GetSingleByCache(convertUserModel.LoginUserID);

            if(queryUser == null || queryUser.Status != NormalActiveStatus.Enable.ToInteger())
            {
                throw new TerproException(ApiErrorCode.InvalidStatus_AccountStatus);
            }

            return convertUserModel;
        }

    }
}
