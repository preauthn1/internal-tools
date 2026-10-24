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
using System.ComponentModel;
using target.Models;
using Newtonsoft.Json;
using System.Text.Json;
using System.Transactions;

namespace target.BL.Services
{

    public class AGGIntegrationService : BaseGameVendorIntegrationService, IGameVendorIntegrationService
    {

        public override GameVendorTypeEnum GameVendor
        {
            get
            {
                return GameVendorTypeEnum.AGG;
            }
        }


        public override string GameVendorPrefix
        {
            get
            {
                return "agg";
            }
        }

        public override double TimeZone
        {
            get
            {
                return 0;
            }
        }


        public BetRecordAGGRepository BetRecordRepository { get; set; }


        [Description("OK")]
        public SingleWalletExceptionHandleModel FormatParsing_ParseExceptionModelToVendorApiResponse(ApiResponseErrorModel model)
        {

            //RS_OK - Successful response.
            //RS_ERROR_UNKNOWN - General error status, for cases without a special error code.
            //RS_ERROR_INVALID_USER_TOKEN - UserToken invalid(or expired) to Integrator's system.
            //RS_ERROR_INVALID_ACCOUNT_TOKEN - AccountToken invalid(or expired) to Integrator's system.
            //RS_ERROR_BANNED_ACCOUNT - User is disabled / locked / banned and can't place bets.
            //RS_ERROR_NOT_ENOUGH_BALANCE - Not enough money(or item) on User's balance.
            //RS_ERROR_WRONG_CURRENCY - User specified currency is invalid or differs from User's wallet currencies.
            //RS_ERROR_DUPLICATE_WAGER - A wager with same id was sent, but there was a different user or game.
            //RS_ERROR_WAGER_DOES_NOT_EXIST - Returned when the wager can't be found on Integrator's side.
            //RS_ERROR_ACCOUNT_DOES_NOT_EXIST - Returned when the account can't be found on Integrator's side.
            //RS_ERROR_INVALID_WAGER_OPERATION - Invalid Operation for this wager's state.
            //RS_ERROR_FORBIDDEN_ITEM - 禁止使用的道具(未達使用條件、未領取、過期、etc..)


            string AGGResponseCode = "RS_ERROR_UNKNOWN";

            //拋錯誤 但是強制回應成功給AGG
            if (model.Code == (int)ApiErrorCode.InvalidStatus_BetRecordDuplicateSettlement)
            {
                return new SingleWalletExceptionHandleModel
                {
                    StatusCode = 200,
                    GameVendorCustomizedResponse = new AGGSingleWalletAPIPlaceBetResponseViewModel
                    {
                        StatusCode = SystemConstant.AGGGameVendorDefaultSuccessResponse,
                        Balances = new List<AGGSingleWalletEachBalanceItemViewModel>
                        {
                            new AGGSingleWalletEachBalanceItemViewModel
                            {
                                Amount = 0,
                                CurrencyCode = "CNY"
                            }
                        }
                    }
                };
            }

            if (model.Code == (int)ApiErrorCode.FunkyApi_GameCodeNotDefine ||
               model.Code == (int)ApiErrorCode.NotFoundGame)
            {
                //遊戲未定義
                AGGResponseCode = "RS_ERROR_FORBIDDEN_ITEM";
            }
            else if (model.Code == (int)ApiErrorCode.BusinessLogicDeny_DuplicateDataExist)
            {
                //投注單號已存在
                AGGResponseCode = "RS_ERROR_DUPLICATE_WAGER";
            }
            else if (IsInvalidStatusError((ApiErrorCode)model.Code))
            {
                //狀態不合法
                AGGResponseCode = "RS_ERROR_INVALID_WAGER_OPERATION";
            }
            else if (IsNotFoundBetRecordError((ApiErrorCode)model.Code))
            {
                //查無注單
                AGGResponseCode = "RS_ERROR_WAGER_DOES_NOT_EXIST";
            }
            else if (IsBalanceNotEnughError((ApiErrorCode)model.Code))
            {
                //錢包餘額不足
                AGGResponseCode = "RS_ERROR_NOT_ENOUGH_BALANCE";
            }
            else if (IsAuthenrizeError((ApiErrorCode)model.Code))
            {
                //身分驗證不合法
                AGGResponseCode = "RS_ERROR_INVALID_USER_TOKEN";
            }

            object debugInfo = null;


            if (SettingHelper.IsDev)
            {
                debugInfo = model;
            }

            return new SingleWalletExceptionHandleModel
            {
                StatusCode = 400,
                PGSDebug = debugInfo,
                GameVendorCustomizedResponse = new AGGAPIResponseModel
                {
                    StatusCode = AGGResponseCode
                }
            };
        }


        [Description("OK")]
        public object FormatParsing_ConvertToPlaceBetResponseResult(SingleWalletPlaceBetDebitAmountResponseParameter placeBetInfo)
        {
            return new AGGSingleWalletAPIPlaceBetResponseViewModel
            {
                StatusCode = SystemConstant.AGGGameVendorDefaultSuccessResponse,
                Balances = new List<AGGSingleWalletEachBalanceItemViewModel>
                {
                    new AGGSingleWalletEachBalanceItemViewModel
                    {
                        Amount = placeBetInfo.AfterGamePoint,
                        CurrencyCode = placeBetInfo.Currency
                    }
                }
            };
        }


        [Description("OK")]
        public object FormatParsing_ConvertToSettleBetResponseResult(SettlementConvertRequiredInfoModel model)
        {
            return new AGGSingleWalletAPIPlaceBetResponseViewModel
            {
                StatusCode = SystemConstant.AGGGameVendorDefaultSuccessResponse,
                Balances = new List<AGGSingleWalletEachBalanceItemViewModel>
                {
                    new AGGSingleWalletEachBalanceItemViewModel
                    {
                        Amount = model.AfterGamePoiint,
                        CurrencyCode = model.Currency
                    }
                }
            };
        }


        [Description("OK")]
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

            return new PerGameDetailLogEntity
            {
                ConvertRate = queryBetRecord.ConvertRate,
                BetAmount = queryBetRecord.BetAmount,
                ValidAmount = queryBetRecord.ValidAmount,
                UpdateAmount = queryBetRecord.WinLossAmount,
                PayoutAmount = queryBetRecord.PayoutAmount,
                JackpotAmount = queryBetRecord.JackpotAmount, //AGG特有的欄位
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



        public void WriteRawBetRecordWhenPlaceBet(InsertRawBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        {
            rawInfo.RuleFor(x => x.AGG).Rule_ModelNotNull();

            var placeBetInfo = rawInfo.AGG.PlaceBetInfo;

            decimal betAmount = model.BetAmount;

            var insert = new BetRecordAGGEntity
            {
                GameID = model.GameID,
                PlayerID = model.PlayerID,
                MerchantID = model.MerchantID,
                BetAmount = betAmount,
                WinLossAmount = 0,
                ValidAmount = 0,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                BetRecordNo = placeBetInfo.WagerId,
                PlayerIPAddress = "",
                PlayerName = model.PlayerName,
                Status = 0,
                ThirdPartyGameCode = placeBetInfo.GameId,
                UpdateTime = DateTimeHelper.GetCurrentTime(),
                SettlementTime = null,
                BetRecordNoPGS = model.BetRecordNoPGS,
                ConvertRate = model.ConvertRate
            };

            BetRecordRepository.Create(insert);
        }


        public string AGGAccountTokenConvertToSessionToken(string token)
        {
            //accountToken是我方自訂的格式 {gameID}__{sesssionToken} 其實只是這樣的格式
            var split = token.Split(new string[] { "__" }, StringSplitOptions.None).LastOrDefault();
            return split;
        }


        public decimal GetBetAmountFromOperationModel(List<AGGSingleWalletOperationRequestModel> operations)
        {
            if (operations.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseDataUnExpected);

            decimal betAmount = operations.FirstOrDefault(x => x.OperationType == "WagerWithdrawOperation")?.Withdraw ?? 0;
            return betAmount;
        }


        public void UpdateRawBetRecordAfterSettlement(SettlementBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        {
            rawInfo.RuleFor(x => x.AGG).Rule_ModelNotNull();

            var settleInfo = rawInfo.AGG.SettleBetInfo;

            decimal payoutPoint = GetPayoutAmount(rawInfo.AGG.SettleBetInfo);
            decimal jackpotPoint = GetJackpotAmount(rawInfo.AGG.SettleBetInfo);

            model.PayoutAmount = ConvertGamePointToMoney(model.MerchantID, payoutPoint);
            model.JackpotAmount = ConvertGamePointToMoney(model.MerchantID, jackpotPoint);

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
                RepositoryHelper.CheckRecordNotFoundAndDuplicateStatus(queryExist, GameVendorBetRecordStatusEnum.SettleComplete, model.PayoutAmount);

                queryExist.UpdateTime = now;
                queryExist.DebugInfo += RepositoryHelper.GenerateDebugInfo(queryExist.Status, memo);
                queryExist.Status = GameVendorBetRecordStatusEnum.SettleComplete.ToInteger();
                queryExist.PayoutAmount = model.PayoutAmount;
                queryExist.JackpotAmount = model.JackpotAmount;

                // 1. 結算回調的時候對方沒有提供有效流水 自己算
                // 2. 後台手動結算的時候 可以自己輸入有效流水
                queryExist.ValidAmount = model.IsAdminSysManualOperate ? model.ValidAmount : queryExist.BetAmount;
                queryExist.WinLossAmount = model.PayoutAmount - queryExist.BetAmount;
                queryExist.SettlementTime = now;
                queryExist.ManualSettleUserID = model.UserID;

                BetRecordRepository.Update(queryExist);
                scope.Complete();
            }

        }

        public ExternalVendorCancelBetResultModel CancelBet(string betRecordNoPGS, string memo)
        {
            int merchantID = 0;
            DateTime now = DateTimeHelper.GetCurrentTime();
            using (var scope = new TransactionScope())
            {
                var queryExist = BetRecordRepository.GetByBetRecordNoPGS(betRecordNoPGS);

                //檢查狀態合法
                RepositoryHelper.CheckBetRecordStatusIsWhiteList(queryExist, GameVendorBetRecordStatusEnum.Pending);

                merchantID = queryExist.MerchantID;

                queryExist.UpdateTime = now;
                queryExist.DebugInfo += RepositoryHelper.GenerateDebugInfo(queryExist.Status, memo);
                queryExist.Status = GameVendorBetRecordStatusEnum.Cancel.ToInteger();

                BetRecordRepository.Update(queryExist);
                scope.Complete();
            }

            var queryMerchant = MerchantRepository.GetSingleByCache(merchantID);

            return new ExternalVendorCancelBetResultModel
            {
                GameVendorCustomizedResponse = new AGGSingleWalletAPIPlaceBetResponseViewModel
                {
                    StatusCode = SystemConstant.AGGGameVendorDefaultSuccessResponse,
                    Balances = new List<AGGSingleWalletEachBalanceItemViewModel>
                    {
                        new AGGSingleWalletEachBalanceItemViewModel
                        {
                            Amount = 0,
                            CurrencyCode = queryMerchant.Currency
                        }
                    }
                }
            };
        }


        [Description("OK")]
        public void CheckIsValidRequestInfo(ExternalVendorInputValidateRequestParameter model)
        {
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

        public decimal GetPayoutAmount(AGGSingleWalletAPITransactionRequestModel model)
        {
            return model.Operations.FirstOrDefault(z => z.OperationType == "WagerDepositOperation")?.Deposit ?? 0;
        }


        public decimal GetJackpotAmount(AGGSingleWalletAPITransactionRequestModel model)
        {
            try
            {
                var parse = model.Operations.FirstOrDefault(z => z.OperationType == "WagerTagOperation" && z.TagType == "Money");

                if (parse == null || parse.Tag == null)
                    return 0;

                var matchJackpot = parse.Tag.Actions?.FirstOrDefault(z => z.Action == "WinJp");

                if (matchJackpot == null)
                    return 0;

                return matchJackpot.MoneyAmount;
            }
            catch
            {

            }

            return 0;
        }



        public void CreateVendorPlayerAccountOrNot(ExternalVendorCreateAccountRequestParameter model)
        {
            //遊戲商會自動建 不用處理
        }


        public string GetRecordReplayHistoryByApi(ExternalVendorGetSingleBetRecordRequestParameter model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.BetRecordNo).Rule_NotEmptyText();

            var req = new RequestParameters
            {
                ApiUrl = "game/record/replayByWager",
                Method = HttpMethodType.HttpPost,
            };

            var queryMerchant = MerchantRepository.GetSingleByCache(model.MerchantID);

            //apiToken、api域名 吃資料庫設置
            SetRequestInfoByDBSetting(req, model.MerchantID, "GetBetDetail", model.DebugReferenceUUID);


            var body = new
            {
                gameId = model.GameCode,
                wagerId = model.BetRecordNo,
                lang = "zh-TW",
                currency = new
                {
                    currencyCode = queryMerchant.Currency
                }
            };


            //對應model AGGVendorApiGetDetailResponseModel.cs
            var res = APIRequestHelper.RequestApiResult<string>(req, body);

            //對方API透過303 + redirect的方式回傳
            if (res.ResponseStatusCode == 303)
            {
                return res.ResponseHeaders["Location"];
            }

            //檢查http錯誤、model解析是否成功
            CheckResponseIsError(res);

            return res.ResponseText;
        }


        /// <summary>
        /// 要求AGG再次發出結算通知
        /// </summary>
        public string TriggerSettlementRequestAgain(AGGForceTriggerSettlementRequestReqeustParameter model)
        {
            model.RuleFor(x => x.BetRecordNoPGS).Rule_NotEmptyText();

            var queryRecord = BetRecordRepository.GetByBetRecordNoPGS(model.BetRecordNoPGS);

            if (queryRecord == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            //參考文件.pdf  5.1  启动游戏 API
            var req = new RequestParameters
            {
                ApiUrl = "game​/command",
                Method = HttpMethodType.HttpPost,
            };

            var queryMerchant = MerchantRepository.GetSingleByCache(queryRecord.MerchantID);

            //apiToken、api域名 吃資料庫設置
            SetRequestInfoByDBSetting(req, queryRecord.MerchantID, "GetBetDetail", model.DebugReferenceUUID);

            var body = new
            {
                command = "user.settleUp",
                userId = queryRecord.PlayerID,
                gameId = model.GameCode,
                args = new
                {
                    wagerId = queryRecord.BetRecordNo
                }
            };

            var res = APIRequestHelper.RequestApiResult<string>(req, body);

            //檢查http錯誤、model解析是否成功
            CheckResponseIsError(res);

            return res.ResponseText;
        }


        public ExternalVendorLaunchUrlViewModel GetGameLaunchUrl(ExternalVendorLaunchUrlRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.PlayerName).Rule_NotEmptyText();

            var queryMerchant = MerchantRepository.GetSingleByCache(model.MerchantID);

            string launchUrl = GetVendorSettingByKey(model.MerchantID, VendorConfigSettingEnum.ApiUrl)?.SettingValue;
            string assignVendorID = GetVendorSettingByKey(model.MerchantID, VendorConfigSettingEnum.MainAssignVendorID)?.SettingValue;

            string relativeUrl = model.IsTrialMode ? "/game/trial" : "/game/launch";
            string sessionToken = [REDACTED];
            string gameCode = model.ThirdPartyGameCode;
            string returnUrl = model.BackUrl;

            //backUrl={returnUrl}

            //組出靜態網址
            string finalUrl = $"{launchUrl}{relativeUrl}?iid={assignVendorID}&usertoken=[REDACTED]";

            return new ExternalVendorLaunchUrlViewModel
            {
                Url = finalUrl
            };
        }


        public ThirdPartyGameVendorBetRecordDetailUrlViewModel GetExternalBetDetailInfo(ExternalVendorBetRecordDetailInfoRequestModel model)
        {
            var queryBetRecord = BetRecordRepository.GetByBetRecordNo(model.BetRecordNo, model.ThirdPartyGameCode);

            if (queryBetRecord == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            var apiResult = GetRecordReplayHistoryByApi(new ExternalVendorGetSingleBetRecordRequestParameter
            {
                BetRecordNo = model.BetRecordNo,
                MerchantID = model.MerchantID,
                GameCode = model.ThirdPartyGameCode
            });

            return new ThirdPartyGameVendorBetRecordDetailUrlViewModel
            {
                InfoType = ThirdPartyVendorInfoTypeEnum.ExternalLink,
                Url = apiResult
            };
        }


        public string BeautifyJson(string unPrettyJson)
        {
            var jsonElement = JsonConvert.DeserializeObject<dynamic>(unPrettyJson);
            return JsonConvert.SerializeObject(jsonElement, Formatting.Indented);
        }



        public List<ExternalVendorGameListConvertModel> GetGames(int merchantID)
        {
            //AGG遊戲商沒有提供此API
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public ExternalVendorGetReportResultModel GetDbAndApiSummaryReport(ExternalVendorGetReportRequestParameter model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.StartTimeUTC).Rule_NotEmptyDateTime();
            model.RuleFor(x => x.EndTimeUTC).Rule_NotEmptyDateTime();


            throw new TerproException(ApiErrorCode.NotImplement);


            ////參考文件.pdf 10.1 游戏赢/输报告 API
            //var req = new RequestParameters
            //{
            //    ApiUrl = "/Funky/Report/GetConsolidateReport",
            //    Method = HttpMethodType.HttpPost,
            //};

            ////apiToken、api域名 吃資料庫設置
            //SetRequestInfoByDBSetting(req, model.MerchantID, "GetReport" , model.DebugRequestUUID);

            //req.BaseUrl = GetVendorSettingByKey(model.MerchantID, VendorConfigSettingEnum.BetReportApiUrl)?.SettingValue;

            ////此遊戲商時區為0
            //var body = new
            //{
            //    startDate = model.StartTimeTargetTimeZone.ToString("yyyy-MM-dd"),  //如果要查單一日期 起始/結束 都帶同樣的時間
            //    endDate = model.EndTimeTargetTimeZone.AddDays(-1).ToString("yyyy-MM-dd"),
            //    includeTestAccount = false,
            //};

            //var res = APIRequestHelper.RequestApiResult<FunkyApiGetReportResponseModel>(req, body);

            ////檢查http錯誤、model解析是否成功
            //CheckResponseIsError(res);

            //if (res.ResultObject.errorCode != 0)
            //    throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseError);

            //var resConvert = res.ResultObject;

            //var queryRawBetRecordSummary = BetRecordRepository.GetTotalSummaryReport(model.StartTimeUTC, model.EndTimeUTC);

            //var availableColumns = new List<GameVendorReportApiAvailableColumnsEnum>
            //    {
            //        GameVendorReportApiAvailableColumnsEnum.TotalValidBetAmount,
            //        GameVendorReportApiAvailableColumnsEnum.TotalBetCount,
            //        GameVendorReportApiAvailableColumnsEnum.TotalWinLoss,
            //    };

            //return new ExternalVendorGetReportResultModel
            //{                
            //    ResponseText = res.ResponseText,
            //    ApiSummary = new BetRecordSummaryReportModel
            //    {
            //        TotalBetCount = res.ResultObject.totalBetCount,
            //        TotalWinLoss = res.ResultObject.data.Sum(x => x.winloss),
            //        TotalValidBetAmount = res.ResultObject.data.Sum(x => x.betAmount),
            //        AvailableColumns = availableColumns
            //    },
            //    DbSummary = new BetRecordSummaryReportModel
            //    {
            //        TotalBetCount = queryRawBetRecordSummary.TotalBetCount,
            //        TotalPayoutAmount = queryRawBetRecordSummary.TotalPayoutAmount,
            //        TotalValidBetAmount = queryRawBetRecordSummary.TotalValidBetAmount,
            //        TotalWinLoss = queryRawBetRecordSummary.TotalWinLoss,
            //        TotalBetAmount = queryRawBetRecordSummary.TotalBetAmount,
            //        AvailableColumns = availableColumns
            //    },
            //};
        }


        private void CheckResponseIsError<T>(ApiRequestResultModel<T> res)
        {
            //檢查http錯誤、model解析是否成功
            if (res.Exception != null)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ResponseError);

            if (res.ResultObject == null)
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ParsingFail);
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

            string apiDomain = GetVendorSettingByKey(merchantID, VendorConfigSettingEnum.BetReportApiUrl)?.SettingValue;
            string apiKey = [REDACTED], VendorConfigSettingEnum.MainSecretKeyOrApiToken)?.SettingValue;
            string vanedorID = GetVendorSettingByKey(merchantID, VendorConfigSettingEnum.MainAssignVendorID)?.SettingValue;

            if (apiDomain.IsNullOrWhiteSpace())
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted, new ErrorDetailModel
                {
                    ColumnName = "apiDomain",
                });
            }

            req.UserAgent = "target";
            req.BaseUrl = apiDomain;
            req.AddHeaderData("X-API-KEY", apiKey);
            req.AddHeaderData("X-INTEGRATOR", vanedorID);
        }


    }
}
