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
using target.Utility.Extensions;
using target.Models.ViewModels.API.AdminSys;
using Microsoft.Extensions.Caching.Memory;
using target.Utility.Utility.ThirdParty;
using target.Utility.Utility;
using target.Models.ViewModels.Base;
using target.Models;

namespace target.API.AdminSys.Controllers
{
    public class InternalTestController : BaseController
    {
        public IGameLaunchService GameLaunchService { get; set; }

        /// <summary>
        /// 測試取得遊戲網址
        /// </summary>
        [Route("internal-test/game-launch")]
        [AdminTokenrAuthorizeFilter]
        [HttpGet]
        public InternalTestResultViewModel LaunchGameAndGetGamePageUrlByPlayerID([FromQuery]InternalTestBundleRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotMinus();

            string resultUrl = GameLaunchService.LaunchGameAndGetGamePageUrlByPlayerID(model.PlayerID , model.GameID);
            return new InternalTestResultViewModel
            {
                Url = resultUrl
            };
        }


        /// <summary>
        /// 檢查dev與正式環境的pipelineID
        /// </summary>
        /// <returns></returns>
        [Route("internal-test/pipeline-id-check")]
        [HttpGet]
        public object CheckOnlineVersion()
        {
            List<string> listDev = new List<string>();
            List<string> list = new List<string>();

            string line = "---------------------------------------";
            string password = [REDACTED];
            string query = $"?password=[REDACTED]";

            listDev.Add(CheckVersion("API.AdminSys", $"http://example.invalid{query}").msg);
            listDev.Add(CheckVersion("API.External", $"http://example.invalid{query}").msg);
            listDev.Add(CheckVersion("API.Job", $"http://example.invalid{query}").msg);
            listDev.Add(CheckVersion("API.Game", $"http://example.invalid{query}").msg);
            listDev.Add(CheckVersion("API.Public", $"http://example.invalid{query}").msg);
            listDev.Add(CheckVersion("API.Vendor", $"http://example.invalid{query}").msg);


            list.Add(CheckVersion("API.AdminSys", $"https://example.invalid{query}").msg);
            list.Add(CheckVersion("API.External", $"https://example.invalid{query}").msg);
            list.Add(CheckVersion("API.Job", $"https://example.invalid{query}").msg);
            list.Add(CheckVersion("API.Game", $"https://example.invalid{query}").msg);
            list.Add(CheckVersion("API.Public", $"https://example.invalid{query}").msg);
            list.Add(CheckVersion("API.Vendor", $"https://example.invalid{query}").msg);

            string slackMsg = "各API版本為" + Environment.NewLine;

            slackMsg += Environment.NewLine + "*Develop環境*" + Environment.NewLine + line + Environment.NewLine;

            foreach (var msg in listDev)
            {
                slackMsg += msg + Environment.NewLine;
            }

            slackMsg += Environment.NewLine + "*Production環境*" + Environment.NewLine + line + Environment.NewLine;

            foreach (var msg in list)
            {
                slackMsg += msg + Environment.NewLine;
            }

            MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
            {
                Channel = SlackChannelEnum.backend_release_k8s,
                Message = slackMsg
            });

            return "訊息已發送";
        }



        private (bool isPipelineExist, string msg) CheckVersion(string category, string homePageUrl)
        {
            var result = APIRequestHelper.RequestApiResult<ApiResponseGenericModel<Dictionary<string, object>>>(new RequestParameters
            {
                BaseUrl = homePageUrl,
                ApiUrl = APIRequestHelper.No_Url_Keyword,
                Method = HttpMethodType.HttpGet
            });

            if (result.ResultObject != null &&
               result.ResultObject.Data != null &&
               result.ResultObject.Data.IsNotEmptyNull() &&
               result.ResultObject.Data.ContainsKey("pipeline_ID")
               )
            {
                var pipeline = result.ResultObject.Data["pipeline_ID"].ToString();
                string url = "http://example.invalid:1688/backend/target/game.backend/pipelines/" + pipeline;
                string msg = $"PipelineID: " + $"<{url}|{pipeline}>" + $" <{homePageUrl}|`{category}`>";
                return (true, msg);
            }

            return (false, $"查無PipelineID `{category}`");
        }






    }
}
