using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using Microsoft.AspNetCore.Cors;
using target.Web.Modules.Filter;
using target.Web.Modules.Filter;

namespace target.API.Public.Controllers
{
    [Route("v{version:int}")]
    public class GameRoomController : BaseController
    {
        public IGameLaunchService GameLaunchService { get; set; }

        /// <summary>
        /// 選廳頁 取得房間&基本資訊 (這邊不是preDemo的 servers API)
        /// </summary>
        [Route("player/game-rooms")]
        [HttpGet]
        [SwaggerDocumentAccessTokenRequired]
        [SignKeyVerifyFilter]
        public GameServerViewModel GetGameServers([FromQuery]GameServerRequestModel model)
        {
            var query = GameLaunchService.GetGameServers(model);
            return query;
        }

    }
}
