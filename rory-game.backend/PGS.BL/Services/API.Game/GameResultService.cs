using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Utility.Helper;
using target.Models;
using target.Models.Attributes;
using target.Utility;
using game.Entity.Game;

namespace target.BL.Services
{

    public class GameResultService : BaseGameService, IGameResultService
    {
        public PerGameLogRepository PerGameLogRepository { get; set; }

        public PerGameDetailLogRepository DetailLogRepository { get; set; }

        public PerGameDetailLogZeroRepository PerGameDetailLogZeroRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }

        public BetRecordRepository BetRecordRepository { get; set; }

        public IPlayerWalletRefService PlayerWalletRefService { get; set; }

        public ISingleWalletMerchantBetProcessRefService SingleWalletMerchantBetProcessRefService { get; set; }


        public int _robotSkipInsertMerchantID = -1;



        public GamePlayResultViewModel UpdateSettlementAmount(GamePlayResultRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();

            var queryPlayer = GetValidPlayerAndMerchant(model.PlayerID);

            var queryServer = ServerListRepository.GetBriefByCache(model.ServerID);

            if (queryServer == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            var queryGame = GameListRepository.GetSingleByCache(queryServer.GameID);

            //目前暫時沒有試玩模式
            WalletTypeEnum walletType = WalletTypeEnum.GameWallet;

            decimal updateGamePoint = model.UpdateAmount;
            decimal updateAmount = ConvertGamePointToMoney(queryPlayer.MerchantID , updateGamePoint);

            var updateModel = model.ReturnInjecteObject<WriteGameResulSPParameterModel>();
            updateModel.WalletType = walletType.ToInteger();
            updateModel.WalletMode = WalletTypeEnum.GameWallet.ToInteger();
            updateModel.GameID = queryServer.GameID;
            updateModel.UpdateAmount = updateAmount;

            bool isSingleWallet = CommonLogicHelper.IsSingleWalletMode(queryPlayer.WalletMode);
            bool isTransferMode = isSingleWallet == false;

            bool isOldVersionProcess = !model.IsPlaceBetProcess;

            //轉帳模式 可以記錄到更新錢金額
            if(isTransferMode)
            {
                var walletAmount = PlayerWalletRefService.GetWalletAmountByMerchantWalletMode(model.PlayerID, walletType);
                updateModel.CurrentAmount = walletAmount.Amount;
            }

            //舊版模式(v1)
            if (isOldVersionProcess)
            {             
                //更新Db錢包金額 並結束
                PlayerRepository.UpdatePlayerWalletAfterGameSettle(updateModel);

                decimal afterAmount = updateModel.CurrentAmount + updateAmount;
                decimal afterGamePoint = ConvertGamePointToMoney(queryPlayer.MerchantID, afterAmount);

                return new GamePlayResultViewModel
                {
                    AfterAmount = afterGamePoint
                };
            }

            //如果是v2版的結算邏輯 (單錢包模式的邏輯)
            var betInfo = model.BetInfo;

            betInfo.RuleFor(z => z.BetAmount).Rule_NotMinus();
            betInfo.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();

            string reqeustUUID = Guid.NewGuid().ToString();
            betInfo.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();

            var queryBetRecord = BetRecordRepository.GetByTransactionNo(betInfo.TransactionNo);

            if (queryBetRecord == null)
            {
                throw new TerproException(ApiErrorCode.NotFoundRecord, new ErrorDetailModel
                {
                    IsAlwaysNotify = true,
                    LogExtraInfo = $"TransactionNo: {betInfo.TransactionNo}"
                });
            }

            decimal betGamePoint = betInfo.BetAmount;
            decimal payoutGamePoint = betInfo.PayoutAmount;
            decimal validGamePoint = betInfo.ValidAmount;
            decimal taxGamePoint = betInfo.TaxAmount;

            decimal betAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, betGamePoint);
            decimal payoutAmount = ConvertGamePointToMoney(queryPlayer.MerchantID , payoutGamePoint);
            decimal validAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, validGamePoint);
            decimal taxAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, taxGamePoint);


            DateTime now = DateTimeHelper.GetCurrentTime();
            queryBetRecord.Status = GameVendorBetRecordStatusEnum.SettleComplete.ToInteger();
            queryBetRecord.UpdateTime = now;
            queryBetRecord.SettlementTime = now;
            queryBetRecord.WinLossAmount = payoutAmount - betAmount - taxAmount;
            queryBetRecord.PayoutAmount = payoutAmount;
            queryBetRecord.ValidAmount = validAmount;
            queryBetRecord.TaxAmount = taxAmount;
            queryBetRecord.ProcessRequestUUID = reqeustUUID;


            // Step 1. 更新betRecord
            BetRecordRepository.Update(queryBetRecord);


            //若投注金額為0元，會走bonus結算流程
            bool isFreeBonus = queryBetRecord.Type == BetRecordTypeEnum.FreeBonus.ToInteger();

            bool isFreeGameWithoutMoneyChanged = isFreeBonus && model.BetInfo.PayoutAmount == 0;

            // Case 1. FreeGame結算、無中獎
            if (isFreeGameWithoutMoneyChanged)
            {
                var walletResult = PlayerWalletRefService.GetWalletAmountByLatestToken(model.PlayerID);
                return new GamePlayResultViewModel
                {
                    AfterAmount = ConvertMoneyToGamePoint(queryPlayer.MerchantID , walletResult.Amount)
                };
            }

            // Case 2. 轉帳模式
            if(!isSingleWallet)
            {
                // [轉帳模式] 更新Db錢包金額
                PlayerRepository.UpdatePlayerWalletAfterGameSettle(updateModel);
                return new GamePlayResultViewModel
                {
                    AfterAmount = ConvertMoneyToGamePoint(queryPlayer.MerchantID, updateModel.CurrentAmount + updateAmount)
                };
            }


            // Case 3. 單錢包模式-一般投注
            if (!isFreeBonus)
            {
                // [單錢包模式] 呼叫對方API
                var swResult = SingleWalletMerchantBetProcessRefService.SettlemenetProcessPassToMerchantApi(new SingleWalletModeSettleApiParameter
                {
                    PlayerID = model.PlayerID,
                    BetAmount = betAmount,
                    PayoutAmount = payoutAmount,
                    TaxAmount = taxAmount,
                    ValidAmount = validAmount,
                    MerchantID = queryPlayer.MerchantID,
                    GameID = queryServer.GameID,
                    BetRecordNoPGS = model.BetInfo.TransactionNo,
                    SettlementTime = now,
                    DebugProcessUUID = reqeustUUID
                });

                return new GamePlayResultViewModel
                {
                    AfterAmount = ConvertMoneyToGamePoint(queryPlayer.MerchantID, swResult.Balance)
                };
            }

            // Case 4. 單錢包模式-FreeGame/Bonus
            var bonusResult = SingleWalletMerchantBetProcessRefService.BonusProcessPassToMerchantApi(new SingleWalletModeSettleApiParameter
            {                                
                BetAmount = 0,
                PlayerID = queryBetRecord.PlayerID,
                PayoutAmount = payoutAmount,
                TaxAmount = taxAmount,
                ValidAmount = validAmount,
                MerchantID = queryPlayer.MerchantID,
                GameID = queryServer.GameID,
                BetRecordNoPGS = model.BetInfo.TransactionNo,
                SettlementTime = now,
                DebugProcessUUID = reqeustUUID
            });

            return new GamePlayResultViewModel
            {
                AfterAmount = ConvertMoneyToGamePoint(queryPlayer.MerchantID, bonusResult.Balance)
            };
        }


        /// <summary>
        /// 寫每局遊戲紀錄
        /// </summary>
        public PerGameLogViewModel WritePerGameLog(PerGameLogRequestModel model)
        {
            var perGameLogModel = model.PerGameLog;
            model.RuleFor(x => x.PerGameLog).Rule_ModelNotNull();
            perGameLogModel.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            perGameLogModel.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            perGameLogModel.RuleFor(x => x.GameData).Rule_NotEmptyText();
            perGameLogModel.RuleFor(x => x.GameBeginTime).Rule_NotEmptyDateTime();
            perGameLogModel.RuleFor(x => x.GameEndTime).Rule_NotEmptyDateTime();

            foreach (var item in model.PerGameLog.UserGameDetailLogs.NullToEmptyCollection())
            {
                item.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
                item.RuleFor(x => x.TradeID).Rule_NotZeroNotMinus();
            }

            var queryGame = GameListRepository.GetSingleByCache(perGameLogModel.GameID);

            if (queryGame == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            decimal totalTaxGamePoint = model.PerGameLog.TaxAmount;
            decimal totalTaxAmount = model.PerGameLog.TaxAmount; //這邊不管了...難以處理

            //step 1. 寫入perGame
            long perGameID = PerGameLogRepository.Create(new PerGameLogEntity
            {
                BeginTime = model.PerGameLog.GameBeginTime,
                EndTime = model.PerGameLog.GameEndTime,
                GameData = model.PerGameLog.GameData,
                ServerID = model.PerGameLog.ServerID,
                GameID = model.PerGameLog.GameID,
                TaxAmount = model.PerGameLog.TaxAmount,
                CreateTime = DateTimeHelper.GetCurrentTime()
            });

            List<PerGameDetailLogEntity> detailLogs = new List<PerGameDetailLogEntity>();
            List<PerGameDetailLogZeroEntity> detailLogsZero = new List<PerGameDetailLogZeroEntity>();

            //舊版寫法(理論上已沒用處了) 傳入字串、拆解解析字串 轉成 PerGameDetailLog
            detailLogs = OldVersionAssignPerGameLogsValues(model, perGameID);

            bool isAlwaysWriteLog = !queryGame.IsSkipInsertZeroBet;
            bool IsSkipInserZeroBet = queryGame.IsSkipInsertZeroBet;

            //新版寫法 傳入array物件
            if (model.PerGameLog.UserGameDetailLogs.IsNotEmptyNull())
            {
                //若此gameID有定義 IsSkipInserZeroBet 的話，當投注行為(掛網觀局)的資料會寫到另一張表
                var logs = model.PerGameLog.UserGameDetailLogs.Where(x => isAlwaysWriteLog || !(x.BetAmount == 0 && x.UpdateAmount == 0)).ToList();
                var logs_zero = model.PerGameLog.UserGameDetailLogs.Where(x => IsSkipInserZeroBet && x.BetAmount == 0 && x.UpdateAmount == 0).ToList();

                foreach (var item in logs)
                {
                    var queryPlayer = PlayerRepository.GetSingleByCache(item.PlayerID);

                    if (queryPlayer == null)
                        throw new TerproException(ApiErrorCode.NotFoundPlayer);

                    decimal betAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, item.BetAmount);
                    decimal payoutAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, item.PayoutAmount);
                    decimal taxAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, item.TaxAmount);
                    decimal validAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, item.ValidAmount);
                    decimal currentAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, item.CurrentAmount);
                    decimal updateAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, item.UpdateAmount);

                    detailLogs.Add(new PerGameDetailLogEntity
                    {
                        GameID = model.PerGameLog.GameID,
                        ServerID = model.PerGameLog.ServerID,
                        CreateTime = DateTimeHelper.GetCurrentTime(),
                        PlayerID = item.PlayerID,
                        TradeID = item.TradeID,

                        ValidAmount = validAmount,
                        CurrentAmount = currentAmount,
                        UpdateAmount = updateAmount,
                        BetAmount = betAmount,
                        PayoutAmount = payoutAmount,
                        TaxAmount = taxAmount,
                        ConvertRate = GetMerchantConvertRate_PGSGame(queryPlayer.MerchantID),

                        GameTime = model.PerGameLog.GameEndTime,
                        MerchantID = 0,
                        PerGameLogID = perGameID,
                        TransactionNo = item.TransactionNo,
                    });
                }

                foreach (var item in logs_zero)
                {
                    detailLogsZero.Add(new PerGameDetailLogZeroEntity
                    {
                        GameID = model.PerGameLog.GameID,
                        ServerID = model.PerGameLog.ServerID,
                        CreateTime = DateTimeHelper.GetCurrentTime(),
                        CurrentAmount = item.CurrentAmount,
                        UpdateAmount = item.UpdateAmount,
                        PlayerID = item.PlayerID,
                        TradeID = item.TradeID,
                        ValidAmount = item.ValidAmount,
                        BetAmount = item.BetAmount,
                        PayoutAmount = item.PayoutAmount,
                        GameTime = model.PerGameLog.GameEndTime,
                        MerchantID = 0,
                        PerGameLogID = perGameID,
                    });
                }
            }

            //補上 merchantIDs & playerName 資訊
            SetMerchantIDAndPlayerName(detailLogs);
            SetMerchantIDAndPlayerName(detailLogsZero);


            // 經討論後決定機器人不寫入detailLog資料 (2020/08/31)
            detailLogs = detailLogs.Where(x => x.MerchantID != _robotSkipInsertMerchantID).ToList();
            detailLogsZero = detailLogsZero.Where(x => x.MerchantID != _robotSkipInsertMerchantID).ToList();


            //step 2. 寫入detailLog
            DetailLogRepository.CreateMany(detailLogs);

            //step 3. 寫入detailLogZero (觀局資料)
            PerGameDetailLogZeroRepository.CreateMany(detailLogsZero);

            return new PerGameLogViewModel
            {
                PerGameLogID = perGameID
            };
        }

        private List<PerGameDetailLogEntity> OldVersionAssignPerGameLogsValues(PerGameLogRequestModel model, long perGameID)
        {
            List<string> splitEachUser = new List<string>();

            if (!string.IsNullOrWhiteSpace(model.PerGameLog.UserGameDetail))
                splitEachUser = model.PerGameLog.UserGameDetail.Split(',').Where(x => !string.IsNullOrWhiteSpace(x)).ToList();

            var list = new List<PerGameDetailLogEntity>();

            foreach (var each in splitEachUser)
            {
                var split = each.Split(':');

                int playerID = Int32.Parse(split[0]);
                decimal currentAmount = decimal.Parse(split[1]);
                decimal updateAmount = decimal.Parse(split[2]);
                long tradeID = Int64.Parse(split[3]);

                list.Add(new PerGameDetailLogEntity
                {
                    GameID = model.PerGameLog.GameID,
                    ServerID = model.PerGameLog.ServerID,
                    CreateTime = DateTimeHelper.GetCurrentTime(),
                    CurrentAmount = currentAmount,
                    UpdateAmount = updateAmount,
                    PlayerID = playerID,
                    TradeID = tradeID,
                    GameTime = model.PerGameLog.GameEndTime,
                    MerchantID = 0,
                    PerGameLogID = perGameID,
                });
            }

            return list;
        }


        private void SetMerchantIDAndPlayerName<T>(List<T> models) where T : IPlayerRowDataColumn
        {
            List<int> playerIDs = models.Select(x => x.PlayerID).ToList();

            if (playerIDs.IsEmptyOrNull())
                return;

            List<PlayerEntity> queryPlayers = GetPlayersByCache(playerIDs);

            //如果完全查不到players 直接拋exception
            if (queryPlayers.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            foreach (var item in models)
            {
                var match = queryPlayers.FirstOrDefault(x => x.PlayerID == item.PlayerID);

                if (match != null)
                {
                    //若判斷此筆資料是機器人
                    if (match.PlayerType == (int)PlayerTypeEnum.Robot)
                    {
                        //這邊填 -1 後面會做些處理
                        item.MerchantID = _robotSkipInsertMerchantID;
                    }
                    else
                    {
                        item.MerchantID = match.MerchantID;
                        item.PlayerName = match.PlayerName;
                    }
                }
            }
        }


        /// <summary>
        /// 如果只有一筆資料要查，使用cache查 (slot應該是這種)
        /// 減少DB存取
        /// </summary>
        private List<PlayerEntity> GetPlayersByCache(List<int> playerIDs)
        {
            if (playerIDs.IsEmptyOrNull())
                return new List<PlayerEntity>();

            playerIDs = playerIDs.Distinct().ToList();

            bool isSingle = playerIDs.Count() == 1;

            if (isSingle)
            {
                var querySingle = PlayerRepository.GetSingleByCache(playerIDs.FirstOrDefault());
                return new List<PlayerEntity>
                {
                    querySingle
                };
            }
            else
            {
                var queryPlayers = PlayerRepository.GetByIds(playerIDs);
                return queryPlayers;
            }

        }


    }
}
