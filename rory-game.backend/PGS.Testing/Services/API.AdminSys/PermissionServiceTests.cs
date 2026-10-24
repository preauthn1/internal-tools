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
    public class PermissionServiceTests : BaseTest
    {

        public IPermissionService PermissionService { get; set; }


        [TestMethod]
        public void GetAvailablePermissions_Test()
        {
            var query = PermissionService.GetAvailablePermissions(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query.AllPermissions);
            return;
        }


        [TestMethod]
        public void GetPermissions_Test()
        {
            var query = PermissionService.GetPermissions(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query);
            return;
        }


        [TestMethod]
        public void GetPermissionsNestedFormat_Test()
        {
            var query = PermissionService.GetPermissionsNestedFormat(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query.AllPermissions);
            return;
        }

    }
}

