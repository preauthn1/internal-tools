using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.Paging;
using target.Web.Modules.Filter;
using target.API.AdminSys;

namespace target.API.Vendor.Controllers
{
    [Route("v{version:int}")]
    public class BetRecordController : BaseController
    {
        public IGameRecordService GameRecordService { get; set; }


        /// <summary>
        /// 3.1.1 Get betRecords (取得投注记录)
        /// </summary>
        [SwaggerApiDescriptionEN("This API can only be access `7200 times` in one hour <BR> Records older than 90 days are not available.")]
        [SwaggerApiDescriptionCN("此API接口 每小时限制存取 `7200次` <BR> 不得查询超过 `90天以前` 的资料")]
        [Route("game-records")]
        [ApiExplorerSettings(GroupName = "tfMode")]
        [HttpGet]
        public PagedList<GameRecordViewModel> GetGameRecords([FromQuery]GameRecordRequestModel model)
        {
            var query = GameRecordService.GetBetRecords(model);
            return query;
        }


        /// <summary>
        /// 3.1.2 Get betRecord by transactionNo (根据投注单号取得投注记录) (SeamlessMode) 
        /// </summary>
        [SwaggerApiDescriptionEN("Get betRecord by transactionNo")]
        [SwaggerApiDescriptionCN("根据投注单号取得投注记录")]
        [Route("game-records/transactionNo")]
        [ApiExplorerSettings(GroupName = "swMode")]
        [ApiConfig(IsResponseLogEnable = false)]
        [HttpGet]
        public SingleWalletGameRecordViewModel GetGameRecordByTransactionNo([FromQuery] SingleWalletGameRecordRequestModel model)
        {
            var query = GameRecordService.GetBetRecordBySingleWalletTransactionNo(model);
            return query;
        }


        /// <summary>
        /// 3.1.3 Get merchant statistic report (取得商户总报表资讯)
        /// </summary>
        [SwaggerApiDescriptionEN("Get merchant statistic report.")]
        [SwaggerApiDescriptionCN("取得投注数据报表")]
        [Route("game-records/report")]
        [HttpGet]
        public MerchantReportSummaryViewModel GetReportSummary([FromQuery] VendorReportSummaryRequestModel model)
        {
            var query = GameRecordService.GetReportSummary(model);
            return query;
        }


        /// <summary>
        /// (這隻API不完整) 先不對外
        /// 產生給外部查看 盤面結果的 查看頁網址
        /// </summary>
        [Route("game-records/detail-info/link")]
        [ApiExplorerSettings(IgnoreApi = true)]
        [HttpGet]
        public BetRecordDetailInfoLinkGenerateViewModel GenerateDetailInfoLink([FromQuery] BetRecordDetailInfoLinkGenerateRequestModel model)
        {
            var query = GameRecordService.GenerateDetailInfoLink(model);
            return query;
        }

    }
}
