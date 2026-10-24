using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using target.DA.Repository;

namespace target.Testing
{
    [TestClass]
    public class GamePlayServiceTests : BaseTest
    {

        public IGameRoomService GamePlayService { get; set; }


        [TestMethod]
        public void GameRoomEnter_And_RoomOut_Test()
        {
            int serverID = _TEST_DATA.Default.ServerID;
            int playerID = _TEST_DATA.Not_Super_User.PlayerID;
            int merchantID = _TEST_DATA.Not_Super_User.MerchantID;
            string playerName = _TEST_DATA.Not_Super_User.PlayerName;

            string jwt = GeneatePlayerJwtToken(playerID, merchantID, playerName);

            //Step 1. 進房
            var result = GamePlayService.GameRoomEnter(playerID, new GameRoomEnterRequestModel
            {
                AccessToken = [REDACTED],
                IPAddress = _TEST_DATA.Default.IPAddress,
                ServerID = serverID,
                IsFishHunter = false
            });

            Assert.IsNotNull(result);
            Assert.IsNotNull(result.PlayerInfo);
            Assert.IsTrue(result.TradeID >= 1000000);

            //Step 2. 退房
            var result2 = GamePlayService.GameRoomLeave(new GameRoomLeaveRequestModel
            {
                PlayerID = playerID,
                DrawCount = 1,
                LoseCount = 3,
                WinCount = 4,
                ServerID = serverID,
                IPAddress = _TEST_DATA.Default.IPAddress,
                LoginTime = DateTimeHelper.GetCurrentTime(),
                LogoutTime = DateTimeHelper.GetCurrentTime(),
                UpdateAmount = 10,
                FlowAmount = 10,
                TaxAmount = 1,
                LastWinLose = 1,
                GameID = _TEST_DATA.Default.GameID,
                TradeID = result.TradeID
            });

            Assert.IsNotNull(result2);
            return;
        }

    }
}

