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
    public class ViewPermissionCheckServiceTests : BaseTest
    {

        public IViewPermissionCheckService ServiceInstance { get; set; }


        [TestMethod]
        public void CheckViewMerchantPermission_Test()
        {
            ServiceInstance.CheckViewMerchantPermission(_TEST_DATA.Default.UserID , _TEST_DATA.Default.MerchantID);
            return;
        }


        [TestMethod]
        public void CheckViewPlayerPermission_Test()
        {
            ServiceInstance.CheckViewPlayerPermission(_TEST_DATA.Default.UserID, _TEST_DATA.Default.PlayerID);
            return;
        }


        [TestMethod]
        public void CheckViewUserPermission_Test()
        {
            ServiceInstance.CheckViewUserPermission(_TEST_DATA.Default.UserID, _TEST_DATA.Default.UserID);
            return;
        }



    }
}

