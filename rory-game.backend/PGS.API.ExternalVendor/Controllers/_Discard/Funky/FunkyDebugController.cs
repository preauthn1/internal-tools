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

namespace target.API.ExternalVendor.Controllers
{
    public class FunkyDebugController : BaseController
    {
        public FunkyIntegrationService FunkyIntegrationService { get; set; }

        public IGameVendorSingleWalletService SingleWalletVendorService { get; set; }

        public IGameVendorRoutineJobService GameVendorRoutineJobService { get; set; }



        public FunkyDebugController()
        {
            FunkyIntegrationService = new FunkyIntegrationService();
        }


        /// <summary>
        /// debug用 取得原始request response  (遊戲列表API)
        /// </summary>
        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.Default)]
        [Route("funky/debug/games")]
        [HttpGet]
        public object GetGames()
        {
            var queryApi = FunkyIntegrationService.GetGames(SystemConstant.DefaultMerchantID ,"EN");
            return queryApi;
        }


        /// <summary>
        /// debug用 取得原始request response (報表API)
        /// </summary>
        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.Default)]
        [Route("funky/debug/report")]
        [HttpGet]
        public object GetDbAndApiSummaryReport(string begin, string end)
        {
            DateTime beginUTC = DateTime.Parse(begin);
            DateTime endUTC = DateTime.Parse(end);

            var queryApi = FunkyIntegrationService.GetDbAndApiSummaryReport(new ExternalVendorGetReportRequestParameter
            {
                StartTimeTargetTimeZone = beginUTC,
                EndTimeTargetTimeZone = endUTC,
                StartTimeUTC = beginUTC,
                EndTimeUTC = endUTC,
                MerchantID = SystemConstant.DefaultMerchantID
            });

            return queryApi;
        }



    }
}
