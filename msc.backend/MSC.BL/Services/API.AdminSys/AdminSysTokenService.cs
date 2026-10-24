using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.Utility;
using MSC.DA.Repository;
using MSC.Models.Table;
using System.Transactions;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.Base;

namespace MSC.BL.Services
{

    public class AdminSysTokenService : BaseAdminService, IAdminSysTokenService
    {

        public RolePermissionRepository RolePermissionRepository { get; set; }

        public UserTokenRepository UserTokenRepository { get; set; }


        public AdminUserTokenJwtModel GetAndVerifyTokenIsValid(string jwtToken)
        {
            if (jwtToken.IsNullOrWhiteSpace())
            {
                throw new MSCException(ErrorCode.Authorize_AccessToken_NotValid);
            }

            var jwtModel = MyJwtHelper.GetAdminUserModelFromJwtString(jwtToken);

            if (jwtModel == null)
            {
                throw new MSCException(ErrorCode.Authorize_AccessToken_NotValid);
            }

            var nowTs = DateTimeHelper.GetCurrentTimeTimestamp();

            if (nowTs > jwtModel.TokenExpiredTime)
            {
                throw new MSCException(ErrorCode.Authorize_AccessToken_Expired);
            }

            string jwtHash = jwtToken.MD5Hash();

            var query = UserTokenRepository.GetCacheDataByHashToken(jwtHash);

            if(query == null)
            {
                throw new MSCException(ErrorCode.Authorize_AccessToken_Expired);
            }

            if (DateTimeHelper.GetCurrentTime() > query.ExpiredTime)
            {
                throw new MSCException(ErrorCode.Authorize_AccessToken_Expired);
            }

            return jwtModel;
        }



        public void CheckPermissionExist(int userID, List<string> requiredTokens)
        {
            var query = UserRepository.GetSingleByCache(userID);

            if (query == null)
                throw new MSCException(ErrorCode.Forbidden_WatchDataPermissionNotAllow);

            var permissions = RolePermissionRepository.GetRolePermissionByCache(query.RoleID);

            bool isAny = permissions.Select(z => z.PermissionKey).Any(x => requiredTokens.Contains(x));

            if (!isAny)
                throw new MSCException(ErrorCode.Forbidden_WatchDataPermissionNotAllow);

        }


    }
}
