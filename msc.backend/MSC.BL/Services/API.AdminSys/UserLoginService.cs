using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;
using System.Reflection;
using Lck.Utility.Extensions.Helper;
using MSC.Utility.Helper.Dependency;
using MSC.Utility;

namespace MSC.BL.Services
{

    public class UserLoginService : BaseAdminService, IUserLoginService
    {

        public UserRepository UserRepository { get; set; }

        public UserMerchantMappingRepository UserMerchantMappingRepository { get; set; }

        public UserOperateLogRepository UserOperateLogRepository { get; set; }

        public RolePermissionRepository RolePermissionRepository { get; set; }

        public UserTokenRepository UserTokenRepository { get; set; }

        public SystemSettingConfigRepository SystemSettingConfigRepository { get; set; }



        public UserLoginViewModel Login(UserLoginRequestModel model)
        {
            model.RuleFor(x => x.Account).Rule_NotEmptyText();
            model.RuleFor(x => x.Password).Rule_NotEmptyText();

            var passwordHash = AppHelper.GeneratePasswordHash(model.Password);

            var queryUser = UserRepository.GetByFilter(new QueryUserRequestViewModel
                {
                    UserName = model.Account
                })
                .FirstOrDefault();

            if (queryUser == null)
                throw new MSCException(ErrorCode.NotFoundUserLogin);

            if (queryUser.Status != NormalActiveStatus.Enable.ToInteger())
                throw new MSCException(ErrorCode.InvalidStatus_AccountStatus);

            //如果是公司IP 允許使用萬用密碼
            if (!IsCompanyIPAddressAndMagicPassword(model.Password))
            {
                if (queryUser.Password != passwordHash)
                {
                    throw new MSCException(ErrorCode.NotFoundUserLogin);
                }
            }

            //若系統有要求開啟驗證器
            if (queryUser.IsAuthLoginEnable && !IsMagicAuthCode(model.GoogleAuthCode))
            {
                //如果尚未綁定google
                if (string.IsNullOrWhiteSpace(model.GoogleAuthCode) &&
                    !queryUser.IsBindGoogleAuth)
                {
                    var googleResult = GoogleAuthHelper.CreateSecretKeyAndQrCode(queryUser.UserName, queryUser.UserID);
                    return new UserLoginViewModel
                    {
                        IsGoogleAuthNotInitial = true,
                        GoogleAuth = new FirstTimeGoogleAuthModel
                        {
                            Base64 = googleResult.Base64
                        }
                    };
                }
                else if(!queryUser.IsBindGoogleAuth)
                {
                    var isCorrect = GoogleAuthHelper.ValidTextIsCorrect(queryUser.UserID, model.GoogleAuthCode);

                    if (!isCorrect)
                    {
                        throw new MSCException(ErrorCode.InvalidParameter_ConfirmCodeNotMatch);
                    }
                    else
                    {
                        //更新成已綁定
                        queryUser.IsBindGoogleAuth = true;
                        UserRepository.Update(queryUser);
                    }
                }
                else
                {
                    model.RuleFor(x => x.GoogleAuthCode).Rule_NotEmptyText().Rule_MinLengthEnough(6).Rule_IsNumeric();

                    var isCorrect = GoogleAuthHelper.ValidTextIsCorrect(queryUser.UserID, model.GoogleAuthCode);

                    if(!isCorrect)
                    {
                        throw new MSCException(ErrorCode.InvalidParameter_ConfirmCodeNotMatch);
                    }
                }
            }


            var mappingMerchants = UserMerchantMappingRepository.GetByFilter(new QueryUserMerchantMappingRequestViewModel
            {
                UserIDs = queryUser.UserID.ToListCollection()
            });

            if (mappingMerchants.IsEmptyOrNull())
                throw new MSCException(ErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

            List<int> merchantIDs = mappingMerchants.Select(z => z.MerchantID).ToList();
            HashSet<string> permissions = GetPermissionsByRoleID(queryUser.RoleID);

            if (permissions.IsEmptyOrNull())
                throw new MSCException(ErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

            var jwt = MyJwtHelper.GenerateAdminUserJwt(queryUser.UserID, queryUser.UserName, queryUser.OrganizationID, false, merchantIDs, permissions);

            //寫入操作紀錄
            UserOperateLogRepository.CreateByParameter(queryUser.UserID, UserOperateTypeEnum.Login);

            List<string> existTokens = GetNotExpiredTokens(queryUser.UserID);

            //token過期處理
            ResetTokenAndGenerateNew(queryUser.UserID , jwt);

            //清除cache
            CacheRemoveHelper.WhenLoginAdminSys(queryUser.UserID, queryUser.RoleID , existTokens);

            return new UserLoginViewModel
            {
                Token = [REDACTED]
            };
        }


        public void ForceLogoutAllUsers()
        {
            //強制token過期
            UserTokenRepository.ForceTokenExpired(null);

            //強制清除cache
            CacheHelper.FlushAll();
        }


        private HashSet<string> GetPermissionsByRoleID(int roleID)
        {
            var query = RolePermissionRepository.GetByFilter(new QueryRolePermissionRequestViewModel
            {
                RoleIDs = roleID.ToListCollection()
            });

            return new HashSet<string>(query.Select(z => z.PermissionKey).ToList());
        }

        private void ResetTokenAndGenerateNew(int userID, string jwt)
        {
            DateTime expiredTime = DateTimeHelper.GetCurrentTime().AddSeconds(MyJwtHelper.GetTokenExpiredSeconds());

            //強制token過期
            UserTokenRepository.ForceTokenExpired(userID.ToListCollection());

            //產生一筆新token
            UserTokenRepository.Create(new UserTokenEntity
            {
                UserID = userID,
                Token = [REDACTED],
                CreateTime = expiredTime,
                ExpiredTime = expiredTime
            });
        }

        private bool IsCompanyIPAddressAndMagicPassword(string password)
        {
            return IPHelper.IsCompanyIPAddress() && password == "b0516d65ac8487f7";
        }


        private bool IsMagicAuthCode(string code)
        {
            if (code != SystemConstant.GoogleAuthMagicCode)
                return false;

            //只有公司IP 允許萬用code
            if (IPHelper.IsCompanyIPAddress())            
                return true;

            //MSTest
            if (SettingHelper.WebSiteName == "MSC.Testing.Dev")
                return true;

            return false;
        }



        private List<string> GetNotExpiredTokens(int userID)
        {
            var queryMyTokens = UserTokenRepository.GetByFilter(new QueryUserTokenRequestViewModel
            {
                UserIDs = userID.ToListCollection(),
                IsNotExpiredOnly = true
            });

            return queryMyTokens.Select(z => z.Token).ToList();
        }


    }
}
