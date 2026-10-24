using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Lck.Utility.Extensions.Helper;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility.Extensions;
using MSC.Web.Modules.Filter;

namespace MSC.API.Game.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    public class FakeTestController : BaseController
    {
        public ICardInputService CardInputService { get; set; }


        /// <summary>
        /// 測試輸入隨機卡號
        /// </summary>
        [DisabledWhenProduction]
        [Route("fake-test/card/input/randomCard")]
        [HttpPost]
        public CardInputViewModel FakeTestRandomInputCard([FromBody]CardInputFakeTestRequestModel model)
        {
            var result = CardInputService.FakeTestRandomInputCard(model);
            return result;
        }

        /// <summary>
        /// 測試輸入指定卡號 (merchantRequestID 直接透過API指定)
        /// </summary>
        [Route("fake-test/card/input")]
        [HttpPost]
        public CardInputViewModel FakeTestInputTargetCardNumber([FromBody] CardInputFakeTestTargetCardNumberRequestModel model)
        {
            var result = CardInputService.FakeTestInputTargetCardNumber(model);
            return result;
        }


    }
}
