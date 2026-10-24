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


namespace target.BL.Services
{

    public class ReportService : BaseAdminSysService, IReportService
    {
        public StatReportStoreProcedureRepository StatReportStoreProcedureRepository = new StatReportStoreProcedureRepository();

        public List<HomePageSummaryViewModel> GetHomePageSummary(int userID, QueryHomePageReportApiRequestParameterModel model)
        {
            if (model.Begin < DateTimeHelper.GetDefaultTime())
                model.Begin = DateTimeHelper.GetCurrentTime().To_00_00_00();
            if (model.End < DateTimeHelper.GetDefaultTime())
                model.End = DateTimeHelper.GetCurrentTime().To_23_59_59();
            model.RuleFor(x => x.Begin).Rule_StartTimeEndTimeRange(model.End, 1);

            var filter = model.ReturnInjecteObject<QueryHomePageReportRequestViewModel>();
            filter.UserID = userID;
            filter.MerchantIDs = GetCurrentOwnPermissionByCache(userID).AvailableMerchantIDs;
            var query = StatReportStoreProcedureRepository.GetHomePageSummary(filter);
            return query;
        }
    }
}
