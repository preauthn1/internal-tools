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
using target.BL.Services;
using target.Models.Enums;
using target.Models.ViewModels.GameVendor;
using target.Models;
using target.Utility.Helper;
using target.BL.Services;
using System.Threading;
using target.DA.Repository;
using System.Transactions;
using target.Utility.Helper.HttpContext;
using System.Threading.Tasks;
using target.Utility.Utility;
using Lck.Utility.Extensions.Helper;
using target.Utility.Utility.ThirdParty;

namespace target.API.ExternalVendor.Controllers
{
    public class InternalUtilityController : BaseController
    {
        public IGameVendorRoutineJobService GameVendorRoutineJobService { get; set; }

        public IGameLaunchService GameLaunchService { get; set; }


        /// <summary>
        /// 提供給AGG技術人員 的測試用網址
        /// </summary>
        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
        [Route("AGG/launch-test")]
        [HttpGet]
        public object LaunchLobbyUrlFroAGGteam(string playerName = "test00112",int merchantID = 175033)
        {
            var url = GameLaunchService.LaunchGameAndGetGamePageUrl(new GameLaunchRequestModel
            {
                Lang = "zh-tw",
                GameID = 0,
                MerchantID = merchantID,
                PlayerName = playerName,
            });

            return new
            {
                url = url
            };
        }

    }
}
