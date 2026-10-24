using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using target.DA.Repository;
using target.Models.Enums;

namespace target.Testing
{
    [TestClass]
    public class UsersServiceTests : BaseTest
    {

        public IUsersService UsersService { get; set; }


        [TestMethod]
        public void GetUsersPagingByFilter_Test()
        {
            var query = UsersService.GetUsersPagingByFilter(new QueryUserApiRequestParameterModel { PageSize = 10});
            Assert_NotEmptyCollection(query.Items);
            return;
        }



        [TestMethod]
        public void GetUserByID_Test()
        {
            var query = UsersService.GetUserByID(_TEST_DATA.Default.UserID);
            Assert.IsNotNull(query);
            Assert.IsTrue(query.UserID == _TEST_DATA.Default.UserID);
            return;
        }



        [TestMethod]
        public void User_Update_Test()
        {
            var query = UsersService.GetUserByID(_TEST_DATA.Default.UserID);

            Assert.IsNotNull(query);
            
            var updateModel = query.ReturnInjecteObject<UserUpdateModel>();
            updateModel.Merchants = query.Merchants;
            updateModel.Status = NormalActiveStatus.Enable.ToInteger(); //永遠保持開啟

            UsersService.Update(updateModel);
            return;
        }



    }
}

