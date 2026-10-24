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
using System.Threading;

namespace target.BL.Services
{

    public partial class GameSlotService
    {

        private static List<string> _list = new List<string>();

        #region SlotGameLog

        public PagedList<GameSlotHistoryLogViewModel> GetSlotGameLog(GameSlotLogRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();

            DateTime defaultTime = DateTimeHelper.GetDefaultTime();
            model.GameTimeBegin = (!model.GameTimeBegin.HasValue)
                ? DateTimeHelper.GetCurrentTime().AddDays(-7).To_00_00_00()
                : (model.GameTimeBegin <= defaultTime) ? defaultTime : model.GameTimeBegin;

            model.GameTimeEnd = (!model.GameTimeEnd.HasValue)
                ? DateTimeHelper.GetCurrentTime().AddDays(1).To_23_59_59()
                : (model.GameTimeEnd <= defaultTime) ? defaultTime : model.GameTimeEnd;

            var filter = new QuerySlotGameHistoryLogRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                GameTimeBegin = model.GameTimeBegin,
                GameTimeEnd = model.GameTimeEnd,
                PageIndex = model.PageIndex,
                PageSize = model.PageSize,
                OrderByColumn = model.OrderByColumn,
                IsASC = model.IsASC
            };

            var queryPlayer = PlayerRepository.GetSingleByCache(model.PlayerID);

            PagedList<SlotGameHistoryLogEntity> query = SlotGameHistoryLogRepository.GetPagingByFilter(filter);

            var paging = ConvertToNewTypePagingObject<SlotGameHistoryLogEntity, GameSlotHistoryLogViewModel>(query);
            var rowNum = (model.PageIndex - 1) * model.PageSize;
            paging.Items = query.Items
                .Select(x =>
                {
                    var convert = ConvertToViewModels(queryPlayer.MerchantID , x);
                    convert.RowNum = ++rowNum;
                    return convert;
                })
                .ToList();

            return paging;
        }



        public void AddSlotGameLog(GameSlotHistoryLogRequestModel model)
        {
            /**
             * 單一筆資料格式：CreateTime; GameID; ServerID; UserID; CurrentAmount; BetAmount; WinAmount; FinalAmount; GameType; LogType; GameData
             * 單一筆資料範例：2019-12-31 10:00:55.960;1005;10104;15267008;1057380.00;280.00;0.00;1057100.00;10;1;p7,p9,p4,p9|p5,p7,p11,p5|p11,p9,p8,p11|p10,p3,p7,p11|p8,p8,p9,p7
             */
            List<SlotGameHistoryLogEntity> historyLog = new List<SlotGameHistoryLogEntity>();

            if (model.Data.IsNotEmptyNull())
            {
                try
                {
                    DateTime defaultTime = DateTimeHelper.GetDefaultTime();
                    DateTime now = DateTimeHelper.GetCurrentTime();
                    DateTime gameTime;

                    foreach (var item in model.Data)
                    {
                        List<string> rowData = item.Split(';').ToList();
                        gameTime = (rowData[1].IsNotEmptyNull()) ? DateTimeOffset.Parse(rowData[0]).UtcDateTime : defaultTime;

                        if (gameTime <= defaultTime)
                        {
                            gameTime = defaultTime;
                        }

                        int playerID = Int32.Parse(rowData[3]);

                        var queryPlayer = PlayerRepository.GetSingleByCache(playerID);

                        if (queryPlayer == null)
                            throw new TerproException(ApiErrorCode.NotFoundPlayer);

                        historyLog.Add(new SlotGameHistoryLogEntity
                        {
                            GameTime = gameTime,
                            GameID = Int32.Parse(rowData[1]),
                            ServerID = Int32.Parse(rowData[2]),
                            PlayerID = playerID,
                            CurrentAmount = ConvertGamePointToMoney(queryPlayer.MerchantID , Decimal.Parse(rowData[4])),
                            BetAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, Decimal.Parse(rowData[5])),
                            WinAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, Decimal.Parse(rowData[6])),
                            FinalAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, Decimal.Parse(rowData[7])),
                            GameType = Int32.Parse(rowData[8]),
                            LogType = Int32.Parse(rowData[9]),
                            GameData = rowData[10].ToString(),
                            CreateTime = now
                        });
                    }

                    SlotGameHistoryLogRepository.CreateMany(historyLog);
                }
                catch (Exception ex)
                {
                    throw new TerproException(ApiErrorCode.InvalidParameter_ParsingError, ex.Message);
                }
            }
        }
        #endregion     

        #region Slot Bonus Record
        public PlayerSlotBonusRecordViewModel GetSlotBonusRecord(GameSlotRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var filter = new QueryPlayerSlotBonusRecordRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };

            var querySlotBonusRecord = PlayerSlotBonusRecordRepository.GetByFilter(filter);

            if (querySlotBonusRecord.IsEmptyOrNull())
            {
                return new PlayerSlotBonusRecordViewModel
                {
                    SlotBonusRecord = new List<PlayerSlotBonusRecordItemViewModel> { }
                };
            }

            var result = new PlayerSlotBonusRecordViewModel
            {
                SlotBonusRecord = querySlotBonusRecord.Select(x => new PlayerSlotBonusRecordItemViewModel
                {
                    PlayerID = x.PlayerID,
                    GameID = x.GameID,
                    FlagID = x.FlagID,
                    OriginID = x.OriginID,
                    TurnID = x.TurnID,
                    Location = x.Location,
                    BetAmount = x.BetAmount,
                    BetRate = x.BetRate,
                    Data = x.Data,
                    PatternID = x.PatternID
                }).ToList()
            };

            return result;
        }

        public void UpdateSlotBonusRecord(GameSlotBonusRecordRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.OriginID).Rule_NotZeroNotMinus();

            var filter = new QueryPlayerSlotBonusRecordRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID,
                OriginID = model.OriginID,
                TurnID = model.TurnID
            };

            //使用lock，同時間只允許處理一筆 相同鍵值的資料，避免時間差/互相覆蓋的問題
            string lockKey = $"UpdateSlotBonusRecord_{model.PlayerID}_{model.GameID}_{model.FlagID}_{model.OriginID}_{model.TurnID}";


            lock (ThreadLockHelper.InternConcurrentSafe(lockKey))
            {
                var queryRecord = PlayerSlotBonusRecordRepository.GetByFilter(filter).FirstOrDefault();

                if (queryRecord == null)
                {
                    filter.TurnID = null;
                    filter.Location = model.Location;

                    queryRecord = PlayerSlotBonusRecordRepository.GetByFilter(filter).FirstOrDefault();

                    if (queryRecord != null)
                        throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);   //新增失败，位置已被选择

                    model.RuleFor(x => x.BetAmount).Rule_NotMinus();
                    model.RuleFor(x => x.BetRate).Rule_NotMinus();

                    PlayerSlotBonusRecordRepository.Create(new PlayerSlotBonusRecordEntity
                    {
                        PlayerID = model.PlayerID,
                        GameID = model.GameID,
                        FlagID = model.FlagID,
                        OriginID = model.OriginID,
                        TurnID = model.TurnID,
                        Location = model.Location,
                        BetAmount = model.BetAmount,
                        BetRate = model.BetRate,
                        Data = (model.Data.IsEmptyOrNull()) ? "" : model.Data,
                        PatternID = model.PatternID
                    });
                }
                else
                {
                    queryRecord.Location = model.Location;
                    queryRecord.PatternID = model.PatternID;
                    PlayerSlotBonusRecordRepository.Update(queryRecord);
                }
            }
        }

        public void ResetSlotBonusRecord(GameSlotBonusRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var filter = new QueryPlayerSlotBonusRecordRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID,
                OriginID = model.OriginID
            };

            var queryRecord = PlayerSlotBonusRecordRepository.GetByFilter(filter);
            if (queryRecord.IsNotEmptyNull())
            {
                foreach (var record in queryRecord)
                {
                    PlayerSlotBonusRecordRepository.Delete(record.ID);
                }
            }
        }
        #endregion

        #region Slot Free Spin Record
        public PlayerSlotFreeSpinRecordViewModel GetSlotFreeSpinRecord(GameSlotRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var filter = new QueryPlayerSlotFreeSpinRecordRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };
            var querySlotFreeSpinRecord = PlayerSlotFreeSpinRecordRepository.GetByFilter(filter);
            if (querySlotFreeSpinRecord.IsEmptyOrNull())
            {
                PlayerSlotFreeSpinRecordRepository.Create(new PlayerSlotFreeSpinRecordEntity
                {
                    PlayerID = model.PlayerID,
                    GameID = model.GameID,
                    FlagID = model.FlagID,
                    OriginID = 0,
                    FreeSpinGroup = 0,
                    ChooseRate = 0,
                    BetAmount = 0,
                    Times = 0,
                    WinAmount = 0,
                    WinTimes = 0,
                    Data = string.Empty,
                    PatternID = 0
                });
                querySlotFreeSpinRecord = PlayerSlotFreeSpinRecordRepository.GetByFilter(filter);
            }

            var result = new PlayerSlotFreeSpinRecordViewModel
            {
                SlotFreeSpinRecord = querySlotFreeSpinRecord.Select(x => new PlayerSlotFreeSpinRecordItemViewModel
                {
                    PlayerID = x.PlayerID,
                    GameID = x.GameID,
                    FlagID = x.FlagID,
                    OriginID = x.OriginID,
                    FreeSpinGroup = x.FreeSpinGroup,
                    ChooseRate = x.ChooseRate,
                    BetAmount = x.BetAmount,
                    Times = x.Times,
                    WinAmount = x.WinAmount,
                    WinTimes = x.WinTimes,
                    Data = x.Data,
                    PatternID = x.PatternID
                }).ToList()
            };
            return result;
        }

        public void UpdateSlotFreeSpinRecord(GameSlotFreeSpinRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var filter = new QueryPlayerSlotFreeSpinRecordRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };
            var queryRecord = PlayerSlotFreeSpinRecordRepository.GetByFilter(filter).FirstOrDefault();

            if (queryRecord == null)
            {
                PlayerSlotFreeSpinRecordRepository.Create(new PlayerSlotFreeSpinRecordEntity
                {
                    PlayerID = model.PlayerID,
                    GameID = model.GameID,
                    FlagID = model.FlagID,
                    OriginID = model.OriginID,
                    FreeSpinGroup = model.FreeSpinGroup,
                    ChooseRate = model.ChooseRate,
                    BetAmount = model.BetAmount,
                    Times = model.Times,
                    WinAmount = model.WinAmount,
                    WinTimes = model.WinTimes,
                    Data = (model.Data.IsEmptyOrNull()) ? "" : model.Data,
                    PatternID = model.PatternID
                });
            }
            else
            {
                queryRecord.OriginID = model.OriginID;
                queryRecord.ChooseRate = model.ChooseRate;
                queryRecord.ChooseRate = model.ChooseRate;
                queryRecord.BetAmount = model.BetAmount;
                queryRecord.Times = model.Times;
                queryRecord.WinAmount = model.WinAmount;
                queryRecord.WinTimes = model.WinTimes;
                queryRecord.Data = model.Data;
                queryRecord.PatternID = model.PatternID;
                PlayerSlotFreeSpinRecordRepository.Update(queryRecord);
            }
        }
        #endregion

        #region Slot Lightning Record
        public PlayerSlotLightingRecordViewModel GetSlotLightningRecord(GameSlotRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var filter = new QueryPlayerSlotLightningRecordRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };
            var querySlotLightingRecord = PlayerSlotLightningRecordRepository.GetByFilter(filter);
            
            if (querySlotLightingRecord.IsEmptyOrNull())
            {
                return new PlayerSlotLightingRecordViewModel
                {
                    SlotLightingRecord = new List<PlayerSlotLightingRecordItemViewModel> { }
                };
            }
            
            var result = new PlayerSlotLightingRecordViewModel
            {
                SlotLightingRecord = querySlotLightingRecord.Select(x => new PlayerSlotLightingRecordItemViewModel
                {
                    PlayerID = x.PlayerID,
                    GameID = x.GameID,
                    FlagID = x.FlagID,
                    Data = x.Data,
                    PatternID = x.PatternID
                }).ToList()
            };


            return result;
        }

        public void UpdateSlotLightningRecord(GameSlotLightningRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var filter = new QueryPlayerSlotLightningRecordRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };
            var queryRecord = PlayerSlotLightningRecordRepository.GetByFilter(filter).FirstOrDefault();

            if (queryRecord == null)
            {
                if (!model.Data.IsEmptyOrNull())
                {
                    PlayerSlotLightningRecordRepository.Create(new PlayerSlotLightingRecordEntity
                    {
                        PlayerID = model.PlayerID,
                        GameID = model.GameID,
                        FlagID = model.FlagID,
                        Data = model.Data,
                        PatternID = model.PatternID
                    });
                }
            }
            else
            {
                //-- No Data, DELETE the record
                if (model.Data.IsEmptyOrNull())
                {
                    PlayerSlotLightningRecordRepository.Delete(queryRecord.ID);
                }
                else
                {
                    queryRecord.Data = model.Data;
                    queryRecord.PatternID = model.PatternID;
                    PlayerSlotLightningRecordRepository.Update(queryRecord);
                }
            }
        }

        #endregion

    }
}
