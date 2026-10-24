using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.Testing
{
    [TestClass]
    public class RoleServiceTests : BaseTest
    {

        public IRoleService ServiceInstance { get; set; }


        [TestMethod]
        public void RoleService_GetRolePermissions_Test()
        {
            var query = ServiceInstance.GetRolePermissions(_TEST_DATA.Default.RoleID);
            Assert_NotEmptyCollection(query.Permissions);
            return;
        }

    }
}

