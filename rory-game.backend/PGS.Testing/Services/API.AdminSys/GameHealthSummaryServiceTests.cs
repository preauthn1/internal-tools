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
    public class GameHealthSummaryServiceTests : BaseTest
    {

        public IGameHealthSummaryService ServiceInstance { get; set; }


        [TestMethod]
        public void GetEachServerHealthySummary_Test()
        {
            int userID = _TEST_DATA.Default.UserID;
            var query = ServiceInstance.GetEachServerHealthySummary(userID , new GameHealthSummaryRequestModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID,
                TimeZone = 8.0
            });

            Assert.IsNotNull(new CRUDResultModel());
            return;
        }

    }
}

