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
using target.Models.Models;

namespace target.API.ExternalVendor.Controllers
{

    public class AGGLiveTokenAuthController : BaseGameVendorController
    {
        protected override GameVendorTypeEnum _gameVendor
        {
            get
            {
                return GameVendorTypeEnum.AGGLive;
            }
        }

        public AGGLiveIntegrationService GameVendorIntegrationService { get; set; }


        public AGGLiveTokenAuthController()
        {
            GameVendorIntegrationService = new AGGLiveIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_AGGApi)]
        [Route("aggLive/api/auth")]
        [HttpPost]
        public AGGLiveSingleWalletAPICheckTokenResponseDataModel CheckTokenAuthAndBalance([FromBody] AGGLiveSingleWalletCheckTokenRequestModel model)
        {
            string token = [REDACTED];

            //step 1.先檢查session是否合法
            GameVendorIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                SessionOrToken = [REDACTED],
            });

            // step 2. 取得錢包餘額
            var walletInfo = SingleWalletVendorService.GetWalletBalance(_gameVendor, new SingleWalletGetWalletBalanceRequestParameter
            {
                SessionToken = [REDACTED],
            });

            return new AGGLiveSingleWalletAPICheckTokenResponseDataModel
            {
               Playerdetailresponse = new AGGLiveStatusOnlyResponseModel
               {
                   status = new AGGLiveStatusResponseModel
                   {
                       success = true,
                       autherror = false,
                       message = ""
                   },                  
               },
               accountid = walletInfo.PlayerID.ToString(),
               accountname = walletInfo.PlayerName,
               balance = walletInfo.BalanceGamePoint,
               currencycode = walletInfo.Currency
            };
        }

    }
}
