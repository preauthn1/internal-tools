using target.Models;
using target.Models.Models;
using target.Models.ViewModels.Base;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Enums;

namespace target.Utility.Helper
{
    public class MyJwtHelper
    {
        public static PlayerJwtGenerateResultModel GeneratePlayerJwt(int playerID , int merchantID , string loginName , string lang , string sessionTokenPGS , string sessionTokenRaw)
        {
            //dev環境 token過期時間延長
            int extraRate = SettingHelper.IsDev ? 10 : 1;

            var jwtModel = new PlayerTokenJwtModel
            {
                TokenVersion = SystemConstant.JwtTokenVersionPlayer,
                MerchantID = merchantID,
                LoginUserID = playerID,
                LoginName = loginName,
                Lang = lang,
                RndID = Guid.NewGuid().ToString(),
                SessionTokenPGS = sessionTokenPGS,
                SessionTokenRaw = sessionTokenRaw,
                JwtTokenType = JwtTokenTypeEnum.PlayerJWtToken.ToInteger(),
                TokenExpiredTime = DateTimeHelper.GetCurrentTimeTimestamp() + SystemConstant.JwtToken_Expired_Seconds * extraRate
            };

            jwtModel.SignKey = GeneratePlayerSignKey(jwtModel);

            var jwt = BaseJWTTokenHelper.Generate_JWT_Token(jwtModel);
            return new PlayerJwtGenerateResultModel
            {
                JwtToken = [REDACTED],
                JwtModel = jwtModel
            };
        }


        public static string GenerateAdminUserJwt(int userID, string loginName , int orgID , bool isSuperAdmin, List<int> merchantIDs)
        {
            //dev環境 token過期時間延長
            int extraRate = SettingHelper.IsDev ? 10 : 1;

            var jwtModel = new AdminUserTokenJwtModel
            {
                TokenVersion = SystemConstant.JwtTokenVersionAdminSys,
                IsSuperAdminOrg = isSuperAdmin,
                MerchantIDs = merchantIDs,
                LoginUserID = userID,
                LoginName = loginName,
                OrganizationID = orgID,
                RndID = Guid.NewGuid().ToString(),
                JwtTokenType = JwtTokenTypeEnum.AdminSysJwtToken.ToInteger(),
                TokenExpiredTime = DateTimeHelper.GetCurrentTimeTimestamp() + SystemConstant.JwtToken_Expired_Seconds * extraRate
            };

            jwtModel.SignKey = GenerateAdminSysSignKey(jwtModel);
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

            //這邊就不檢查signKey了 拿到的 LoginUserID 只是純粹log使用，即使偽造影響也不大
            return jwtModel.LoginUserID.ToString();
        }


        public static PlayerTokenJwtModel GetPlayerModelFromJwtString(string jwtTokenText)
        {
            if (string.IsNullOrWhiteSpace(jwtTokenText))
                return null;

            var jwtModel = BaseJWTTokenHelper.GetObjectFromJWT<PlayerTokenJwtModel>(jwtTokenText);

            if (jwtModel == null)
                return null;

            string expectSignKey = GeneratePlayerSignKey(jwtModel);
            CheckSignIsValid(jwtModel.SignKey , expectSignKey);

            return jwtModel;
        }
        
     
        public static AdminUserTokenJwtModel GetAdminUserModelFromJwtString(string jwtTokenText)
        {
            if (string.IsNullOrWhiteSpace(jwtTokenText))
                return null;

            var jwtModel = BaseJWTTokenHelper.GetObjectFromJWT<AdminUserTokenJwtModel>(jwtTokenText);

            if (jwtModel == null)
                return null;

            string expectSignKey = GenerateAdminSysSignKey(jwtModel);
            CheckSignIsValid(jwtModel.SignKey, expectSignKey);

            return jwtModel;
        }


        private static void CheckSignIsValid(string originalSign, string expectSignKey)
        {
            //加密簽再驗證一次 檢查是否偽造
            if (expectSignKey != originalSign)
            {
                throw new TerproException(ApiErrorCode.Authorize_JwtTokenParsingFail);
            }
        }

        private static string GeneratePlayerSignKey(PlayerTokenJwtModel model)
        {
            return ($"{model.LoginUserID}_{model.MerchantID}_{model.LoginName}_{SystemConstant.JwtToken_SaltKey}_{model.TokenExpiredTime}_{model.RndID}_{model.TokenVersion}").MD5Hash();
        }

        private static string GenerateAdminSysSignKey(AdminUserTokenJwtModel model)
        {
            return ($"{model.LoginUserID}_{model.OrganizationID}_{model.LoginName}_{SystemConstant.JwtToken_SaltKey}_{model.TokenExpiredTime}_{model.RndID}_{model.TokenVersion}").MD5Hash();
        }


    }
}
