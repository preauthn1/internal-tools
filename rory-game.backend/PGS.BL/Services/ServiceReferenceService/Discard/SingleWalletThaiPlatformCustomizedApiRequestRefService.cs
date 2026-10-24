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
using System.Security.Cryptography;
using System.Reflection.Emit;
using System.Text;

namespace target.BL.Services
{

    public class SingleWalletThaiPlatformCustomizedApiRequestRefService : BaseService //, ISingleWalletPlatformApiRequestRefService
    {
        public MerchantRepository MerchantRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public SingleWalletMerchantSessionTokenRepository SingleWalletMerchantSessionTokenRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public IExternalApiRequestLogRefService ExternalApiRequestLogRefService { get; set; }


        private static string PGS_TOKEN = [REDACTED]";



        /// <summary>
        /// 取得錢包餘額
        /// </summary>
        public SingleWalletModeWalletApiResponseModel GetWalletBalance(SingleWalletModeWalletApiRequestParameter model)
        {
            var querySessionToken = [REDACTED];

            var queryPlayer = PlayerRepository.GetSingleByCache(model.PlayerID);

            //取得API配置
            string requestID = Guid.NewGuid().ToString();

            var para = InstanceApiRequestPara(querySessionToken.MerchantID, "balance", requestID);

            var body = new ThaiPlatformSeamlessGetBalanaceRequestModel
            {
                appid = GetAppID(querySessionToken.MerchantID),
                timestamp = DateTimeHelper.GetCurrentTimeTimestamp(),
                username = queryPlayer.PlayerName,
            };

            body.hash = GenerateMD5Hash(body, GetSecretKey(querySessionToken.MerchantID));

            string formUrlEncode = ConvertToFormUrlEncode(body);
            var apiResult = APIRequestHelper.RequestApiResult<ThaiPlatformSeamlessGetBalanaceResponseModel>(para, formUrlEncode);

            CheckIsErrorResponse(apiResult, requestID, isThrowException: true);

            return new SingleWalletModeWalletApiResponseModel
            {
                Balance = apiResult.ResultObject.Balance,                
            };
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

            var para = InstanceApiRequestPara(querySessionToken.MerchantID, "bet", param.DebugProcessUUID);

            var body = new ThaiPlatformSeamlessPlaceBetRequestModel
            {
                appid = GetAppID(querySessionToken.MerchantID),
                timestamp = DateTimeHelper.GetCurrentTimeTimestamp(),
                amount = param.BetAmount,
                gamecode = param.GameID.ToString(),
                id = param.BetRecordNoPGS,
                username = queryPlayer.PlayerName,
                roundid = ""
            };

            body.hash = GenerateMD5Hash(body, GetSecretKey(querySessionToken.MerchantID));

            string formUrlEncode = ConvertToFormUrlEncode(body);
            var apiResult = APIRequestHelper.RequestApiResult<ThaiPlatformSeamlessPlaceBetResponseModel>(para, formUrlEncode);

            var checkResult = CheckIsErrorResponse(apiResult, param.DebugProcessUUID, isThrowException: true);

            return new SingleWalletModePlaceBetApiResponseModel
            {

            };
        }


        /// <summary>
        /// 結算注單 (一般投注)
        /// </summary>
        public SingleWalletModeSettleApiResponseModel Settlement(SingleWalletModeSettleApiParameter param)
        {
            param.RuleFor(z => z.BetRecordNoPGS).Rule_NotEmptyText();
            param.RuleFor(z => z.GameID).Rule_NotZeroNotMinus();

            var para = InstanceApiRequestPara(param.MerchantID, "settle-bet", param.DebugProcessUUID);

            var queryPlayer = PlayerRepository.GetSingleByCache(param.PlayerID);

            var body = new ThaiPlatformSeamlessSettleRequestModel
            {
                appid = GetAppID(param.MerchantID),
                timestamp = DateTimeHelper.GetCurrentTimeTimestamp(),
                amount = param.PayoutAmount,
                gamecode = param.GameID.ToString(),
                username = queryPlayer.PlayerName,
                id = param.BetRecordNoPGS,
                description = "",
                roundid = "",
            };

            body.hash = GenerateMD5Hash(body, GetSecretKey(param.MerchantID));

            string formUrlEncode = ConvertToFormUrlEncode(body);
            var apiResult = APIRequestHelper.RequestApiResult<ThaiPlatformSeamlessSettleResponseModel>(para, formUrlEncode);

            var checkResult = CheckIsErrorResponse(apiResult, param.DebugProcessUUID, isThrowException: true);

            return new SingleWalletModeSettleApiResponseModel
            {

            };

        }


        /// <summary>
        /// 結算注單 (獎勵/免費遊戲)
        /// </summary>
        public SingleWalletModeSettleApiResponseModel BonusSettlement(SingleWalletModeSettleApiParameter param)
        {
            param.RuleFor(z => z.BetRecordNoPGS).Rule_NotEmptyText();
            param.RuleFor(z => z.GameID).Rule_NotZeroNotMinus();

            var para = InstanceApiRequestPara(param.MerchantID, "bonus-win", param.DebugProcessUUID);

            var queryPlayer = PlayerRepository.GetSingleByCache(param.PlayerID);

            var body = new ThaiPlatformSeamlessBonusRequestModel
            {
                appid = GetAppID(param.MerchantID),
                timestamp = DateTimeHelper.GetCurrentTimeTimestamp(),
                amount = param.PayoutAmount,
                gamecode = param.GameID.ToString(),
                username = queryPlayer.PlayerName,
                id = param.BetRecordNoPGS,
                description = "",
                roundid = "",
                type = ""
            };

            body.hash = GenerateMD5Hash(body, GetSecretKey(param.MerchantID));

            string formUrlEncode = ConvertToFormUrlEncode(body);
            var apiResult = APIRequestHelper.RequestApiResult<ThaiPlatformSeamlessBonusResponseModel>(para, formUrlEncode);

            var checkResult = CheckIsErrorResponse(apiResult, param.DebugProcessUUID, isThrowException: true);

            return new SingleWalletModeSettleApiResponseModel
            {

            };
        }


        /// <summary>
        /// 通知注單取消
        /// </summary>
        public SingleWalletModeCancelApiResponseModel Cancel(SingleWalletModeCancelApiRequestModel param)
        {
            param.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();
            param.RuleFor(z => z.MerchantID).Rule_NotZeroNotMinus();
            param.RuleFor(z => z.GameID).Rule_NotZeroNotMinus();

            var para = InstanceApiRequestPara(param.MerchantID, "cancel-bet", param.DebugProcessUUID);

            var queryPlayer = PlayerRepository.GetSingleByCache(param.PlayerID);

            var body = new ThaiPlatformSeamlessCancelRequestModel
            {
                appid = GetAppID(param.MerchantID),
                timestamp = DateTimeHelper.GetCurrentTimeTimestamp(),
                id = param.TransactionNo,
                betid = param.TransactionNo,
                gamecode = param.GameID.ToString(),
                username = queryPlayer.PlayerName,
                roundid = ""
            };

            body.hash = GenerateMD5Hash(body, GetSecretKey(param.MerchantID));

            string formUrlEncode = ConvertToFormUrlEncode(body);
            var apiResult = APIRequestHelper.RequestApiResult<ThaiPlatformSeamlessCancelResponseModel>(para, formUrlEncode);

            return new SingleWalletModeCancelApiResponseModel
            {

            };
        }



        private RequestParameters InstanceApiRequestPara(int merchantID, string url, string requestUUID)
        {
            //取得API配置
            string domainUrl = GetMerchantApiDomainUrl(merchantID);

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
                    Log_SubCategory = url
                }
            };

            para.AddHeaderData("X-target-IDENTITY", PGS_TOKEN);
            para.AddHeaderData("X-target-REQUEST-ID", requestUUID);

            return para;
        }


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
                throw new TerproException(ApiErrorCode.ThirdPartyError_SettingDataNotCompleted, new ErrorDetailModel("Integration_SingleWalletApiDomain"));

            return querySetting.SettingValue;
        }


        private string GetAppID(int merchantID)
        {
            return "XXXXXXXX";

            //var querySetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.Integration_SingleWalletApiDomain);

            //if (querySetting == null || querySetting.SettingValue.IsNullOrWhiteSpace())
            //    throw new TerproException(ApiErrorCode.ThirdPartyError_SettingDataNotCompleted, new ErrorDetailModel("Integration_SingleWalletApiDomain"));

            //return querySetting.SettingValue;
        }

        private string GetSecretKey(int merchantID)
        {
            return "XXXXXXXX";

            //var querySetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.Integration_SingleWalletApiDomain);

            //if (querySetting == null || querySetting.SettingValue.IsNullOrWhiteSpace())
            //    throw new TerproException(ApiErrorCode.ThirdPartyError_SettingDataNotCompleted, new ErrorDetailModel("Integration_SingleWalletApiDomain"));

            //return querySetting.SettingValue;
        }




        private (bool isException, ApiErrorCode errorCode) CheckIsErrorResponse<T>(ApiRequestResultModel<T> response, string requestID, bool isThrowException = true) where T : BaseThaiPlatformSeamlessApiResponseModel
        {
            ApiErrorCode errorCode = ApiErrorCode.NotImplement;

            if (response.Exception == null)
            {
                return (false, ApiErrorCode.NotImplement);
            }

            if (response.IsTimeoutException)
            {
                errorCode = ApiErrorCode.ThirdPartyApiFail_Timeout;
            }
            else if (response.Exception != null)
            {
                errorCode = ApiErrorCode.ThirdPartyApiFail_ResponseError;
            }
            else if (response.ResultObject == null)
            {
                errorCode = ApiErrorCode.ThirdPartyApiFail_ParsingFail;
            }

            if (response.ResultObject != null &&
               response.ResultObject.Status != 0)
            {                
                errorCode = ApiErrorCode.ThirdPartyApiFail_ResponseError;
            }

            if (isThrowException)
            {
                throw new TerproException(errorCode, new ErrorDetailModel
                {
                    LogExtraInfo = requestID + "_" + response.Exception.Message + " . " + response.ResponseText
                });
            }

            return (true, errorCode);
        }


        public static string GenerateMD5Hash(object post, string secretKey)
        {
            Dictionary<string, object> dic = new Dictionary<string, object>();

            foreach (var item in post.GetType().GetProperties().OrderBy(z => z.Name))
            {
                string columnName = item.Name.ToLower();
                var value = item.GetValue(post);

                if (columnName == "hash")
                    continue;

                dic.Add(columnName, value.ToString());
            }

            string beforeMD5 = string.Join("&", dic.Select(z => $"{z.Key}={z.Value}")) + secretKey;
            string finalResult = GetMD5Signature(beforeMD5);
            return finalResult;
        }


        private string ConvertToFormUrlEncode(object model)
        {
            string form = string.Empty;

            var keyvalues = model.GetType().GetProperties()
                .Select(p => $"{p.Name} = {p.GetValue(model)}")
                .ToArray();

            return string.Join('&', keyvalues);
        }


        public static string GetMD5Signature(string rawData)
        {
            var md5Content = string.Empty;
            using (MD5 md5Hash = MD5.Create())
            {
                md5Content = GetMd5Hash(md5Hash, rawData);
            }
            return md5Content;
        }


        private static string GetMd5Hash(MD5 md5Hash, string rawData)
        {
            // Convert the input string to a byte array and compute the hash.
            byte[] data = md5Hash.ComputeHash(Encoding.UTF8.GetBytes(rawData));
            // Create a new Stringbuilder to collect the bytes
            // and create a string.
            var sBuilder = new StringBuilder();
            // Loop through each byte of the hashed data
            // and format each one as a hexadecimal string.
            for (int i = 0; i < data.Length; i++)
            {
                sBuilder.Append(data[i].ToString("x2"));
            }
            // Return the hexadecimal string.
            return sBuilder.ToString();
        }


    }
}
