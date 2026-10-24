using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using example.invalid;
using System.Linq;
using System.Reflection;
using System.Runtime.Versioning;
using System.Threading;
using System.Threading.Tasks;
using game.Entity.Game;
using Lck.Utility.Extensions.Helper;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Utility.Helper.Dependency;
using target.Utility.IOC;

namespace target.API.Public.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    public class FakeResponseController : BaseController
    {

        /// <summary>
        /// 假資料API，假裝是 廠商的session確認API
        /// </summary>
        [HttpPost]
        [Route("target/validate-session")]
        [Consumes("application/x-www-form-urlencoded")]
        public MerchantIntergrationApiResultModel ValidateUserSession([FromForm] MerchantIntergrationApiRequestModel model)
        {
            return new MerchantIntergrationApiResultModel
            {
                IsValidOK = true,
                PlayerName = model.PlayerName
            };
        }


        /// <summary>
        /// 測試用API
        /// </summary>
        [Route("fake/player/game-rooms")]
        [HttpGet]
        public object TestSignKey([FromQuery] GameServerRequestModel model)
        {
            string token = [REDACTED];

            var text = SignKeyRuleHelper.GenerateF2E_HttpGet_SignKeyBeforeMD5(HttpContext.Request);
            return new
            {
                text = text,
                textMD5 = text.MD5Hash()
            };
        }


        /// <summary>
        /// 模擬單錢包的假API
        /// </summary>
        [Route("target/players/wallet")]
        [HttpPost]
        public SingleWalletModeWalletApiResponseModel PGSFakeSingle_GetWallet([FromBody] SingleWalletModeWalletApiRequestModel model)
        {
            //int rndValue = RandomHelper.GetRandIntegerNumbers(1, 99);

            //if (50 > rndValue)
            //{
            //    throw new TerproException(ApiErrorCode.ThirdPartyApiFail_IsNotValid);
            //}

            return new SingleWalletModeWalletApiResponseModel
            {
                Balance = RandomHelper.GetRandIntegerNumbers(1001, 1999),
            };
        }


        [Route("target/bet/placeBet")]
        [HttpPost]
        public SingleWalletModePlaceBetApiResponseModel PGSFakeSingle_PlaceBet([FromBody] SingleWalletModePlaceBetApiRequestModel model)
        {
            RandomRatioThrowException(model.TransactionNo , SystemConfigSettingEnum.FakeSingleWalletApi_ErrorRate_PlaceBet);

            decimal afterBalance = RandomHelper.GetRandIntegerNumbers(1001, 1999);

            //讓捕魚遊戲placeBet後 一定是餘額0元
            if (model.GameID == 310001)
                afterBalance = 0;

            return new SingleWalletModePlaceBetApiResponseModel
            {
                Balance = afterBalance,
                TransactionNo = model.TransactionNo
            };
        }



        [Route("target/bet/settle")]
        [HttpPost]
        public SingleWalletModeSettleApiResponseModel PGSFakeSingle_Settle([FromBody] SingleWalletModeSettleApiRequestModel model)
        {
            RandomRatioThrowException(model.TransactionNo , SystemConfigSettingEnum.FakeSingleWalletApi_ErrorRate_Settlement);
            return new SingleWalletModeSettleApiResponseModel
            {
                Balance = RandomHelper.GetRandIntegerNumbers(1001, 1999),
                TransactionNo = model.TransactionNo,
            };
        }



        [Route("target/bonus/settle")]
        [HttpPost]
        public SingleWalletModeSettleApiResponseModel PGSFakeSingle_BonusSettle([FromBody] SingleWalletModeSettleApiRequestModel model)
        {
            RandomRatioThrowException(model.TransactionNo , SystemConfigSettingEnum.FakeSingleWalletApi_ErrorRate_Settlement);
            return new SingleWalletModeSettleApiResponseModel
            {
                Balance = RandomHelper.GetRandIntegerNumbers(1001, 1999),
                TransactionNo = model.TransactionNo,
            };
        }



        [Route("target/bet/cancel")]
        [HttpPost]
        public SingleWalletModeCancelApiResponseModel PGSFakeSingle_Cancel([FromBody] SingleWalletModeCancelApiRequestModel model)
        {
            return new SingleWalletModeCancelApiResponseModel
            {
                TransactionNo = model.TransactionNo
            };
        }


        /// <summary>
        /// 隨機機率 產生exception 提供模擬
        /// </summary>
        private void RandomRatioThrowException(string transactionNo , SystemConfigSettingEnum type)
        {
            string cacheKeyExceptionTranNo = "Keep_Transaction_No_" + transactionNo;

            var keepTranNo = CacheHelper.GetCacheValue<string>(cacheKeyExceptionTranNo);

            //如果這張單模擬過錯誤 第二次也會讓他錯誤
            if (!string.IsNullOrWhiteSpace(keepTranNo))
            {
                throw new TerproException(ApiErrorCode.InvalidStatus);
            }

            var service = IOCHelper.Resolve<ISystemConfigSettingService>();

            var querySetting = service.GetSystemConfigSetting(type);
            decimal requiredValue = querySetting == null ? 0 : querySetting.SettingValueOfDecimal;

            int rndValue = RandomHelper.GetRandIntegerNumbers(1, 99);


            if(requiredValue >= 10000)
            {
                int delaySeconds = (int)requiredValue - 10000;
                delaySeconds = delaySeconds >= 60 ? 60 : delaySeconds;
                Thread.Sleep(1000 * delaySeconds);
                return;
            }

            if (requiredValue > rndValue)
            {
                //記住這張單號 讓retry下一次 也發生錯誤
                CacheHelper.SetCacheValue(cacheKeyExceptionTranNo, transactionNo, new TimeSpan(0, 10, 0));
                throw new TerproException(ApiErrorCode.InvalidStatus);
            }
        }


    }
}
