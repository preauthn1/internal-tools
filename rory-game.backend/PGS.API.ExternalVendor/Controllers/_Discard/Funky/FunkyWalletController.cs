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

    public class FunkyWalletController : BaseController
    {

        public FunkyIntegrationService FunkyIntegrationService { get; set; }

        public IGameVendorSingleWalletService SingleWalletVendorService { get; set; }


        public FunkyWalletController()
        {
            FunkyIntegrationService = new FunkyIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_FunkyApi)]
        [Route("Funky/User/GetBalance")]
        [HttpPost]
        public FunkyAPIResponseModel<FunkySingleWalletAPIGetBalanceResponseDataModel> GetPlayerWallet([FromBody]FunkySingleWalletAPIGetBalanceRequestModel model)
        {
            // step 1. 先檢查session是否合法
            FunkyIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                SessionOrToken = [REDACTED],
                PlayerID = model?.PlayerId ?? 0
            });

            // step 2. 取得錢包餘額
            var walletInfo = SingleWalletVendorService.GetWalletBalance(new SingleWalletGetWalletBalanceRequestParameter
            {
                SessionToken = [REDACTED],
                PlayerID = model.PlayerId
            });
            
            return new FunkyAPIResponseModel<FunkySingleWalletAPIGetBalanceResponseDataModel>
            {
                ErrorCode = 0,
                Data = new FunkySingleWalletAPIGetBalanceResponseDataModel
                {
                    Balance = walletInfo.Amount,
                    Currency = walletInfo.Currency
                }
            };
        }

    }
}
