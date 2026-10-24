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
    public class ConfigSettingServiceTests : BaseTest
    {

        public IConfigSettingService ServiceInstance { get; set; }


        [TestMethod]
        public void ConfigSettingService_GetShoppingsIframeUrl_Test()
        {
            var query = ServiceInstance.GetShoppingsIframeUrl();
            return;
        }


        [TestMethod]
        public void ConfigSettingService_GetCardInputPageSetting_Test()
        {
            var query = ServiceInstance.GetCardInputPageSetting(_TEST_DATA.Default.MerchantID);
            Assert.IsNotNull(query);
            return;
        }


    }
}

