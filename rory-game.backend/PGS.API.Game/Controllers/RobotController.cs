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
   
    public class RobotController : BaseController
    {
        /// <summary>
        /// Game Server Setting Service
        /// </summary>
        public IGameServerSettingService GameServerSettingService { get; set; }
        

        /// <summary>
        /// 取得機器人設定
        /// </summary>
        [Route("system/server/robots")]
        [ApiConfig(IsHeaderLogEnable = false)]
        [HttpGet]
        public RobotListResultModel GetRobots([FromQuery]RobotListRequestModel model)
        {
            var result = GameServerSettingService.GetRobotList(model);
            return result;
        }

        
    }
}
