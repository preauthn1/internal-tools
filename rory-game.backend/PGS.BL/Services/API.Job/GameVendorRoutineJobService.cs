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
using target.DA.Repository;
using System.Collections.Concurrent;
using target.Models;
using target.Utility.Helper;
using target.Models.ViewModels.GameVendor;
using target.Utility.Utility.ThirdParty;

namespace target.BL.Services
{

    public class GameVendorRoutineJobService : BaseService, IGameVendorRoutineJobService
    {
        public GameListRepository GameListRepository { get; set; }

        public MerchantGameSummaryDataRepository MerchantGameSummaryDataRepository { get; set; }

        public LanguageTextRepository LanguageTextRepository { get; set; }


        public void GetVendorGames(GameVendorTypeEnum vendorType)
        {
            int vendorID = vendorType.ToInteger();

            var vendorService = ServiceInstanceHelper.InstanceVendorService(vendorID);

            //透過三方API 取得遊戲清單
            var queryGamesByApi = vendorService.GetGames(SystemConstant.DefaultMerchantID);

            if (queryGamesByApi.IsEmptyOrNull())
                return;

            var queryVendorGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                GameVendorTypes = vendorID.ToListCollection()
            });

            //如果遊戲商沒有任何遊戲 或 沒有預設遊戲編號 預設建一筆
            if (queryVendorGames.IsEmptyOrNull() || !queryVendorGames.Any(x => x.GameID == vendorID))
            {
                string gameCode = ((GameVendorTypeEnum)vendorID).ToString().ToUpper() + "_default";

                string languageCode = AppHelper.GetGameListLanguageCode(gameCode);

                GameListRepository.Create(new GameListEntity
                {
                    LanguageCode = languageCode,
                    GameID = vendorID, //預設第一筆直接跟vendor編號一樣
                    CreateTime = DateTimeHelper.GetCurrentTime(),
                    GameCode = gameCode,
                    ThirdPartyGameCode = gameCode,
                    GameName = gameCode,
                    GameVendorType = vendorID,
                    GameType = GameTypeEnum.None.ToInteger(),
                    Status = NormalActiveStatus.Enable.ToInteger(),
                    IsSkipInsertZeroBet = false
                });
            }

            List<string> existGameCodes = queryVendorGames.Select(x => x.ThirdPartyGameCode).ToList();

            var notExistGames = queryGamesByApi.Where(x => !existGameCodes.Contains(x.ThirdPartyGameCode)).ToList();

            //TODO 暫不考慮遊戲下架
            foreach (var item in notExistGames)
            {
                int gameID = GameListRepository.GetMaxIdByVendor(vendorID);

                string languageCode = AppHelper.GetGameListLanguageCode(item.GameCode);

                GameListRepository.Create(new GameListEntity
                {
                    LanguageCode = languageCode,
                    GameID = gameID + 1,
                    CreateTime = DateTimeHelper.GetCurrentTime(),
                    GameCode = item.GameCode,
                    ThirdPartyGameCode = item.ThirdPartyGameCode,
                    GameName = item.GameName,
                    GameVendorType = vendorID,
                    GameType = GameTypeEnum.None.ToInteger(),
                    Status = NormalActiveStatus.Enable.ToInteger(),
                    IsSkipInsertZeroBet = false
                });
            }

            List<LanguageTextEntity> insertItems = new List<LanguageTextEntity>();

            foreach (var item in queryGamesByApi)
            {
                string gameCode = AppHelper.GetGameListLanguageCode(item.GameCode);

                foreach (var lang in item.LanguageText)
                {
                    insertItems.Add(new LanguageTextEntity
                    {
                        Language = lang.Key.ToString(),
                        LanguageCode = gameCode,
                        Text = lang.Value
                    });
                }
            }

            //新增多語系文字
            LanguageTextRepository.CreateNotExist(insertItems);
        }


        public void CalculateYesterdayReportIsCorrect(GameVendorTypeEnum vendorType)
        {
            var vendorService = ServiceInstanceHelper.InstanceVendorService(vendorType);

            double timeZone = vendorService.GetVendorTimeZone();

            //找出特定時區的昨日
            DateTime beginDate = DateTimeHelper.GetCurrentTime(timeZone).AddDays(-1).Date;
            DateTime endTime = beginDate.AddDays(1);


            CalculateDateRangeReportIsCorrect(vendorType, beginDate , endTime);
        }

      
        public void CalculateTargetDateReportIsCorrect(GameVendorTypeEnum vendorType, DateTime targetDate , double timeZone)
        {
            //將指定時區的特定日期 轉成UTC
            //e.g  +8時區的 11/10 00:00:00 等於 +0時區的 11/09 16:00:00
            DateTime beginUTC = targetDate.Date.AddHours(timeZone * -1);
            DateTime endUTC = beginUTC.AddDays(1);

            CalculateDateRangeReportIsCorrect(vendorType, beginUTC, endUTC);
        }


        public void CalculateDateRangeReportIsCorrect(GameVendorTypeEnum vendorType, DateTime beginUTC, DateTime endUTC)
        {
            int vendorID = vendorType.ToInteger();

            var vendorService = ServiceInstanceHelper.InstanceVendorService(vendorID);

            var queryGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                GameVendorTypes = vendorID.ToListCollection()
            });

            List<int> merchantIDs = new List<int>();
            List<int> gameIDs = queryGames.Select(z => z.GameID).ToList();

            //每家遊戲商的時區依據不同
            double vendorTimeZone = vendorService.GetVendorTimeZone();

            // step 1. 先算自己DB的統計數據
            var queryPgsSummary = MerchantGameSummaryDataRepository.GetSummaryForVendorReportDataCheck(merchantIDs, gameIDs, beginUTC, endUTC);


            string requestUUID = Guid.NewGuid().ToString();

            // step 2. 從API取得統計數據
            ExternalVendorGetReportResultModel queryApiReport = vendorService.GetDbAndApiSummaryReport(new ExternalVendorGetReportRequestParameter
            {
                MerchantID = SystemConstant.DefaultMerchantID,
                StartTimeUTC = beginUTC,
                EndTimeUTC = endUTC,
                StartTimeTargetTimeZone = beginUTC.AddHours(vendorTimeZone),
                EndTimeTargetTimeZone = endUTC.AddHours(vendorTimeZone),
                DebugRequestUUID = requestUUID,
            });


            //檢查PGS統計數據 是否跟 API對得起來
            bool isCorrect1 = CheckIsDataCorrect(queryPgsSummary, queryApiReport.ApiSummary);

            //檢查原始投注資料 是否跟 API對得起來
            bool isCorrect2 = CheckIsDataCorrect(queryApiReport.DbSummary, queryApiReport.ApiSummary);

            //如果判斷到帳對不起來
            if (!isCorrect1 || !isCorrect2)
            {
                var vendor = (GameVendorTypeEnum)vendorID;
                string msg = $" {vendor.ToString()} 對帳失敗，UUID: `{requestUUID}` " + Environment.NewLine;
                msg += $"時間區間 {beginUTC.ToString("yyyy-MM-dd HH:mm:ss")} - {endUTC.ToString("yyyy-MM-dd HH:mm:ss")}" + Environment.NewLine;
                msg += $"- TotalBetCount：{queryPgsSummary.TotalBetCount}｜{queryApiReport.DbSummary.TotalBetCount}｜{queryApiReport.ApiSummary.TotalBetCount}" + Environment.NewLine;
                msg += $"- TotalWinLoss：{queryPgsSummary.TotalWinLoss}｜{queryApiReport.DbSummary.TotalWinLoss}｜{queryApiReport.ApiSummary.TotalWinLoss}" + Environment.NewLine;
                msg += $"- TotalBetAmount：{queryPgsSummary.TotalBetAmount}｜{queryApiReport.DbSummary.TotalBetAmount}｜{queryApiReport.ApiSummary.TotalBetAmount}" + Environment.NewLine;
                msg += $"- TotalValidBetAmount：{queryPgsSummary.TotalValidBetAmount}｜{queryApiReport.DbSummary.TotalValidBetAmount}｜{queryApiReport.ApiSummary.TotalValidBetAmount}" + Environment.NewLine;

                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                {
                    Channel = SlackChannelEnum.n_error_public,
                    Message = msg
                });
            }

            //TODO 重跑API數據....?
            //TODO 重跑stat統計數據
        }


        private bool CheckIsDataCorrect(BetRecordSummaryReportModel betSummary, BetRecordSummaryReportModel compareTarget)
        {
            bool isEqual = true;

            //取得兩者交集
            List<GameVendorReportApiAvailableColumnsEnum> intersectColumns = betSummary.AvailableColumns.Intersect(compareTarget.AvailableColumns).ToList();

            //AvailableColumns 定義了這家遊戲商 對帳行為要檢查那些欄位 (可能A遊戲商 只能檢查輸贏、B遊戲商可檢查三個欄位)
            foreach (var each_column in intersectColumns)
            {
                switch (each_column)
                {
                    case GameVendorReportApiAvailableColumnsEnum.TotalBetCount:
                        isEqual = betSummary.TotalBetCount == compareTarget.TotalBetCount;
                        break;
                    case GameVendorReportApiAvailableColumnsEnum.TotalBetAmount:
                        isEqual = betSummary.TotalBetAmount == compareTarget.TotalBetAmount;
                        break;
                    case GameVendorReportApiAvailableColumnsEnum.TotalWinLoss:
                        isEqual = betSummary.TotalWinLoss == compareTarget.TotalWinLoss;
                        break;
                    case GameVendorReportApiAvailableColumnsEnum.TotalValidBetAmount:
                        isEqual = betSummary.TotalValidBetAmount == compareTarget.TotalValidBetAmount;
                        break;
                    default:
                        throw new TerproException(ApiErrorCode.NotImplement);
                }

                if (!isEqual)
                    break;
            }

            return isEqual;
        }





    }
}
