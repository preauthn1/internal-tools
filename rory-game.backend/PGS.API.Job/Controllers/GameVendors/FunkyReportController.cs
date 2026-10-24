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

namespace target.API.Job.Controllers
{
    public class FunkyReportController : BaseController
    {
        public FunkyIntegrationService FunkyIntegrationService { get; set; }

        public IGameVendorSingleWalletService SingleWalletVendorService { get; set; }

        public IGameVendorRoutineJobService GameVendorRoutineJobService { get; set; }


        private static GameVendorTypeEnum _gameVendor = GameVendorTypeEnum.Funky;


        public FunkyReportController()
        {
            FunkyIntegrationService = new FunkyIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.Default)]
        [Route("Report-calcualte/Funky/Yesterday")]
        [HttpGet]
        public void CalculateYesterdayReportIsCorrect()
        {
            GameVendorRoutineJobService.CalculateYesterdayReportIsCorrect(_gameVendor);
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.Default)]
        [Route("Report-calcualte/Funky/DateRange")]
        [HttpGet]
        public void CalculateDateRangeReportIsCorrect(string begin , string end)
        {
            DateTime beginUTC = DateTime.Parse(begin);
            DateTime endUTC = DateTime.Parse(end);

            GameVendorRoutineJobService.CalculateDateRangeReportIsCorrect(_gameVendor , beginUTC , endUTC);
        }


    }
}
