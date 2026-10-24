using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.API.Game.Controllers
{
    /// <summary>
    /// 超端(水池相關)
    /// </summary>    
    public class SuperToolController : BaseController
    {

        public ISuperToolPoolService SuperToolPoolService { get; set; }

        [Route("super-tool/pool/win-loss")]
        [HttpGet]
        public SuperToolPoolViewModel GetWinLossPool([FromQuery]SuperToolPoolRequestModel model)
        {
            var query = SuperToolPoolService.GetWinLossPool(model);
            return query;
        }
        
        [Route("super-tool/pool/win-loss")]
        [HttpPost]
        public SuperToolPoolViewModel UpdateWinLossPool([FromBody]SuperToolPoolUpdateModel model)
        {
            var query = SuperToolPoolService.UpdateWinLossPool(model);
            return query;
        }
    }
}
