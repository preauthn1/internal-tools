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
using target.Utility.Utility;
using target.Utility;
using game.Entity.Game;
using System.Threading.Tasks;
using System.Threading;
using target.Utility.Utility.ThirdParty;

namespace target.BL.Services
{

    public class SingleWalletPlatformApiRequestRefService : BaseService, ISingleWalletPlatformApiRequestRefService
    {
        public MerchantRepository MerchantRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public SingleWalletMerchantSessionTokenRepository SingleWalletMerchantSessionTokenRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public IExternalApiRequestLogRefService ExternalApiRequestLogRefService { get; set; }


        /// <summary>
        /// 取得錢包餘額
        /// </summary>
        public SingleWalletModeWalletApiResponseModel GetWalletBalance(SingleWalletModeWalletApiRequestParameter model)
        {
            var querySessionToken = [REDACTED];         

            //取得API配置
            string requestID = Guid.NewGuid().ToString();
            var para = InstanceApiRequestPara(querySessionToken.MerchantID, "target/players/wallet", requestID);

            var body = new SingleWalletModeWalletApiRequestModel
            {
                SessionToken = [REDACTED],
            };

            var apiResult = APIRequestHelper.RequestApiResult<BaseSingleWalletModeResponseModel<SingleWalletModeWalletApiResponseModel>>(para, body);

            CheckIsErrorResponse(apiResult, requestID, isThrowException: true);

            return apiResult.ResultObject.Data;
        }



        /// <summary>
        /// 投注
        /// </summary>
        public SingleWalletModePlaceBetApiResponseModel PlaceBet(SingleWalletModePlaceBetParameter param)
        {
            param.RuleFor(z => z.BetRecordNoPGS).Rule_NotEmptyText();
            param.RuleFor(z => z.SessionToken).Rule_NotEmptyText();
            param.RuleFor(z => z.GameID).Rule_NotZeroNotMinus();

            var querySessionToken = [REDACTED];

            var queryPlayer = PlayerRepository.GetSingleByCache(querySessionToken.PlayerID);

            var para = InstanceApiRequestPara(querySessionToken.MerchantID, "target/bet/placeBet", param.DebugProcessUUID , debugInfo1: param.BetRecordNoPGS , debugInfo2: param.PlayerID.ToString());

            var body = new SingleWalletModePlaceBetApiRequestModel
            {
                PlayerUniqueID = queryPlayer.PlayerName,
                SessionToken = [REDACTED],
                BetAmount = param.BetAmount,
                GameID = param.GameID,
                TransactionNo = param.BetRecordNoPGS
            };

            var apiResult = RequestApiResult<BaseSingleWalletModeResponseModel<SingleWalletModePlaceBetApiResponseModel>>(para, body);

            var checkResult = CheckIsErrorResponse(apiResult, param.DebugProcessUUID, isThrowException: true);

            // Case 1. 正常拿到回應
            if (!checkResult.isException && apiResult.ResultObject.Data?.TransactionNo == param.BetRecordNoPGS)
                return apiResult.ResultObject.Data;

            throw new TerproException(checkResult.errorCode);
        }


        /// <summary>
        /// 結算注單 (一般投注)
        /// </summary>
        public SingleWalletModeSettleApiResponseModel Settlement(SingleWalletModeSettleApiParameter param)
        {
            param.RuleFor(z => z.BetRecordNoPGS).Rule_NotEmptyText();
            param.RuleFor(z => z.GameID).Rule_NotZeroNotMinus();

            var para = InstanceApiRequestPara(param.MerchantID, "target/bet/settle", param.DebugProcessUUID, debugInfo1: param.BetRecordNoPGS, debugInfo2: param.PlayerID.ToString());

            var body = new SingleWalletModeSettleApiRequestModel
            {                
                TransactionNo = param.BetRecordNoPGS,
                BetAmount = param.BetAmount,
                GameID = param.GameID,
                TaxAmount = param.TaxAmount,
                PayoutAmount = param.PayoutAmount,
                ValidAmount = param.ValidAmount,
                SettlementTime = param.SettlementTime.HasValue ? param.SettlementTime.Value.ToISO8601() : null
            };

            var apiResult = RequestApiResult<BaseSingleWalletModeResponseModel<SingleWalletModeSettleApiResponseModel>>(para, body);

            var checkResult = CheckIsErrorResponse(apiResult, param.DebugProcessUUID, isThrowException: true);
            
            if (apiResult.ResultObject.Data?.TransactionNo != param.BetRecordNoPGS)
            {
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ParsingFail);
            }

            return apiResult.ResultObject.Data;
        }


        /// <summary>
        /// 結算注單 (獎勵/免費遊戲)
        /// </summary>
        public SingleWalletModeSettleApiResponseModel BonusSettlement(SingleWalletModeSettleApiParameter param)
        {
            param.RuleFor(z => z.BetRecordNoPGS).Rule_NotEmptyText();
            param.RuleFor(z => z.GameID).Rule_NotZeroNotMinus();

            var para = InstanceApiRequestPara(param.MerchantID, "target/bonus/settle", param.DebugProcessUUID, debugInfo1: param.BetRecordNoPGS, debugInfo2: param.PlayerID.ToString());

            var queryPlayer = PlayerRepository.GetSingleByCache(param.PlayerID);

            var body = new SingleWalletModeSettleApiRequestModel
            {
                PlayerUniqueID = queryPlayer.PlayerName,
                TransactionNo = param.BetRecordNoPGS,
                BetAmount = param.BetAmount,
                GameID = param.GameID,
                TaxAmount = param.TaxAmount,
                PayoutAmount = param.PayoutAmount,
                ValidAmount = param.ValidAmount,
                SettlementTime = param.SettlementTime.HasValue ? param.SettlementTime.Value.ToISO8601() : null
            };

            var apiResult = RequestApiResult<BaseSingleWalletModeResponseModel<SingleWalletModeSettleApiResponseModel>>(para, body);

            var checkResult = CheckIsErrorResponse(apiResult, param.DebugProcessUUID, isThrowException: true);

            if (apiResult.ResultObject.Data?.TransactionNo != param.BetRecordNoPGS)
            {
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ParsingFail);
            }

            return apiResult.ResultObject.Data;
        }


        /// <summary>
        /// 通知注單取消
        /// </summary>
        public SingleWalletModeCancelApiResponseModel Cancel(SingleWalletModeCancelApiRequestModel param)
        {
            param.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();
            param.RuleFor(z => z.MerchantID).Rule_NotZeroNotMinus();
            param.RuleFor(z => z.GameID).Rule_NotZeroNotMinus();
         
            var para = InstanceApiRequestPara(param.MerchantID, "target/bet/cancel", param.DebugProcessUUID, debugInfo1: param.TransactionNo, debugInfo2: param.PlayerID.ToString());

            var body = new SingleWalletModeCancelApiRequestModel
            {
                TransactionNo = param.TransactionNo,
                GameID = param.GameID,
            };

            var apiResult = RequestApiResult<BaseSingleWalletModeResponseModel<SingleWalletModeCancelApiResponseModel>>(para, body);

            List<int> ignoreErrorCode = new List<int>
            {
                ApiErrorCode.NotFoundBetRecord.ToInteger()
            };

            var checkResult = CheckIsErrorResponse(apiResult, param.DebugProcessUUID, isThrowException: true , ignoreErrorCode);

            //取消注單 若對方回應not found 也視為合法
            if(ignoreErrorCode.Contains(checkResult.errorCode.ToInteger()))
            {
                return new SingleWalletModeCancelApiResponseModel();
            }

            if (apiResult.ResultObject.Data?.TransactionNo != param.TransactionNo)
            {
                throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ParsingFail);
            }

            return apiResult.ResultObject.Data;
        }

        /// <summary>
        /// 因為要寫MSTest 所以特別包成一層 virtual function
        /// </summary>
        protected virtual ApiRequestResultModel<T> RequestApiResult<T>(RequestParameters para, object body = null)
        {
            var apiResult = APIRequestHelper.RequestApiResult<T>(para, body);
            return apiResult;
        }


        private RequestParameters InstanceApiRequestPara(int merchantID, string url, string requestUUID , string debugInfo1 = "" , string debugInfo2 = "")
        {
            //取得API配置
            string domainUrl = GetMerchantApiDomainUrl(merchantID);
            string PgsToken = [REDACTED];

            var para = new RequestParameters
            {
                BaseUrl = domainUrl,
                ApiUrl = url,
                Method = HttpMethodType.HttpPost,
                LogAction = WriteRequestLog,
                TimeoutMillisecond = 15000, //等太久反而會拖垮我方自己的效能?
                AdvancedSetting = new AdvancedSetting
                {
                    Log_RequestUUID = requestUUID,
                    Log_ApiCategory = merchantID.ToString(),
                    Log_SubCategory = url,
                    Log_DebugInfo1 = debugInfo1,
                    Log_DebugInfo2 = debugInfo2,
                    IsCamelCase = true
                }
            };

            para.AddHeaderData("X-target-IDENTITY", PgsToken);
            para.AddHeaderData("X-target-REQUEST-ID", requestUUID);

            return para;
        }

        /// <summary>
        /// 因為要寫MSTest 所以特別包成一層function
        /// </summary>
        public SingleWalletMerchantSessionTokenEntity CheckIsValidSessionToken(string sessionTokenPGS)
        {
            //TODO 也許可以改用cache ?
            //TODO 也許可以改用cache ?
            //TODO 也許可以改用cache ?
            //TODO 也許可以改用cache ?
            //TODO 也許可以改用cache ?
            //TODO 也許可以改用cache ?
            //TODO 也許可以改用cache ?
            var query = SingleWalletMerchantSessionTokenRepository.GetBySessionTokenPGS(sessionTokenPGS);

            if (query == null)
                throw new TerproException(ApiErrorCode.NotFoundSessionToken);

            return query;
        }

        protected void WriteRequestLog(ApiRequestLogTemporaryInfoModel model)
        {
            ExternalApiRequestLogRefService.CreateSingleWalletApiLog(model);
        }


        private string GetMerchantApiDomainUrl(int merchantID)
        {
            var querySetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.Integration_SingleWalletApiDomain);

            if (querySetting == null || querySetting.SettingValue.IsNullOrWhiteSpace())
                throw new TerproException(ApiErrorCode.ThirdPartyError_SettingDataNotCompleted , new ErrorDetailModel("Integration_SingleWalletApiDomain"));
            
            return querySetting.SettingValue;
        }


        private string GetMerchantSettingPGSValidToken(int merchantID)
        {
            var querySetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.Integration_SingleWalletApiPgsValidToken);

            if (querySetting == null || querySetting.SettingValue.IsNullOrWhiteSpace())
                throw new TerproException(ApiErrorCode.ThirdPartyError_SettingDataNotCompleted, new ErrorDetailModel("Integration_SingleWalletApiPgsValidToken"));

            return querySetting.SettingValue;
        }


        private (bool isException, ApiErrorCode errorCode) CheckIsErrorResponse<T>(
            ApiRequestResultModel<BaseSingleWalletModeResponseModel<T>> response, 
            string requestID, 
            bool isThrowException = true,
            List<int> ignoreErrorCodes = null)
        {
            ApiErrorCode errorCode = ApiErrorCode.Success;

           
            if (response.IsTimeoutException)
            {
                errorCode = ApiErrorCode.ThirdPartyApiFail_Timeout;
            }
            else if (response.Exception != null)
            {
                //400 statusCode 也是走這邊
                errorCode = ApiErrorCode.ThirdPartyApiFail_ResponseError;

                var tryParse = response.ResponseText.TryJsonConvertToModel<BaseSingleWalletModeResponseModel<string>>();

                //如果能解析的話 
                if(tryParse != null && tryParse.Code > 400000 && ExtensionOfEnum.IsDefinedValue((ApiErrorCode)tryParse.Code))
                {
                    string msg = MultiLanguageTextHelper.GetEnumsCultureText(ApiErrorCode.ThirdPartyApiFail_ResponseError);
                    errorCode = (ApiErrorCode)tryParse.Code;
                }

            }
            else if (response.ResultObject == null)
            {
                errorCode = ApiErrorCode.ThirdPartyApiFail_ParsingFail;
            }

            if (response.ResultObject != null &&
               response.ResultObject.Code != 0)
            {
                var convertCode = (ApiErrorCode)response.ResultObject.Code;
                if (ExtensionOfEnum.IsDefinedValue(convertCode))
                {
                    //如果有定義特定的errorCode不拋錯誤
                    if(ignoreErrorCodes.IsNotEmptyNull() && ignoreErrorCodes.Contains(response.ResultObject.Code))
                    {
                        return (false, convertCode);
                    }
                    else
                    {
                        errorCode = convertCode;
                    }
                }
                else
                {
                    errorCode = ApiErrorCode.ThirdPartyApiFail_ResponseError;
                }
            }

            if (response.Exception == null && errorCode == ApiErrorCode.Success)
            {
                return (false, ApiErrorCode.Success);
            }

            if (isThrowException)
            {
                throw new TerproException(errorCode, new ErrorDetailModel
                {
                    LogExtraInfo = requestID + "_" + response.Exception?.Message + " . " + response.ResponseText
                });
            }
            
            return (true, errorCode);
        }



    }
}
