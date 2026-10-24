using target.Core.Repository;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using target.Utility.Extensions;
using target.DA.Repository;
using target.Models.Table;
using target.Utility.Helper;
using target.Models.ViewModels.Base;
using Lck.Utility.Extensions.Helper;
using target.Utility.IOC;
using target.Utility;

namespace target.BL.Services
{
    public class BaseAdminSysService : BaseService
    {

        /// <summary>
        /// 寫入後台修改日誌 (其實apiLog可以處理大部分紀錄了)
        /// 但如果是要讓後台可以查的記錄 寫入這張表
        /// </summary>
        protected void WriteOperateLog(int userID, UserOperateTypeEnum type, int targetID, UserOperateChangeLogModel changeLog = null)
        {
            var rep = InstanceRepository<UserOperateRecordRepository>();

            var model = InstanceOperateLog(userID, type, targetID, changeLog);

            rep.Create(model);
        }


        protected UserOperateRecordEntity InstanceOperateLog(int userID, UserOperateTypeEnum type, int targetID, UserOperateChangeLogModel changeLog = null)
        {
            var rep = InstanceRepository<UserOperateRecordRepository>();

            return new UserOperateRecordEntity
            {
                CreateTime = DateTimeHelper.GetCurrentTime(),
                CreateIPAddress = IPHelper.GetIP(),
                ChangeLogInfo = changeLog.ConvertToJsonString(),
                UserID = userID,
                OperateType = type.ToInteger(),
                TargetID = targetID,
                MerchantID = 0,
                Memo = type.ToString()
            };
        }


        protected virtual AdminUserTokenJwtModel GetCurrentUser()
        {
            string jwt = HeaderValueHelper.GetJwtToken();

            var convert = MyJwtHelper.GetAdminUserModelFromJwtString(jwt);

            if (convert == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            return convert;
        }


        protected bool IsInternalTestAccountPlayerOrMerchantPlayer(int playerID)
        {
            var playerRepository = new PlayerRepository();

            var queryPlayer = playerRepository.GetSingleByCache(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var rep1 = new MerchantCustomizedSettingRepository();

            var queryMerchantSetting = rep1.GetMerchantSettingByCache(queryPlayer.MerchantID, MerchantCustomizedSettingKeyEnum.IsInternalTestMerchant);

            if (queryMerchantSetting != null && queryMerchantSetting.SettingValueOfBoolean)
                return true;

            var rep2 = new SystemConfigSettingRepository();

            var querySettingPlayerIDs = rep2.GetParseModelByEnum(SystemConfigSettingEnum.InternalTestAccountPlayerIDs);

            if (querySettingPlayerIDs.SettingValueOfListInt.Contains(playerID))
                return true;

            return false;
        }




        protected List<int> FilterAvailablePermissionMerchantIDs(int userID, List<int> merchantIDs)
        {
            var queryCurrentUser = GetCurrentOwnPermissionByCache(userID);

            //如果是超級管理員的組織
            if (queryCurrentUser.IsSuperAdminOrg)
                return merchantIDs ?? new List<int>();

            //如果一般等級的廠商，沒篩選merchantID 強制要求要篩選
            if (merchantIDs.IsEmptyOrNull())
                return queryCurrentUser.AvailableMerchantIDs;

            //檢查是否有帶 違反權限的merchantIDs
            if (merchantIDs.Any(x => !queryCurrentUser.AvailableMerchantIDs.Contains(x)))
                throw new TerproException(ApiErrorCode.Forbidden_WatchDataPermissionNotAllow);

            return merchantIDs;
        }


        protected AdminUserCacheModelAdminUpdateViewModel GetCurrentOwnPermissionByCache(int userID)
        {
            if (userID <= 0)
                return null;

            string cacheKey = CacheKey.AdminSys_UserOwnPermission(userID);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                var merchantRep = InstanceRepository<MerchantRepository>();
                var userRep = InstanceRepository<UserRepository>();
                var orgRep = InstanceRepository<OrganizationRepository>();
                var orgMappingRep = InstanceRepository<OrganizationMerchantMappingRepository>();
                var roleMappingRep = InstanceRepository<UserRoleMappingRepository>();
                var userMerchantRep = InstanceRepository<UserMerchantMappingRepository>();

                var queryUser = userRep.GetById(userID);
                var queryOrg = orgRep.GetById(queryUser.OrganizationID);
                var queryMyMerchants = userMerchantRep.GetByUserID(userID);

                bool isSuperAdmin = queryOrg.OrganizationType == OrganizationTypeEnum.SuperAdmin.ToInteger();

                var result = new AdminUserCacheModelAdminUpdateViewModel
                {
                    UserID = queryUser.UserID,
                    UserName = queryUser.UserName,
                    OrganizationID = queryUser.OrganizationID,
                    IsSuperAdminOrg = isSuperAdmin,
                    Roles = roleMappingRep.GetByUserID(queryUser.UserID).Select(x => x.RoleID).ToList(),
                    AvailableOrganizationIDs = queryUser.OrganizationID.ToListCollection(),
                    AvailableMerchantIDs = queryMyMerchants.Select(x => x.MerchantID).ToList(),
                    AllEnableMerchantsOrAvailableMerchants = queryMyMerchants.Select(x => x.MerchantID).ToList()
                };

                //如果是超級管理員組織
                if (isSuperAdmin)
                {
                    var queryEnableMerchants = merchantRep.GetByFilter(new QueryMerchantRequestViewModel
                    {
                        Status = NormalActiveStatus.Enable.ToInteger()
                    });

                    result.AllEnableMerchantsOrAvailableMerchants = queryEnableMerchants.Select(z => z.MerchantID).ToList();
                }

                return result;

            }, 10);

            return queryCache;
        }


        protected void CheckViewMerchantPermission(int userID, int merchantID)
        {
            var queryUser = GetCurrentOwnPermissionByCache(userID);

            if (queryUser == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            if (queryUser.IsSuperAdminOrg)
                return;

            if (queryUser.AvailableMerchantIDs.Contains(merchantID))
                return;

            throw new TerproException(ApiErrorCode.Forbidden_WatchDataPermissionNotAllow);
        }



    }
}
