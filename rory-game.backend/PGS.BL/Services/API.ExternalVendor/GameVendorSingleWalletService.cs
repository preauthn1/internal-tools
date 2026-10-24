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
using System.Collections.Concurrent;
using target.Utility.Helper;
using game.Entity.Game;
using System.Threading.Tasks;
using System.Threading;
using target.Utility.Utility.ThirdParty;
using target.Utility.Helper.HttpContext;
using target.Models;
using target.Models.ViewModels.GameVendor;

namespace target.BL.Services
{

    public class GameVendorSingleWalletService : BaseService, IGameVendorSingleWalletService
    {
        public VendorMaintenanceRepository VendorMaintenanceRepository { get; set; }

        public VendorSessionTokenRepository VendorSessionTokenRepository { get; set; }

        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public PerGameLogRepository PerGameLogRepository { get; set; }

        public PerGameDetailLogRepository PerGameDetailLogRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public BetRecordTransactionNoMappingRepository BetRecordTransactionNoMappingRepository { get; set; }

        public PerGameDetailLogTagMappingRepository PerGameDetailLogTagMappingRepository { get; set; }


        public IPlayerWalletRefService PlayerWalletRefService { get; set; }

        public ISingleWalletMerchantBetProcessRefService SingleWalletMerchantBetProcessRefService { get; set; }

        public ISingleWalletPlatformApiRequestRefService SingleWalletPlatformApiRequestRefService { get; set; }




        public SingleWalletGetWalletBalanceResultViewModel GetWalletBalance(GameVendorTypeEnum vendorID , SingleWalletGetWalletBalanceRequestParameter model)
        {
            int playerID = GetPlayerIdFromRequest(model);

            var queryPlayer = PlayerRepository.GetBriefByCache(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var queryMerchant = MerchantRepository.GetSingleByCache(queryPlayer.MerchantID);

            var queryWallet = PlayerWalletRefService.GetWalletAmountByMerchantWalletMode(queryPlayer.PlayerID, model.SessionToken);

            string currency = AppHelper.ConvertCurrencyBETA_To_CNY(queryMerchant.Currency);

            decimal gamePoint = ConvertMoneyToGamePoint(queryPlayer.MerchantID , vendorID , queryWallet.Amount);

            return new SingleWalletGetWalletBalanceResultViewModel
            {
                PlayerID = playerID,
                PlayerName = queryPlayer.PlayerName,
                BalanceAmount = queryWallet.Amount,
                BalanceGamePoint = gamePoint,
                Currency = currency
            };
        }


        public SingleWalletVendorCustomizedResponseConvertModel PlaceBet(GameVendorTypeEnum vendorID, SingleWalletBetDebitAmountRequestParameter model)
        {
            string debugUUID = Guid.NewGuid().ToString();

            
            model.RuleFor(z => z.BetRecordNo).Rule_NotEmptyText();

            int playerID = GetPlayerIdFromRequest(model);

            var queryPlayer = PlayerRepository.GetBriefByCache(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var queryMerchant = MerchantRepository.GetSingleByCache(queryPlayer.MerchantID);

            IGameVendorIntegrationService vendorService = ServiceInstanceHelper.InstanceVendorService(vendorID);

            //檢查遊戲存在/遊戲開啟
            var queryGame = GameListRepository.GetGameCacheByVendorGameCode(vendorID.ToInteger(), model.GameCode);

            if (queryGame == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            if (queryGame.Status != NormalActiveStatus.Enable.ToInteger())
                throw new TerproException(ApiErrorCode.NotFoundGame);

            //檢查遊戲商維護
            var queryMaintain = VendorMaintenanceRepository.GetVendorMaintenanceByCache(vendorID.ToInteger());

            if (queryMaintain != null)
                throw new TerproException(ApiErrorCode.SystemMaintain);

            // step 1. 檢查單號是否重複
            var queryBetRecord = vendorService.GetBetRecordByNo(model.BetRecordNo, model.GameCode);

            if (queryBetRecord != null)
            {
                if(vendorID == GameVendorTypeEnum.CrazyFishing)
                {
                    //捕魚這家機制不夠完善 直接進入取消注單程序
                    FishingGameSecondTimesCancelProcess(vendorID, model, queryPlayer, queryBetRecord);
                }
                else
                {
                    throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateOrderNo);
                }
            }

            string transactionNoPGS = Guid.NewGuid().ToString("N");

            string processUUID = Guid.NewGuid().ToString("N");

            UpdateBalanceSPResult updateResult = null;

            //需要經過 遊戲幣與法幣轉換
            decimal betPoint = model.BetAmount;
            decimal betAmount = ConvertGamePointToMoney(queryPlayer.MerchantID , vendorID , betPoint);

            // step 2. 更新錢包金額(禁止負數)
            if (CommonLogicHelper.IsSingleWalletMode(queryMerchant.WalletMode))
            {
                // 2-1 單錢包模式 請求對方API
                var swResult = SingleWalletMerchantBetProcessRefService.PlaceBetProcessPassToMerchantApi(queryMerchant.MerchantID, new SingleWalletModePlaceBetParameter
                {
                    PlayerID = queryPlayer.PlayerID,
                    DebugProcessUUID = processUUID,
                    SessionToken = [REDACTED],
                    BetAmount = betAmount,
                    GameID = queryGame.GameID,
                    BetRecordNoPGS = transactionNoPGS
                });

                decimal swAfterAmount = swResult.Balance;

                updateResult = new UpdateBalanceSPResult
                {
                    AfterAmount = swAfterAmount, //後面的程式碼會轉成gamePoint
                    Code = ApiErrorCode.SuccessOK
                };
            }
            else
            {
                // 2-2 轉帳模式，更新PGS自己的錢包金額
                var transactionType = TransactionTypeEnum.SingleWallet_PlaceBetDebit;
                decimal updateAmount = betAmount > 0 ? betAmount * -1 : betAmount;
                updateResult = PlayerWalletRepository.UpdateWalletAmountForbiddenMinus(playerID, queryGame.GameID, updateAmount, transactionType, model.BetRecordNo, model.PlayerIPAddress);
            }

            if (!updateResult.IsSuccess)
            {
                throw new TerproException((ApiErrorCode)updateResult.Code);
            }

            try
            {
                // step 3. 寫入一筆單號對應Mapping
                BetRecordTransactionNoMappingRepository.Create(new BetRecordTransactionNoMappingEntity
                {
                    BetRecordNoPGS = transactionNoPGS,
                    BetRecordNo = model.BetRecordNo,
                    GameVendorID = vendorID.ToInteger(),
                    CreateTime = DateTimeHelper.GetCurrentTime()
                });

                var requiredInfo = new InsertRawBetRecordRequiredInfoModel
                {
                    GameID = queryGame.GameID,
                    MerchantID = queryPlayer.MerchantID,
                    PlayerID = queryPlayer.PlayerID,
                    PlayerName = queryPlayer.PlayerName,
                    BetRecordNoPGS = transactionNoPGS,
                    BetAmount = betAmount,
                    BetGamePoint = betPoint,
                    ConvertRate = GetMerchantConvertRate(queryPlayer.MerchantID , vendorID).ConvertRate
                };

                // step 4. 寫入一筆投注紀錄 (raw betRecord)
                vendorService.WriteRawBetRecordWhenPlaceBet(requiredInfo, model.GameVendorRawRequestInfo);

                string currency = AppHelper.ConvertCurrencyBETA_To_CNY(queryMerchant.Currency);

                decimal aftetAmount = updateResult.AfterAmount;
                decimal aftetGamePoint = ConvertMoneyToGamePoint(queryPlayer.MerchantID, vendorID, aftetAmount);

                // step 5. 轉成遊戲商要求的API格式
                var converResponse = vendorService.FormatParsing_ConvertToPlaceBetResponseResult(new SingleWalletPlaceBetDebitAmountResponseParameter
                {
                    Currency = currency,
                    AfterAmount = aftetAmount,
                    AfterGamePoint = aftetGamePoint
                });

                return new SingleWalletVendorCustomizedResponseConvertModel
                {
                    GameVendorCustomizedResponse = converResponse
                };

            }
            catch (Exception ex)
            {
                string msg = $"{vendorID.ToString()} PlaceBet 後半段DB寫入錯誤、可能需要補救處理?，三方單號:{model.BetRecordNo}";

                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                {
                    Channel = SlackChannelEnum.n_error_external_api,
                    Message = msg
                });

                throw ex;
            }
        }

        private void FishingGameSecondTimesCancelProcess(GameVendorTypeEnum vendorID, SingleWalletBetDebitAmountRequestParameter model, PlayerBriefModel queryPlayer, Models.ViewModels.GameVendor.BetRecordGeneralColumnModel queryBetRecord)
        {
            if (queryBetRecord.Status == (int)GameVendorBetRecordStatusEnum.Pending)
            {
                CancelBet(vendorID, new SingleWalletCancelBetRequestParameter
                {
                    BetRecordNo = model.BetRecordNo,
                    PlayerID = queryPlayer.PlayerID,
                    GameCode = model.GameCode,
                    PlayerName = queryPlayer.PlayerName,
                });
            }

            throw new TerproException(ApiErrorCode.InvalidStatus);
        }

        public SingleWalletVendorCustomizedResponseConvertModel SettleBet(GameVendorTypeEnum vendorID, SingleWalletBetSettleRequestParameter model)
        {
            int playerID = GetPlayerIdFromRequest(model);

            var queryPlayer = PlayerRepository.GetBriefByCache(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var queryGame = GameListRepository.GetGamebyVendorGameCode(vendorID.ToInteger(), model.GameCode);
            
            if (queryGame == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            var queryMerchant = MerchantRepository.GetSingleByCache(queryPlayer.MerchantID);

            IGameVendorIntegrationService vendorService = ServiceInstanceHelper.InstanceVendorService(vendorID);

            // step 1. 檢查注單存在 (檢查狀態是在 UpdateRawBetRecordAfterSettlement 處理)
            var queryBetRecord = vendorService.GetBetRecordByNo(model.BetRecordNo, queryGame?.ThirdPartyGameCode);

            if (queryBetRecord == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            string requestUUID = Guid.NewGuid().ToString();

            var requiredInfo = new SettlementBetRecordRequiredInfoModel
            {
                MerchantID = queryPlayer.MerchantID,
                BetRecordNoPGS = queryBetRecord.BetRecordNoPGS,
                DebugExtraMemo = requestUUID
            };

            // step 2. 更新 raw records (輸贏金額、狀態)
            vendorService.UpdateRawBetRecordAfterSettlement(requiredInfo, model.GameVendorRawRequestInfo);

            // step 3. 轉成perGameLog格式
            var convertPerGame = vendorService.FormatParsing_ConvertBetRecordToPerGameDetailLog(queryBetRecord.BetRecordNoPGS);

            // step 4. 寫入perGameLog
            long perGameDetailID = WritePerGameLog(convertPerGame, queryBetRecord);

            UpdateBalanceSPResult updateResult = null;

            // step 5. 更新錢包金額
            if (CommonLogicHelper.IsSingleWalletMode(queryMerchant.WalletMode))
            {
                //單錢包模式
                var swResult = SingleWalletMerchantBetProcessRefService.SettlemenetProcessPassToMerchantApi(new SingleWalletModeSettleApiParameter
                {
                    PlayerID = convertPerGame.PlayerID,
                    PayoutAmount = convertPerGame.PayoutAmount,
                    BetAmount = convertPerGame.BetAmount,
                    GameID = convertPerGame.GameID,
                    MerchantID = convertPerGame.MerchantID,
                    TaxAmount = convertPerGame.TaxAmount,
                    ValidAmount = convertPerGame.ValidAmount,
                    BetRecordNoPGS = queryBetRecord.BetRecordNoPGS,
                    SettlementTime = convertPerGame.CreateTime,
                    DebugProcessUUID = requestUUID
                });

                decimal swAfterAmount = swResult.Balance;

                updateResult = new UpdateBalanceSPResult
                {
                    AfterAmount = swAfterAmount, //後面的程式碼會轉成gamePoint
                    Code = ApiErrorCode.Success
                };
            }
            else
            {
                //轉帳模式
                var transactionType = TransactionTypeEnum.SingleWallet_BetSettlement;

                decimal payoutPoint = model.PayoutAmount;
                decimal payoutAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, vendorID, payoutPoint);

                updateResult = PlayerWalletRepository.UpdateWalletAmountForbiddenMinus(playerID, queryGame.GameID, payoutAmount, transactionType, model.BetRecordNo, "");
            }

            //PGS後台中 [當前金幣] 慣例上是顯示 未下注前(不含預扣)的金幣 (在三方遊戲架構下，BeforeAmountc會拿到 已預扣betAmount的金額，因此統一轉換一致的格式)
            decimal currentAmount = updateResult.BeforeAmount + queryBetRecord.BetAmount;

            // step 6. 更新金額
            PerGameDetailLogRepository.UpdateCurrentAmount(perGameDetailID, currentAmount);

            string currency = AppHelper.ConvertCurrencyBETA_To_CNY(queryMerchant.Currency);

            decimal afterAmount = updateResult.AfterAmount;
            decimal afterGamePoint = ConvertMoneyToGamePoint(queryPlayer.MerchantID , vendorID , afterAmount);

            // step 7. 轉成每一家遊戲商要求的 response格式
            var converResponse = vendorService.FormatParsing_ConvertToSettleBetResponseResult(new SettlementConvertRequiredInfoModel
            {
                Currency = currency,
                AfterAmount = updateResult.AfterAmount,
                AfterGamePoiint = afterGamePoint,
                BetRecordNo = model.BetRecordNo,
                PlayerName = queryPlayer.PlayerName,
                PlayerID = playerID,
            });

            return new SingleWalletVendorCustomizedResponseConvertModel
            {
                GameVendorCustomizedResponse = converResponse
            };
        }


        public SingleWalletVendorCustomizedResponseConvertModel CancelBet(GameVendorTypeEnum vendorID, SingleWalletCancelBetRequestParameter model)
        {
            int playerID = GetPlayerIdFromRequest(model);

            var queryPlayer = PlayerRepository.GetBriefByCache(model.PlayerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var queryMerchant = MerchantRepository.GetSingleByCache(queryPlayer.MerchantID);

            IGameVendorIntegrationService vendorService = ServiceInstanceHelper.InstanceVendorService(vendorID);

            var queryBetRecord = vendorService.GetBetRecordByNo(model.BetRecordNo, model.GameCode);

            if (queryBetRecord == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            //Step 1. 更新狀態
            var cancelResult = vendorService.CancelBet(queryBetRecord.BetRecordNoPGS, "");

            if (CommonLogicHelper.IsSingleWalletMode(queryMerchant.WalletMode))
            {

                // [單錢包模式] 請求三方API 告知取消注單
                var queryGame = GameListRepository.GetGamebyVendorGameCode(vendorID.ToInteger(), model.GameCode);
                SingleWalletMerchantBetProcessRefService.CancelProcessPassToMerchantApi(new SingleWalletModeCancelApiRequestModel
                {
                    MerchantID = queryBetRecord.MerchantID,
                    GameID = queryBetRecord.GameID,
                    TransactionNo = queryBetRecord.BetRecordNoPGS
                });
            }
            else
            {
                // [轉帳模式] 退還金幣
                PlayerWalletRepository.UpdateWalletAmountForbiddenMinus(queryBetRecord.PlayerID, queryBetRecord.GameID, queryBetRecord.BetAmount, TransactionTypeEnum.SingleWallet_CancelRefundWallet, queryBetRecord.BetRecordNoPGS, "");
            }

            return new SingleWalletVendorCustomizedResponseConvertModel
            {
                GameVendorCustomizedResponse = cancelResult.GameVendorCustomizedResponse
            };
        }


        private long WritePerGameLog(PerGameDetailLogEntity convertPerGame , BetRecordGeneralColumnModel betRecordInfo)
        {
            //step 1. 寫入perGame
            long perGameID = PerGameLogRepository.Create(new PerGameLogEntity
            {
                BeginTime = convertPerGame.GameTime,
                EndTime = convertPerGame.GameTime,
                GameData = "",
                ServerID = convertPerGame.ServerID,
                GameID = convertPerGame.GameID,
                TaxAmount = 0,
                CreateTime = DateTimeHelper.GetCurrentTime()
            });

            convertPerGame.PerGameLogID = perGameID;

            //step 2. 寫入perGameDetail
            long resultPerGameDetailID = PerGameDetailLogRepository.Create(convertPerGame);

            //step 3. 如果是jackpot 標記Tag來記錄 (用此方法可以很快的查到jackpot資料) 
            if(convertPerGame.JackpotAmount > 0)
            {
                PerGameDetailLogTagMappingRepository.Create(new PerGameDetailLogTagMappingEntity
                {
                    PerGameDetailLogID = resultPerGameDetailID,
                    TagID = SystemTagEnum.PerGameDetailLog_Jackpot.ToInteger()
                });
            }

            return resultPerGameDetailID;
        }



        private (bool IsDefine, decimal ConvertRate) GetMerchantConvertRate(int merchantID, GameVendorTypeEnum gameVendor)
        {
            var settingKeyResult = AppHelper.GetConvertRateCustomizedSettingKey(gameVendor);

            if (!settingKeyResult.IsDefine)
                return (false , 1.0M);

            decimal rate = MerchantCustomizedSettingRepository.GetMerchantGamePointConvertRateByCache(merchantID, settingKeyResult.SettingKey);
            return (true, rate);
        }


        private decimal ConvertMoneyToGamePoint(int merchantID, GameVendorTypeEnum gameVendor, decimal money)
        {
            var settingKeyResult = GetMerchantConvertRate(merchantID , gameVendor);
            return money * settingKeyResult.ConvertRate;
        }


        private decimal ConvertGamePointToMoney(int merchantID , GameVendorTypeEnum gameVendor, decimal point)
        {
            var settingKeyResult = GetMerchantConvertRate(merchantID, gameVendor);
            return point / settingKeyResult.ConvertRate;
        }


        private int GetPlayerIdFromRequest<T>(T model) where T : SingleWalletBasePlayerInfoModel
        {
            if (model.PlayerID != 0)
            {
                AddPlayerIDInfoToHeader(model.PlayerID);
                return model.PlayerID;
            }

            if (model.SessionToken.IsNotNullEmptyText())
            {
                //ps. controller層有檢查過期了
                var queryToken = [REDACTED];
                int playerID = queryToken == null ? 0 : queryToken.PlayerID;
                AddPlayerIDInfoToHeader(playerID);
                return playerID;
            }

            //其他遊戲商可能要用playerName反查 (有其他遊戲商在根據狀況實作)
            return 0;
        }


        /// <summary>
        /// 透過此方式，ApiLog會獨立某個欄位記錄是哪個 playerID
        /// </summary>
        private void AddPlayerIDInfoToHeader(int playerID)
        {
            if (MyHttpContext.Current.Request.Headers.ContainsKey(SystemConstant.HeaderKey_IdentityMemo))
                return;

            MyHttpContext.Current.Request.Headers.Add(SystemConstant.HeaderKey_IdentityMemo, playerID.ToString());
        }


        /// <summary>
        /// (沒啥意義 還是算了) 如果注單結算有跨日 且 時間差很多，標註下
        /// </summary>
        private void AddTagIfSettleCrossTooMuchTime(PerGameDetailLogEntity convertPerGame, BetRecordGeneralColumnModel betRecordInfo, long resultPerGameDetailID)
        {
            int dayOfBet = betRecordInfo.CreateTime.DayOfYear;
            int dayOfSettle = convertPerGame.GameTime.DayOfYear;

            if (dayOfBet == dayOfSettle)
                return;

            int totalMinutes = (int)(convertPerGame.GameTime - betRecordInfo.CreateTime).TotalMinutes;

            if (totalMinutes <= 30)
                return;

            PerGameDetailLogTagMappingRepository.Create(new PerGameDetailLogTagMappingEntity
            {
                PerGameDetailLogID = resultPerGameDetailID,
                TagID = SystemTagEnum.PerGameDetailLog_CrossDateSettle.ToInteger()
            });
        }


    }
}
