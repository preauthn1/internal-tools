using game.Entity.Game;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using NSubstitute;
using NSubstitute.ExceptionExtensions;
using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class SingleWalletFailProcessTests : BaseTest
    {


        [TestCleanup]
        public void TestDone()
        {
            //清空還原mock
            IOCHelper.MockInstanceMapping = new Dictionary<string, object>();
        }



        [TestMethod]
        public void SingleWallet_PlaceBet_Fail_Test()
        {
            string jwt = GenerateFakeJwt();

            //使用mock
            MockServiceMethod(jwt);

            decimal betAmount = 20;

            //使用mockInstance的話，instance service必須在註冊mock之後
            var GameBetService = IOCHelper.Resolve<IGameBetService>();

            var expectCode = ApiErrorCode.Success;

            try
            {
                var result = GameBetService.PlaceBet(new GamePlaceBetRequestModel
                {
                    ServerID = _TEST_DATA.Default.ServerID,
                    AccessJwtToken = [REDACTED],
                    IPAddress = _TEST_DATA.Default.IPAddress,
                    BetAmount = betAmount,
                    PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                    TradeID = _TEST_DATA.Default.TradeID
                });

            }
            catch(TerproException ex)
            {
                expectCode = ex.ErrroCode;
            }

            Assert.IsTrue(expectCode == ApiErrorCode.BetFail_WalletAmountNotEnough);
            return;
        }



        private void MockServiceMethod(string jwt)
        {
            //使用mock
            var substitute = Substitute.For<ISingleWalletPlatformApiRequestRefService>();
          
            substitute.PlaceBet(null).ThrowsForAnyArgs(new TerproException(ApiErrorCode.BetFail_WalletAmountNotEnough));

            IOCHelper.AddMockInstance<ISingleWalletPlatformApiRequestRefService>(substitute);
        }


        private string GenerateFakeJwt()
        {
            var resultJWT = MyJwtHelper.GeneratePlayerJwt(
                _TEST_DATA.SingleWallet.PlayerID,
                _TEST_DATA.SingleWallet.MerchantID,
                _TEST_DATA.SingleWallet.PlayerName,
                SystemConstant.DefaultLang,
                Guid.NewGuid().ToString(),
                Guid.NewGuid().ToString()
             );

            return resultJWT.JwtToken;
        }


    }
}

