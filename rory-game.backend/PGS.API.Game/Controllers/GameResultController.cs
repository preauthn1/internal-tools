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
using target.Models.Models;
using target.Utility.Helper;
using target.Utility.Extensions;

namespace target.API.Game.Controllers
{
    public class GameResultController : BaseController
    {
        public IGameResultService GamePlayResultService { get; set; }


        /// <summary>
        /// 每局結算完寫分、更新錢包
        /// </summary>
        [Route("game/result/settlement")]
        [ApiConfig(IsHeaderLogEnable = false)]
        [LoggingPlayerID]
        [HttpPost]
        public GamePlayResultViewModel UpdateSettlementResult(GamePlayResultRequestModel model)
        {
            var result = GamePlayResultService.UpdateSettlementAmount(model);
            return result;
        }


        /// <summary>
        /// 每局結算完寫分、更新錢包
        /// </summary>
        [Route("v2/game/result/settlement")]
        [ApiConfig(IsHeaderLogEnable = false)]
        [LoggingPlayerID]
        [HttpPost]
        public GamePlayResultViewModel UpdateSettlementResultV2Process(GamePlayResultRequestModel model)
        {
            model = model ?? new GamePlayResultRequestModel();
            model.IsPlaceBetProcess = true;

            var result = GamePlayResultService.UpdateSettlementAmount(model);
            return result;
        }


        /// <summary>
        /// 寫入每局遊戲紀錄
        /// </summary>
        [Route("game/result/per-game-log")]
        [ApiConfig(IsHeaderLogEnable = false)]
        [HttpPost]
        public PerGameLogViewModel WritePerGameLog(PerGameLogRequestModel model)
        {
            var result = GamePlayResultService.WritePerGameLog(model);
            return result;
        }



        /// <summary>
        /// 寫入每局遊戲紀錄
        /// </summary>
        [Route("v2/game/result/per-game-log")]
        [ApiConfig(IsHeaderLogEnable = false)]
        [HttpPost]
        public PerGameLogViewModel WritePerGameLogV2Process(PerGameLogRequestModel model)
        {
            model = model ?? new PerGameLogRequestModel();
            model.IsPlaceBetProcess = true;

            var result = GamePlayResultService.WritePerGameLog(model);
            return result;
        }



    }
}
