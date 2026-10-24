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
    public class SystemSettingConfigServiceTests : BaseTest
    {

        public ISystemSettingConfigService ServiceInstance { get; set; }


        [TestMethod]
        public void SystemSettingConfigService_GetGlobalSystemSettings_Test()
        {
            var query = ServiceInstance.GetGlobalSystemSettings();
            Assert.IsNotNull(query);
            Assert_NotEmptyCollection(query.Groups);
            return;
        }

    }
}

