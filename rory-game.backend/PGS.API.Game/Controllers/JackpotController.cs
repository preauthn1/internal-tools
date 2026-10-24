using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Web.Modules.Filter;

namespace target.API.Game.Controllers
{
    /// <summary>
    /// 老虎機
    /// </summary>
    /// <remarks>
    /// 老虎機 Jackpot
    /// </remarks>
    public class JackpotController : BaseController
    {

        public IGameSlotService GameSlotService { get; set; }

        /// <summary>
        /// 更新老虎機彩池金額, 後返回當前數值
        /// </summary>
        [Route("game/{gameID:int}/jackpot")]
        [ApiConfig(IsHeaderLogEnable = false , IsResponseLogEnable = false)]
        [HttpPost]
        public GameJackpotViewModel UpdatetJackpot(int gameID, [FromBody]GameSlotJackpotUpdateRequestModel model)
        {
            var query = GameSlotService.UpdateJackpot(gameID, model);
            return query;
        }
    }
}
