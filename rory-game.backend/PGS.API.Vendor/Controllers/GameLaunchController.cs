using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.API.AdminSys;

namespace target.API.Vendor.Controllers
{
    [Route("v{version:int}")]
    public class GameLaunchController : BaseController
    {
        public IGameLinkService GameLinkService { get; set; }

        /// <summary>
        /// 2.1.1 Get game launch url (取得游戏启动连结) (TransferMode)
        /// </summary>
        [SwaggerApiDescriptionEN("Get the link of a game.")]
        [SwaggerApiDescriptionCN("取得游戏网址")]
        [Route("game/link")]
        [ApiExplorerSettings(GroupName = "tfMode")]
        [HttpGet]
        public GameLinkGenerateResultViewModel GenerateGameUrl([FromQuery] GameLaunchVendorApiTransferModeRequestModel model)
        {
            var result = GameLinkService.GenerateTransferModeLaunchLink(model);
            return result;
        }

        /// <summary>
        /// 2.1.2 Get game launch url (取得游戏启动连结) (SeamlessMode)
        /// </summary>
        [SwaggerApiDescriptionEN("Get the link of a game.")]
        [SwaggerApiDescriptionCN("取得游戏网址")]
        [Route("game/launch-link")]
        [ApiExplorerSettings(GroupName = "swMode")]
        [HttpGet]
        public GameLinkGenerateResultViewModel GenerateSingleWalletGameLaunchUrl([FromQuery] GameLaunchVendorApiSingleWalletModeRequestModel model)
        {
            var result = GameLinkService.GenerateSingleWalletLaunchLink(model);
            return result;
        }
    }
}
