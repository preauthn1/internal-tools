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
    public class GameInfoServiceTests : BaseTest
    {

        public IGameInfoService ServiceInstance { get; set; }


        [TestMethod]
        public void GetAllGames_Test()
        {
            var query = ServiceInstance.GetAllGames(new GameInfoRequestModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
            });

            Assert_NotEmptyCollection(query);
            return;
        }

    }
}

