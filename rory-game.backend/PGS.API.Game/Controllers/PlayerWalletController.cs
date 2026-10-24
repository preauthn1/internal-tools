using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Web.Modules.Filter;
using target.Models.ViewModels;
using target.Utility.Helper;

namespace target.API.Game.Controllers
{
    public class PlayerWalletController : BaseController
    {
        public IPlayerWalletInfoService PlayerWalletInfoService { get; set; }


        /// <summary>
        /// (access_Token*) 用戶基本資訊
        /// </summary>
        [PlayerTokenAuthorizeFilter]
        [Route("players/info")]
        [HttpGet]
        public PlayerBriefInfoViewModel GetPlayerInfo()
        {
            int playerID = GetCurrentPlayerID();
            var query = PlayerWalletInfoService.GetPlayerInfo(playerID);
            return query;
        }


        /// <summary>
        /// (access_Token*) 用戶錢包 (單純取錢包金額)
        /// </summary>
        [PlayerTokenAuthorizeFilter]
        [Route("players/wallet")]
        [HttpGet]
        public PlayerWalletViewModel GetPlayerWalletCash()
        {
            string accessToken = [REDACTED];

            int playerID = GetCurrentPlayerID();
            var query = PlayerWalletInfoService.GetWalletAmount(playerID , accessToken);
            return query;
        }

    }
}
