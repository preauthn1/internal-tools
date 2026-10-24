using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;

namespace MSC.BL.Services
{

    public class StatReportService : BaseAdminService, IStatReportService
    {
        public UserMerchantMappingRepository UserMerchantMappingRepository { get; set; }

        public StatSPRepository StatSPRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public MerchantWalletRepository MerchantWalletRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }


        public StatReportViewModel GetMerchantDailySummary(int userID, StatReportRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.BeginDate).Rule_NotEmptyDateTime();
            model.RuleFor(x => x.EndDate).Rule_NotEmptyDateTime();

            var timeConvert = ConvertToUTCTimeParameter(model);

            var spResult = StatSPRepository.GetMerchantDailySummaryV2(new MerchantDailySummarySPRequestParameter
            {
                BeginTime = timeConvert.BeginUTC,
                EndTime = timeConvert.EndUTC,
                MerchantID = model.MerchantID,
                TimeZone = model.TimeZone
            });

            var resultItems = new List<StatReportMerchantDailySummaryViewModel>();

            foreach (DateTime date in timeConvert.AllDates)
            {
                var match = spResult.FirstOrDefault(x => x.Date == date);

                var newInstance = new StatReportMerchantDailySummaryViewModel();

                if (match != null)
                {
                    newInstance = new StatReportMerchantDailySummaryViewModel
                    {
                        TotalSuccessCardCount = match.TotalSuccessCount,
                        TotalMemberCount = match.TotalMembers,
                        TotalPoint = match.TotalAmount,
                        TotalOrderCount = match.TotalOrderCount,
                        DepositCompletedRatio = CalculateRate(match.TotalOrderCount , match.TotalSuccessCount)
                    };
                }

                newInstance.Date = date.ToString("yyyy-MM-dd");
                newInstance.MerchantID = model.MerchantID;


                var matchMerchant = MerchantRepository.GetSingleByCache(model.MerchantID);

                newInstance.MerchantName = matchMerchant?.DisplayName;
                newInstance.Currency = matchMerchant?.Currency;

                resultItems.Add(newInstance);
            }

            resultItems = resultItems.OrderByDescending(z => z.Date).ToList();
            resultItems = resultItems ?? new List<StatReportMerchantDailySummaryViewModel>();

            var summary = new StatReportMerchantDailySummaryViewModel
            {
                TotalMemberCount = resultItems.Sum(z => z.TotalMemberCount),
                TotalOrderCount = resultItems.Sum(z => z.TotalOrderCount),
                TotalSuccessCardCount = resultItems.Sum(z => z.TotalSuccessCardCount),
                TotalPoint = resultItems.Sum(z => z.TotalPoint),
                DepositCompletedRatio = resultItems.Where(z => z.TotalOrderCount > 0).DefaultIfEmpty().Average(z => z.DepositCompletedRatio)
            };

            return new StatReportViewModel
            {
                Items = resultItems,
                SummaryInfo = summary
            };
        }

        private decimal CalculateRate(int totalOrderCount , int totalSuccessOrderCount)
        {
            return totalOrderCount <= 0 ? 0 : (((decimal)totalSuccessOrderCount) / ((decimal)totalOrderCount)) * 100;
        }

        public MerchantSummaryInfoViewModel GetDefaultMerchantSummaryInfo(int userID)
        {
            var queryMyMerchants = UserMerchantMappingRepository.GetByFilter(new QueryUserMerchantMappingRequestViewModel
            {
                UserIDs = userID.ToListCollection()
            });

            if(queryMyMerchants.Count() > 1)
            {
                throw new MSCException(ErrorCode.BusinessLogicDeny);
            }

            var viewModel = GetMerchantSummaryInfo(userID , queryMyMerchants.FirstOrDefault().MerchantID);
            return viewModel;
        }


        public MerchantSummaryInfoViewModel GetMerchantSummaryInfo(int userID, int merchantID)
        {
            if(!UserMerchantMappingRepository.IsMerchantPermission(userID , merchantID))
            {
                throw new MSCException(ErrorCode.Forbidden_WatchDataPermissionNotAllow);
            }

            var queryMerchant = MerchantRepository.GetById(merchantID);

            var queryWallet = MerchantWalletRepository.GetByMerchantID(merchantID);

            var querySettings = MerchantCustomizedSettingRepository.GetByFilter(new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
            });

            var amountSetting = querySettings.FirstOrDefault(x => x.CustomizedSettingKey == (int)MerchantCustomizedSettingKeyEnum.Basic_AvailableAmounts);

            List<decimal> avaibleAmounts = GetAvailableAmounts(amountSetting);

            return new MerchantSummaryInfoViewModel
            {
                WalletBalance = queryWallet.Amount,
                MerchantID = merchantID,
                MerchantName = queryMerchant.DisplayName,
                Currency = queryMerchant.Currency,
                Status = queryMerchant.Status,
                StatusText = MultiLanguageTextHelper.GetEnums_i18nText((NormalActiveStatus)queryMerchant.Status),
                TodayDepositAmounts = 0,
                TodayDepositCardCount = 0,
                AvailableAmounts = avaibleAmounts,
                AssignDomain = querySettings.FirstOrDefault(x=> x.CustomizedSettingKey == (int)MerchantCustomizedSettingKeyEnum.Basic_CardInputPage_DomainUrl)?.SettingValue,
            };
        }

        private List<decimal> GetAvailableAmounts(MerchantCustomizedSettingEntity model)
        {
            if (model == null)
                return new List<decimal>();

            var convert = model.ReturnInjecteObject<MerchantCustomizedSettingConvertModel>();

            return convert.SettingValueListDecimal;
        }


        public (DateTime BeginUTC, DateTime EndUTC, List<DateTime> AllDates) ConvertToUTCTimeParameter(StatReportRequestModel model)
        {
            List<DateTime> allDates = ExtensionOfDatetime.GetAllDateBetweenTwoDate(model.BeginDate.Date, model.EndDate.Date);

            if(allDates.Count() > 1)
            {
                //不包含最後1筆  (e.g. 11/1 ~ 12/1  只顯示 11/1 ~ 11/30 的資料)
                allDates.RemoveAt(allDates.Count() - 1);
            }

            DateTime beginUTC = model.BeginDate.Date.AddHours(model.TimeZone * -1);
            DateTime endUTC = model.EndDate.Date.AddHours(model.TimeZone * -1);

            return (beginUTC, endUTC, allDates);
        }


    }
}
