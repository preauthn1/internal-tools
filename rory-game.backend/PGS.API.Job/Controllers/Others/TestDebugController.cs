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

namespace target.API.Job.Controllers
{
    public class TestDebugController : BaseController
    {
        public IGameVendorRoutineJobService GameVendorRoutineJobService { get; set; }

        public IGameLaunchService GameLaunchService { get; set; }


        public FunkyIntegrationService FunkyIntegrationService = new FunkyIntegrationService();




        public static List<string> TempMsg = new List<string>();



        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
        [Route("Debug/Funky/Testttttt")]
        [HttpPost]
        public object DebugTest(int gameID = 250006)
        {
            //GameVendorRoutineJobService.CalculateDateRangeReportIsCorrect(GameVendorTypeEnum.Funky, begin , end);




            //var result = GameLaunchService.LaunchGameAndGetGamePageUrl(new GameLaunchRequestModel
            //{
            //    PlayerName = "marshall006677",
            //    GameID = gameID,
            //    MerchantID = 175001,
            //    SignKey = SystemConstant.PublicAPI_SignKey_MagicPassword,
            //    Ts = DateTimeHelper.GetCurrentTimeTimestamp(),
            //});


            //return result;
            return "XXXXXXXXXX";

            //GameVendorRoutineJobService.GetVendorGames((int)GameVendorTypeEnum.Funky);
        }



        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
        [Route("Debug/Report")]
        [HttpPost]
        public object DebugTestReport(string begin = "" , string end = "")
        {
            //請傳入 2020-11-01  這樣
            DateTime beginDate = DateTime.Parse(begin).Date;
            DateTime endDate = DateTime.Parse(end).Date;

            GameVendorRoutineJobService.CalculateDateRangeReportIsCorrect(GameVendorTypeEnum.Funky, beginDate, endDate);
            return "XXXXXXXXXX";
        }


        [Route("Debug/Funky/Test002")]
        [HttpPost]
        public object Test002()
        {
            return new
            {
                text1 = AppHelper.SimplifiedToTraditional("这是测试档案"),
                text2 = "这是测试档案",
            };
        }



        [Route("Debug/Funky/Test003")]
        [HttpGet , HttpPost]
        public object Test003(string userID)
        {
            //lock(userID)
            //{
            //    TempMsg.Add(userID + "_" + DateTime.Now.ToString("HH:mm:ss"));
            //    Thread.Sleep(5000);

            //    if(userID != "6666")
            //    {
            //        throw new Exception("xxxxxxxx");
            //    }
            //}

            return TempMsg;
        }



        [Route("Debug/Funky/Test004")]
        [HttpGet, HttpPost]
        public object Test004()
        {
            var rep = new PlayerRepository();
            rep.TestDebug_TransactionIn();
            return TempMsg;
        }



    }
}
