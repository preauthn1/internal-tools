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

namespace target.API.ExternalVendor.Controllers
{
    public class GameVendorDebugController : BaseController
    {
        public AGGIntegrationService AGGIntegrationService { get; set; }

        public FishingIntegrationService FishingIntegrationService { get; set; }

        public IGameVendorSingleWalletService SingleWalletVendorService { get; set; }

        public IGameVendorRoutineJobService GameVendorRoutineJobService { get; set; }

        public IGameLaunchService GameLaunchService { get; set; }


        public GameVendorDebugController()
        {
            AGGIntegrationService = new AGGIntegrationService();
            FishingIntegrationService = new FishingIntegrationService();
        }


        [Route("agg/debug/lannch")]
        [HttpGet]
        public object GetAGGGameUrl(int gameID = 290001)
        {
            var qqqqqq = GameLaunchService.LaunchGameAndGetGamePageUrl(new GameLaunchRequestModel
            {
                GameID = gameID,
                MerchantID = 175001,
                PlayerName = "qa_002",
            });

            return qqqqqq;
        }


        [Route("fishing/debug/lannch")]
        [HttpGet]
        public object GetFishingGameUrl(int gameID = 310001 , string playerName = "qa_004")
        {
            int merchantID = 175027; //175027: 100倍比例

            var qqqqqq = GameLaunchService.LaunchGameAndGetGamePageUrl(new GameLaunchRequestModel
            {
                Lang = "zh-tw",
                GameID = gameID,
                MerchantID = merchantID,
                PlayerName = playerName,
            });

            return qqqqqq;
        }


        [Route("fishing/debug/getRecords")]
        [HttpGet]
        public object GetFishingRecords(int gameID = 310001, string playerName = "qa_004" , string no = "2000002374-171")
        {
            int merchantID = 175027; //175027: 100倍比例

            var qqqqqq = FishingIntegrationService.GetExternalBetDetailInfo(new ExternalVendorBetRecordDetailInfoRequestModel
            {
                BetRecordNo = no,
            });

            return qqqqqq;
        }



        /// <summary>
        /// debug用 取得原始request response  (遊戲列表API)
        /// </summary>
        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.Default)]
        [Route("agg/debug/betDetail")]
        [HttpGet]
        public object GetBetDetail(int gameID = 290001)
        {
            var response = AGGIntegrationService.GetRecordReplayHistoryByApi(new ExternalVendorGetSingleBetRecordRequestParameter
            {
                BetRecordNo = "cfaf76e3-ae47-4660-9095-cd5f02b094c3",
                GameCode = "ht-emperorqin",
                MerchantID = 175001
            });

            return response;
        }



    }
}
