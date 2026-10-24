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
using target.Models.Models;
using target.Utility.Helper;
using target.Utility.Extensions;

namespace target.API.Game.Controllers
{
    public class GameBetController : BaseController
    {
     
        public IGameBetService GameBetService { get; set; }


        /// <summary>
        /// 投注
        /// </summary>
        [PlayerTokenAuthorizeFilter]
        [Route("game/placeBet")]
        [HttpPost]
        public GamePlaceBetResultViewModel PlaceBet([FromBody]GamePlaceBetRequestModel model)
        {
            model = model ?? new GamePlaceBetRequestModel();
            model.AccessJwtToken = [REDACTED];
            model.PlayerID = GetCurrentPlayerID();

            var result = GameBetService.PlaceBet(model);
            return result;
        }


    }
}
