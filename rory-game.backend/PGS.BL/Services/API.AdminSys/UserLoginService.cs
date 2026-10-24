using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;
using target.Utility;
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{

    public class UserLoginService : BaseAdminSysService, IUserLoginService
    {

        public UserMerchantMappingRepository UserMerchantMappingRepository { get; set; }

        public UserRepository UserRepository { get; set; }

        public OrganizationMerchantMappingRepository OrganizationMerchantMappingRepository { get; set; }

        public OrganizationRepository OrganizationRepository { get; set; }


        public UserLoginViewModel Login(UserLoginRequestModel model)
        {
            model.RuleFor(x => x.Account).Rule_NotEmptyText().Rule_IsValidAccountFormat();
            model.RuleFor(x => x.Password).Rule_NotEmptyText();

            string passwordMD5 = AppHelper.GeneratePasswordHash(model.Password);

            var queryUser = UserRepository.GetByAccountPassword(model.Account , passwordMD5);

            if (queryUser == null)
                throw new TerproException(ApiErrorCode.NotFoundUserLogin);

            if (queryUser.Status != NormalActiveStatus.Enable.ToInteger())
                throw new TerproException(ApiErrorCode.InvalidStatus_AccountStatus);

            //查出此用戶擁有哪幾家商戶的權限
            var queryMapping = UserMerchantMappingRepository.GetByFilter(new QueryUserMerchantMappingRequestViewModel
            {
                UserIDs = queryUser.UserID.ToListCollection()
            });

            //如果完全沒有分配任何商戶直接拋錯誤 (強制要求至少要一家)
            if (queryMapping.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

            List<int> availableMerchantIDs = queryMapping.Select(x => x.MerchantID).ToList();

            var queryOrg = OrganizationRepository.GetById(queryUser.OrganizationID);

            if (queryOrg.Status != NormalActiveStatus.Enable.ToInteger())
                throw new TerproException(ApiErrorCode.InvalidStatus_OrganizationStatus);
         
            //更新最後登入時間
            UserRepository.UpdateLastLoginInfo(queryUser.UserID, IPHelper.GetIP());

            //寫入登入紀錄
            WriteOperateLog(queryUser.UserID, UserOperateTypeEnum.User_UserLogin, 0);

            //刷新Cache
            CacheRemoveHelper.RemoveCache_WhenUserLogin(queryUser.UserID);


            bool isSuperAdmin = queryOrg.OrganizationType == OrganizationTypeEnum.SuperAdmin.ToInteger();
            var jwt = MyJwtHelper.GenerateAdminUserJwt(queryUser.UserID, queryUser.UserName  , queryOrg.OrganizationID, isSuperAdmin, availableMerchantIDs);
            return new UserLoginViewModel
            {
                Access_Token = [REDACTED]
            };
        }


    }
}
