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
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{

    public class UsersService : BaseAdminSysService, IUsersService
    {
        public UserRepository UserRepository { get; set; }

        public OrganizationRepository OrganizationRepository { get; set; }

        public UserRoleMappingRepository UserRoleMappingRepository { get; set; }

        public RoleRepository RoleRepository { get; set; }

        public UserMerchantMappingRepository UserMerchantMappingRepository { get; set; }



        public PagedList<UsersViewModel> GetUsersPagingByFilter(QueryUserApiRequestParameterModel filter)
        {
            filter.OrganizationIDs = filter.OrganizationID.ToIntegerListCollection();

            var query = UserRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<UserEntity , UsersViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels<UsersViewModel>(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }

        
        public UsersUpdatePageViewModel GetUserByID(int userID)
        {
            var queryUser = UserRepository.GetById(userID);

            if (queryUser == null)
                throw new TerproException(ApiErrorCode.NotFoundUser);

            var convert = ConvertToViewModels<UsersUpdatePageViewModel>(queryUser);

            var queryRoleSetting = UserRoleMappingRepository.GetByUserID(userID).FirstOrDefault();
            convert.RoleID = queryRoleSetting?.RoleID ?? 0;

            var queryMerchantMapping = UserMerchantMappingRepository.GetByUserID(userID);
            convert.Merchants = queryMerchantMapping.Select(x => x.MerchantID).ToList();

            return convert;
        }


        public void Create(UserCreateModel model)
        {
            model.RuleFor(x => x.UserName).Rule_NotEmptyText().Rule_IsValidAccountFormat();
            model.RuleFor(x => x.DisplayName).Rule_NotEmptyText().Rule_MaxLengthOK(maxLength: 50);
            model.RuleFor(x => x.Password).Rule_NotEmptyText().Rule_IsValidPasswordFormat();
            model.RuleFor(x => x.ConfirmPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.RoleID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.OrganizationID).Rule_NotZeroNotMinus();

            if (!RoleRepository.IsExist(model.RoleID))
                throw new TerproException(ApiErrorCode.NotFoundRole);

            if (model.Password != model.ConfirmPassword)
                throw new TerproException(ApiErrorCode.InvalidParameter_ConfirmCodeNotMatch);

            var queryOrgs = OrganizationRepository.GetById(model.OrganizationID);

            if (queryOrgs == null)
                throw new TerproException(ApiErrorCode.NotFoundOrganization);

            model.UserName = $"{model.UserName}";

            var queryDuplicate = UserRepository.GetByAccount(model.UserName);

            //檢查資料重複
            if (queryDuplicate != null)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);

            var newEntity = model.ReturnInjecteObject<UserEntity>();

            var currentUser = GetCurrentUser();
            newEntity.Password = [REDACTED];
            newEntity.Status = NormalActiveStatus.Enable.ToInteger();
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.LastLogonTime = DateTimeHelper.GetDefaultTime();
            newEntity.UpdateBy = currentUser.LoginName;
            newEntity.CreateBy = currentUser.LoginName;
            newEntity.DefaultLang = queryOrgs.Lang;

            if (newEntity.TimeZoneAreaID == 0)
                newEntity.TimeZoneAreaID = TimeZoneAreaEnum.China_0800.ToInteger();

            //step 1. 新增user資料
            var newUserID = UserRepository.Create(newEntity);

            //step 2. 新增role資料
            UserRoleMappingRepository.Create(new UserRoleMappingEntity
            {
                RoleID = model.RoleID,
                UserID = newUserID,
            });

            //stpe 3. 新增merchantMapping
            UserMerchantMappingRepository.DeleteAndInsertNew(newUserID, model.Merchants);
        }        
        

        public void Update(UserUpdateModel model)
        {
            model.RuleFor(x => x.UserID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.DisplayName).Rule_NotEmptyText().Rule_MaxLengthOK(maxLength: 50);
            model.RuleFor(x => x.RoleID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.TimeZoneAreaID).Rule_NotZeroNotMinus();

            //如果有輸入password才檢查
            if (!string.IsNullOrWhiteSpace(model.Password))
            {
                if (model.Password != model.ConfirmPassword)
                {
                    throw new TerproException(ApiErrorCode.InvalidParameter_ConfirmCodeNotMatch);
                }
            }

            if (!RoleRepository.IsExist(model.RoleID))
                throw new TerproException(ApiErrorCode.NotFoundRole);

            var queryExist = UserRepository.GetById(model.UserID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundUser);

            queryExist.DisplayName = model.DisplayName;
            queryExist.Status = model.Status;
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateBy = GetCurrentUser().LoginName;
            queryExist.TimeZoneAreaID = model.TimeZoneAreaID;


            //如果有輸入password才更新
            if (!string.IsNullOrWhiteSpace(model.Password))
            {
                queryExist.Password = [REDACTED];
            }
                       
            //step 1. 更新 user資料
            UserRepository.Update(queryExist);

            //step 2. 更新 role資料
            UserRoleMappingRepository.DeleteAndInsertNew(model.UserID,  model.RoleID.ToListCollection());

            //stpe 3. 更新 merchantMapping
            UserMerchantMappingRepository.DeleteAndInsertNew(model.UserID , model.Merchants);


            //清除cache
            CacheRemoveHelper.RemoveCache_WhenUpdateUser(model.UserID);
        }


        public void Delete(int userID)
        {
            var queryExist = UserRepository.GetById(userID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundUser);

            queryExist.Status = NormalActiveStatus.PermanentDelete.ToInteger();
            UserRepository.Update(queryExist);
        }


        private T ConvertToViewModels<T>(UserEntity queryUser) where T: UsersViewModel , new()
        {
            var convert = queryUser.ReturnInjecteObject<T>();

            var queryUserRoleMapping = UserRoleMappingRepository.GetByUserID(convert.UserID).FirstOrDefault();
            var roleName = string.Empty;
            if (queryUserRoleMapping != null) {
                var queryRole = RoleRepository.GetById(queryUserRoleMapping.RoleID);
                roleName = queryRole.RoleName;
            }
            convert.RoleName_CN = roleName;
            convert.RoleName_EN = roleName;
            convert.OrganizationName = OrganizationRepository.GetSingleOrgByCache(queryUser.OrganizationID)?.OrganizationName;
            convert.TimeZoneAreaName = MultiLanguageTextHelper.GetTimeZoneDisplayName(queryUser.TimeZoneAreaID);

            return convert;
        }

    }
}
