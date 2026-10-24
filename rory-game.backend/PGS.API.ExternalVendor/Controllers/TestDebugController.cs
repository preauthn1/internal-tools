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
    public class TestDebugController : BaseController
    {
        public IGameVendorRoutineJobService GameVendorRoutineJobService { get; set; }

        public IGameLaunchService GameLaunchService { get; set; }

        public FunkyIntegrationService FunkyIntegrationService = new FunkyIntegrationService();

        public BetRecordRepository BetRecordRepository = new BetRecordRepository();


        public static List<string> TempMsg = new List<string>();

        public static List<string> resultTimes = new List<string>();


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
        [Route("Debug/Test/TTTTTTTTTTTTT")]
        [HttpGet]
        public object TTTTTTTTTT()
        {
            return "";
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
        [Route("Debug/Test/DDOS")]
        [HttpGet]
        public object TestDDOSRequest(string token = [REDACTED]")
        {
            List<Task> allTask = new List<Task>();

            for (int i = 0; i < 1500; i++)
            {
                Task.Factory.StartNew(() =>
                {
                    var para = new AGGSingleWalletAPITransactionRequestModel
                    {
                        AccountToken = $"290001__{token}",
                        GameId = "ht-emperorqin",
                        GameSessionId = Guid.NewGuid().ToString(),
                        TransactionId = Guid.NewGuid().ToString(),
                        UserID = 0,
                        WagerId = Guid.NewGuid().ToString(),
                        Operations = new List<AGGSingleWalletOperationRequestModel>
                            {
                                new AGGSingleWalletOperationRequestModel
                                {
                                    Action = "create",
                                    OperationType = "WagerLifecycleOperation",
                                },
                                new AGGSingleWalletOperationRequestModel
                                {
                                    OperationType = "WagerWithdrawOperation",
                                    Withdraw = 10
                                }
                            }
                    };

                    try
                    {
                        //var ppppppp = new RequestParameters
                        //{
                        //    BaseUrl = "https://example.invalid",
                        //    ApiUrl = "/agg/wager/transaction?v=v5",
                        //    Method = HttpMethodType.HttpPost
                        //};

                        var ppppppp = new RequestParameters
                        {
                            BaseUrl = "https://example.invalid",
                            ApiUrl = "/Debug/Test/Cache/Table",
                            Method = HttpMethodType.HttpGet
                        };

                        var qqq = APIRequestHelper.RequestApiResult<string>(ppppppp, para);

                        bool isOK = qqq.ResponseText?.Contains("RS_OK") ?? false;
                        string memoText = isOK ? "success" : "fail";


                        string msg = $"memoText_{qqq.ExecutedTime}";
                        resultTimes.Add(msg);
                    }
                    catch (Exception ex)
                    {
                        string qqq = ex.Message;
                    }
                });

            }

            return "TestDDOSRequest";
        }



        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
        [Route("Debug/Test/Cache/Table")]
        [HttpGet]
        public object TestTableValueRedis(int gameID = 1023 )
        {
            var rep = new GameListRepository();
            var qqq = rep.GetSingleByCache(gameID);

            //似乎有出現 flushKey 會timeout的問題
            string testKey = "XXXX_YYYY_TTTT";
            CacheHelper.SetCacheValue(testKey, "XXXXXXXX" , new TimeSpan(0,5,0));
            CacheHelper.RemoveCacheAsync(testKey);
            
            return qqq;
        }



        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
        [Route("Debug/Header/LogTest")]
        [HttpGet]
        public object HeaderLogTest()
        {
            MyHttpContext.Current.Request.Headers.Add("DebugInfo2", "tttttttttttt");
            return "XXXXXXXXXX";
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
        [Route("Debug/Funky/TestTransaction")]
        [HttpGet]
        public object TestTransaction()
        {
            using (var tran = new TransactionScope())
            {
                var query = BetRecordRepository.GetByTransactionNo("DFDC210127094239791rz0fh");
                query.Memo = DateTimeHelper.GetCurrentTime().ToString("yyyy-MM-dd HH:mm:ss");

                BetRecordRepository.Update(query);
                tran.Complete();
            }

            return "XXXXXXXXXX";
        }

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
        public object DebugTestReport(string begin = "", string end = "")
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


        [Route("Debug/Delay/Test")]
        [HttpPost]
        public object TestDelay(int milliseconds = 200)
        {
            Thread.Sleep(milliseconds);
            return new
            {
                now = DateTimeHelper.GetCurrentTime().ToString()
            };
        }



        [Route("Debug/TestException")]
        [HttpGet, HttpPost]
        public object TestException()
        {
            throw new Exception("ttttttttest");
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
