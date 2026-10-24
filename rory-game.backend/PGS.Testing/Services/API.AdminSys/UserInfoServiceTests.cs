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
using System.Collections.Generic;
using System;
using target.Models;

namespace target.Testing
{
    [TestClass]
    public class UserInfoServiceTests : BaseTest
    {

        public IUserInfoService ServiceInstance { get; set; }


        [TestMethod]
        public void GetUserProfile_Test()
        {
            var query = ServiceInstance.GetUserProfile(_TEST_DATA.Default.UserID);
            Assert.IsNotNull(query);
            return;
        }


        [TestMethod]
        public void UpdateUserLanguage_Test()
        {
            ServiceInstance.UpdateUserLanguage(_TEST_DATA.Default.UserID , SystemConstant.DefaultLang);
            return;
        }


    }
}

