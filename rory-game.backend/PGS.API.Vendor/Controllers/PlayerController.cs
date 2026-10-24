using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.API.Vendor.Controllers;
using target.API.AdminSys;

namespace target.API.Vendor.Controllers
{
    [Route("v{version:int}")]
    public class PlayerController : BaseController
    {
        public IPlayerVendorService PlayerVendorService { get; set; }


        /// <summary>
        /// 2.2.1 Check player exist (检查用户是否存在)
        /// </summary>
        [SwaggerApiDescriptionEN("Check player account exist in game provider's system.")]
        [SwaggerApiDescriptionCN("检查运营商之特定玩家账号、是否已存在于游戏商")]
        [Route("player/exist")]
        [ApiExplorerSettings(GroupName = "tfMode")]
        [HttpGet]
        public PlayerVendorViewModel CheckIsPlayerAccountExist([FromQuery]QueryPlayerAccountExistRequestModel model)
        {
            var query = PlayerVendorService.IsAccountExist(model);
            return query;
        }


        /// <summary>
        /// 2.2.2 Create player account (建立玩家帐户)
        /// </summary>
        [SwaggerApiDescriptionEN("Create player game account.")]
        [SwaggerApiDescriptionCN("运营商可透过此接口建立玩家账号")]
        [Route("player")]
        [ApiExplorerSettings(GroupName = "tfMode")]
        [HttpPost]
        public PlayerCreateResultVenodrApiViewModel CreateNewPlayer([FromBody]CreatePlayerAccountExistRequestModel model)
        {
            var result = PlayerVendorService.CreateNewPlayer(model);
            return result;
        }

    }
}
