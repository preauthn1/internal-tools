using game.Entity.Game;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using NSubstitute;
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
    public class SingleWalletProcessTests : BaseTest
    {

        [TestCleanup]
        public void TestDone()
        {
            //清空還原mock
            IOCHelper.MockInstanceMapping = new Dictionary<string, object>();
        }


        [TestMethod]
        public void GetWalletAmount_Test()
        {
            string jwt = GenerateFakeJwt();
            MockServiceMethod(jwt);

            //使用mockInstance的話，instance service必須在註冊mock之後
            var PlayerWalletInfoService = IOCHelper.Resolve<IPlayerWalletInfoService>();
            var result = PlayerWalletInfoService.GetWalletAmount(_TEST_DATA.SingleWallet.PlayerID, jwt);

            Assert.IsTrue(result.Amount == _TEST_DATA.SingleWallet.BalanceAmount);
            ResetMockRegister();
            return;
        }


        [TestMethod]
        public void SingleWallet_PlaceBet_Test()
        {
            string jwt = GenerateFakeJwt();

            //使用mock
            MockServiceMethod(jwt);

            decimal betAmount = 20;

            //使用mockInstance的話，instance service必須在註冊mock之後
            var GameBetService = IOCHelper.Resolve<IGameBetService>();
            var result = GameBetService.PlaceBet(new GamePlaceBetRequestModel
            {
                ServerID = _TEST_DATA.Default.ServerID,
                AccessJwtToken = [REDACTED],
                IPAddress = _TEST_DATA.Default.IPAddress,
                BetAmount = betAmount,
                PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                TradeID = _TEST_DATA.Default.TradeID
            });

            Assert.IsTrue(_TEST_DATA.SingleWallet.BalanceAmount > result.AfterAmount);
            ResetMockRegister();
            return;
        }



        [TestMethod]
        public void SingleWallet_UpdateSettlementAmount_Test()
        {
            string jwt = GenerateFakeJwt();

            //使用mock
            MockServiceMethod(jwt);

            string transactionNo = _TEST_DATA.Default.TransactionNo_PgsBet;

            //使用mockInstance的話，instance service必須在註冊mock之後
            var GameResultService = IOCHelper.Resolve<IGameResultService>();
            var result = GameResultService.UpdateSettlementAmount(new GamePlayResultRequestModel
            {
                ServerID = _TEST_DATA.Default.ServerID,
                IPAddress = _TEST_DATA.Default.IPAddress,
                PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                TradeID = _TEST_DATA.Default.TradeID,
                IsPlaceBetProcess = true,
                UpdateAmount = 10,
                BetInfo = new SettlementBetInfoRequestModel
                {
                    BetAmount = 10,
                    PayoutAmount = 10,
                    TransactionNo = transactionNo
                }
            });

            ResetMockRegister();
            return;
        }




        private void MockServiceMethod(string jwt)
        {
            //使用mock
            var substitute = Substitute.For<ISingleWalletPlatformApiRequestRefService>();
            substitute.GetWalletBalance(null).ReturnsForAnyArgs(new SingleWalletModeWalletApiResponseModel
            {
                Balance = _TEST_DATA.SingleWallet.BalanceAmount
            });

            substitute.PlaceBet(null).ReturnsForAnyArgs(new SingleWalletModePlaceBetApiResponseModel
            {
                Balance = _TEST_DATA.SingleWallet.BalanceAmount - 100,
                TransactionNo = RandomHelper.GetOrderNo("AA")
            });

            substitute.Settlement(null).ReturnsForAnyArgs(new SingleWalletModeSettleApiResponseModel
            {
                Balance = _TEST_DATA.SingleWallet.BalanceAmount,
                TransactionNo = RandomHelper.GetOrderNo("AA")
            });

            IOCHelper.AddMockInstance<ISingleWalletPlatformApiRequestRefService>(substitute);
        }


        private void ResetMockRegister()
        {
            //清空其他測試注冊的mock
            IOCHelper.MockInstanceMapping = new Dictionary<string, object>();
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

