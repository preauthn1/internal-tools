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
    public class SuperToolPoolServiceTests : BaseTest
    {

        public ISuperToolPoolService ServiceInstance { get; set; }


        [TestMethod]
        public void GetWinLossPool_Test()
        {
            var query = ServiceInstance.GetWinLossPool(new SuperToolPoolRequestModel
            {
                FlagID = 0,
                GameID = _TEST_DATA.Default.GameID
            });

            Assert.IsNotNull(query);
            return;
        }

        [TestMethod]
        public void UpdateWinLossPool_Test()
        {
            var query = ServiceInstance.UpdateWinLossPool(new SuperToolPoolUpdateModel
            {
                GameID = _TEST_DATA.Default.GameID,
                FlagID = 0,
                UpdatePoolAmount = 1.0M
            });

            Assert.IsNotNull(query);
            return;
        }


    }
}

