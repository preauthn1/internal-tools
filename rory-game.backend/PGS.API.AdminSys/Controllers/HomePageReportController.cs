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
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class HomePageReportController : BaseController
    {
        public IReportService ReportService { get; set; }

        /// <summary>
        /// 當日基本訊息
        /// </summary>
        /// <param name="model"></param>
        /// <returns></returns>
        [Route("report/home-page-summary")]
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.Home)]
        [HttpGet]
        public List<HomePageSummaryViewModel> GetHomePageSummary([FromQuery]QueryHomePageReportApiRequestParameterModel model)
        {
            int userID = GetCurrentUserID();
            var query = ReportService.GetHomePageSummary(userID, model);
            return query;
        }
    }
}
