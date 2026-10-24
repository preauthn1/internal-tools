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
using example.invalid;

namespace target.BL.Services
{

    public class StatReportService : BaseAdminSysService, IStatReportService
    {

        public StatReportStoreProcedureRepository StatReportStoreProcedureRepository = new StatReportStoreProcedureRepository();

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }


        public StatGameReportViewModel GetGameStatReport(int userID, QueryStatGameReportApiRequestParameterModel filter)
        {
            var queryCurrentUser = GetCurrentOwnPermissionByCache(userID);

            if (!filter.Begin.HasValue)
                filter.Begin = DateTimeHelper.GetCurrentTime().To_00_00_00();

            if (!filter.End.HasValue)
                filter.End = DateTimeHelper.GetCurrentTime().AddDays(1).To_23_59_59();

            filter.MerchantIDs = FilterAvailablePermissionMerchantIDs(userID, filter.MerchantIDs);
            filter.GameIDs = filter.GameIDs ?? new List<int>();

            List<Stat_GameSummarySPResultModel> query = StatReportStoreProcedureRepository.GetGameSummaryStatReport(filter.Begin.Value, filter.End.Value, filter.MerchantIDs, filter.GameIDs);

            //如果沒篩選merchantIDs 顯示當前用戶能看到的商戶(s)
            if(filter.MerchantIDs.IsEmptyOrNull())
            {
                filter.MerchantIDs = queryCurrentUser.AllEnableMerchantsOrAvailableMerchants;
            }

            var queryMyMerchants = MerchantRepository.GetByIds(filter.MerchantIDs)
                .Select(x => new StatGameReportMerchantItemViewModel
                {
                    MerchantID = x.MerchantID,
                    MerchantName = x.Name,
                    Items = new List<StatGameReportItemViewModel>()
                })
                .ToList();

            var queryAllGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                IsNotDelete = true,
            });

            var queryEachMerchantGames = MerchantGameSettingRepository.GetByFilter(new QueryMerchantGameSettingRequestViewModel
            {
                Status = NormalActiveStatus.Enable.ToInteger(),
                MerchantIDs = filter.MerchantIDs
            });

            foreach (var eachMerchant in queryMyMerchants)
            {
                var matchMyEnableGames = queryEachMerchantGames.Where(x => x.MerchantID == eachMerchant.MerchantID).ToList();

                foreach (var eachGames in queryAllGames)
                {
                    var matchSummaryData = query.FirstOrDefault(x => x.MerchantID == eachMerchant.MerchantID && x.GameID == eachGames.GameID);
                    var matchGameEnable = matchMyEnableGames.FirstOrDefault(x => x.GameID == eachGames.GameID);

                    //如果有資料 或 有開放遊戲 都會顯示數據
                    if (matchSummaryData != null || matchGameEnable != null)
                    {
                        var convert = new StatGameReportItemViewModel();
                        convert.GameID = eachGames.GameID;
                        convert.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(convert.GameID);


                        if (matchSummaryData != null)
                        {
                            convert = matchSummaryData.ReturnInjecteObject<StatGameReportItemViewModel>();
                        }

                        convert.RTP = CommonLogicHelper.CalculateRTP_Percentage(convert.BetAmount, convert.PayoutAmount);

                        eachMerchant.Items.Add(convert);
                    }
                }

                eachMerchant.Items = eachMerchant.Items.OrderBy(x => x.GameID).ToList();
            }

            return new StatGameReportViewModel
            {
                Merchants = queryMyMerchants
            };
        }


        public MerchantSummaryStatReportViewModel GetMerchantSummaryStatReport(int userID, QueryMerchantSummaryStatReportReportApiRequestParameterModel filter)
        {
            if (!filter.Begin.HasValue)
                filter.Begin = DateTimeHelper.GetCurrentTime().To_00_00_00();

            if (!filter.End.HasValue)
                filter.End = DateTimeHelper.GetCurrentTime().AddDays(1).To_23_59_59();

            filter.MerchantIDs = FilterAvailablePermissionMerchantIDs(userID, filter.MerchantIDs);

            List<Stat_MerchantSummarySPResultModel> query = StatReportStoreProcedureRepository.GetMerchantSummaryStatReport(filter.Begin.Value, filter.End.Value, filter.MerchantIDs);

            return new MerchantSummaryStatReportViewModel
            {
                Merchants = query.Select(x => new MerchantSummaryStatReportItemViewModel
                {
                    MerchantID = x.MerchantID,
                    MerchantName = x.MerchantName,
                    BasicSummary = x.ReturnInjecteObject<MerchantSummaryStatReportItemBaseInfoViewModel>(),
                    LoginSummary = x.ReturnInjecteObject<MerchantSummaryStatReportItemLoginInfoViewModel>(),
                    GameSummary = x.ReturnInjecteObject<MerchantSummaryStatReportItemGameInfoViewModel>(),
                })
                .ToList()
            };
        }


        public MerchantDailySummaryStatReportInfoViewModel GetMerchantDailySummaryStatReport(int userID, QueryMerchantDailySummaryStatReportReportApiRequestParameterModel filter)
        {
            filter = filter ?? new QueryMerchantDailySummaryStatReportReportApiRequestParameterModel();

            if (!filter.Begin.HasValue)
                filter.Begin = DateTimeHelper.GetCurrentTime(filter.TimeZone).AddDays(-1).To_00_00_00();

            if (!filter.End.HasValue)
                filter.End = filter.Begin.Value.To_23_59_59();

            //不開放查詢太長
            if ((filter.End.Value - filter.Begin.Value).TotalDays > 93)
            {
                throw new TerproException(ApiErrorCode.InvalidParameter_DateTimeRangeTooLarge);
            }

            List<int> merchantIDs = FilterAvailablePermissionMerchantIDs(userID, filter.MerchantID.ToIntegerListCollection());

            var mercahnts = string.Join(",", merchantIDs);

            var sp_filter = new Stat_MerchantDailySummaryParameterModel
            {
                DateStart = filter.Begin.Value,
                DateEnd = filter.End.Value,
                MerchantList = mercahnts.IsNullOrWhiteSpace() ? null : mercahnts,
                Timezone = filter.TimeZone
            };

            List<MerchantDailySummaryResultModel> query = StatReportStoreProcedureRepository.GetMerchantDailySummary(sp_filter);

            //有時候查 8/31 ~ 8/31 23:59 竟然會吐 9/1 
            query = query
                .Where(x => x.Date <= filter.End.Value)
                .ToList();

            var rowItems = query.Select(x =>
                {
                    var convert = x.ReturnInjecteObject<MerchantDailySummaryStatReportViewModel>();
                    convert.RTP = CommonLogicHelper.CalculateRTP_Percentage(convert.BetAmount, convert.PayoutAmount);
                    convert.Date = x.Date.ToString("yyyy-MM-dd");
                    convert.MerchantWinLossAmount = x.WinLossAmount * -1;

                    return convert;
                })
                .ToList();

            var totalSummary = rowItems.CalculateTotalSumMergeToOneObject();
            totalSummary.RTP = rowItems.Where(z => z.BetAmount > 0).Select(z => z.RTP).DefaultIfEmpty().Average();

            return new MerchantDailySummaryStatReportInfoViewModel
            {
                TotalSummary = totalSummary,
                Items = rowItems
            };
        }

        public Stream ExportDailySummaryExcelByFilter(int userID, QueryMerchantDailySummaryStatReportReportApiRequestParameterModel filter)
        {
            var paging = GetMerchantDailySummaryStatReport(userID, filter);
            var excelStream = ExcelDataConvertHelper.ConvertToMerchantDailySummaryExcelStream(paging.Items);
            return excelStream;
        }


        public ChartDataStatReportViewModel GetStatReportChartData(int userID, ChartDataTypeEnum chartType, QueryChartDataStatReportApiRequestParameterModel model)
        {
            if (!model.Begin.HasValue)
                model.Begin = DateTimeHelper.GetCurrentTime().To_00_00_00();

            if (!model.End.HasValue)
                model.End = DateTimeHelper.GetCurrentTime();

            model.MerchantIDs = FilterAvailablePermissionMerchantIDs(userID, model.MerchantIDs);

            var query = StatReportStoreProcedureRepository.GetStatReportChartData(model.Begin.Value, model.End.Value, model.MerchantIDs, chartType);

            List<ChartDataStatReportItemViewModel> groupBy;

            if (chartType == ChartDataTypeEnum.Game_OnlineUser || chartType == ChartDataTypeEnum.Game_TotalGame)
            {
                groupBy = query.GroupBy(x => x.GameID).Select(x =>
                {
                    int? gameID = x.FirstOrDefault().GameID;
                    var convert = new ChartDataStatReportItemViewModel
                    {
                        GameID = gameID,
                        GameName = gameID.HasValue ? DbMultiLanguageTextHelper.GetGameDisplayName(gameID.Value) : "",
                        Data = query
                            .Where(y => y.GameID == x.FirstOrDefault().GameID)
                            .Select(z => new ChartDataViewModel
                            {
                                StatTime = z.StatTime,
                                DataCount = z.DataCount
                            }).ToList()
                    };

                    return convert;
                }


              ).ToList();
            }
            else
            {
                groupBy = query.GroupBy(x => x.MerchantID).Select(x => new ChartDataStatReportItemViewModel
                {
                    MerchantID = x.FirstOrDefault().MerchantID,
                    MerchantName = x.FirstOrDefault().MerchantName,
                    Data = query
                            .Where(y => y.MerchantID == x.FirstOrDefault().MerchantID)
                            .Select(z => new ChartDataViewModel
                            {
                                StatTime = z.StatTime,
                                DataCount = z.DataCount
                            }).ToList()
                }).ToList();
            }

            var result = new ChartDataStatReportViewModel
            {
                ChartName = MultiLanguageTextHelper.GetEnumsCultureText((ChartDataTypeEnum)chartType),
                Begin = model.Begin.GetValueOrDefault(),
                End = model.End.GetValueOrDefault(),
                ChartDatas = groupBy
            };

            return result;
        }
    }
}
