using MSC.Models;
using MSC.Models.Models;
using MSC.Models.ViewModels.Base;
using System;
using System.Linq;
using System.Collections.Generic;
using System.Text;
using MSC.Utility.Extensions;

namespace MSC.Utility.Helper
{
    public class MyJwtHelper
    {
        

        public static string GenerateAdminUserJwt(int userID, string loginName , int orgID , bool isSuperAdmin, List<int> merchantIDs , HashSet<string> permissions = null)
        {
        
            var jwtModel = new AdminUserTokenJwtModel
            {
                TokenVersion = SystemConstant.JwtTokenVersionAdminSys,
                IsSuperAdminOrg = isSuperAdmin,
                MerchantIDs = merchantIDs,
                LoginUserID = userID,
                LoginName = loginName,
                OrganizationID = orgID,
                RndID = Guid.NewGuid().ToString(),
                Permissions = permissions,
                TokenExpiredTime = DateTimeHelper.GetCurrentTimeTimestamp() + GetTokenExpiredSeconds()
            };

            jwtModel.SignKey = GenerateSignKey(userID, jwtModel.TokenExpiredTime, jwtModel.RndID);
            jwtModel.SignKeyPermission = GeneratePermissionSignKey(permissions);


            var jwt = BaseJWTTokenHelper.Generate_JWT_Token(jwtModel);
            return jwt;
        }



        public static string GetUserIdentityTextFromJwt(string jwtTokenText)
        {
            if (string.IsNullOrWhiteSpace(jwtTokenText))
                return null;

            var jwtModel = BaseJWTTokenHelper.GetObjectFromJWT<BaseJwtModel>(jwtTokenText);

            if (jwtModel == null)
                return string.Empty;

            CheckSignIsValid(jwtModel.SignKey, jwtModel , new List<int> 
            { 
                SystemConstant.JwtTokenVersionPlayer, 
                SystemConstant.JwtTokenVersionAdminSys 
            });

            return jwtModel.LoginUserID.ToString();
        }


        public static PlayerTokenJwtModel GetPlayerModelFromJwtString(string jwtTokenText)
        {
            if (string.IsNullOrWhiteSpace(jwtTokenText))
                return null;

            var jwtModel = BaseJWTTokenHelper.GetObjectFromJWT<PlayerTokenJwtModel>(jwtTokenText);

            if (jwtModel == null)
                return null;

            CheckSignIsValid(jwtModel.SignKey , jwtModel , new List<int>
            {
                SystemConstant.JwtTokenVersionPlayer
            });

            return jwtModel;
        }
        
     
        public static AdminUserTokenJwtModel GetAdminUserModelFromJwtString(string jwtTokenText)
        {
            if (string.IsNullOrWhiteSpace(jwtTokenText))
                return null;

            var jwtModel = BaseJWTTokenHelper.GetObjectFromJWT<AdminUserTokenJwtModel>(jwtTokenText);

            if (jwtModel == null)
                return null;

            CheckSignIsValid(jwtModel.SignKey, jwtModel , new List<int>
            {
                SystemConstant.JwtTokenVersionAdminSys,
            });

            CheckSignPermissionIsValid(jwtModel.Permissions , jwtModel.SignKeyPermission);
            return jwtModel;
        }


        private static void CheckSignIsValid<T>(string originalSign, T model , List<int> availableVersions) where T : ITokenJwtModel
        {
            string signKeyAgain = GenerateSignKey(model.LoginUserID, model.TokenExpiredTime, model.RndID);

            //加密簽再驗證一次 檢查是否偽造
            if (signKeyAgain != originalSign)
            {
                throw new MSCException(ErrorCode.Authorize_JwtTokenParsingFail);
            }

            if(!availableVersions.Contains(model.TokenVersion))
            {
                throw new MSCException(ErrorCode.Authorize_JwtTokenVersionNotMatch);
            }
        }

        public static long GetTokenExpiredSeconds()
        {
            //dev環境 token過期時間延長
            int extraRate = SettingHelper.IsDev ? 10 : 1;

            return SystemConstant.JwtToken_Expired_Seconds * extraRate;
        }


        private static void CheckSignPermissionIsValid(HashSet<string> permissions , string mySignKey)
        {
            var expectSignKey = GeneratePermissionSignKey(permissions);

            if(expectSignKey != mySignKey)
            {
                throw new MSCException(ErrorCode.Authorize_JwtTokenVersionNotMatch);
            }
        }


        private static string GenerateSignKey(int userID , long expiredTime , string rndID)
        {
            return $"{userID}_{SystemConstant.JwtToken_SaltKey}_{expiredTime}_{rndID}".MD5Hash();
        }

        public static string GeneratePermissionSignKey(HashSet<string> permissions)
        {
            string concateText = string.Join(",", permissions.OrderBy(z => z).ToList());
            return $"{SystemConstant.JwtToken_SaltKey}_{concateText}_{SystemConstant.JwtToken_SaltKey}".MD5Hash();
        }
    }
}
