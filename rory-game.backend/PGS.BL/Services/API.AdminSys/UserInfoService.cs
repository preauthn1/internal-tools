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

namespace target.BL.Services
{

    public class UserInfoService : BaseService, IUserInfoService
    {
        public UserRepository UserRepository { get; set; }

        public OrganizationRepository OrganizationRepository { get; set; }


        public UserInfoViewModel GetUserProfile(int userID)
        {
            var queryUser = UserRepository.GetById(userID);

            if (queryUser == null)
                throw new TerproException(ApiErrorCode.NotFoundUser);

            var queryOrgs = OrganizationRepository.GetById(queryUser.OrganizationID);

            TimeZoneAreaEnum timeZone = (TimeZoneAreaEnum)queryUser.TimeZoneAreaID;

            return new UserInfoViewModel
            {
                TimeZone = DataModelConvertHelper.GetTimeZoneByAreaID(queryUser.TimeZoneAreaID),
                TimeZoneAreaID = queryUser.TimeZoneAreaID,
                TimeZoneAreaName = MultiLanguageTextHelper.GetTimeZoneDisplayName(queryUser.TimeZoneAreaID),
                UserID = userID,
                UserAccount = queryUser.UserName,
                OrganizationID = queryUser.OrganizationID,
                OrganizationName = queryOrgs.OrganizationName,
                DefaultLang = queryUser.DefaultLang,
                IsSuperAdminOrganization = queryOrgs.OrganizationType == OrganizationTypeEnum.SuperAdmin.ToInteger()
            };
        }

        public void UpdateUserLanguage(int userID , string lang)
        {
            bool isSupport = RegexRuleHelper.IsSupportLanguage(lang);

            if (!isSupport)
                throw new TerproException(ApiErrorCode.InvalidParameter_NotSupportLanguageCodet);

            var queryExist = UserRepository.GetById(userID);
            queryExist.DefaultLang = lang;

            UserRepository.Update(queryExist);
        }

        public void UpdateUserPassword(int userID, UpdateCurrentUserPasswordRequestModel model)
        {
            model.RuleFor(x => x.OriginalPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.NewPassword).Rule_NotEmptyText();
            model.RuleFor(x => x.NewPassword).Rule_IsValidPasswordFormat();
            model.RuleFor(x => x.ConfirmPassword).Rule_NotEmptyText();

            if (model.NewPassword != model.ConfirmPassword)
                throw new TerproException(ApiErrorCode.InvalidParameter_ConfirmCodeNotMatch);
            
            var queryUser = UserRepository.GetById(userID);

            var originalPasswordHash = AppHelper.GeneratePasswordHash(model.OriginalPassword);

            if(queryUser.Password != originalPasswordHash)
            {
                throw new TerproException(ApiErrorCode.InvalidParameter_ConfirmCodeNotMatch);
            }
                            
            queryUser.Password = [REDACTED];
                       
            UserRepository.Update(queryUser);
        }

    }
}
