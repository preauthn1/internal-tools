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
using target.Models.Enums;
using target.Models.ViewModels.GameVendor;

namespace target.API.ExternalVendor.Controllers
{
    public class FunkyBetController : BaseController
    {
        public FunkyIntegrationService FunkyIntegrationService { get; set; }

        public IGameVendorSingleWalletService SingleWalletVendorService { get; set; }


        private static GameVendorTypeEnum _gameVendor = GameVendorTypeEnum.Funky;


        public FunkyBetController()
        {
            FunkyIntegrationService = new FunkyIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_FunkyApi)]
        [Route("Funky/Bet/PlaceBet")]
        [HttpPost]
        public FunkyAPIResponseModel<FunkySingleWalletAPIPlaceBetResultViewModel> PlaceBet([FromBody]FunkySingleWalletAPIPlaceBetRequestModel model)
        {
            // step 1. 先檢查session是否合法
            FunkyIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                SessionOrToken = [REDACTED],
            });

            // step 2. 投注預先扣錢
            var betResult = SingleWalletVendorService.PlaceBet(_gameVendor, new SingleWalletBetDebitAmountRequestParameter
            {
                GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
                {
                    Funky = new FunkyApiRawRequestInfoVendor
                    {
                        PlaceBetInfo = model
                    }
                },
                BetAmount = model.Bet.Stake,
                GameCode = model.Bet.GameCode,
                BetRecordNo = model.Bet.RefNo,
                SessionToken = [REDACTED],
                PlayerIPAddress = model.PlayerIp
            });

            return betResult.GameVendorCustomizedResponse;
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_FunkyApi)]
        [Route("Funky/Bet/SettleBet")]
        [HttpPost]
        public FunkyAPIResponseModel<FunkySingleWalletAPISettleBetResponseModel> SettleBet([FromBody]FunkySingleWalletAPISettleBetRequestModel model)
        {
            // step 1. 先檢查session是否合法
            FunkyIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                IsSkipCheckSessionToken = [REDACTED],
            });

            // step 2. 結算、更新錢包、寫投注紀錄
            var betResult = SingleWalletVendorService.SettleBet(_gameVendor, new SingleWalletBetSettleRequestParameter
            {
                BetRecordNo = model.refNo,
                PayoutAmount = model.betResultReq.winAmount,
                GameCode = model.betResultReq.gameCode,
                PlayerID = model.betResultReq.playerId,
                GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
                {
                    Funky = new FunkyApiRawRequestInfoVendor
                    {
                        SettleBetInfo = model
                    }
                }
            });

            return betResult.GameVendorCustomizedResponse;
        }

       
        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_FunkyApi)]
        [Route("Funky/Bet/CancelBet")]
        [HttpPost]
        public FunkyAPIResponseModel<FunkySingleWalletCancelBetResponseModel> CancelBet([FromBody] FunkySingleWalletCancelBetRequestModel model)
        {
            // step 1. 先檢查session是否合法
            FunkyIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                IsSkipCheckSessionToken = [REDACTED],
            });

            // step 2. 進行取消注單的行為實作
            var betResult = SingleWalletVendorService.CancelBet(_gameVendor, new SingleWalletCancelBetRequestParameter
            {
                PlayerID = model.PlayerId,
                BetRecordNo = model.RefNo,
            });

            return betResult.GameVendorCustomizedResponse;
        }

        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_FunkyApi)]
        [Route("Funky/Bet/CheckBet")]
        [HttpPost]
        public FunkyAPIResponseModel<FunkySingleWalletAPICheckBetResponseModel> CheckBet([FromBody] FunkySingleWalletCheckBetRequestModel model)
        {
            // step 1. 先檢查session是否合法
            FunkyIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                IsSkipCheckSessionToken = [REDACTED],
            });

            // step 2. 檢查注單
            var checkResult = FunkyIntegrationService.CheckBet(model?.Id);

            return checkResult;
        }


    }
}
