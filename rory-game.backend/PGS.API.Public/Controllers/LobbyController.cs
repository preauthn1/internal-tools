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

namespace target.API.Public.Controllers
{
    [Route("v{version:int}")]
    public class LobbyController : BaseController
    {
        public IGameLaunchService GameLaunchService { get; set; }

        public IInternalUtilityService InternalUtilityService { get; set; }


        [SwaggerDocumentAccessTokenRequired]
        [LobbyTokenAuthorizeFilter]
        [Route("lobby/games")]
        [HttpGet]
        public LobbyApiGetGamesInfoViewModel GetLobbyGames([FromQuery] LobbyApiGetGamesRequestModel model)
        {
            var query = InternalUtilityService.GetLobbyGames(model);
            return query;
        }


        [SwaggerDocumentAccessTokenRequired]
        [LobbyTokenAuthorizeFilter]
        [Route("lobby/game-link")]
        [HttpGet]
        public LobbyGenerateLinkResultModel GenerateGameLinkForLobby([FromQuery] LobbyApiGenerateLinkParameterModel model)
        {
            var query = GameLaunchService.GenerateGameLinkForLobby(model);
            return query;
        }

    }
}
