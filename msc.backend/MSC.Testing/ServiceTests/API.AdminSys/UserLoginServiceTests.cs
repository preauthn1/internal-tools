using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.Models.Models;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System.Linq;
using MSC.Models.CoreLibs;
using MSC.Models.ViewModels;
using MSC.DA.Repository;
using MSC.Models.Enums;
using System.Collections.Generic;
using System;
using MSC.Models;

namespace MSC.Testing
{
    [TestClass]
    public class UserLoginServiceTests : BaseTest
    {

        public IUserLoginService ServiceInstance { get; set; }


        [TestMethod]
        public void UserLoginService_Login_Test()
        {
            var query = ServiceInstance.Login(new UserLoginRequestModel
            {
                GoogleAuthCode = SystemConstant.GoogleAuthMagicCode,
                Account = _TEST_DATA.Default.UserName,
                Password = [REDACTED]
            });

            Assert.IsTrue(query.Token.IsNotNullNotWhiteSpace());
            return;
        }

    }
}

