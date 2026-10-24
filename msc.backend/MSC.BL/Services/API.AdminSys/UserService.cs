using Lck.Utility.Extensions.Helper;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.BL.Services
{
    public class UserService : BaseAdminService, IUserService
    {
        public RoleRepository RoleRepository { get; set; }

        public UserMerchantMappingRepository UserMerchantMappingRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public UserOperateLogRepository UserOperateLogRepository { get; set; }

        public PagedList<UserViewModel> GetByPagingFilter(int userID, QueryUserApiRequestParameterModel filter)
        {
            var query = UserRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<UserEntity, UserViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            SetExtraAttribute(viewModels.Items);
            return viewModels;
        }


        public void Create(int userID, UserCreateRequestModel model)
        {
            model.RuleFor(x => x.UserName).Rule_NotEmptyText();
            model.RuleFor(x => x.DisplayName).Rule_NotEmptyText();
            model.RuleFor(x => x.Password).Rule_NotEmptyText();
            model.RuleFor(x => x.ConfirmPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();

            if (model.Password != model.ConfirmPassword)
                throw new MSCException(ErrorCode.InvalidParameter_ConfirmCodeNotMatch);

            var queryExist = UserRepository.GetByFilter(new QueryUserRequestViewModel
            {
                UserName = model.UserName
            });

            if (queryExist.IsNotEmptyNull())
                throw new MSCException(ErrorCode.BusinessLogicDeny_DuplicateDataExist);

            string userName = GetUserNameByUserID(userID);

            var newEntity = model.ReturnInjecteObject<UserEntity>();
            newEntity.IsAuthLoginEnable = true;
            newEntity.Password = [REDACTED];
            newEntity.Status = NormalActiveStatus.Enable.ToInteger();
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.CreateBy = userName;
            newEntity.UpdateBy = userName;

            // Step 1. user資料
            long resultUserID = UserRepository.Create(newEntity);

            // Step 2. Merchant Mapping
            UserMerchantMappingRepository.Create(new UserMerchantMappingEntity
            {
                UserID = (int)resultUserID,
                MerchantID = model.MerchantID,

            });
        }


        public void Update(int userID, UserCreateUpdateModel model)
        {
            model.RuleFor(x => x.UserID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.DisplayName).Rule_NotEmptyText();
            model.RuleFor(x => x.Password).Rule_NotEmptyText();
            model.RuleFor(x => x.ConfirmPassword).Rule_NotEmptyText();

            var queryExist = UserRepository.GetById(model.UserID);

            if (queryExist == null)
                throw new MSCException(ErrorCode.NotFoundUser);

            string userName = GetUserNameByUserID(userID);

            queryExist.DisplayName = model.DisplayName;
            queryExist.Status = model.Status;
            queryExist.CreateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.CreateBy = userName;
            queryExist.UpdateBy = userName;

            UserRepository.Update(queryExist);

            //更新cache
            CacheRemoveHelper.WhenUpdateUser(queryExist.UserID);
        }


        public void ChangePassword(int currentUserID, UserChangePasswordModel model)
        {
            model.RuleFor(x => x.TargetUserID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.NewPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.ConfirmPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.NewPassword).Rule_IsValidPasswordFormat();

            if (model.ConfirmPassword != model.NewPassword)
                throw new MSCException(ErrorCode.InvalidParameter_ConfirmCodeNotMatch);

            var queryUser = UserRepository.GetById(model.TargetUserID);

            if (queryUser == null)
                throw new MSCException(ErrorCode.NotFoundUser);

            queryUser.Password = [REDACTED];
            queryUser.UpdateTime = DateTimeHelper.GetCurrentTime();

            // Step 1. 更新
            UserRepository.Update(queryUser);

            // Step 2. 寫入一筆紀錄
            UserOperateLogRepository.CreateByParameter(currentUserID, UserOperateTypeEnum.UpdateOtherPassword, model.TargetUserID);
        }


        public void ChangeSelfPassword(int targetUserID, UserChangeSelfPasswordModel model)
        {
            model.RuleFor(x => x.NewPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.OldPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.ConfirmPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.NewPassword).Rule_IsValidPasswordFormat();

            if (model.ConfirmPassword != model.NewPassword)
                throw new MSCException(ErrorCode.InvalidParameter_ConfirmCodeNotMatch);

            var queryExist = UserRepository.GetById(targetUserID);

            if (queryExist == null)
                throw new MSCException(ErrorCode.NotFoundUser);

            string expectPassword = [REDACTED];

            if (queryExist.Password != expectPassword)
                throw new MSCException(ErrorCode.InvalidParameter_OriginalPasswordNotMatch);

            queryExist.Password = [REDACTED];
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();

            // Step 1. 更新
            UserRepository.Update(queryExist);

            // Step 2. 寫入一筆紀錄
            UserOperateLogRepository.CreateByParameter(targetUserID, UserOperateTypeEnum.UpdateSelfPassword, targetUserID);
        }



        private UserViewModel ConvertToViewModels(UserEntity query)
        {
            var convert = query.ReturnInjecteObject<UserViewModel>();
            convert.StatusText = MultiLanguageTextHelper.GetEnums_i18nText((NormalActiveStatus)query.Status);
            return convert;
        }


        private void SetExtraAttribute(IEnumerable<UserViewModel> models)
        {
            List<int> userIDs = models.Select(z => z.UserID).ToList();

            var queryUserMerchants = UserMerchantMappingRepository.GetByFilter(new QueryUserMerchantMappingRequestViewModel
            {
                UserIDs = userIDs
            });

            List<int> roleIDs = models.Select(x => x.RoleID).Distinct().ToList();

            var queryRoles = RoleRepository.GetByIds(roleIDs);

            foreach (var item in models)
            {
                item.RoleName = queryRoles.FirstOrDefault(z => z.RoleID == item.RoleID)?.DisplayName;

                var matchOfThis = queryUserMerchants.Where(z => z.UserID == item.UserID).ToList();

                item.AvailableMerchantCount = matchOfThis.Count();
                item.AvailableMerchantIDs = matchOfThis.Select(z => z.MerchantID).ToList();
                item.AvailableMerchantNames = matchOfThis.Select(z => GetMerchantName(z.MerchantID)).ToList();
                item.AvailableMerchantName = string.Join(",", item.AvailableMerchantIDs.Select(z => GetMerchantName(z)));
            }
        }

        private string GetMerchantName(int merchantID)
        {
            return MerchantRepository.GetSingleByCache(merchantID)?.DisplayName;
        }


    }
}
