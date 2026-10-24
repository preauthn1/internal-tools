using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Collections.Generic;
using System;
using MSC.DA.Repository;
using MSC.Models.ViewModels;
using MSC.Testing;
using MSC.Utility.Extensions;

namespace MSC.Testing
{
    [TestClass]
    public class UserRepositoryTests : BaseTest
    {

        public UserRepository UserRepository { get; set; }


        [TestMethod]
        public void UserRepository_GetPagingByFilter_Test()
        {
            var query = UserRepository.GetPagingByFilter(new QueryUserRequestViewModel
            {
                PageSize = 10
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }




        [TestMethod]
        public void UserRepository_GetUserNameByUserID_Test()
        {
            string userName = UserRepository.GetUserNameByUserID(_TEST_DATA.Default.UserID);
            Assert.IsTrue(userName.IsNotNullNotWhiteSpace());
            return;
        }



    }
}

