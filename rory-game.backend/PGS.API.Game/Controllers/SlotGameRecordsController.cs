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
   
    public class SlotGameRecordsController : BaseController
    {

        /// <summary>
        /// Game Slot Service
        /// </summary>
        public IGameSlotService GameSlotService { get; set; }


        /// <summary>
        /// Slot 取遊戲歷程
        /// </summary>
        [Route("game/slot/log")]
        [LoggingPlayerID]
        [HttpGet]
        public PagedList<GameSlotHistoryLogViewModel> GetSlotGameLog([FromQuery]GameSlotLogRequestModel model)
        {
            var query = GameSlotService.GetSlotGameLog(model);
            return query;
        }

        /// <summary>
        /// 寫入Slot遊戲歷程
        /// </summary>
        [Route("game/slot/log")]
        [HttpPost]
        public void AddSlotGameLog([FromBody]GameSlotHistoryLogRequestModel model)
        {
            GameSlotService.AddSlotGameLog(model);
        }
       

        /// <summary>
        /// 獲取老虎機資料2 Slot Bonus Record
        /// </summary>
        /// <remarks>
        /// 讀取PlayerSlotBonusRecord資料 (原PrGs_Slot_BonusGameRead)
        /// </remarks>
        [Route("game/slot/bonus")]
        [LoggingPlayerID]
        [HttpGet]
        public PlayerSlotBonusRecordViewModel GetSlotBonusRecord([FromQuery]GameSlotRequestModel model)
        {
            var query = GameSlotService.GetSlotBonusRecord(model);
            return query;
        }

        /// <summary>
        /// 更新老虎機資料2
        /// </summary>
        /// <remarks>
        /// 更新Slot_BonusGameRecord資料 (原PrGs_Slot_BonusGameRecord)
        /// </remarks>
        [Route("game/slot/bonus")]
        [LoggingPlayerID]
        [HttpPost]
        public void UpdateSlotBonusRecord([FromBody]GameSlotBonusRecordRequestModel model)
        {
            GameSlotService.UpdateSlotBonusRecord(model);
        }

        /// <summary>
        /// 重置老虎機資料2
        /// </summary>
        /// <remarks>
        /// 重置Slot_BonusGameRecord資料 (原PrGs_Slot_BonusGameReset)
        /// </remarks>
        [Route("game/slot/bonus")]
        [LoggingPlayerID]
        [HttpDelete]
        public void ResetSlotBonusRecord([FromQuery]GameSlotBonusRequestModel model)
        {
            GameSlotService.ResetSlotBonusRecord(model);
        }


        /// <summary>
        /// 獲取老虎機資料3
        /// </summary>
        /// <remarks>
        /// 讀取PlayerSlotFreeSpinRecord資料 (原PrGs_Slot_FreeSpinRead)
        /// </remarks>
        [Route("game/slot/free-spin")]
        [LoggingPlayerID]
        [HttpGet]
        public PlayerSlotFreeSpinRecordViewModel GetSlotFreeSpinRecord([FromQuery]GameSlotRequestModel model)
        {
            var query = GameSlotService.GetSlotFreeSpinRecord(model);
            return query;
        }

        /// <summary>
        /// 更新老虎機資料3
        /// </summary>
        /// <remarks>
        /// 寫入FreeSpin紀錄 (原PrGs_Slot_FreeSpinRecord)
        /// </remarks>
        [Route("game/slot/free-spin")]
        [HttpPost]
        public void UpdateSlotFreeSpinRecord([FromBody]GameSlotFreeSpinRequestModel model)
        {
            GameSlotService.UpdateSlotFreeSpinRecord(model);
        }



        /// <summary>
        /// 獲取老虎機資料 (金吉報喜 相關的遊戲會用到)
        /// </summary>
        /// <remarks>
        /// 讀取PlayerSlotLightningRecord資料 (原PrGs_Slot_LightningGameRead)
        /// </remarks>
        [Route("game/slot/lightning")]
        [HttpGet]
        public PlayerSlotLightingRecordViewModel GetSlotLightningRecord([FromQuery]GameSlotRequestModel model)
        {
            var query = GameSlotService.GetSlotLightningRecord(model);
            return query;
        }

        /// <summary>
        /// 獲取老虎機資料4
        /// </summary>
        /// <remarks>
        /// 寫入PlayerSlotLightningRecord紀錄 (原PrGs_Slot_LightningRecord)
        /// </remarks>
        [Route("game/slot/lightning")]
        [ApiConfig(IsRequestBodyLogEnable = false)]
        [HttpPost]
        public void UpdateSlotLightningRecord([FromBody]GameSlotLightningRequestModel model)
        {
            GameSlotService.UpdateSlotLightningRecord(model);
        }



    }
}
