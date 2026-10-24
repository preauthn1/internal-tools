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
    public class AGGBetController : BaseGameVendorController
    {
        protected override GameVendorTypeEnum _gameVendor
        {
            get
            {
                return GameVendorTypeEnum.AGG;
            }
        }

        public AGGIntegrationService GameVendorIntegrationService { get; set; }


        public AGGBetController()
        {
            GameVendorIntegrationService = new AGGIntegrationService();
        }


        [ApiConfig(ApiCategory = ApiConfigCategoryEnum.SingleWalletApi_AGGApi , IsHeaderLogEnable = false)]
        [Route("agg/wager/transaction")]
        [HttpPost]
        public AGGSingleWalletAPIPlaceBetResponseViewModel PlaceBet([FromBody]AGGSingleWalletAPITransactionRequestModel model)
        {
            model.RuleFor(z => z.AccountToken).Rule_NotEmptyText();

            // Case 1. 如果是結算、執行結算加錢的邏輯
            if (IsSettleProcessRequest(model))
            {
                var settleResult = SettlementProcess(model);
                return settleResult;
            }

            // Case 2. 投注流程
            if (IsPlaceBetProcessRequest(model))
            {
                var result = PlaceBetProcess(model);
                return result;
            }


            if(IsCancelProcessRequest(model))
            {
                var result = CancelProcess(model);
                return result;
            }

            throw new TerproException(ApiErrorCode.NotImplement);
        }


        private AGGSingleWalletAPIPlaceBetResponseViewModel PlaceBetProcess(AGGSingleWalletAPITransactionRequestModel model)
        {
            DateTime time0 = DateTimeHelper.GetCurrentTime();

            string sessionToken = [REDACTED];

            decimal betAmount = GameVendorIntegrationService.GetBetAmountFromOperationModel(model?.Operations);


            // Step 1. 先檢查session是否合法
            GameVendorIntegrationService.CheckIsValidRequestInfo(new ExternalVendorInputValidateRequestParameter
            {
                SessionOrToken = [REDACTED],
            });

            // Step 2. 如果是投注、投注預先扣錢
            var betResult = SingleWalletVendorService.PlaceBet(_gameVendor, new SingleWalletBetDebitAmountRequestParameter
            {
                SessionToken = [REDACTED],
                PlayerID = model.UserID,
                BetRecordNo = model.WagerId,
                GameCode = model.GameId,
                BetAmount = betAmount,
                GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
                {
                    AGG = new AGGApiRawRequestInfoVendor
                    {
                        PlaceBetInfo = model
                    }
                },
            });

            Develop_Simulate_DelayTimeout(SystemConfigSettingEnum.GameVendor_PlaceBet_SuccessUpdateBut_DelayTimeout);
            return betResult.GameVendorCustomizedResponse;
        }

        private AGGSingleWalletAPIPlaceBetResponseViewModel SettlementProcess(AGGSingleWalletAPITransactionRequestModel model)
        {
            string sessionToken = [REDACTED];

            decimal payoutAmount = GameVendorIntegrationService.GetPayoutAmount(model); 

            // 結算、更新錢包、寫投注紀錄
            var betResult = SingleWalletVendorService.SettleBet(_gameVendor, new SingleWalletBetSettleRequestParameter
            {
                SessionToken = [REDACTED],
                BetRecordNo = model.WagerId,
                PayoutAmount = payoutAmount,
                GameCode = model.GameId,
                PlayerID = model.UserID,
                GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
                {
                    AGG = new AGGApiRawRequestInfoVendor
                    {
                        SettleBetInfo = model
                    }
                }
            });

            Develop_Simulate_DelayTimeout(SystemConfigSettingEnum.GameVendor_Settle_SuccessUpdateBut_DelayTimeout);
            return betResult.GameVendorCustomizedResponse;
        }


        private AGGSingleWalletAPIPlaceBetResponseViewModel CancelProcess([FromBody] AGGSingleWalletAPITransactionRequestModel model)
        {
            //進行取消注單的行為實作
            var betResult = SingleWalletVendorService.CancelBet(_gameVendor, new SingleWalletCancelBetRequestParameter
            {
                PlayerID = model.UserID,
                GameCode = model.GameId,
                BetRecordNo = model.WagerId,
                GameVendorRawRequestInfo = new EachGameVendorRawRequestInfoModel
                {
                    AGG = new AGGApiRawRequestInfoVendor
                    {
                        CancelBet = model
                    }
                }
            });

            return betResult.GameVendorCustomizedResponse;
        }



        private string GetWagerLifecycleOperation(AGGSingleWalletAPITransactionRequestModel model)
        {
            model.RuleFor(z => z.Operations).Rule_NotEmptyNull();

            var actionType = model.Operations.FirstOrDefault(z => z.OperationType == "WagerLifecycleOperation")?.Action;

            return actionType;
        }


        private bool IsCancelProcessRequest(AGGSingleWalletAPITransactionRequestModel model)
        {
            var actionType = GetWagerLifecycleOperation(model);
            return actionType == "cancel";
        }


        private bool IsSettleProcessRequest(AGGSingleWalletAPITransactionRequestModel model)
        {
            var actionType = GetWagerLifecycleOperation(model);
            return actionType == "confirm";
        }


        private bool IsPlaceBetProcessRequest(AGGSingleWalletAPITransactionRequestModel model)
        {
            var actionType = GetWagerLifecycleOperation(model);
            return actionType == "create";
        }



    }
}
