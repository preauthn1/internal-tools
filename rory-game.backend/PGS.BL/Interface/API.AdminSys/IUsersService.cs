using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IUsersService
    {
        PagedList<UsersViewModel> GetUsersPagingByFilter(QueryUserApiRequestParameterModel filter);


        UsersUpdatePageViewModel GetUserByID(int userID);


        void Create(UserCreateModel model);


        void Update(UserUpdateModel model);


        void Delete(int userID);

    }
}
