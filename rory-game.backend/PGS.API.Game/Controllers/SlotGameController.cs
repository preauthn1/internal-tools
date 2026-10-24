using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.ViewModels;
using target.Web.Modules.Filter;

namespace target.API.Game.Controllers
{
    /// <summary>
    /// Game Slot Play
    /// </summary>
    public class SlotGameController : BaseController
    {

        /// <summary>
        /// Game Slot Service
        /// </summary>
        public IGameSlotService GameSlotService { get; set; }

       


        /// <summary>
        /// 獲取老虎機設置 GameSlotProbabilitySetting
        /// </summary>
        [Route("game/slot/prob")]
        [LoggingPlayerID]
        [HttpGet]
        public GameSlotProbabilitySettingViewModel GetSlotProbRecord([FromQuery]GameSlotRequestModel model)
        {
            var query = GameSlotService.GetSlotProbabilitySetting(model);
            return query;
        }

        /// <summary>
        /// 更新老虎機資料  GameSlotProbabilitySetting
        /// </summary>
        [Route("game/slot/prob")]
        [LoggingPlayerID]
        [HttpPost]
        public void UpdateSlotProbRecord([FromBody]GameSlotProbRequestModel model)
        {
            GameSlotService.UpdateSlotProbRecord(model);
        }

        /// <summary>
        /// 重置 GameSlotProbabilitySetting
        /// </summary>
        [Route("game/slot/prob")]
        [LoggingPlayerID]
        [HttpDelete]
        public void DeleteSlotProbRecord([FromQuery]GameSlotRequestModel model)
        {
            GameSlotService.DeleteSlotProbRecord(model);
        }


        /// <summary>
        /// 讀取PlayerSlotDeferredAward資料
        /// </summary>
        [Route("game/slot/deferred-award")]
        [LoggingPlayerID]
        [HttpGet]
        public GameSlotDeferredAwardViewModel GetSlotDeferredAward([FromQuery]GameSlotDeferredAwardRequestModel model)
        {
            var query = GameSlotService.GetSlotDeferredAward(model);
            return query;
        }

        /// <summary>
        /// 寫入 SlotDeferredAward資料
        /// </summary>
        [Route("game/slot/deferred-award")]
        [LoggingPlayerID]
        [HttpPost]
        public void AddSlotDeferredAward([FromBody]GameSlotDeferredAwardUpdateRequestModel model)
        {
            GameSlotService.AddSlotDeferredAward(model);
        }

        /// <summary>
        /// 讀取 SlotDeferredAward第一筆記錄，並刪除該紀錄
        /// </summary>
        [Route("game/slot/deferred-award")]
        [LoggingPlayerID]
        [HttpDelete]
        public GameSlotDeferredAwardItemViewModel DeleteSlotDeferredAward([FromQuery]GameSlotDeferredAwardRequestModel model)
        {
            var query = GameSlotService.DeleteSlotDeferredAward(model);
            return query;
        }
        

        /// <summary>
        /// 讀取 Slot_UserAssignedMode 資料
        /// </summary>
        [Route("game/slot/user-assigned-mode")]
        [LoggingPlayerID]
        [HttpGet]
        public GameSlotPlayerAssignedModeViewModel GetUserAssignedMode([FromQuery]GameSlotRequestModel model)
        {
            var query = GameSlotService.GetSlotAssignedMode(model);
            return query;
        }

        /// <summary>
        /// 寫入 Slot_UserAssignedMode資料
        /// </summary>
        [Route("game/slot/user-assigned-mode")]
        [LoggingPlayerID]
        [HttpPost]
        public void AddSlotAssignedMode([FromBody]GameSlotPlayerAssignedRequestModel model)
        {
            GameSlotService.AddSlotAssignedMode(model);
        }

        /// <summary>
        /// 刪除 Slot_UserAssignedMode指定紀錄
        /// </summary>
        [Route("game/slot/user-assigned-mode")]
        [LoggingPlayerID]
        [HttpDelete]
        public void DeleteSlotAssignedMode([FromQuery]GameSlotRequestModel model)
        {
            GameSlotService.DeleteSlotAssignedMode(model);
        }
       

    }
}
