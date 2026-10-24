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
using target.Models.Paging;
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class BetRecordsController : BaseController
    {

        public IBetRecordService BetRecordService { get; set; }

        public IBetRecordStatusChangeNotifyService BetRecordStatusChangeNotifyService { get; set; }

        public IThirdPartyGameVendorService ThirdPartyGameVendorService { get; set; }



        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagement_BetRecordSearch_SearchList)]
        [Route("bet-records/union/list")]
        [HttpGet]
        public PagedList<UnionBetRecordViewModel> GetUnionBetRecordsByPagingFilter([FromQuery] QueryUnionBetRecordRequestViewModel model)
        {
            int userID = GetCurrentUserID();
            var query = BetRecordService.GetUnionBetRecordsByPagingFilter(userID, model);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagement_BetRecord_ManualRefund)]
        [Route("bet-records/notify-queue")]
        [HttpGet]
        public BetRecordNotifyQueueDetailInfoViewModel GetNotifyLogInfo([FromQuery] long notifyQueueID)
        {
            int userID = GetCurrentUserID();
            var query = BetRecordService.GetNotifyQueueInfo(userID, notifyQueueID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagement_BetRecord_ManualRefund)]
        [Route("bet-records/notify-queue/trigger")]
        [HttpPost]
        public void TriggerNotifyQueue([FromQuery] long notifyQueueID)
        {
            int userID = GetCurrentUserID();
            BetRecordStatusChangeNotifyService.TriggerRetryByID( notifyQueueID);
        }



        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagement_BetRecord_ManualRefund)]
        [Route("bet-records/manual/refund")]
        [HttpPost]
        public void ManualRefundBetRecord([FromBody] BetRecordManualRefundRequestModel model)
        {
            int userID = GetCurrentUserID();
            BetRecordService.ManualRefundBetRecord(userID, model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagement_BetRecord_ManualRefund)]
        [Route("bet-records/manual/settle")]
        [HttpPost]
        public void ManualSettleBetRecord([FromBody] BetRecordManualSettleRequestModel model)
        {
            int userID = GetCurrentUserID();
            BetRecordService.ManualSettleBetRecord(userID, model);
        }


        /// <summary>
        /// 三方遊戲商，產生外部連結 (也許並非每個廠商都有)
        /// </summary>
        [Route("game-vendor/betRecord/{perGameDetailID}/url")]
        [Route("game-vendor/betRecord/{perGameDetailID}/detailInfo")]
        [HttpGet]
        [AdminTokenrAuthorizeFilter]
        public ThirdPartyGameVendorBetRecordDetailUrlViewModel GetDetailUrlByPegGameDetailLogID(int perGameDetailID)
        {
            int userID = GetCurrentUserID();
            var result = ThirdPartyGameVendorService.GetDetailUrlByPegGameDetailLogID(userID, perGameDetailID);
            return result;
        }


        /// <summary>
        /// 三方遊戲商，產生外部連結 (也許並非每個廠商都有)
        /// </summary>
        [Route("game-vendor/betRecord/by-transaction-no/{transactionNo}/detailInfo")]
        [HttpGet]
        [AdminTokenrAuthorizeFilter]
        public ThirdPartyGameVendorBetRecordDetailUrlViewModel GetDetailUrlByTransactionNo(string transactionNo)
        {
            int userID = GetCurrentUserID();
            var result = ThirdPartyGameVendorService.GetDetailUrlByTransactionNo(userID, transactionNo);
            return result;
        }


    }
}
