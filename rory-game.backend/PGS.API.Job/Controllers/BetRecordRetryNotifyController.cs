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
    public class BetRecordRetryNotifyController : BaseController
    {
        public IBetRecordStatusChangeNotifyService BetRecordStatusChangeNotifyService { get; set; }

        /// <summary>
        /// 補救某些未結算/未通知完成的注單
        /// </summary>
        [Route("bet-record/retry/queue")]
        [HttpGet]
        public void BetReocrdNotifyQueueRetryProcess()
        {
            BetRecordStatusChangeNotifyService.BetReocrdNotifyQueueRetryProcess();
        }

    }
}
