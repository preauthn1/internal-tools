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
    public class SingleWalletVendorServiceTests : BaseTest
    {

        public IGameVendorSingleWalletService ServiceInstance { get; set; }


        [TestMethod]
        public void GetWalletBalance_Test()
        {
            var query = ServiceInstance.GetWalletBalance(GameVendorTypeEnum.AGG, new SingleWalletGetWalletBalanceRequestParameter
            {
                PlayerID = _TEST_DATA.Default.PlayerID
            });

            Assert.IsNotNull(query);
            return;
        }


        [TestMethod]
        public void GameVendor_PlaceBet_Test()
        {
            //var query = ServiceInstance.PlaceBet(GameVendorTypeEnum.Funky , new SingleWalletBetDebitAmountRequestParameter
            //{
            //    SessionToken = [REDACTED]
            //});

            //Assert.IsNotNull(query);
            return;
        }


    }
}

