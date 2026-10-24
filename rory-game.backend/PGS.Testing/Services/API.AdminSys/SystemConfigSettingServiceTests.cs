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
    public class SystemConfigSettingServiceTests : BaseTest
    {

        public ISystemConfigSettingService ServiceInstance { get; set; }


        [TestMethod]
        public void GetAllSystemConfigSetting_Test()
        {
            var query = ServiceInstance.GetAllSystemConfigSetting(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query.Groups);
            Assert.IsNotNull(query);
            return;
        }

    }
}

