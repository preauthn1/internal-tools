using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.ViewModels.Base;
using target.Utility.Utility;
using target.Utility.Extensions;
using System.Linq;
using target.Models.Models;
using target.Utility;
using target.DA.Repository;
using target.BL.Interface;
using target.Models.ViewModels.GameVendor;
using target.Utility.Helper;
using target.Utility.Helper.HttpContext;
using System.Transactions;

namespace target.BL.Services
{

    public class FunkyIntegrationService : BaseGameVendorIntegrationService, IGameVendorIntegrationService
    {

        public override GameVendorTypeEnum GameVendor
        {
            get
            {
                return GameVendorTypeEnum.Funky;
            }
        }


        public override string GameVendorPrefix
        {
            get
            {
                return "fun";
            }
        }

        public override double TimeZone
        {
            get
            {
                return 0;
            }
        }


        public BetRecordFunkyRepository BetRecordRepository { get; set; }


        public SingleWalletExceptionHandleModel FormatParsing_ParseExceptionModelToVendorApiResponse(ApiResponseErrorModel model)
        {
            //funky的API錯誤碼 有定義了10組左右的errorCode 末三碼會轉成Funky文件定義的編號 (請參考文件 错误代码列表)
            ApiErrorCode defaultError = ApiErrorCode.FunkyApi_InvalidInput;

            int errorCode = model.Code % 1000;

            if (errorCode < 400 || errorCode > 415)
            {
                errorCode = defaultError.ToInteger() % 1000;
            }

            if (model.Code == (int)ApiErrorCode.FunkyApi_GameCodeNotDefine ||
               model.Code == (int)ApiErrorCode.NotFoundGame)
            {
                //遊戲未定義
                errorCode = 10005;
            }
            else if (model.Code == (int)ApiErrorCode.NotFoundBetRecord)
            {
                 //遊戲未定義
                errorCode = 404;
            }
            else if (model.Code == (int)ApiErrorCode.WalletTranferFail_WalletAmountNotEnough)
            {
                //錢包餘額不足
                errorCode = 402;
            }

            return new SingleWalletExceptionHandleModel
            {
                StatusCode = 200,
                GameVendorCustomizedResponse = new FunkyAPIResponseModel<string>
                {
                    ErrorCode = errorCode,
                    ErrorMessage = ((ApiErrorCode)model.Code).ToString(),
                    Data = null
                }
            };
        }

        public object FormatParsing_ConvertToPlaceBetResponseResult(SingleWalletPlaceBetDebitAmountResponseParameter placeBetInfo)
        {
            return new FunkyAPIResponseModel<FunkySingleWalletAPIPlaceBetResultViewModel>
            {
                Data = new FunkySingleWalletAPIPlaceBetResultViewModel
                {
                    Balance = placeBetInfo.AfterAmount
                }
            };
        }


        public object FormatParsing_ConvertToSettleBetResponseResult(SettlementConvertRequiredInfoModel model)
        {
            return new FunkyAPIResponseModel<FunkySingleWalletAPISettleBetResponseModel>
            {
                Data = new FunkySingleWalletAPISettleBetResponseModel
                {
                    RefNo = model.BetRecordNo,
                    Balance = model.AfterAmount,
                    PlayerId = model.PlayerID.ToString(),
                    StatementDate = DateTimeHelper.GetCurrentTime().ToString("yyyy-MM-dd")
                }
            };
        }


        public PerGameDetailLogEntity FormatParsing_ConvertBetRecordToPerGameDetailLog(string betRecordNoPGS)
        {
            var queryBetRecord = BetRecordRepository.GetByBetRecordNoPGS(betRecordNoPGS);

            if (queryBetRecord == null)
            {
                throw new TerproException(ApiErrorCode.NotFoundData);
            }

            if (queryBetRecord.Status != (int)GameVendorBetRecordStatusEnum.SettleComplete)
            {
                throw new TerproException(ApiErrorCode.NotFoundData);
            }

            var queryPlayer = PlayerRepository.GetBriefByCache(queryBetRecord.PlayerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var convert = new PerGameDetailLogEntity
            {
                BetAmount = queryBetRecord.BetAmount,
                ValidAmount = queryBetRecord.ValidAmount,
                UpdateAmount = queryBetRecord.WinLossAmount,
                PayoutAmount = queryBetRecord.PayoutAmount,
                TaxAmount = 0,
                ServerID = 0,
                GameID = queryBetRecord.GameID,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                CurrentAmount = 0, //後半段的code再補上
                GameTime = queryBetRecord.SettlementTime.Value,
                MerchantID = queryBetRecord.MerchantID,
                PlayerID = queryBetRecord.PlayerID,
                PlayerName = queryPlayer.PlayerName,
                TradeID = 0,
                ThirdPartyBetRecordNo = queryBetRecord.BetRecordNo,
                TransactionNo = queryBetRecord.BetRecordNoPGS
            };

            return convert;
        }



        public void WriteRawBetRecordWhenPlaceBet(InsertRawBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        {
            rawInfo.RuleFor(x => x.Funky).Rule_ModelNotNull();

            var placeBetInfo = rawInfo.Funky.PlaceBetInfo;

            var insert = new BetRecordFunkyEntity
            {
                GameID = model.GameID,
                PlayerID = model.PlayerID,
                MerchantID = model.MerchantID,
                BetAmount = placeBetInfo.Bet.Stake,
                WinLossAmount = 0,
                ValidAmount = 0,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                BetRecordNo = placeBetInfo.Bet.RefNo,
                PlayerIPAddress = placeBetInfo.PlayerIp,
                PlayerName = model.PlayerName,
                Status = 0,
                ThirdPartyGameCode = placeBetInfo.Bet.GameCode,
                UpdateTime = DateTimeHelper.GetCurrentTime(),
                SettlementTime = null,
                BetRecordNoPGS = model.BetRecordNoPGS
            };

            BetRecordRepository.Create(insert);
        }


        public void UpdateRawBetRecordAfterSettlement(SettlementBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        {
            rawInfo.RuleFor(x => x.Funky).Rule_ModelNotNull();

            var settleInfo = rawInfo.Funky.SettleBetInfo;

            model.PayoutAmount = settleInfo.betResultReq.winAmount;
            model.ValidAmount = settleInfo.betResultReq.effectiveStake;            

            UpdateRawBetRecordAfterSettlement(model);
        }


        public void UpdateRawBetRecordAfterSettlement(SettlementBetRecordRequiredInfoModel model)
        {
            model.RuleFor(z => z.BetRecordNoPGS).Rule_NotEmptyText();

            DateTime now = DateTimeHelper.GetCurrentTime();
            string memo = $"processID:{model.DebugExtraMemo}";

            using (var scope = new TransactionScope())
            {
                var queryExist = BetRecordRepository.GetByBetRecordNoPGS(model.BetRecordNoPGS);

                //檢查狀態合法
                RepositoryHelper.CheckRecordNotFoundAndDuplicateStatus(queryExist, GameVendorBetRecordStatusEnum.SettleComplete , model.PayoutAmount);

                queryExist.UpdateTime = now;
                queryExist.DebugInfo += RepositoryHelper.GenerateDebugInfo(queryExist.Status, memo);
                queryExist.Status = GameVendorBetRecordStatusEnum.SettleComplete.ToInteger();
                queryExist.PayoutAmount = model.PayoutAmount;
                queryExist.ValidAmount = model.ValidAmount;
                queryExist.WinLossAmount = model.PayoutAmount - queryExist.BetAmount;
                queryExist.SettlementTime = now;
                queryExist.ManualSettleUserID = model.UserID;

                //特殊遊戲處理，如果是捕魚遊戲 對方會預扣所有錢包金額 導致投注金額/派彩金額 等於 錢包金額
                //因此投注金額 改成使用 validAmount & 派彩金額則是 輸贏+有效流水
                if (queryExist.ThirdPartyGameCode == "10211" ||
                    queryExist.ThirdPartyGameCode == "10212")
                {
                    decimal payout = queryExist.PayoutAmount - queryExist.BetAmount;
                    queryExist.BetAmount = model.ValidAmount;
                    queryExist.PayoutAmount = model.ValidAmount + queryExist.WinLossAmount;
                }

                BetRecordRepository.Update(queryExist);
                scope.Complete();
            }

        }


        public FunkyAPIResponseModel<FunkySingleWalletAPICheckBetResponseModel> CheckBet(string betRecordNo)
        {
            var queryRecord = BetRecordRepository.GetByBetRecordNo(betRecordNo);

            if (queryRecord == null)
            {
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);
            }

            //R：正在进行
            //W：赢了
            //L：输了
            //C：取消
            //D：平局

            string statusText = string.Empty;

            if (queryRecord.Status == GameVendorBetRecordStatusEnum.Pending.ToInteger())
            {
                statusText = "R";
            }
            else if (queryRecord.Status == GameVendorBetRecordStatusEnum.Cancel.ToInteger())
            {
                statusText = "C";
            }
            else if (queryRecord.WinLossAmount > 0)
            {
                statusText = "W";
            }
            else if (queryRecord.WinLossAmount < 0)
            {
                statusText = "L";
            }
            else
            {
                statusText = "D";
            }

            return new FunkyAPIResponseModel<FunkySingleWalletAPICheckBetResponseModel>
            {
                Data = new FunkySingleWalletAPICheckBetResponseModel
                {
                    Stake = queryRecord.BetAmount,
                    RefNo = queryRecord.BetRecordNo,
                    WinAmount = queryRecord.WinLossAmount,
                    Status = statusText,
                    StatementDate = queryRecord.SettlementTime.HasValue ? queryRecord.SettlementTime.Value.ToString("yyyy-MM-dd") : ""
                }
            };
        }



        public ExternalVendorCancelBetResultModel CancelBet(string betRecordNoPGS, string memo)
        {           
            var queryRecord = BetRecordRepository.GetByBetRecordNoPGS(betRecordNoPGS);

            if (queryRecord == null)
            {
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);
            }

            //根據遊戲商要求 實作不同的errorCode
            if (queryRecord.Status == GameVendorBetRecordStatusEnum.SettleComplete.ToInteger())
            {
                throw new TerproException(ApiErrorCode.FunkyApi_Bet_Was_Already_Settled);
            }

            if (queryRecord.Status == GameVendorBetRecordStatusEnum.Cancel.ToInteger())
            {
                throw new TerproException(ApiErrorCode.FunkyApi_Bet_Was_Already_Cancelled);
            }

            if (queryRecord.Status != GameVendorBetRecordStatusEnum.Pending.ToInteger())
            {
                throw new TerproException(ApiErrorCode.FunkyApi_Bet_Was_Already_Cancelled);
            }

            //更新狀態
            queryRecord.Status = GameVendorBetRecordStatusEnum.Cancel.ToInteger();
            queryRecord.UpdateTime = DateTimeHelper.GetCurrentTime();

            BetRecordRepository.Update(queryRecord);

            return new ExternalVendorCancelBetResultModel
            {
                GameVendorCustomizedResponse = new FunkyAPIResponseModel<FunkySingleWalletCancelBetResponseModel>
                {
                    ErrorCode = 0,
                    ErrorMessage = "Success",
                    Data = new FunkySingleWalletCancelBetResponseModel
                    {
                        RefNo = queryRecord.BetRecordNo,
                    }
                }
            };
        }


        public BetRecordGeneralColumnModel GetBetRecordByNo(string betRecordNo, string gameCode)
        {
            var queryRecord = BetRecordRepository.GetByBetRecordNo(betRecordNo, gameCode);
            return queryRecord.ReturnInjecteObject<BetRecordGeneralColumnModel>();
        }


        public BetRecordGeneralColumnModel GetBetRecordByPGSTransactionNo(string betRecordNoPGS)
        {
            var queryRecord = BetRecordRepository.GetByBetRecordNoPGS(betRecordNoPGS);
            return queryRecord.ReturnInjecteObject<BetRecordGeneralColumnModel>();
        }



        public void CheckIsValidRequestInfo(ExternalVendorInputValidateRequestParameter model)
        {
            string headers = HeaderValueHelper.GetHeaderValueByKey("Authentication");

            //檢查遊戲商的header是否有密鑰
            if (headers != "23d828a8-0bba-49e8-90a6-9007bb50de1f")
                throw new TerproException(ApiErrorCode.Forbidden_NotValidDefineValue);

            if (model.IsSkipCheckSessionToken)
                return;

            model.RuleFor(x => x.SessionOrToken).Rule_NotEmptyText();

            var query = VendorSessionTokenRepository.GetTokenByCache(model.SessionOrToken);

            if (query == null)
            {
                throw new TerproException(ApiErrorCode.NotFoundSessionToken);
            }

            if (model.PlayerID != 0 && query.PlayerID != model.PlayerID)
            {
                throw new TerproException(ApiErrorCode.NotFoundSessionToken);
            }

            if (query.ExpiredTime <= DateTimeHelper.GetCurrentTime())
            {
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);
            }
        }


        public void CreateVendorPlayerAccountOrNot(ExternalVendorCreateAccountRequestParameter model)
        {
            //funky會自動建 不用處理
        }



        public ExternalVendorLaunchUrlViewModel GetGameLaunchUrl(ExternalVendorLaunchUrlRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.PlayerName).Rule_NotEmptyText();

            //參考文件.pdf  5.1  启动游戏 API
            var req = new RequestParameters
            {
                ApiUrl = "/Funky/Game/LaunchGame",
                Method = HttpMethodType.HttpPost,
            };

            var queryMerchant = MerchantRepository.GetSingleByCache(model.MerchantID);

            //apiToken、api域名 吃資料庫設置
            SetRequestInfoByDBSetting(req, model.MerchantID, "LaunchGame", model.DebugReferenceUUID);

            if(model.Language == LanguageEnum.zh_tw)
            {
                model.Language = LanguageEnum.zh_cn;
            }

            var body = new
            {
                gameCode = model.ThirdPartyGameCode,
                userName = model.PlayerName,
                playerId = model.PlayerID,
                currency = model.Currency,
                language = model.Language.ToString().ToUpper(),
                playerIp = IPHelper.GetIP(),
                sessionId = model.ThirdPartyLaunchSessionToken,
                isTestAccount = model.IsTrialMode,
            };

            var res = APIRequestHelper.RequestApiResult<FunkyVendorApiResponseModel<FunkyVendorApiLaunchResponseModel>>(req, body);

            //檢查http錯誤、model解析是否成功
            CheckResponseIsError(res);

            if (res.ResultObject.errorCode != 0)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseError , InstanceErrorLogInfo(res));

            if (res.ResultObject.data.gameUrl.IsNullOrWhiteSpace())
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseDataUnExpected, InstanceErrorLogInfo(res));

            string gameUrl = res.ResultObject.data.gameUrl;
            string token = [REDACTED];

            string finalUrl = $"{gameUrl}?token=[REDACTED]";

            return new ExternalVendorLaunchUrlViewModel
            {
                Url = finalUrl
            };
        }


        private ErrorDetailModel InstanceErrorLogInfo<T>(ApiRequestResultModel<T> res)
        {
            return new ErrorDetailModel
            {
                IsAlwaysNotify = true,
                LogExtraInfo = "FunkyAPI_" + res.RelativeUrl + "_" + res.ResponseText,
            };
        }


        public ThirdPartyGameVendorBetRecordDetailUrlViewModel GetExternalBetDetailInfo(ExternalVendorBetRecordDetailInfoRequestModel model)
        {
            model.RuleFor(z => z.BetRecordNo).Rule_NotEmptyText();
            model.RuleFor(z => z.MerchantID).Rule_NotZeroNotMinus();

            var queryBetRecord = BetRecordRepository.GetByBetRecordNo(model.BetRecordNo , model.ThirdPartyGameCode);

            if (queryBetRecord == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            //參考文件.pdf  8.1  获取赌注详情 API
            var req = new RequestParameters
            {
                ApiUrl = "/Funky/Bet/GetBetDetail",
                Method = HttpMethodType.HttpPost,
            };

            //apiToken、api域名 吃資料庫設置
            SetRequestInfoByDBSetting(req, model.MerchantID, "GetBetDetail");

            var body = new
            {
                gameCode = queryBetRecord.ThirdPartyGameCode,
                playerId = queryBetRecord.PlayerID,
                refNo = queryBetRecord.BetRecordNo,
                lang = "EN"
            };

            var res = APIRequestHelper.RequestApiResult<FunkyVendorApiResponseModel<FunkyVendorApiGetDetailUrlResponseModel>>(req, body);

            //檢查http錯誤、model解析是否成功
            CheckResponseIsError(res);

            if (res.ResultObject.errorCode != 0)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseError);

            return new ThirdPartyGameVendorBetRecordDetailUrlViewModel
            {
                InfoType = ThirdPartyVendorInfoTypeEnum.ExternalLink,
                Url = res.ResultObject.data.url
            };
        }



        public List<ExternalVendorGameListConvertModel> GetGames(int merchantID)
        {
            var queryEN = GetGames(merchantID, "EN");
            var queryCN = GetGames(merchantID, "ZH_CN");

            List<ExternalVendorGameListConvertModel> listResult = new List<ExternalVendorGameListConvertModel>();

            foreach (var item in queryEN)
            {
                string gameCode = GameVendorPrefix + "_" + item.gameCode;

                var newInstance = new ExternalVendorGameListConvertModel
                {
                    GameName = item.gameName,
                    GameCode = GameVendorPrefix + "_" + item.gameCode,
                    GameVendorType = this.GameVendorID,
                    ThirdPartyGameCode = item.gameCode,
                    LanguageText = new Dictionary<LanguageEnum, string>()
                };


                newInstance.LanguageText.Add(LanguageEnum.en_us, item.gameName);

                var matchCN = queryCN.FirstOrDefault(x => x.gameCode == item.gameCode);

                if (matchCN != null)
                {
                    newInstance.LanguageText.Add(LanguageEnum.zh_tw, matchCN.gameName);
                    newInstance.LanguageText.Add(LanguageEnum.zh_cn, AppHelper.SimplifiedToTraditional(matchCN.gameName));
                }

                listResult.Add(newInstance);
            }

            return listResult;
        }


        public List<FunkyVendorApiGetGamesItemModel> GetGames(int merchantID, string languageCode)
        {
            //參考文件.pdf  4.1  获取游戏列表 API
            var req = new RequestParameters
            {
                ApiUrl = "/Funky/Game/GetGameList",
                Method = HttpMethodType.HttpPost,
            };

            //apiToken、api域名 吃資料庫設置
            SetRequestInfoByDBSetting(req, merchantID, "GetGames");

            var body = new
            {
                language = languageCode,
                gameType = 0
            };

            var res = APIRequestHelper.RequestApiResult<FunkyVendorApiRootResponseModel>(req, body);

            //檢查http錯誤、model解析是否成功
            CheckResponseIsError(res);

            if (res.ResultObject.errorCode != 0)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseError);

            if (res.ResultObject.gameList.IsEmptyOrNull())
            {
                return new List<FunkyVendorApiGetGamesItemModel>();
            }

            return res.ResultObject.gameList;
        }


        public ExternalVendorGetReportResultModel GetDbAndApiSummaryReport(ExternalVendorGetReportRequestParameter model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.StartTimeUTC).Rule_NotEmptyDateTime();
            model.RuleFor(x => x.EndTimeUTC).Rule_NotEmptyDateTime();

            //參考文件.pdf 10.1 游戏赢/输报告 API
            var req = new RequestParameters
            {
                ApiUrl = "/Funky/Report/GetConsolidateReport",
                Method = HttpMethodType.HttpPost,
            };

            //apiToken、api域名 吃資料庫設置
            SetRequestInfoByDBSetting(req, model.MerchantID, "GetReport", model.DebugRequestUUID);

            req.BaseUrl = GetVendorSettingByKey(model.MerchantID, VendorConfigSettingEnum.BetReportApiUrl)?.SettingValue;

            //此遊戲商時區為0
            var body = new
            {
                startDate = model.StartTimeTargetTimeZone.ToString("yyyy-MM-dd"),  //如果要查單一日期 起始/結束 都帶同樣的時間
                endDate = model.EndTimeTargetTimeZone.AddDays(-1).ToString("yyyy-MM-dd"),
                includeTestAccount = false,
            };

            var res = APIRequestHelper.RequestApiResult<FunkyApiGetReportResponseModel>(req, body);

            //檢查http錯誤、model解析是否成功
            CheckResponseIsError(res);

            if (res.ResultObject.errorCode != 0)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseError);

            var resConvert = res.ResultObject;

            var queryRawBetRecordSummary = BetRecordRepository.GetTotalSummaryReport(model.StartTimeUTC, model.EndTimeUTC);

            var availableColumns = new List<GameVendorReportApiAvailableColumnsEnum>
                {
                    GameVendorReportApiAvailableColumnsEnum.TotalValidBetAmount,
                    GameVendorReportApiAvailableColumnsEnum.TotalBetCount,
                    GameVendorReportApiAvailableColumnsEnum.TotalWinLoss,
                };

            return new ExternalVendorGetReportResultModel
            {
                ResponseText = res.ResponseText,
                ApiSummary = new BetRecordSummaryReportModel
                {
                    TotalBetCount = res.ResultObject.totalBetCount,
                    TotalWinLoss = res.ResultObject.data.Sum(x => x.winloss),
                    TotalValidBetAmount = res.ResultObject.data.Sum(x => x.betAmount),
                    AvailableColumns = availableColumns
                },
                DbSummary = new BetRecordSummaryReportModel
                {
                    TotalBetCount = queryRawBetRecordSummary.TotalBetCount,
                    TotalPayoutAmount = queryRawBetRecordSummary.TotalPayoutAmount,
                    TotalValidBetAmount = queryRawBetRecordSummary.TotalValidBetAmount,
                    TotalWinLoss = queryRawBetRecordSummary.TotalWinLoss,
                    TotalBetAmount = queryRawBetRecordSummary.TotalBetAmount,
                    AvailableColumns = availableColumns
                },
            };
        }


        private void CheckResponseIsError<T>(ApiRequestResultModel<T> res)
        {
            //檢查http錯誤、model解析是否成功
            if (res.Exception != null)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseError  , InstanceErrorLogInfo(res));

            if (res.ResultObject == null)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ParsingFail , InstanceErrorLogInfo(res));
        }



        private void SetRequestInfoByDBSetting(RequestParameters req, int merchantID, string apiCategory, string requestUUID = "")
        {
            req.LogAction = WriteRequestLog;
            req.AdvancedSetting = new AdvancedSetting
            {
                Log_ApiCategory = apiCategory,
                Log_SubCategory = GameVendor.ToString(),
                Log_RequestUUID = requestUUID
            };

            string apiDomain = GetVendorSettingByKey(merchantID, VendorConfigSettingEnum.ApiUrl)?.SettingValue;
            string opeatorID = GetVendorSettingByKey(merchantID, VendorConfigSettingEnum.MainAssignVendorID)?.SettingValue;
            string apiToken = [REDACTED], VendorConfigSettingEnum.MainSecretKeyOrApiToken)?.SettingValue;

            if (apiDomain.IsNullOrWhiteSpace() || opeatorID.IsNullOrWhiteSpace() || apiToken.IsNullOrWhiteSpace())
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted);
            }

            req.BaseUrl = apiDomain;
            req.UserAgent = opeatorID;
            req.AddHeaderData("Authentication", apiToken);
            req.AddHeaderData("X-Request-ID", Guid.NewGuid().ToString());
        }


    }
}
