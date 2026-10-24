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
using target.BL.Services;

namespace target.Testing
{
    [TestClass]
    public class BaseAdminSysServiceTests : BaseTest
    {
        public IAdminSysDropdownService ServiceInstance { get; set; }


        public class MockBaseAdminSysService : BaseAdminSysService
        {
            public AdminUserCacheModelAdminUpdateViewModel MockGetCurrentOwnPermissionByCache(int userID)
            {
                return this.GetCurrentOwnPermissionByCache(userID);
            }
        }



        [TestMethod]
        public void MockGetCurrentOwnPermissionByCache_Test()
        {
            var service = new MockBaseAdminSysService();
            var query = service.MockGetCurrentOwnPermissionByCache(_TEST_DATA.Default.UserID);
            Assert.IsTrue(query != null);
            return;
        }


    }
}

