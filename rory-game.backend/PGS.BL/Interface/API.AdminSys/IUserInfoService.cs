using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IUserInfoService
    {        
        UserInfoViewModel GetUserProfile(int userID);

        void UpdateUserLanguage(int userID, string lang);

        void UpdateUserPassword(int userID, UpdateCurrentUserPasswordRequestModel model);

    }
}
