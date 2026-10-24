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
    public class PlayerWalletRefServiceTests : BaseTest
    {

        public IPlayerWalletRefService ServiceInstance { get; set; }


        [TestMethod]
        public void GetWalletAmountByMerchantWalletMode_Test()
        {
            var query = ServiceInstance.GetWalletAmountByMerchantWalletMode(_TEST_DATA.Default.PlayerID , WalletTypeEnum.GameWallet);
            Assert.IsNotNull(query);
            return;
        }

    }
}

