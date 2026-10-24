using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.Enums;
using target.API.AdminSys;

namespace target.API.Vendor.Controllers
{
    [Route("v{version:int}")]
    public class GamesController : BaseController
    {
        public IGameInfoService GameInfoService { get; set; }



        /// <summary>
        /// 4.6.1 Get game list (取得游戏列表)
        /// </summary>
        [SwaggerApiDescriptionEN("Get enable game list.")]
        [SwaggerApiDescriptionCN("取得运营商已开放的游戏清单")]
        [Route("games")]
        [HttpGet]
        public List<GameInfoViewModel> GetAllGames([FromQuery]GameInfoRequestModel model)
        {
            var query = GameInfoService.GetAllGames(model);
            return query;
        }

    }
}
