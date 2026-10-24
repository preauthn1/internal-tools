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
    public class StatReportServiceTests : BaseTest
    {

        public IStatReportService ServiceInstance { get; set; }


        [TestMethod]
        public void StatReportService_GetMerchantDailySummary_Test()
        {
            var query = ServiceInstance.GetMerchantDailySummary(_TEST_DATA.Default.UserID, new StatReportRequestModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID,
                BeginDate = DateTimeHelper.GetCurrentTime().AddMonths(-1),
                EndDate = DateTimeHelper.GetCurrentTime()
            });

            return;
        }


        [TestMethod]
        public void StatReportService_GetMerchantSummaryInfo_Test()
        {
            var query = ServiceInstance.GetMerchantSummaryInfo(_TEST_DATA.Default.UserID, _TEST_DATA.Default.MerchantID);
            Assert.IsNotNull(query);
            Assert.IsTrue(query.MerchantID > 0);
            return;
        }

    }
}

