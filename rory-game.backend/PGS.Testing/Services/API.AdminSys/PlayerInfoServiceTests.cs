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
using game.Entity.Game;

namespace target.Testing
{
    [TestClass]
    public class PlayerInfoServiceTests : BaseTest
    {

        public IPlayerInfoService ServiceInstance { get; set; }


        [TestMethod]
        public void PlayerInfo_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryPlayerApiRequestParameterModel { PageSize = 10});
            Assert_NotEmptyCollection(query.Items);
            return;
        }

        
        [TestMethod]
        public void PlayerInfo_GetPlayerSummaryInfo_Test()
        {
            var query = ServiceInstance.GetPlayerSummaryInfo(_TEST_DATA.Default.PlayerID);
            Assert.IsNotNull(query);
            return;
        }



        [TestMethod]
        public void UpdatePlayerStatus_Test()
        {
            ServiceInstance.UpdatePlayerStatus(_TEST_DATA.Default.UserID , new PlayerUpdateStatusRequestModel
            {
                Memo = "unitTest",
                Status = PlayerStatusEnum.Enable.ToInteger(),
                PlayerID = _TEST_DATA.Default.PlayerID
            });

            return;
        }


        [TestMethod]
        public void ManualDepositPlayerWallet_Test()
        {
            var result = ServiceInstance.ManualDepositPlayerWallet(_TEST_DATA.Default.UserID,  new PlayerManualDepositRequestModel
            {
                Amount = 10,
                PlayerID = _TEST_DATA.Default.PlayerID,
                PlayerName = _TEST_DATA.Default.PlayerName,
                Memo = "MSTest"
            });

            Assert.IsTrue(result.BeforeAmount != 0 || result.AfterAmount != 0);
            return;
        }


        [TestMethod]
        public void InternalTestAccountWalletDepositAmount_Test()
        {
            var result = ServiceInstance.InternalTestAccountWalletDepositAmount(_TEST_DATA.Default.UserID, new InternalTestAccountPlusWalletRequestModel
            {
                Amount = 10,
                PlayerID = _TEST_DATA.Default.PlayerID
            });

            Assert.IsTrue(result.BeforeAmount != 0 || result.AfterAmount != 0);
            return;
        }



    }
}

