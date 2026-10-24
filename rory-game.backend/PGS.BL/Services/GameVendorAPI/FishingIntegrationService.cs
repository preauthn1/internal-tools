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

    public class FishingIntegrationService : BaseGameVendorIntegrationService, IGameVendorIntegrationService
    {

        public override GameVendorTypeEnum GameVendor
        {
            get
            {
                return GameVendorTypeEnum.CrazyFishing;
            }
        }


        public override string GameVendorPrefix
        {
            get
            {
                return "fish";
            }
        }

        public override double TimeZone
        {
            get
            {
                return 0;
            }
        }

        private static int CrazyFishingGameID = 310001;


        private static string CrazyFishingGameCode = SystemConstant.CrazyFishingDefaultGameCode;


        public BetRecordFishRepository BetRecordRepository { get; set; }



        public SingleWalletExceptionHandleModel FormatParsing_ParseExceptionModelToVendorApiResponse(ApiResponseErrorModel model)
        {
            string msg = "Error! " + model.Message;
            return new SingleWalletExceptionHandleModel
            {
                StatusCode = 200,
                PGSDebug = model,
                GameVendorCustomizedResponse = new FishingAPIResponseModel
                {
                    Code = 9999,
                    Msg = msg
                }
            };
        }


        public object FormatParsing_ConvertToPlaceBetResponseResult(SingleWalletPlaceBetDebitAmountResponseParameter placeBetInfo)
        {
            return new FishingSingleWalletPlaceBetResponseModel
            {
                Code = 0,
                Msg = "OK",
                Balance = placeBetInfo.AfterGamePoint,
            };
        }


        [Description("OK")]
        public object FormatParsing_ConvertToSettleBetResponseResult(SettlementConvertRequiredInfoModel model)
        {
            return new FishingSingleWalletSettleBetResponseModel
            {
                Code = 0,
                Msg = "OK",
                Balance = model.AfterGamePoiint,
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
            rawInfo.RuleFor(x => x.Fishing).Rule_ModelNotNull();

            var placeBetInfo = rawInfo.Fishing.PlaceBetInfo;

            decimal betAmount = model.BetAmount;

            var insert = new BetRecordFishEntity
            {
                GameID = model.GameID,
                PlayerID = model.PlayerID,
                MerchantID = model.MerchantID,
                BetAmount = betAmount,
                WinLossAmount = 0,
                ValidAmount = 0,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                BetRecordNo = placeBetInfo.OrderId,
                PlayerIPAddress = "",
                PlayerName = model.PlayerName,
                Status = 0,
                ThirdPartyGameCode = CrazyFishingGameCode,
                UpdateTime = DateTimeHelper.GetCurrentTime(),
                SettlementTime = null,
                BetRecordNoPGS = model.BetRecordNoPGS,
                ConvertRate = model.ConvertRate
            };

            BetRecordRepository.Create(insert);
        }



        public void UpdateRawBetRecordAfterSettlement(SettlementBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        {
            rawInfo.RuleFor(x => x.Fishing).Rule_ModelNotNull();

            var settleInfo = rawInfo.Fishing.SettleBetInfo;

            decimal payoutPoint = settleInfo.Stake;
            decimal validPoint = settleInfo.Cost;

            model.PayoutAmount = ConvertGamePointToMoney(model.MerchantID , payoutPoint);
            model.ValidAmount = ConvertGamePointToMoney(model.MerchantID, validPoint);

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

                queryExist.ValidAmount = model.ValidAmount;
                queryExist.WinLossAmount = model.PayoutAmount - queryExist.BetAmount;
                queryExist.SettlementTime = now;
                queryExist.ManualSettleUserID = model.UserID;

                BetRecordRepository.Update(queryExist);
                scope.Complete();
            }

        }

        public ExternalVendorCancelBetResultModel CancelBet(string betRecordNoPGS, string memo)
        {
            DateTime now = DateTimeHelper.GetCurrentTime();
            using (var scope = new TransactionScope())
            {
                var queryExist = BetRecordRepository.GetByBetRecordNoPGS(betRecordNoPGS);

                //檢查狀態合法
                RepositoryHelper.CheckBetRecordStatusIsWhiteList(queryExist, GameVendorBetRecordStatusEnum.Pending);

                queryExist.UpdateTime = now;
                queryExist.DebugInfo += RepositoryHelper.GenerateDebugInfo(queryExist.Status, memo);
                queryExist.Status = GameVendorBetRecordStatusEnum.Cancel.ToInteger();

                BetRecordRepository.Update(queryExist);
                scope.Complete();
            }

            return new ExternalVendorCancelBetResultModel
            {
                GameVendorCustomizedResponse = new FishingSingleWalletPlaceBetResponseModel
                {
                    Code = 0,
                    Msg = "OK",
                    Balance = 0,
                }
            };
        }


        [Description("OK")]
        public void CheckIsValidRequestInfo(ExternalVendorInputValidateRequestParameter model)
        {
            model.RuleFor(z => z.EncryptSign).Rule_NotEmptyText();

            int playerID = CheckSessionTokenAndPlayer(model);

            if (playerID <= 0)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var queryPlayer = PlayerRepository.GetSingleByCache(playerID);

            string secretKey = GetVendorSettingByKey(queryPlayer.MerchantID, VendorConfigSettingEnum.MainSecretKeyOrApiToken)?.SettingValue;

            if (model?.GameVendorRawRequestInfo?.Fishing?.PlaceBetInfo != null)
            {
                var mySign = GenerateSignKey(model?.GameVendorRawRequestInfo?.Fishing?.PlaceBetInfo , secretKey);

                if(mySign.afterMD5 != model.EncryptSign)
                {
                    throw new TerproException(ApiErrorCode.InvalidParameter_SignKeyError);
                }
            }


            if (model?.GameVendorRawRequestInfo?.Fishing?.SettleBetInfo != null)
            {
                var mySign = GenerateSignKey(model?.GameVendorRawRequestInfo?.Fishing?.SettleBetInfo, secretKey);

                if (mySign.afterMD5 != model.EncryptSign)
                {
                    throw new TerproException(ApiErrorCode.InvalidParameter_SignKeyError);
                }
            }
        }


        public void CreateVendorPlayerAccountOrNot(ExternalVendorCreateAccountRequestParameter model)
        {
            //遊戲商會自動建 不用處理
        }


        public string GetSingleBetRecordByApi(ExternalVendorGetSingleBetRecordRequestParameter model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.BetRecordNo).Rule_NotEmptyText();

            //參考文件.pdf  5.1  启动游戏 API
            var req = new RequestParameters
            {
                ApiUrl = "fishdata/userbetflow_pgs",
                Method = HttpMethodType.HttpPost,
            };

            var queryMerchant = MerchantRepository.GetSingleByCache(model.MerchantID);

            string secretKey = GetVendorSettingByKey(model.MerchantID, VendorConfigSettingEnum.MainSecretKeyOrApiToken)?.SettingValue;
            string vendorSourceID = GetVendorSettingByKey(model.MerchantID, VendorConfigSettingEnum.MainAssignVendorID)?.SettingValue;

            //apiToken、api域名 吃資料庫設置
            SetRequestInfoByDBSetting(req, model.MerchantID, "GetBetDetail", model.DebugReferenceUUID);

            long beginTs = new DateTime(2021, 5, 1).CovertToUnixTimestamp();
            long endTs = new DateTime(2021, 5, 12).CovertToUnixTimestamp();

            string beforeMD5 = $"{beginTs}{endTs}{secretKey}";
            var body = new
            {
                source = vendorSourceID,
                start_timestamp = beginTs,
                end_timestamp = endTs,
                sign = beforeMD5.MD5Hash()
            };

            //對應model AGGVendorApiGetDetailResponseModel.cs
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
            string secretKey = GetVendorSettingByKey(model.MerchantID, VendorConfigSettingEnum.MainSecretKeyOrApiToken)?.SettingValue;
            string vendorSourceID = GetVendorSettingByKey(model.MerchantID, VendorConfigSettingEnum.MainAssignVendorID)?.SettingValue;

            string relativeUrl = model.IsTrialMode ? "/game/trial" : "/game/launch";

            Dictionary<string, string> keyValues = new Dictionary<string, string>();

            string playerID = model.PlayerID.ToString();
            long ts = DateTimeHelper.GetCurrentTimeTimestamp();

            keyValues.Add("uid", playerID);
            keyValues.Add("ncode", model.PlayerName.UrlEncode());
            keyValues.Add("token", model.ThirdPartyLaunchSessionToken);
            keyValues.Add("timestamp", ts.ToString());
            keyValues.Add("qwsource", vendorSourceID);

            //我们目前有5种语言&lang=zh-cn/en-us/vi-vn/th-th/in-id
            keyValues.Add("lang", AppHelper.IsChineseLang(model.Language) ? "zh-cn" : "en-us");

            string md5 = $"{playerID}&{ts}&{secretKey}".MD5Hash();
            keyValues.Add("sign", md5);

            //組出靜態網址
            string queryString = string.Join("&", keyValues.Select(z => $"{z.Key}={z.Value}"));
            string finalUrl = $"{launchUrl}?{queryString}";

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

            var apiResult = GetSingleBetRecordByApi(new ExternalVendorGetSingleBetRecordRequestParameter
            {
                BetRecordNo = model.BetRecordNo,
                MerchantID = queryBetRecord.MerchantID,
                GameCode = model.ThirdPartyGameCode
            });

            return new ThirdPartyGameVendorBetRecordDetailUrlViewModel
            {
                InfoType = ThirdPartyVendorInfoTypeEnum.ExternalLink,
                Url = apiResult
            };
        }


        public List<ExternalVendorGameListConvertModel> GetGames(int merchantID)
        {
            //遊戲商只有一款遊戲
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


        private (string beforeMD5 , string afterMD5) GenerateSignKey(object model , string secretKey)
        {
            Dictionary<string, string> keyValues = new Dictionary<string, string>();

            foreach(var item in model.GetType().GetProperties().OrderBy(z=> z.Name))
            {
                string nameCamel = item.Name.ToCamelCase();

                if (nameCamel == "sign")
                    continue;                        

                var value = item.GetValue(model);

                //小數點?
                keyValues.Add(nameCamel, value.ToString());
            }

            string beforeMD5 = string.Join("&" , keyValues.Select(z=> $"{z.Key}={z.Value}"));
            beforeMD5 = beforeMD5 + "&key="+ secretKey;

            return (beforeMD5, beforeMD5.MD5Hash());
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

            if (apiDomain.IsNullOrWhiteSpace())
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted , new ErrorDetailModel
                {
                    ColumnName = "apiDomain",
                });
            }

            req.UserAgent = "target";
            req.BaseUrl = apiDomain;
        }


    }
}
