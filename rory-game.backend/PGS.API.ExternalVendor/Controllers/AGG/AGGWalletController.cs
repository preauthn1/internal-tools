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

    public class AGGWalletController : BaseGameVendorController
    {
        protected override GameVendorTypeEnum _gameVendor
        {
            get
            {
                return GameVendorTypeEnum.AGG;
            }
        }


        public AGGIntegrationService GameVendorIntegrationService { get; set; }


        public AGGWalletController()
        {
            GameVendorIntegrationService = new AGGIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_AGGApi)]
        [Route("agg/account/balance")]
        [HttpPost]
        public AGGSingleWalletAPIGetBalanceResponseDataModel GetPlayerWallet([FromBody] AGGSingleWalletAPIGetBalanceRequestModel model)
        {         
            //取得錢包餘額
            var walletInfo = SingleWalletVendorService.GetWalletBalance(_gameVendor , new SingleWalletGetWalletBalanceRequestParameter
            {
                PlayerID = model.UserID
            });
            
            return new AGGSingleWalletAPIGetBalanceResponseDataModel
            {
                Balances = new List<AGGSingleWalletEachBalanceItemViewModel>
                {
                    new AGGSingleWalletEachBalanceItemViewModel
                    {
                        Amount = walletInfo.BalanceGamePoint,
                        CurrencyCode = walletInfo.Currency
                    }
                },
                StatusCode = SystemConstant.AGGGameVendorDefaultSuccessResponse
            };
        }

    }
}
