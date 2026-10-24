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

namespace target.Testing
{
    [TestClass]
    public class AdminSysDropdownServiceTests : BaseTest
    {

        public IAdminSysDropdownService ServiceInstance { get; set; }


        [TestMethod]
        public void GetAvailableMerchants_Test()
        {
            var query = ServiceInstance.GetAvailableMerchants(_TEST_DATA.Default.UserID , new QueryDropdownMerchantRequestViewModel { });
            Assert_NotEmptyCollection(query);
            return;
        }



        [TestMethod]
        public void GetAvailableMerchants_FilterOrgID_Test()
        {
            var filter = new QueryDropdownMerchantRequestViewModel
            {
                OrganizationID = _TEST_DATA.Default.OrganizationID
            };

            var query = ServiceInstance.GetAvailableMerchants(_TEST_DATA.Default.UserID, filter);
            Assert_NotEmptyCollection(query);
            return;
        }



        [TestMethod]
        public void GetAvailableOrgs_Test()
        {
            var query = ServiceInstance.GetAvailableOrgs(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query);
            return;
        }



        [TestMethod]
        public void GetAvailableRoles_Test()
        {
            var query = ServiceInstance.GetAvailableRoles(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query);
            return;
        }


        [TestMethod]
        public void GetGameTypes_Test()
        {
            var query = ServiceInstance.GetGameTypes();
            Assert_NotEmptyCollection(query);
            return;
        }


        [TestMethod]
        public void GetGameVendorTypes_Test()
        {
            var query = ServiceInstance.GetGameVendorTypes();
            Assert_NotEmptyCollection(query);
            return;
        }



        [TestMethod]
        public void GetGames_Test()
        {
            var query = ServiceInstance.GetGames(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query);
            return;
        }



        [TestMethod]
        public void GetServers_Test()
        {
            var query = ServiceInstance.GetServers(_TEST_DATA.Default.UserID , null);
            Assert_NotEmptyCollection(query);
            return;
        }


        [TestMethod]
        public void GetServers_FilterGameID_Test()
        {
            var query = ServiceInstance.GetServers(_TEST_DATA.Default.UserID, new QueryServerListApiRequestParameterModel { GameID = _TEST_DATA.Default.GameID });
            Assert_NotEmptyCollection(query);
            return;
        }



        [TestMethod]
        public void GetTimeZoneAreas_Test()
        {
            var query = ServiceInstance.GetTimeZoneAreas();
            Assert_NotEmptyCollection(query);
            return;
        }

    }
}

