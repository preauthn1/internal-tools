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
    public class StatReportController : BaseController
    {
        /// <summary>
        /// 统计资料
        /// </summary>
        public IStatReportService StatReportService { get; set; }

        /// <summary>
        /// 游戏报表
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.StatisticsGame)]
        [Route("report/stat/game")]
        [HttpGet]
        public StatGameReportViewModel GetStatReportPagingByFilter([FromQuery]QueryStatGameReportApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = StatReportService.GetGameStatReport(userID , filter);
            return query;
        }

        /// <summary>
        /// 数据报表
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.StatisticsMerchantSummary)]
        [Route("report/stat/merchant-summary")]
        [HttpGet]
        public MerchantSummaryStatReportViewModel GetMerchantSummaryStatReport([FromQuery]QueryMerchantSummaryStatReportReportApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = StatReportService.GetMerchantSummaryStatReport(userID , filter);
            return query;
        }



        /// <summary>
        /// 数据报表
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.StatisticsMerchantSummary)]
        [Route("report/stat/merchant-daily-summary")]
        [HttpGet]
        public MerchantDailySummaryStatReportInfoViewModel GetMerchantDailySummaryStatReport([FromQuery]QueryMerchantDailySummaryStatReportReportApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = StatReportService.GetMerchantDailySummaryStatReport(userID, filter);
            return query;
        }



        /// <summary>
        /// 数据报表-圖
        /// </summary>
        /// <param name="filter"></param>
        /// <returns></returns>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.Statistics)]
        [Route("report/stat/chartdata/{chartType}")]
        [HttpGet]  
        public ChartDataStatReportViewModel GetStatReportChartData(ChartDataTypeEnum chartType, [FromQuery]QueryChartDataStatReportApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = StatReportService.GetStatReportChartData(userID, chartType, filter);
            return query;
        }
    }
}
