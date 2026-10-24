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
using Lck.Utility.Extensions.Helper;
using target.Utility.Extensions;
using target.Utility.Helper;

namespace target.API.Game.Controllers
{

    /// <summary>
    /// 專門提供測試環境使用的API
    /// </summary>

    [DoNotCheckCertificate]
    [DisabledWhenProduction]
    public class QATestUtilityController : BaseController
    {
        public IFakeDataGenerateService FakeDataGenerateService { get; set; }


        /// <summary>
        /// 解析token
        /// </summary>
        [Route("fake-test/token/parsing")]
        [HttpPost]
        public object ParsingJwtTokenToModel()
        {
            string body = HttpContext.Request.GetRequestBodyString();
            var query = MyJwtHelper.GetPlayerModelFromJwtString(body);
            return query;
        }


        /// <summary>
        /// 提供測試用，快速產生jwtToken
        /// </summary>
        [Route("fake-test/player/token/generate")]
        [HttpGet]
        public GameLaunchViewModel GenerateToken([FromQuery]FakeDataGenerateRequestModel model)
        {
            int userID = GetCurrentPlayerID();
            var query = FakeDataGenerateService.GenerateJwtTokenForQATesting(model);
            return query;
        }


        /// <summary>
        /// 提供測試用，快速加錢
        /// </summary>
        [Route("fake-test/player/add-money")]
        [HttpGet]
        public PlayerWalletAddResult AddPlayerWalletAmount([FromQuery]int playerID)
        {
            int userID = GetCurrentPlayerID();
            var query = FakeDataGenerateService.AddPlayerWallet(playerID);
            return query;
        }

    }
}
