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

namespace target.API.ExternalVendor.Controllers
{
    public class AGGLiveBetController : BaseGameVendorController
    {
        protected override GameVendorTypeEnum _gameVendor
        {
            get
            {
                return GameVendorTypeEnum.AGGLive;
            }
        }

        public AGGLiveIntegrationService GameVendorIntegrationService { get; set; }


        public AGGLiveBetController()
        {
            GameVendorIntegrationService = new AGGLiveIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_AGGApi , IsHeaderLogEnable = false)]
        [Route("aggLive/wager/transaction")]
        [HttpPost]
        public AGGLiveSingleWalletAPITransactionResponseViewModel PlaceBet([FromBody] AGGLiveSingleWalletAPITransactionRequestModel model)
        {
            model.RuleFor(z => z.FundTransferRequest.token).Rule_NotEmptyText();

            // Case 1. 投注流程
            if (IsPlaceBetProcessRequest(model))
            {
                var result = PlaceBetProcess(model);
                return result;
            }

            // Case 2. 如果是結算、執行結算加錢的邏輯
            if (IsSettleProcessRequest(model))
            {
                var settleResult = SettlementProcess(model);
                return settleResult;
            }

            // Case 3. 取消注單
            if (IsCancelProcessRequest(model))
            {
                var result = CancelProcess(model);
                return result;
            }

            throw new TerproException(ApiErrorCode.NotImplement);
        }


        private AGGLiveSingleWalletAPITransactionResponseViewModel PlaceBetProcess(AGGLiveSingleWalletAPITransactionRequestModel model)
        {
            DateTime time0 = DateTimeHelper.GetCurrentTime();

            string sessionToken = [REDACTED];

            var placeBetInfo = model.FundTransferRequest.Funds.FundInfo.FirstOrDefault();

            decimal betAmount = placeBetInfo.amount.ToAbs();

            // Step 1. 先檢查session是否合法
            GameVendorIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                SessionOrToken = [REDACTED],
            });

            // Step 2. 如果是投注、投注預先扣錢
            var betResult = SingleWalletVendorService.PlaceBet(_gameVendor, new SingleWalletBetDebitAmountRequestParameter
            {
                SessionToken = [REDACTED],
                BetRecordNo = placeBetInfo.transferid,
                GameCode = model.BaseGame.brandgameid , // "還沒寫_還沒寫_還沒寫_還沒寫_還沒寫_還沒寫_還沒寫_還沒寫_還沒寫_還沒寫_還沒寫_",
                BetAmount = betAmount,
                GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
                {
                    AGGLive = new AGGLiveApiRawRequestInfoVendor
                    {
                        PlaceBetInfo = model
                    }
                },
            });

            return betResult.GameVendorCustomizedResponse;
        }

        private AGGLiveSingleWalletAPITransactionResponseViewModel SettlementProcess(AGGLiveSingleWalletAPITransactionRequestModel model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
            //string sessionToken = [REDACTED];

            //decimal payoutAmount = GameVendorIntegrationService.GetPayoutAmount(model);

            //var placeBetInfo = model.FundTransferRequest.Funds.FundInfo.FirstOrDefault();

            //decimal betAmount = placeBetInfo.amount.ToAbs();

            //// step 2. 結算、更新錢包、寫投注紀錄
            //var betResult = SingleWalletVendorService.SettleBet(_gameVendor, new SingleWalletBetSettleRequestParameter
            //{
            //    SessionToken = [REDACTED],
            //    BetRecordNo = model.WagerId,
            //    PayoutAmount = payoutAmount,
            //    GameCode = model.GameId,
            //    PlayerID = model.UserID,
            //    GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
            //    {
            //        AGG = new AGGApiRawRequestInfoVendor
            //        {
            //            SettleBetInfo = model
            //        }
            //    }
            //});

            //return betResult.GameVendorCustomizedResponse;
        }


        private AGGLiveSingleWalletAPITransactionResponseViewModel CancelProcess([FromBody] AGGLiveSingleWalletAPITransactionRequestModel model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);

            ////進行取消注單的行為實作
            //var betResult = SingleWalletVendorService.CancelBet(_gameVendor, new SingleWalletCancelBetRequestParameter
            //{
            //    PlayerID = model.FundTransferRequest.accountid,
            //    GameCode = model.GameId,
            //    BetRecordNo = model.WagerId,
            //    GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
            //    {
            //        AGG = new AGGApiRawRequestInfoVendor
            //        {
            //            CancelBet = model
            //        }
            //    }
            //});

            //return betResult.GameVendorCustomizedResponse;
        }



        private bool IsCancelProcessRequest(AGGLiveSingleWalletAPITransactionRequestModel model)
        {
            return true;
        }


        private bool IsSettleProcessRequest(AGGLiveSingleWalletAPITransactionRequestModel model)
        {
            return true;
        }


        private bool IsPlaceBetProcessRequest(AGGLiveSingleWalletAPITransactionRequestModel model)
        {
            bool isSingleFundInfo = model.FundTransferRequest.Funds.FundInfo.Count() == 1;
            bool isDebitAndCreaditFalse = model.FundTransferRequest.Funds.debitandcredit == false;
            var matchFirst = model.FundTransferRequest.Funds.FundInfo.FirstOrDefault();

            return isSingleFundInfo && isDebitAndCreaditFalse && matchFirst.gamestatemode == 1;
        }


    }
}
