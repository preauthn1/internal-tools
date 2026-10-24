using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class GameServerDistributeRefServiceTests : BaseTest
    {

        public IGameServerDistributeRefService ServiceInstance { get; set; }


        public class MockGameServerDistributeRefService : GameServerDistributeRefService
        {
            protected override List<ServerConcurrentCountModel> GetServerOnlinePlayers(List<int> serverIDs)
            {
                return new List<ServerConcurrentCountModel>
                {
                    new ServerConcurrentCountModel
                    {
                        ServerID = _TEST_DATA.Default.ServerID,
                        TotalCount = 9999
                    }
                };
            }
        }


        [TestMethod]
        public void ChooseAvailableGameServer_PlayerCountOutOfRange_Test()
        {
            var mockService = new MockGameServerDistributeRefService();

            try
            {
                var result = mockService.ChooseAvailableGameServer(
               _TEST_DATA.Default.PlayerID,
               _TEST_DATA.Default.MerchantID,
               _TEST_DATA.Default.GameID,
               new ServerDistributeOptionalParameter
               {

               });

            }
            catch (Exception ex)
            {
                AssertExpectException(ex, ApiErrorCode.ServerPlayerOutOfRange);
            }

            return;
        }

    }
}

