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
    public class StatReportServiceTests : BaseTest
    {

        public IStatReportService ServiceInstance { get; set; }


        [TestMethod]
        public void GetGameStatReport_Test()
        {
            var query = ServiceInstance.GetGameStatReport(_TEST_DATA.Default.UserID , new QueryStatGameReportApiRequestParameterModel());
            Assert.IsNotNull(query);
            return;
        }

        [TestMethod]
        public void GetMerchantDailySummaryStatReport_Test()
        {
            var query = ServiceInstance.GetMerchantDailySummaryStatReport(_TEST_DATA.Default.UserID , new QueryMerchantDailySummaryStatReportReportApiRequestParameterModel
            {
                Begin = DateTimeHelper.GetCurrentTime().AddDays(-30),
                End = DateTimeHelper.GetCurrentTime(),
                MerchantID = _TEST_DATA.Default.MerchantID,
                TimeZone = 0
            });

            Assert.IsNotNull(query);
            Assert_NotEmptyCollection(query.Items);
            return;
        }



        



    }
}

