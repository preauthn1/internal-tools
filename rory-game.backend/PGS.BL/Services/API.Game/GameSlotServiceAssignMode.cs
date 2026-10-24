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

        public GameSlotPlayerAssignedModeViewModel GetSlotAssignedMode(GameSlotRequestModel model)
        {         
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotMinus();

            var filter = new QueryPlayerSlotAssignedModeRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };

            var query = GameSlotAssignedModeRepository.GetByFilter(filter);

            var result = new GameSlotPlayerAssignedModeViewModel
            {
                SlotUserAssignedMode = query.Select(x => new GameSlotPlayerAssignedModeItemViewModel
                {
                    PlayerID = x.PlayerID,
                    GameID = x.GameID,
                    FlagID = x.FlagID,
                    CreateTime = x.CreateTime,
                    Mode = x.Mode,
                    ExpireTime = x.ExpireTime
                }).ToList()
            };

            return result;
        }

        public void AddSlotAssignedMode(GameSlotPlayerAssignedRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotMinus();
            model.RuleFor(x => x.Mode).Rule_NotMinus();
            model.RuleFor(x => x.ExpireTime).Rule_NotEmptyDateTime();

            var filter = new QueryPlayerSlotAssignedModeRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };

            //使用lock，同時間只允許處理一筆 相同鍵值的資料，避免時間差/互相覆蓋的問題
            string lockKey = $"AddSlotAssignedMode_{filter.PlayerID}_{filter.GameID}_{filter.FlagID}";

            lock(ThreadLockHelper.InternConcurrentSafe(lockKey))
            {
                var query = GameSlotAssignedModeRepository.GetByFilter(filter).FirstOrDefault();

                if (query != null)
                {
                    DateTime defaultTime = DateTimeHelper.GetDefaultTime();
                    query.Mode = model.Mode;
                    query.ExpireTime = (model.ExpireTime <= defaultTime) ? defaultTime : model.ExpireTime.ToUniversalTime();  //to utc

                    GameSlotAssignedModeRepository.Update(query);
                }
                else
                {
                    GameSlotAssignedModeRepository.Create(new GameSlotPlayerAssignedModeEntity
                    {
                        PlayerID = model.PlayerID,
                        GameID = model.GameID,
                        FlagID = model.FlagID,
                        Mode = model.Mode,
                        CreateTime = DateTimeHelper.GetCurrentTime(),
                        ExpireTime = model.ExpireTime
                    });
                }
            }
        }

        public void DeleteSlotAssignedMode(GameSlotRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotMinus();

            var filter = new QueryPlayerSlotAssignedModeRequestViewModel   //QueryPlayerSlotDeferredAwardRequestViewModel
            {
                PlayerID = model.PlayerID,
                GameID = model.GameID,
                FlagID = model.FlagID
            };

            var query = GameSlotAssignedModeRepository.GetByFilter(filter);

            if (query.IsNotEmptyNull())
            {
                List<int> ids = query.Select(x => x.ID).ToList();
                GameSlotAssignedModeRepository.Delete(ids);
            }
        }

    }
}
