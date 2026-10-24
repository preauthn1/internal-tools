using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.BL.Interface
{
    public interface IUserService
    {
        PagedList<UserViewModel> GetByPagingFilter(int userID, QueryUserApiRequestParameterModel filter);

        void Create(int userID, UserCreateRequestModel model);

        void Update(int userID, UserCreateUpdateModel model);

        void ChangePassword(int currentUserID, UserChangePasswordModel model);

        void ChangeSelfPassword(int targetUserID, UserChangeSelfPasswordModel model);
    }
}
