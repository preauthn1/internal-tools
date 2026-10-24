using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;
using example.invalid;

namespace target.BL.Interface
{

    public interface IStatReportService
    {

        StatGameReportViewModel GetGameStatReport(int userID, QueryStatGameReportApiRequestParameterModel filter);

        MerchantDailySummaryStatReportInfoViewModel GetMerchantDailySummaryStatReport(int userID, QueryMerchantDailySummaryStatReportReportApiRequestParameterModel filter);

        Stream ExportDailySummaryExcelByFilter(int userID, QueryMerchantDailySummaryStatReportReportApiRequestParameterModel filter);

        MerchantSummaryStatReportViewModel GetMerchantSummaryStatReport(int userID, QueryMerchantSummaryStatReportReportApiRequestParameterModel filter);

        ChartDataStatReportViewModel GetStatReportChartData(int userID, ChartDataTypeEnum chartType, QueryChartDataStatReportApiRequestParameterModel filter);
    }
}
