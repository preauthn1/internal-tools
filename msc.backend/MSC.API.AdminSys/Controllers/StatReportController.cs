using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using MSC.BL.Interface;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Web.Modules.Filter;
using MSC.Models.Paging;
using MSC.Models.Enums;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class StatReportController : BaseController
    {
        public IStatReportService StatReportService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Report_DailySummary_List)]
        [Route("report/stat/merchant-daily-summary")]
        [HttpGet]
        public StatReportViewModel GetMerchantDailySummary([FromQuery]StatReportRequestModel filter)
        {
            int userID = GetCurrentUserID();
            var query = StatReportService.GetMerchantDailySummary(userID, filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Merchant_MerchantSummary)]
        [Route("report/merchant/summary/my")]
        [HttpGet]
        public MerchantSummaryInfoViewModel GetMyMerchantSummaryReport()
        {
            int userID = GetCurrentUserID();
            var query = StatReportService.GetDefaultMerchantSummaryInfo(userID);
            return query;
        }



        [Route("report/merchant/{merchantID}/summary")]
        [HttpGet]
        public MerchantSummaryInfoViewModel GetMerchantSummaryReport(int merchantID)
        {
            int userID = GetCurrentUserID();
            var query = StatReportService.GetMerchantSummaryInfo(userID, merchantID);
            return query;
        }


    }
}
