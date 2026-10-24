using System;
using System.Collections.Generic;
using target.BL.Interface;
using target.BL.Services;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Utility;
using target.Utility.Helper;
using target.Models.ViewModels.GameVendor;
using target.Models.Enums;

namespace target.BL.Services
{

    public class GameRecordService : BaseVendorService, IGameRecordService
    {
        public PerGameDetailLogRepository PerGameDetailLogRepository { get; set; }

        public PerGameLogRepository PerGameLogRepository { get; set; }

        public BetRecordRepository BetRecordRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public BetRecordTransactionNoMappingRepository BetRecordTransactionNoMappingRepository { get; set; }

        public MerchantGameSummaryDataRepository MerchantGameSummaryDataRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }



        public IGameVendorIntegrationService GameVendorIntegrationService { get; set; }


        private const int _maxQuerySize = 10000;


        /// <summary>
        /// 取得投注紀錄、流水紀錄
        /// </summary>
        /// <param name="model"></param>
        /// <returns></returns>
        public PagedList<GameRecordViewModel> GetBetRecords(GameRecordRequestModel model)
        {
            model.RuleFor(x => x.StartTime).Rule_NotEmptyDateTime();
            model.RuleFor(x => x.EndTime).Rule_NotEmptyDateTime();
            
            if (!SettingHelper.IsDev)
            {
                model.RuleFor(x => x.StartTime).Rule_StartTimeEndTimeRange(model.EndTime, 30);
                model.RuleFor(x => x.StartTime).Rule_DateTimeMustLargeThanPreviousDays(90);
            }
                
            this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);

            //最多一次撈N筆
            if (model.PageSize >= _maxQuerySize)
                model.PageSize = _maxQuerySize;

            var filter = model.ReturnInjecteObject<QueryDetailLogRequestViewModel>();
            filter.MerchantIDs = model.MerchantID.ToListCollection();
            filter.PlayerName = model.PlayerName;
            filter.GameStartTimeBegin = model.StartTime;
            filter.GameStartTimeEnd = model.EndTime;

            var query = PerGameDetailLogRepository.GetPagingByFilter(filter);

            var pageModels = ConvertToNewTypePagingObject<PerGameDetailLogEntity, GameRecordViewModel>(query);
            pageModels.Items = query.Items.Select(x => new GameRecordViewModel
                {
                    BetTime = x.GameTime,
                    PlayerName = x.PlayerName,
                    GameID = x.GameID,
                    BeforeAmount = x.CurrentAmount,
                    WinLossAmount = x.UpdateAmount,
                    PerGameLogID = x.PerGameLogID,
                    BetAmount = x.BetAmount,
                    ValidAmount = x.ValidAmount,
                    PayoutAmount = x.PayoutAmount,
                    TaxAmount = x.TaxAmount,
                    BetID = x.Id.ToString(),
                    RoundID = x.PerGameLogID.ToString()
                })
                .ToList();

            //從其他Table補充資訊
            SerPerGameInformation(pageModels.Items);
            return pageModels;
        }



        /// <summary>
        /// 產生給外部查看 盤面結果的 查看頁網址
        /// </summary>
        public BetRecordDetailInfoLinkGenerateViewModel GenerateDetailInfoLink(BetRecordDetailInfoLinkGenerateRequestModel model)
        {
            model.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();

            this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);

            var queryPerGame = PerGameDetailLogRepository.GetByFilter(new QueryDetailLogRequestViewModel
                {
                    TransactionNo = model.TransactionNo
                })
                .FirstOrDefault();

            if (queryPerGame == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            if (queryPerGame.MerchantID != model.MerchantID)
                throw new TerproException(ApiErrorCode.Forbidden_WatchDataPermissionNotAllow);

            var queryMapping = BetRecordTransactionNoMappingRepository.GetByTransactionNo(model.TransactionNo);

            if (queryMapping == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            var simpleToken = [REDACTED] SimpleExternalTokenModel
            {
                TransactionNo = model.TransactionNo,
                RequestID = Guid.NewGuid().ToString("N"),
                ExpiredTime = DateTimeHelper.GetCurrentTime().AddHours(4)
            };

            simpleToken.Sign = AppHelper.GenerateSimpleTokenModelSignKey(simpleToken);

            string tokenBase64 = simpleToken.ConvertToJsonString().Base64Encode();

            var domainSetting = SystemConfigSettingRepository.GetParseModelByEnum(SystemConfigSettingEnum.AdminSysDomain_ExternalDetailInfoUrl);

            if (domainSetting == null)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted);


            bool isExternalVendorGame = AppHelper.IsExternalVendorGame(queryPerGame.GameID);

            if(!isExternalVendorGame)
            {
                //PGS自製遊戲                 
                string pageUrl = $"{domainSetting.SettingValue}/view/betHistory".ReplaceDuplicateUrlChar();
                string finalUrl = $"{pageUrl}?transactionNo={model.TransactionNo}&token=[REDACTED]";

                return new BetRecordDetailInfoLinkGenerateViewModel
                {
                    Url = finalUrl
                };
            }
            else
            {
                //三方遊戲
                var queryGame = GameListRepository.GetSingleByCache(queryPerGame.GameID);

                var service = ServiceInstanceHelper.InstanceVendorService(queryMapping.GameVendorID);

                var apiResult = service.GetExternalBetDetailInfo(new ExternalVendorBetRecordDetailInfoRequestModel
                {
                    BetRecordNo = queryMapping.BetRecordNo,
                    MerchantID = queryPerGame.MerchantID,
                    ThirdPartyGameCode = queryGame.ThirdPartyGameCode
                });

                if(apiResult == null || apiResult.Url.IsNullOrWhiteSpace())
                {
                    throw new TerproException(ApiErrorCode.Forbidden_NotSupportExternalLinkDisplay);
                }

                return new BetRecordDetailInfoLinkGenerateViewModel
                {
                    Url = apiResult.Url
                };
            }          
        }



        /// <summary>
        /// (單錢包模式) 根據投注單號取得資料
        /// </summary>
        public SingleWalletGameRecordViewModel GetBetRecordBySingleWalletTransactionNo(SingleWalletGameRecordRequestModel model)
        {
            model.RuleFor(x => x.TransactionNo).Rule_NotEmptyText();       

            this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);

            var queryTransactionNoMapping = BetRecordTransactionNoMappingRepository.GetByTransactionNo(model.TransactionNo);

            if (queryTransactionNoMapping == null)
                return null;

            var query = GetByPgsTransactionNo(queryTransactionNoMapping.GameVendorID, queryTransactionNoMapping.BetRecordNoPGS);

            //權限不合法
            if (query.MerchantID != model.MerchantID)
                return null;

            var result = new SingleWalletGameRecordViewModel
            {
                BetAmount = query.BetAmount,
                GameID = query.GameID,
                BetTime = query.CreateTime,
                GameName = DbMultiLanguageTextHelper.GetGameDisplayName(query.GameID),
                ValidAmount = query.ValidAmount,
                PayoutAmount = query.PayoutAmount,
                TransactionNo = query.BetRecordNoPGS,
                WinLossAmount = query.WinLossAmount,
                Status = query.Status,
                SettlementTime = query.SettlementTime
            };

            return result;
        }


      
        /// <summary>
        /// 取得總報表
        /// </summary>
        public MerchantReportSummaryViewModel GetReportSummary(VendorReportSummaryRequestModel model)
        {
            model.RuleFor(x => x.StartTime).Rule_NotEmptyDateTime();
            model.RuleFor(x => x.EndTime).Rule_NotEmptyDateTime();

            //正式環境才需要檢查range不得太大
            if(!SettingHelper.IsDev)
            {
                model.RuleFor(x => x.StartTime).Rule_StartTimeEndTimeRange(model.EndTime, 32);
            }

            this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);

            var query = MerchantGameSummaryDataRepository.GetEachGameTotalSumOnDuration(model.MerchantID.ToListCollection() , null, model.StartTime, model.EndTime);

            return new MerchantReportSummaryViewModel
            {
                TotalBetAmount = query.Select(x => x.BetAmount).Sum(),
                TotalValidAmount = query.Select(x => x.ValidAmount).Sum(),
                TotalTaxAmount = query.Select(x => x.TaxAmount).Sum(),
                TotalPayoutAmount = query.Select(x => x.PayoutAmount).Sum(),
                TotalWinLossAmount = query.Select(x => x.WinLossAmount).Sum(),

                GameSummary = query
                    .Select(x=> new MerchantTotalSummaryItemViewModel
                    { 
                        GameID = x.GameID,
                        TotalBetAmount = x.BetAmount,
                        TotalTaxAmount = x.TaxAmount,
                        TotalWinLossAmount = x.WinLossAmount,
                        TotalPayoutAmount = x.PayoutAmount,
                        TotalValidAmount = x.ValidAmount
                    })
                    .ToList()
            };
        }


        private BetRecordGeneralColumnModel GetByPgsTransactionNo(int gameVendorID, string pgsTransactionNo)
        {           
            var service = ServiceInstanceHelper.InstanceVendorService(gameVendorID);
            var result = service.GetBetRecordByPGSTransactionNo(pgsTransactionNo);
            return result;
        }



        private void SerPerGameInformation(List<GameRecordViewModel> viewModels)
        {
            if (viewModels.IsEmptyOrNull())
                return;

            List<long> perGameIDs = viewModels.Select(x => x.PerGameLogID).ToList();

            var queryPerGameLogs = PerGameLogRepository.GetByLargeAmountIDs(perGameIDs);

            foreach(var item in viewModels)
            {
                var match = queryPerGameLogs.FirstOrDefault(x => x.PerGameLogID == item.PerGameLogID);

                if(match != null)
                {
                    item.RoundStartTime = match.BeginTime;
                    item.RoundEndTime = match.EndTime;
                }
            }
        }


    }
}
