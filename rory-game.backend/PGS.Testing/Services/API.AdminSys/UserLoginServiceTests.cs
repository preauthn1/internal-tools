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
    public class UserLoginServiceTests : BaseTest
    {

        public IUserLoginService UserLoginService { get; set; }


        [TestMethod]
        public void User_Login_Test()
        {
            var query = UserLoginService.Login(new UserLoginRequestModel 
            { 
                Account = _TEST_DATA.Default.UserName , 
                Password = [REDACTED]"
            });
            Assert.IsNotNull(query);
            Assert.IsNotNull(query.Access_Token);
            return;
        }

    }
}

