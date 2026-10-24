using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using Microsoft.AspNetCore.Cors;
using target.Web.Modules.Filter;
using target.BL.Services;
using target.Utility.Extensions;
using target.Models.Enums;
using target.Models.ViewModels.GameVendor;
using System.Linq;
using target.Models.Models;
using target.Utility.Helper;
using target.Models;

namespace target.API.ExternalVendor.Controllers
{
    public class CrazyFishingBetController : BaseGameVendorController
    {

        protected override GameVendorTypeEnum _gameVendor
        {
            get
            {
                return GameVendorTypeEnum.CrazyFishing;
            }
        }


        public FishingIntegrationService GameVendorIntegrationService { get; set; }


        public CrazyFishingBetController()
        {
            GameVendorIntegrationService = new FishingIntegrationService();
        }



        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_FishApi, IsHeaderLogEnable = false)]
        [Route("crazyFishing/place_bet")]
        [HttpPost]
        public FishingSingleWalletPlaceBetResponseModel PlaceBet([FromBody] FishingSingleWalletPlaceBetRequestModel model)
        {
            var rawRequestInfo = InstanceRawRequestModel(model, null);

            // step 1. 先檢查session是否合法
            GameVendorIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                EncryptSign = model.Sign,
                SessionOrToken = [REDACTED],
                GameVendorRawRequestInfo = rawRequestInfo
            });

            // step 2. 投注預先扣錢
            var betResult = SingleWalletVendorService.PlaceBet(_gameVendor, new SingleWalletBetDebitAmountRequestParameter
            {
                GameVendorRawRequestInfo = rawRequestInfo,
                BetAmount = model.Stake,
                GameCode = SystemConstant.CrazyFishingDefaultGameCode,
                BetRecordNo = model.OrderId,
                SessionToken = [REDACTED],
            });

            Develop_Simulate_DelayTimeout(SystemConfigSettingEnum.GameVendor_PlaceBet_SuccessUpdateBut_DelayTimeout);
            return betResult.GameVendorCustomizedResponse;
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_FishApi, IsHeaderLogEnable = false)]
        [Route("crazyFishing/settle_bet")]
        [HttpPost]
        public FishingSingleWalletSettleBetResponseModel Settlement([FromBody] FishingSingleWalletSettleBetRequestModel model)
        {
            var rawRequestInfo = InstanceRawRequestModel(null, model);

            // step 1. 先檢查session是否合法
            GameVendorIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                IsSkipCheckSessionToken = [REDACTED],
                EncryptSign = model.Sign,
                PlayerID = model.Account,
                GameVendorRawRequestInfo = rawRequestInfo
            });

            // step 2. 結算、更新錢包、寫投注紀錄
            var betResult = SingleWalletVendorService.SettleBet(_gameVendor, new SingleWalletBetSettleRequestParameter
            {
                BetRecordNo = model.OrderId,
                PayoutAmount = model.Stake,
                GameCode = SystemConstant.CrazyFishingDefaultGameCode,
                PlayerID = model.Account,
                GameVendorRawRequestInfo = rawRequestInfo
            });


            Develop_Simulate_DelayTimeout(SystemConfigSettingEnum.GameVendor_Settle_SuccessUpdateBut_DelayTimeout);
            return betResult.GameVendorCustomizedResponse;
        }



        private EachGameVendorRawRequestInfoModel InstanceRawRequestModel(
            FishingSingleWalletPlaceBetRequestModel placeBetParam,
            FishingSingleWalletSettleBetRequestModel settleParam)
        {
            return new EachGameVendorRawRequestInfoModel
            {
                Fishing = new FishingApiRawRequestInfoVendor
                {
                    PlaceBetInfo = placeBetParam,
                    SettleBetInfo = settleParam
                }
            };
        }




    }
}
