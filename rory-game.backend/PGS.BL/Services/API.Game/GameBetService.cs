using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Utility.Helper;
using target.Models;
using target.Models.Attributes;
using target.Utility;
using game.Entity.Game;
using target.Utility.Utility.ThirdParty;

namespace target.BL.Services
{

    public class GameBetService : BaseGameService, IGameBetService
    {

        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }

        public BetRecordRepository BetRecordRepository { get; set; }

        public BetRecordTransactionNoMappingRepository BetRecordTransactionNoMappingRepository { get; set; }

        public ISingleWalletMerchantBetProcessRefService SingleWalletMerchantBetProcessRefService { get; set; }

        public ISingleWalletPlatformApiRequestRefService SingleWalletPlatformApiRequestRefService { get; set; }

        public IPlayerWalletRefService PlayerWalletRefService { get; set; }


        public GamePlaceBetResultViewModel PlaceBet(GamePlaceBetRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.TradeID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.BetAmount).Rule_NotMinus();

            var queryPlayer = GetValidPlayerAndMerchant(model.PlayerID);

            var queryServer = ServerListRepository.GetBriefByCache(model.ServerID);

            if (queryServer == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            var queryGame = GameListRepository.GetSingleByCache(queryServer.GameID);

            //三方遊戲不可能打這支API
            if (AppHelper.IsExternalVendorGame(queryGame.GameID))
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny);
            }

            string transactionNo = Guid.NewGuid().ToString("N");

            DateTime now = DateTimeHelper.GetCurrentTime();


            bool isFreeBonus = model.BetAmount <= 0;
            var recordType = isFreeBonus ? BetRecordTypeEnum.FreeBonus : BetRecordTypeEnum.Normal;

            decimal betGamePoint = model.BetAmount;
            decimal betAmount = ConvertGamePointToMoney(queryPlayer.MerchantID , betGamePoint);

            if (RegexRuleHelper.IsDecimalPointOverThanPlacesTwo(betAmount))
                throw new TerproException(ApiErrorCode.InvalidParameter_NumberValueMustPlacesTwo, new ErrorDetailModel("BetAmount"));

            var insertBetRecord = new BetRecordEntity
            {
                Type = recordType.ToInteger(),
                Status = GameVendorBetRecordStatusEnum.Pending.ToInteger(),
                BetAmount = betAmount,
                CreateTime = now,
                GameID = queryServer.GameID,
                MerchantID = queryPlayer.MerchantID,
                UpdateTime = now,
                PlayerID = queryPlayer.PlayerID,
                PlayerIPAddress = model.IPAddress,
                PlayerName = queryPlayer.PlayerName,
                BetRecordNo = transactionNo,
                BetRecordNoPGS = transactionNo,
                Memo = "",
                TradeID = model.TradeID,
                ServerID = model.ServerID,
                DebugInfo = model.Memo,
                WalletMode = queryPlayer.WalletMode.ToInteger(),
                ConvertRate = GetMerchantConvertRate_PGSGame(queryPlayer.MerchantID)
            };

            var insertMapping = new BetRecordTransactionNoMappingEntity
            {
                CreateTime = now,
                GameVendorID = 0,
                BetRecordNoPGS = transactionNo,
            };

            string processUUID = Guid.NewGuid().ToString("N");

            var queue1 = BetRecordRepository.GenerateSQL_Create(insertBetRecord);
            var queue2 = BetRecordTransactionNoMappingRepository.GenerateSQL_Create(insertMapping);


            // Case 1. freeGame / bonuse / lighting
            if(isFreeBonus)
            {
                // 僅建立一筆投注資料
                BetRecordRepository.ExecuteTransactionSQLs(queue1, queue2);

                // 跟server溝通，若freeGame請不要使用此欄位 (不主動刷新錢包金額)
                decimal zeroAmount = 0;

                return new GamePlaceBetResultViewModel
                {
                    AfterAmount = zeroAmount,
                    TransactionNo = transactionNo,
                };
            }

            // Case 2. 一般投注
            if (CommonLogicHelper.IsSingleWalletMode(queryPlayer.WalletMode))
            {
                var parseJwtToken = [REDACTED];

                // [單錢包模式] 
                // 1.建投注單、單號mapping
                // 2.請求扣款API
                BetRecordRepository.ExecuteTransactionSQLs(queue1, queue2);

                var swResult = SingleWalletMerchantBetProcessRefService.PlaceBetProcessPassToMerchantApi(queryPlayer.MerchantID, new SingleWalletModePlaceBetParameter
                {
                    DebugProcessUUID = processUUID,
                    BetAmount = betAmount,
                    SessionToken = [REDACTED],
                    GameID = queryServer.GameID,
                    BetRecordNoPGS = transactionNo
                });

                var apiResult = InstanceResultViewModelAndRateConvert(queryPlayer.MerchantID , transactionNo, swResult.Balance);
                return apiResult;
            }
            else
            {
                // [轉帳模式] 
                // 1. 先嘗試扣款  
                // 2. 建投注單、單號mapping

                decimal targetAmount = betAmount * -1;

                var result = PlayerWalletRepository.UpdateWalletAmountForbiddenMinus(queryPlayer.PlayerID, queryServer.GameID, targetAmount, TransactionTypeEnum.SingleWallet_PlaceBetDebit, transactionNo, model.IPAddress);

                if (!result.IsSuccess)
                    throw new TerproException(result.Code);

                try
                {
                    BetRecordRepository.ExecuteTransactionSQLs(queue1, queue2);
                }
                catch(Exception ex)
                {
                    decimal refundAmount = targetAmount.ToAbs();

                    //退還金幣
                    PlayerWalletRepository.UpdateWalletAmountForbiddenMinus(queryPlayer.PlayerID, queryServer.GameID, refundAmount, TransactionTypeEnum.SingleWallet_CancelRefundWallet, transactionNo, model.IPAddress);
                    
                    //發個通知
                    SlackOrTelegramNotify(model, queryPlayer, transactionNo , ex.Message);

                    throw ex;
                }

                var apiResult = InstanceResultViewModelAndRateConvert(queryPlayer.MerchantID, transactionNo, result.AfterAmount);
                return apiResult;
            }

        }

        private GamePlaceBetResultViewModel InstanceResultViewModelAndRateConvert(int merchantID , string transactionNo, decimal afterAmount)
        {
            decimal afterGamePoint = ConvertMoneyToGamePoint(merchantID , afterAmount);

            return new GamePlaceBetResultViewModel
            {
                AfterAmount = afterGamePoint, //Gs已對接好API，欄位已經無法改名了
                TransactionNo = transactionNo
            };
        }
               

        private static void SlackOrTelegramNotify(GamePlaceBetRequestModel model, PlayerMerchantInfoModel queryPlayer, string transactionNo , string errorMsg)
        {
            string msg = $"先扣款再建單失敗，可能需要補錢?。{errorMsg}{Environment.NewLine}player:{queryPlayer.PlayerID}_amount:{model.BetAmount}_orderNo:{transactionNo}";
            MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
            {
                Message = msg,
                TelegramChannelEnum = TelegramChannelEnum.pgs_single_wallet_erorr
            });
        }
    }
}
