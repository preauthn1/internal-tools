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
    /// 針對某些三方遊戲有特別的客製API
    /// </summary>
    public class ThirdPartyCustomizeController : BaseController
    {
       
        public IThirdPartyGameLogicService ThirdPartyGameLogicService { get; set; }

        /// <summary>
        /// 根據slot88的MD5 反查某些玩家資訊
        /// </summary>      
        [Route("third-party/slot88/md5")]
        [HttpGet]
        public Slot88MD5InfoResultModel GetPlayerInformationByMD5([FromQuery]QuerySlot88MD5RequestModel model)
        {
            var query = ThirdPartyGameLogicService.GetPlayerInformationByMD5(model);
            return query;
        }


    }
}
