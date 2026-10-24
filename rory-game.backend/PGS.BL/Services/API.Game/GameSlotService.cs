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
using target.Utility;
using System.Data.SqlClient;
using System.Threading;

namespace target.BL.Services
{

    public partial class GameSlotService : BaseGameService, IGameSlotService
    {

        public GameListRepository GameListRepository { get; set; }
        public JackpotPoolRepository JackpotPoolRepository { get; set; }
        public GameSlotProbabilitySettingRepository GameSlotProbabilitySettingRepository { get; set; }
        public GameSlotDeferredAwardRepository GameSlotDeferredAwardRepository { get; set; }
        public GameSlotPlayerAssignedModeRepository GameSlotAssignedModeRepository { get; set; }
        public PlayerSlotBonusRecordRepository PlayerSlotBonusRecordRepository { get; set; }
        public PlayerSlotFreeSpinRecordRepository PlayerSlotFreeSpinRecordRepository { get; set; }
        public PlayerSlotLightningRecordRepository PlayerSlotLightningRecordRepository { get; set; }
        public SlotGameHistoryLogRepository SlotGameHistoryLogRepository { get; set; }



        public List<GameTypeEnum> JackpotAvailableTypes = new List<GameTypeEnum>
        {
            GameTypeEnum.BoardGame,
            GameTypeEnum.Slot
        };


        public GameJackpotViewModel UpdateJackpot(int gameID, GameSlotJackpotUpdateRequestModel model)
        {
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var queryGameList = GameListRepository.GetSingleByCache(gameID);

            if (queryGameList == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            GameTypeEnum gameType = (GameTypeEnum)queryGameList.GameType;

            bool isDefineGameType = JackpotAvailableTypes.Contains(gameType);

            if(!isDefineGameType)
            {
                throw new TerproException(ApiErrorCode.NotImplement);
            }

            string lockKey = $"UpdateJackpot_{gameID}_{model.FlagID}_{model.TableNo}_{model.IsLink}";
                             
            //使用lock，同時間只允許處理一筆 相同鍵值的資料，避免時間差/互相覆蓋的問題
            lock(ThreadLockHelper.InternConcurrentSafe(lockKey))
            {
                var filter = new QueryJackpotPoolRequestViewModel
                {
                    GameID = gameID,
                    FlagID = model.FlagID,
                    TableNo = model.TableNo,
                    IsLink = model.IsLink
                };

                var queryAllLevels = JackpotPoolRepository.GetByFilter(filter).ToList();
                var queryRecord = queryAllLevels.Where(x => x.LevelID == model.LevelID).FirstOrDefault();

                if (queryRecord == null)
                {
                    JackpotPoolEntity newRecord = new JackpotPoolEntity
                    {
                        GameID = gameID,
                        FlagID = model.FlagID,
                        TableNo = model.TableNo,
                        LevelID = model.LevelID,
                        Contribution = model.Contribution,
                        TotalContribution = model.TotalContribution,
                        TotalAward = model.TotalAward,
                        IsLink = model.IsLink
                    };
                    JackpotPoolRepository.Create(newRecord);
                    queryAllLevels.Add(newRecord);
                }
                else
                {
                    queryRecord.Contribution += model.Contribution;
                    queryRecord.TotalContribution += model.TotalContribution;
                    queryRecord.TotalAward += model.TotalAward;
                    JackpotPoolRepository.Update(queryRecord);
                }

                //返回相關資料
                filter.LevelID = null;
                var result = new GameJackpotViewModel
                {
                    Jackpots = queryAllLevels.Select(x => new GameJackpotItemViewModel
                    {
                        GameID = x.GameID,
                        FlagID = x.FlagID,
                        LevelID = x.LevelID,
                        TableNo = x.TableNo,
                        Contribution = x.Contribution,
                        TotalContribution = x.TotalContribution,
                        TotalAward = x.TotalAward
                    }).ToList()
                };
                return result;
            }
        }


        public GameSlotDeferredAwardViewModel GetSlotDeferredAward(GameSlotDeferredAwardRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.AwardType).Rule_NotEmptyText();

            var filter = new QueryPlayerSlotDeferredAwardRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID,
                AwardType = model.AwardType
            };

            var query = GameSlotDeferredAwardRepository.GetByFilter(filter);

            var result = new GameSlotDeferredAwardViewModel
            {
                PlayerSlotDeferredAward = query.Select(x => new GameSlotDeferredAwardItemViewModel
                {
                    PlayerID = x.PlayerID,
                    GameID = x.GameID,
                    FlagID = x.FlagID,
                    OriginID = x.OriginID,
                    AwardType = x.AwardType,
                    CreateTime = x.CreateTime,
                    PatternID = x.PatternID,
                    Data = x.Data
                }).ToList()
            };

            return result;
        }

        public void AddSlotDeferredAward(GameSlotDeferredAwardUpdateRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotMinus();
            model.RuleFor(x => x.OriginID).Rule_NotMinus();
            model.RuleFor(x => x.AwardType).Rule_NotEmptyText();

            GameSlotDeferredAwardRepository.Create(new GameSlotDeferredAwardEntity
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID,
                OriginID = model.OriginID,
                AwardType = model.AwardType,
                Data = model.Data,
                PatternID = model.PatternID,
                CreateTime = DateTimeHelper.GetCurrentTime()
            });
        }

        public GameSlotDeferredAwardItemViewModel DeleteSlotDeferredAward(GameSlotDeferredAwardRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotMinus();
            model.RuleFor(x => x.AwardType).Rule_NotEmptyText();

            var filter = new QueryPlayerSlotDeferredAwardRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID,
                AwardType = model.AwardType
            };
            var query = GameSlotDeferredAwardRepository.GetByFilter(filter).FirstOrDefault();

            if (query == null)
                return new GameSlotDeferredAwardItemViewModel();

            var result = query.ReturnInjecteObject<GameSlotDeferredAwardItemViewModel>();
            GameSlotDeferredAwardRepository.Delete(query.ID);

            return result;
        }

        private GameSlotHistoryLogViewModel ConvertToViewModels(int merchantID , SlotGameHistoryLogEntity query)
        {
            var convert = query.ReturnInjecteObject<GameSlotHistoryLogViewModel>();
            convert.LogID = query.ID;
            convert.FinalAmount = query.FinalAmount - query.CurrentAmount;

            //顯示給玩家看到的 還是 [遊戲幣]
            convert.FinalAmount = ConvertMoneyToGamePoint(merchantID , convert.FinalAmount);
            convert.BetAmount = ConvertMoneyToGamePoint(merchantID, convert.BetAmount);
            return convert;
        }
    }
}
