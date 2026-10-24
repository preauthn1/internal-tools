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

    public partial class GameSlotService
    {

        public GameSlotProbabilitySettingViewModel GetSlotProbabilitySetting(GameSlotRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var filter = new QueryGameSlotProbabilitySettingRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };
            var querySlotProbRecord = GameSlotProbabilitySettingRepository.GetByFilter(filter);

            if (querySlotProbRecord.IsEmptyOrNull())
            {
                return new GameSlotProbabilitySettingViewModel
                {
                    SlotProbRecord = new List<GameSlotProbabilitySettingItemViewModel> { }
                };
            }

            var result = new GameSlotProbabilitySettingViewModel
            {
                SlotProbRecord = querySlotProbRecord.Select(x => new GameSlotProbabilitySettingItemViewModel
                {
                    PlayerID = x.PlayerID,
                    GameID = x.GameID,
                    FlagID = x.FlagID,
                    OriginID = x.OriginID,
                    Times = x.Times,
                    Data = x.Data,
                    PatternID = x.PatternID
                }).ToList()
            };

            return result;
        }


        public void UpdateSlotProbRecord(GameSlotProbRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.OriginID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.PatternID).Rule_NotZeroNotMinus();

            var filter = new QueryGameSlotProbabilitySettingRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };

            string lockKey = $"UpdateSlotProbRecord_{model.PlayerID}_{ model.GameID}_{model.FlagID}";

            //使用lock，同時間只允許處理一筆 相同鍵值的資料，避免時間差/互相覆蓋的問題
            lock (ThreadLockHelper.InternConcurrentSafe(lockKey))
            {
                var queryRecord = GameSlotProbabilitySettingRepository.GetByFilter(filter).FirstOrDefault();

                if (queryRecord == null)
                {
                    GameSlotProbabilitySettingRepository.Create(new GameSlotProbabilitySettingEntity
                    {
                        PlayerID = model.PlayerID,
                        GameID = model.GameID,
                        FlagID = model.FlagID,
                        OriginID = model.OriginID,
                        Times = model.Times,
                        Data = (model.Data.IsEmptyOrNull()) ? "" : model.Data,
                        PatternID = model.PatternID
                    });
                }
                else
                {
                    queryRecord.OriginID = model.OriginID;
                    queryRecord.Times = model.Times;
                    queryRecord.PatternID = model.PatternID;
                    GameSlotProbabilitySettingRepository.Update(queryRecord);
                }
            }
        }

        public void DeleteSlotProbRecord(GameSlotRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotZeroNotMinus();

            var filter = new QueryGameSlotProbabilitySettingRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };
            var queryRecord = GameSlotProbabilitySettingRepository.GetByFilter(filter).FirstOrDefault();
            if (queryRecord != null)
            {
                GameSlotProbabilitySettingRepository.Delete(queryRecord.ID);
            }
        }


    }
}
