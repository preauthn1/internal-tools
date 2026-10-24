using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using target.DA.Repository;
using System.Collections.Generic;

namespace target.Testing
{
    [TestClass]
    public class GameResultServiceTests : BaseTest
    {

        public IGameResultService GameResultService { get; set; }


        [TestMethod]
        public void WriteSettlementResult_Test()
        {
            var result = GameResultService.UpdateSettlementAmount(new GamePlayResultRequestModel
            {
                PlayerID = _TEST_DATA.Not_Super_User.PlayerID,
                ServerID = _TEST_DATA.Default.ServerID,
                TradeID = _TEST_DATA.Default.TradeID,
                UpdateAmount = 10,
            });

            Assert.IsNotNull(result);
            return;
        }


        [TestMethod]
        public void UpdateSettlementAmount_V2_Test()
        {
            var result = GameResultService.UpdateSettlementAmount(new GamePlayResultRequestModel
            {
                PlayerID = _TEST_DATA.Not_Super_User.PlayerID,
                ServerID = _TEST_DATA.Default.ServerID,
                TradeID = _TEST_DATA.Default.TradeID,
                UpdateAmount = 10,
                IsPlaceBetProcess = true,
                IPAddress = _TEST_DATA.Default.IPAddress,
                BetInfo = new SettlementBetInfoRequestModel
                {
                    TransactionNo = _TEST_DATA.Default.TransactionNo_PgsBet,
                    BetAmount = 10,
                    PayoutAmount = 10,
                    ValidAmount = 10
                }
            });

            Assert.IsNotNull(result);
            return;
        }


        [TestMethod]
        public void WritePerGameLog_Test()
        {
            decimal amount = _TEST_DATA.Default.BetAmount_MustDelete;
            var result = GameResultService.WritePerGameLog(new PerGameLogRequestModel
            {
                PerGameLog = new PerGameLogItemRequestModel
                {
                    TaxAmount = 0,
                    GameData = "遊戲歷程記錄blabla",
                    GameID = _TEST_DATA.Default.GameID,
                    GameBeginTime = DateTimeHelper.GetCurrentTime(),
                    GameEndTime = DateTimeHelper.GetCurrentTime().AddMinutes(2),
                    ServerID = _TEST_DATA.Default.ServerID,
                    UserGameDetailLogs = new List<PerGameDetailLogInfoViewModel>
                    {
                        new PerGameDetailLogInfoViewModel
                        {
                            BetAmount = amount,
                            CurrentAmount = amount,
                            UpdateAmount = amount,
                            ValidAmount = amount ,
                            PayoutAmount = amount,
                            PlayerID = _TEST_DATA.Not_Super_User.PlayerID,
                            TradeID = _TEST_DATA.Default.TradeID_MustDelete
                        }
                    },
                    UserGameDetail = $"{_TEST_DATA.Not_Super_User.PlayerID}:{amount}:{amount}:{_TEST_DATA.Default.TradeID_MustDelete}"
                }
            });

            Assert.IsNotNull(result);
            return;
        }



        [TestMethod]
        public void PerGameLog_DeleteTestData_Test()
        {
            var rep_per = InstanceRepository<PerGameLogRepository>();
            var rep_detail = InstanceRepository<PerGameDetailLogRepository>();
            
            var query = rep_detail.GetByFilter(new QueryDetailLogRequestViewModel
            {
                PlayerIDs = _TEST_DATA.Not_Super_User.PlayerID.ToListCollection()
            });

            if (query.IsEmptyOrNull())
                return;

            List<long> perGameLogIDs = query.Select(x => x.PerGameLogID).ToList();


            //刪除測試產生的資料 (初期frank測試會造成困惑、先不執行此處理)

            //foreach(var perGameID in perGameLogIDs)
            //{
            //    rep_per.Delete(perGameID);
            //}

            //foreach (var item in query)
            //{
            //    rep_detail.Delete(item.Id);
            //}

            return;
        }




    }
}

