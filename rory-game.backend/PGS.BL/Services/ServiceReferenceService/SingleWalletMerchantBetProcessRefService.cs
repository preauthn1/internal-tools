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
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;
using target.Utility.Utility;
using target.Utility;
using game.Entity.Game;
using System.Threading.Tasks;
using System.Threading;
using target.Utility.Utility.ThirdParty;

namespace target.BL.Services
{

    public class SingleWalletMerchantBetProcessRefService : BaseService, ISingleWalletMerchantBetProcessRefService
    {

        public ISingleWalletPlatformApiRequestRefService SingleWalletPlatformApiRequestRefService { get; set; }

        public BetRecordRepository BetRecordRepository { get; set; }

        public BetRecordMerchantNotifyQueueRepository BetRecordMerchantNotifyQueueRepository { get; set; }


        public SingleWalletModePlaceBetApiResponseModel PlaceBetProcessPassToMerchantApi(int merchantID, SingleWalletModePlaceBetParameter model)
        {
            try
            {
                var apiResult = SingleWalletPlatformApiRequestRefService.PlaceBet(model);
                return apiResult;
            }
            catch (Exception ex)
            {
                string msg = $"SingleWalletPlaceBetPassToMerchantApi、將取消注單_TransactionNo: {model.BetRecordNoPGS}{Environment.NewLine}{ex.Message}";

                MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
                {
                    Message = msg,
                    TelegramChannelEnum = TelegramChannelEnum.pgs_single_wallet_erorr
                });

                //取消此單
                if (!AppHelper.IsExternalVendorGame(model.GameID))
                {
                    BetRecordRepository.UpdateStatus(model.BetRecordNoPGS, GameVendorBetRecordStatusEnum.Cancel, msg);
                }                

                //強制cancel這張單
                CancelProcessPassToMerchantApi(new SingleWalletModeCancelApiRequestModel
                {
                    PlayerID = model.PlayerID,
                    DebugProcessUUID = model.DebugProcessUUID,
                    MerchantID = merchantID,
                    GameID = model.GameID,
                    TransactionNo = model.BetRecordNoPGS,
                });

                //強制拋錯誤 讓此投注失敗
                throw ex;
            }
        }


        public SingleWalletModeSettleApiResponseModel SettlemenetProcessPassToMerchantApi(SingleWalletModeSettleApiParameter model)
        {         
            try
            {
                var apiResult = SingleWalletPlatformApiRequestRefService.Settlement(model);

                //最理想的狀況 有拿到回應
                return apiResult;

            }
            catch (Exception ex)
            {
                //失敗 再嘗試retry一次
                Task.Factory.StartNew(() =>
                {
                    Thread.Sleep(3000);

                    try
                    {
                        SingleWalletPlatformApiRequestRefService.Settlement(model);
                    }
                    catch(Exception ex)
                    {
                        string msg = GenerateErrorMsg("Settlement", model.BetRecordNoPGS, model.DebugProcessUUID, model.GameID , ex.Message);

                        MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
                        {
                            Message = msg,
                            TelegramChannelEnum = TelegramChannelEnum.pgs_single_wallet_erorr
                        });

                        //加入補救排程佇列
                        BetRecordMerchantNotifyQueueRepository.Create(new BetRecordMerchantNotifyQueueEntity
                        {
                            PlayerID = model.PlayerID,
                            MerchantID = model.MerchantID,
                            GameID = model.GameID,
                            LastRequestUUID = model.DebugProcessUUID,
                            BetRecordNoPGS = model.BetRecordNoPGS,
                            IsSuccess = false,
                            CreateTime = DateTimeHelper.GetCurrentTime(),
                            RetryTimes = 1,
                            NextNotifyTime = DateTimeHelper.GetCurrentTime().AddMinutes(2),
                            Type = BetRecordMerchantNotifyQueueTypeEnum.Settle.ToInteger(),
                        });                    
                    }

                });
            }

            //若失敗過一次 強制拋錯誤(讓玩家強制需要登出/重連)
            throw new TerproException(ApiErrorCode.ThirdPartyApiFail_Timeout);

            //若失敗過一次 強制進行下一個步驟，balance直接返回0元 (不希望處理太久 導致整個thread timeout掛掉)
            return new SingleWalletModeSettleApiResponseModel
            {
                Balance = 0,
            };
        }



        public SingleWalletModeSettleApiResponseModel BonusProcessPassToMerchantApi(SingleWalletModeSettleApiParameter model)
        {
            try
            {
                var apiResult = SingleWalletPlatformApiRequestRefService.BonusSettlement(model);

                //最理想的狀況 有拿到回應
                return apiResult;

            }
            catch (Exception ex)
            {
                //失敗 再嘗試retry一次
                Task.Factory.StartNew(() =>
                {
                    Thread.Sleep(3000);

                    try
                    {
                        SingleWalletPlatformApiRequestRefService.BonusSettlement(model);
                    }
                    catch (Exception ex)
                    {
                        string msg = GenerateErrorMsg("BonusSettlement", model.BetRecordNoPGS, model.DebugProcessUUID, model.GameID, ex.Message);

                        MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
                        {
                            Message = msg,
                            TelegramChannelEnum = TelegramChannelEnum.pgs_single_wallet_erorr
                        });

                        //加入補救排程佇列
                        BetRecordMerchantNotifyQueueRepository.Create(new BetRecordMerchantNotifyQueueEntity
                        {
                            PlayerID = model.PlayerID,
                            MerchantID = model.MerchantID,
                            GameID = model.GameID,
                            LastRequestUUID = model.DebugProcessUUID,
                            BetRecordNoPGS = model.BetRecordNoPGS,
                            IsSuccess = false,
                            CreateTime = DateTimeHelper.GetCurrentTime(),
                            RetryTimes = 1,
                            NextNotifyTime = DateTimeHelper.GetCurrentTime().AddMinutes(2),
                            Type = BetRecordMerchantNotifyQueueTypeEnum.BonusSettle.ToInteger(),
                        });
                    }
                });
            }

            //若失敗過一次 強制進行下一個步驟，balance直接返回0元 (不希望處理太久 導致整個thread timeout掛掉)
            return new SingleWalletModeSettleApiResponseModel
            {
                Balance = 0,
            };
        }



        public SingleWalletModeSettleApiResponseModel CancelProcessPassToMerchantApi(SingleWalletModeCancelApiRequestModel model)
        {
            try
            {
                SingleWalletPlatformApiRequestRefService.Cancel(model);
            }
            catch(Exception ex1)
            {         
                //N秒後再嘗試重新發送
                Task.Factory.StartNew(() =>
                {
                    Thread.Sleep(3000); 

                    try
                    {
                        SingleWalletPlatformApiRequestRefService.Cancel(model);
                    }
                    catch(Exception ex)
                    {
                        string msg = GenerateErrorMsg("Cancel", model.TransactionNo, model.DebugProcessUUID, model.GameID , ex.Message);

                        MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
                        {
                            Message = msg,
                            TelegramChannelEnum = TelegramChannelEnum.pgs_single_wallet_erorr
                        }); 
                        
                        //加入補救排程佇列
                        BetRecordMerchantNotifyQueueRepository.Create(new BetRecordMerchantNotifyQueueEntity
                        {
                            PlayerID = model.PlayerID,
                            MerchantID = model.MerchantID,
                            GameID = model.GameID,
                            LastRequestUUID = model.DebugProcessUUID,
                            BetRecordNoPGS = model.TransactionNo,
                            IsSuccess = false,
                            CreateTime = DateTimeHelper.GetCurrentTime(),
                            RetryTimes = 1,
                            NextNotifyTime = DateTimeHelper.GetCurrentTime().AddMinutes(2),
                            Type = BetRecordMerchantNotifyQueueTypeEnum.Cancel.ToInteger(),
                        });
                    }

                });
            }

            //若失敗過一次 強制進行下一個步驟，balance直接返回0元 (不希望處理太久 導致整個thread timeout掛掉)
            return new SingleWalletModeSettleApiResponseModel
            {
                Balance = 0,
            };
        }

        public SingleWalletModePlaceBetApiResponseModel ManualSettlementNotifyByNo(int merchantID, SingleWalletModePlaceBetParameter model)
        {
            return null; 
        }

        private string GenerateErrorMsg(string type, string betRecordNoPGS, string requestUUID, int gameID , string exMessage)
        {
            string newline = Environment.NewLine;
            return $"SingleWalletPlatformApiRequestRefService_{type} 通知二度失敗，{newline}"+
                   $"GameID：{gameID}{newline}TransactionNo：{betRecordNoPGS}{newline}RequestUUID：{requestUUID}{newline}" +
                   $"Msg：{exMessage}";
        }


    }
}
