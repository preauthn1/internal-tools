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
using target.Models;
using target.Utility.Helper;
using target.Models.ViewModels.GameVendor;
using target.Utility.Utility.ThirdParty;
using System.Threading.Tasks;
using game.Entity.Game;

namespace target.BL.Services
{

    public class BetRecordStatusChangeNotifyService : BaseService , IBetRecordStatusChangeNotifyService
    {     
        public BetRecordMerchantNotifyQueueRepository BetRecordMerchantNotifyQueueRepository { get; set; }

        public BetRecordTransactionNoMappingRepository BetRecordTransactionNoMappingRepository { get; set; }

        public ISingleWalletPlatformApiRequestRefService SingleWalletPlatformApiRequestRefService { get; set; }


        public void BetReocrdNotifyQueueRetryProcess()
        {
            //取得發送時間小於現在且...
            var queryQueue = BetRecordMerchantNotifyQueueRepository.GetByFilter(new QueryBetRecordMerchantNotifyQueueRequestViewModel
            {
                IsMustNotify = true
            });

            if (queryQueue.IsEmptyOrNull())
                return;

            //先更新下一次時間
            foreach (var item in queryQueue)
            {
                DateTime newNextTime = CalculateNextTime(item.NextNotifyTime , item.RetryTimes);
                BetRecordMerchantNotifyQueueRepository.UpdateNextTime(item.Id, newNextTime);
            }

            foreach (var item in queryQueue)
            {
                //採用非同步處理 避免某一家回應太慢卡住其他queue
                Task.Factory.StartNew(() =>
                {
                    //觸發指定的API
                    TriggerMerchantApiByType(item);
                });
            }
        }


        public void TriggerRetryByID(long id)
        {
            var queryQueue = BetRecordMerchantNotifyQueueRepository.GetById(id);

            if (queryQueue == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            if (queryQueue.IsSuccess)
                throw new TerproException(ApiErrorCode.InvalidStatus);

            //先更新下次時間
            DateTime newNextTime = CalculateNextTime(queryQueue.NextNotifyTime, queryQueue.RetryTimes);
            BetRecordMerchantNotifyQueueRepository.UpdateNextTime(queryQueue.Id, newNextTime);

            //觸發指定的API
            TriggerMerchantApiByType(queryQueue);
        }


        private void TriggerMerchantApiByType(BetRecordMerchantNotifyQueueEntity item)
        {
            var type = (BetRecordMerchantNotifyQueueTypeEnum)item.Type;
            switch (type)
            {
                case BetRecordMerchantNotifyQueueTypeEnum.Cancel:
                    CancelProcess(item);
                    break;

                case BetRecordMerchantNotifyQueueTypeEnum.Settle:
                    SettleProcess(item);
                    break;

                case BetRecordMerchantNotifyQueueTypeEnum.BonusSettle:
                    BonusSettleProcess(item);
                    break;

                default:
                    throw new TerproException(ApiErrorCode.NotImplement);
            }
        }

        private void CancelProcess(BetRecordMerchantNotifyQueueEntity item)
        {
            try
            {
                SingleWalletPlatformApiRequestRefService.Cancel(new SingleWalletModeCancelApiRequestModel
                {
                    PlayerID = item.PlayerID,
                    GameID = item.GameID,
                    DebugProcessUUID = item.LastRequestUUID,
                    MerchantID = item.MerchantID,
                    TransactionNo = item.BetRecordNoPGS
                });

                BetRecordMerchantNotifyQueueRepository.UpdateToSuccess(item.Id);
            }
            catch(Exception ex)
            {

            }
        }


        private void SettleProcess(BetRecordMerchantNotifyQueueEntity item)
        {
            try
            {
                var queryMapping = BetRecordTransactionNoMappingRepository.GetByTransactionNo(item.BetRecordNoPGS);

                if (queryMapping == null)
                    throw new TerproException(ApiErrorCode.NotFoundRecord);

                var service = ServiceInstanceHelper.InstanceVendorService(queryMapping.GameVendorID);

                var queryBetReocrd = service.GetBetRecordByPGSTransactionNo(item.BetRecordNoPGS);

                SingleWalletPlatformApiRequestRefService.Settlement(new SingleWalletModeSettleApiParameter
                {
                    PlayerID = item.PlayerID,
                    TaxAmount = queryBetReocrd.TaxAmount,
                    SettlementTime = queryBetReocrd.SettlementTime,
                    BetRecordNoPGS = queryBetReocrd.BetRecordNoPGS,
                    BetAmount = queryBetReocrd.BetAmount,
                    GameID = queryBetReocrd.GameID,
                    MerchantID = queryBetReocrd.MerchantID,
                    PayoutAmount = queryBetReocrd.PayoutAmount,
                    ValidAmount = queryBetReocrd.ValidAmount,
                    DebugProcessUUID = item.LastRequestUUID,
                });

                BetRecordMerchantNotifyQueueRepository.UpdateToSuccess(item.Id);
            }
            catch (Exception ex)
            {

            }
        }

        private void BonusSettleProcess(BetRecordMerchantNotifyQueueEntity item)
        {
            try
            {
                var queryMapping = BetRecordTransactionNoMappingRepository.GetByTransactionNo(item.BetRecordNoPGS);

                if (queryMapping == null)
                    throw new TerproException(ApiErrorCode.NotFoundRecord);

                var service = ServiceInstanceHelper.InstanceVendorService(queryMapping.GameVendorID);

                var queryBetReocrd = service.GetBetRecordByPGSTransactionNo(item.BetRecordNoPGS);

                SingleWalletPlatformApiRequestRefService.BonusSettlement(new SingleWalletModeSettleApiParameter
                {
                    PlayerID = queryBetReocrd.PlayerID,
                    TaxAmount = queryBetReocrd.TaxAmount,
                    SettlementTime = queryBetReocrd.SettlementTime,
                    BetRecordNoPGS = queryBetReocrd.BetRecordNoPGS,
                    BetAmount = queryBetReocrd.BetAmount,
                    GameID = queryBetReocrd.GameID,
                    MerchantID = queryBetReocrd.MerchantID,
                    PayoutAmount = queryBetReocrd.PayoutAmount,
                    ValidAmount = queryBetReocrd.ValidAmount,
                    DebugProcessUUID = item.LastRequestUUID,
                });

                BetRecordMerchantNotifyQueueRepository.UpdateToSuccess(item.Id);
            }
            catch (Exception ex)
            {

            }
        }


        private DateTime CalculateNextTime(DateTime originalNextTime , int retryTimes)
        {
            retryTimes = retryTimes + 1;

            //超過N次不再處理
            if(retryTimes >= 7)
            {
                return new DateTime(2999, 1, 1);
            }

            //3的N次方 3 9 27 81 243 
            double extendMinutes = Math.Pow(3, (retryTimes));
            return originalNextTime.AddMinutes(extendMinutes);
        }


    }
}
