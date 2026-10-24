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

namespace target.API.ExternalVendor.Controllers
{

    public class AGGTokenAuthController : BaseGameVendorController
    {
        protected override GameVendorTypeEnum _gameVendor
        {
            get
            {
                return GameVendorTypeEnum.AGG;
            }
        }

        public AGGIntegrationService GameVendorIntegrationService { get; set; }



        public AGGTokenAuthController()
        {
            GameVendorIntegrationService = new AGGIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_AGGApi)]
        [Route("agg/user/authenticate")]
        [HttpPost]
        public AGGSingleWalletAPICheckTokenResponseDataModel CheckTokenIsValid([FromBody] AGGSingleWalletCheckTokenRequestModel model)
        {
            // step 1. 先檢查session是否合法
            GameVendorIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                SessionOrToken = [REDACTED],
            });

            // step 2. 取得錢包餘額
            var walletInfo = SingleWalletVendorService.GetWalletBalance(_gameVendor , new SingleWalletGetWalletBalanceRequestParameter
            {
                SessionToken = [REDACTED],
            });

            // Case 3. 取消注單
            return new AGGSingleWalletAPICheckTokenResponseDataModel
            {
                StatusCode = SystemConstant.AGGGameVendorDefaultSuccessResponse,
                UserID = walletInfo.PlayerID,
                Balances = new List<AGGSingleWalletEachBalanceItemViewModel>
                {
                    new AGGSingleWalletEachBalanceItemViewModel
                    {
                        Amount = walletInfo.BalanceGamePoint,
                        CurrencyCode = walletInfo.Currency
                    }
                }
            };
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_AGGApi)]
        [Route("agg/account/token")]
        [HttpPost]
        public AGGSingleWalletGenerateGameTokenResponseViewModel GenerateGameToken([FromBody] AGGSingleWalletGenerateGameTokenRequestModel model)
        {
            // step 1. 先檢查session是否合法
            GameVendorIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                SessionOrToken = [REDACTED],
            });

            // step 2. 取得gameID
            var queryGameID = GameVendorIntegrationService.GetGameIdByGameCode(model.GameId);

            // step 3. 取得錢包餘額
            var walletInfo = SingleWalletVendorService.GetWalletBalance(_gameVendor , new SingleWalletGetWalletBalanceRequestParameter
            {
                SessionToken = [REDACTED],
            });

            //有點鳥的作法 因為不想再開新table，就直接用 gameID + 我方的token 當作對方的game_token
            string gameToken = $"{queryGameID}__{model.UserToken}";

            return new AGGSingleWalletGenerateGameTokenResponseViewModel
            {
                AccountToken = [REDACTED],
                StatusCode = SystemConstant.AGGGameVendorDefaultSuccessResponse,
                Balances = new List<AGGSingleWalletEachBalanceItemViewModel>
                {
                    new AGGSingleWalletEachBalanceItemViewModel
                    {
                        Amount = walletInfo.BalanceGamePoint,
                        CurrencyCode = walletInfo.Currency
                    }
                }
            };
        }


    }
}
