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
using target.Models.Models.StoreProcedures;

namespace target.BL.Services
{

    public class GameHealthSummaryService : BaseAdminSysService, IGameHealthSummaryService
    {
        public GameListRepository GameListRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }


        public GameHealthSummaryViewModel GetEachServerHealthySummary(int userID , GameHealthSummaryRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();

            var queryMerchant = MerchantRepository.GetById(model.MerchantID);

            if (queryMerchant == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            var timeZoneToday = DateTimeHelper.GetTimeZoneToday(model.TimeZone);

            var query = GameListRepository.GetGameHealthResult(timeZoneToday.TodayBeginUTC , timeZoneToday.TodayEndUTC , model.MerchantID.ToListCollection());

            List<GameHealthSummaryItemViewModel> viewModels = ConvertToViewModels(query);        

            return new GameHealthSummaryViewModel
            {
                Items = viewModels,
                RowTotal = viewModels.CalculateTotalSumMergeToOneObject()
            };
        }


        private List<GameHealthSummaryItemViewModel> ConvertToViewModels(List<Stat_GameHealthSummarySPResultModel> queryGameHealthSummary)
        {
            return queryGameHealthSummary.Select(x =>
                {
                    var convert = x.ReturnInjecteObject<GameHealthSummaryItemViewModel>();
                    return convert;
                })
                .ToList();
        }

    }
}
