using game.Entity.Game;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using NSubstitute;
using target.BL.Interface;
using target.DA.Repository;
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
using System.Threading.Tasks;
using NSubstitute.ExceptionExtensions;
using target.Models.Table;

namespace target.Testing
{
    [TestClass]
    public class GameBetServiceTests : BaseTest
    {

        public IGameBetService ServiceInstance { get; set; }



        [TestCleanup]
        public void TestDone()
        {
            //清空還原mock
            IOCHelper.MockInstanceMapping = new Dictionary<string, object>();
        }



        [TestMethod]
        public void PlaceBet_Test()
        {
            var result = ServiceInstance.PlaceBet(new GamePlaceBetRequestModel
            {
                ServerID = _TEST_DATA.Default.ServerID,
                IPAddress = _TEST_DATA.Default.IPAddress,
                BetAmount = 10,
                PlayerID = _TEST_DATA.Default.PlayerID,
                TradeID = _TEST_DATA.Default.TradeID,
                AccessJwtToken = ""
            });

            Assert.IsNotNull(result);
            return;
        }



        [TestMethod]
        public void PlaceBet_SingleWallet_SessionNotFound()
        {
            string sessionToken = [REDACTED];

            string jwtToken = [REDACTED]
                _TEST_DATA.SingleWallet.PlayerID,
                _TEST_DATA.SingleWallet.MerchantID,
                _TEST_DATA.SingleWallet.PlayerName
            );

            try
            {
                var result = ServiceInstance.PlaceBet(new GamePlaceBetRequestModel
                {
                    ServerID = _TEST_DATA.Default.ServerID,
                    IPAddress = _TEST_DATA.Default.IPAddress,
                    BetAmount = 10,
                    PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                    TradeID = _TEST_DATA.Default.TradeID,
                    AccessJwtToken = [REDACTED]
                });
            }
            catch (TerproException ex)
            {
                Assert.IsTrue(ex.ErrroCode == ApiErrorCode.NotFoundSessionToken);
            }

            return;
        }



        [TestMethod]
        public void PlaceBet_SingleWallet_Test()
        {
            string sessionToken = [REDACTED];

            string jwtToken = [REDACTED]
                _TEST_DATA.SingleWallet.PlayerID,
                _TEST_DATA.SingleWallet.MerchantID,
                _TEST_DATA.SingleWallet.PlayerName,
                sessionToken
            );

            //使用mock
            var substitute = Substitute.For<ISingleWalletPlatformApiRequestRefService>();

            substitute.PlaceBet(null).ReturnsForAnyArgs(new SingleWalletModePlaceBetApiResponseModel
            {
                Balance = _TEST_DATA.SingleWallet.BalanceAmount - 100,
                TransactionNo = RandomHelper.GetOrderNo("AA")
            });

            substitute.CheckIsValidSessionToken(null).ReturnsForAnyArgs(new SingleWalletMerchantSessionTokenEntity
            {
                PlayerDisplayName = _TEST_DATA.SingleWallet.PlayerName,
                PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                MerchantID = _TEST_DATA.Default.MerchantID
            });

            IOCHelper.AddMockInstance<ISingleWalletPlatformApiRequestRefService>(substitute);

            //若使用mockService 註冊service之後
            ServiceInstance = IOCHelper.Resolve<IGameBetService>();

            var result = ServiceInstance.PlaceBet(new GamePlaceBetRequestModel
            {
                ServerID = _TEST_DATA.Default.ServerID,
                IPAddress = _TEST_DATA.Default.IPAddress,
                BetAmount = 10,
                PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                TradeID = _TEST_DATA.Default.TradeID,
                AccessJwtToken = [REDACTED]
            });

            Assert.IsNotNull(result);
            return;
        }



        [TestMethod]
        public void PlaceBet_SingleWallet_Test_Fail()
        {
            //使用mock
            var substitute = Substitute.For<ISingleWalletPlatformApiRequestRefService>();
            substitute.PlaceBet(null).ThrowsForAnyArgs(new TerproException(ApiErrorCode.BetFail_WalletAmountNotEnough));

            IOCHelper.AddMockInstance<ISingleWalletPlatformApiRequestRefService>(substitute);

            //若使用mockService 註冊service之後
            ServiceInstance = IOCHelper.Resolve<IGameBetService>();

            string jwtToken = [REDACTED]
                _TEST_DATA.SingleWallet.PlayerID, 
                _TEST_DATA.SingleWallet.MerchantID, 
                _TEST_DATA.SingleWallet.PlayerName
            );

            try
            {
                var result = ServiceInstance.PlaceBet(new GamePlaceBetRequestModel
                {
                    ServerID = _TEST_DATA.Default.ServerID,
                    IPAddress = _TEST_DATA.Default.IPAddress,
                    BetAmount = 10,
                    PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                    TradeID = _TEST_DATA.Default.TradeID,
                    AccessJwtToken = [REDACTED]
                });
            }
            catch(TerproException ex)
            {
                Assert.IsTrue(ex.ErrroCode == ApiErrorCode.BetFail_WalletAmountNotEnough);
            }

            return;
        }

    }
}

