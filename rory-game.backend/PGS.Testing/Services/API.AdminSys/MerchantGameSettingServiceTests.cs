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
    public class MerchantGameSettingServiceTests : BaseTest
    {

        public IMerchantGameSettingService ServiceInstance { get; set; }


        [TestMethod]
        public void GetGameSettingByMerchantID_Test()
        {
            var query = ServiceInstance.GetGameSettingByMerchantID(_TEST_DATA.Default.MerchantID);

            Assert.IsNotNull(query);
            Assert_NotEmptyCollection(query.GameSettingItems);
            return;
        }

    }
}

