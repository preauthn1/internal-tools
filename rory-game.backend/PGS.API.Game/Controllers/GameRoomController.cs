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
using target.Utility.Helper;

namespace target.API.Game.Controllers
{
    public class GameRoomController : BaseController
    {
        public IGameRoomService GamePlayService { get; set; }


        /// <summary>
        /// (access_Token*) 遊戲進房
        /// </summary>
        [PlayerTokenAuthorizeFilter]
        [Route("game/room/in")]
        [HttpPost]
        public GameRoomEnterViewModel GameRoomEnter(GameRoomEnterRequestModel model)
        {
            model = model ?? new GameRoomEnterRequestModel();
            model.AccessToken = [REDACTED];

            int playerID = GetCurrentPlayerID();

            var query = GamePlayService.GameRoomEnter(playerID, model);
            return query;
        }
        

        /// <summary>
        /// 遊戲退房
        /// </summary>
        [Route("game/room/out")]
        [LoggingPlayerID]
        [HttpPost]
        public GameRoomLeaveViewModel GameRoomLeave(GameRoomLeaveRequestModel model)
        {
            var query = GamePlayService.GameRoomLeave(model);
            return query;
        }


    }
}
