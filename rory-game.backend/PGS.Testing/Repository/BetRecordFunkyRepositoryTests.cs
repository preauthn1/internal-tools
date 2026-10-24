using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class BetRecordFunkyRepositoryTests : BaseTest
    {

        public BetRecordFunkyRepository BetRecordFunkyRepository { get; set; }


        [TestMethod]
        public void GetTotalSummaryReport_Test()
        {
            DateTime begin = DateTimeHelper.GetCurrentTime().AddDays(-2);
            DateTime end = DateTimeHelper.GetCurrentTime();

            var query = BetRecordFunkyRepository.GetTotalSummaryReport(begin , end);
            return;
        }

    }
}

