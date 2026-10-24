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
using target.Models;
using target.Utility.IOC;

namespace target.API.ExternalVendor.Controllers
{

    public class CrazyFishingTokenController : BaseGameVendorController
    {

        protected override GameVendorTypeEnum _gameVendor
        {
            get
            {
                return GameVendorTypeEnum.CrazyFishing;
            }
        }


        public FishingIntegrationService GameVendorIntegrationService { get; set; }


        public CrazyFishingTokenController()
        {
            GameVendorIntegrationService = new FishingIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_FishApi)]
        [Route("crazyFishing/account_auth")]
        [HttpPost]
        public FishingSingleWalletCheckAccountTokenResponseModel CheckTokenIsValid([FromBody] FishingSingleWalletCheckAccountTokenRequestModel model)
        {
            // step 1. 先檢查session是否合法
            GameVendorIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                EncryptSign = model.Sign,
                SessionOrToken = [REDACTED],
                GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
                {
                    Fishing = new FishingApiRawRequestInfoVendor
                    {
                        TokenAuth = model
                    }
                }
            });

            // step 2. 取得錢包餘額
            var walletInfo = SingleWalletVendorService.GetWalletBalance(_gameVendor , new SingleWalletGetWalletBalanceRequestParameter
            {
                SessionToken = [REDACTED],
            });
            
            return new FishingSingleWalletCheckAccountTokenResponseModel
            {
                Code = 0,
                Msg = "OK",
                Nick = walletInfo.PlayerName,
                Account = walletInfo.PlayerID.ToString(),
                Icon = "https://example.invalid/public/f51954c7-df31-4bab-a367-e60b33fcb35d.png",
                Wealth = walletInfo.BalanceGamePoint,
            };
        }


    }
}
