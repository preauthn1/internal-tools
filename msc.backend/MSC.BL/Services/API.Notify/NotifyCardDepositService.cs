using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;
using MSC.Utility.Utility;
using System.Threading.Tasks;

namespace MSC.BL.Services
{

    public class NotifyCardDepositService : BaseVendorService, INotifyCardDepositService
    {

        public CardConsumeMerchantNotifyQueueRepository CardConsumeMerchantNotifyQueueRepository { get; set; }

        public CardConsumeRecordRepository CardConsumeRecordRepository { get; set; }
        

        public IExternalApiRequestLogRefService ExternalApiRequestLogRefService { get; set; }


        private const int MAX_RETRY_TIMES = 5;


        public NotifyCardDepositViewModel WriteQueueAndTriggerApiImmediately(NotifyCardDepositRequestModel model)
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            CardConsumeRecordEntity queryConsumeRecord = CardConsumeRecordRepository.GetById(model.CardConsumeRecordID);

            if (queryConsumeRecord == null)
                throw new MSCException(ErrorCode.NotFoundCardConsume);

            //寫入queue (第一次立即執行、狀態為處理中)
            long queueJob_ID = CardConsumeMerchantNotifyQueueRepository.Create(new CardConsumeMerchantNotifyQueueEntity
            {
                MerchantRequestID = queryConsumeRecord.MerchantRequestID,
                MerchantID = queryConsumeRecord.MerchantID,
                CardNo = queryConsumeRecord.CardNo,
                ReferenceOrderNo = queryConsumeRecord.OrderNo,
                CreateTime = now,
                UpdateTime = now,
                Status = NotifyQueueStatusEnum.On_Executed.ToInteger(),
                MemberCode = queryConsumeRecord.MemberCode,
                ReferenceID = queryConsumeRecord.ID,
                RetryTimes = 0,
                NotifyUrl = model.NotifyApiUrl,
                NotifyType = MerchantNotifyTypeEnum.NotifyCardInputSuccess.ToInteger(),
            });

            //通知三方API (非同步處理)
            TriggerExternalAPIAndUpdateStatusAsync(model, queueJob_ID, queryConsumeRecord);
            return null;
        }


        public void TriggerNotifyAgainByAdminSys(int userID , TriggerNotifyAgainRequestModel model)
        {
            //TODO 檢查權限之類的

            //觸發通知
            TriggerNotifyAgain(model);
        }



        public void TriggerNotifyAgain(TriggerNotifyAgainRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.MerchantRequestID).Rule_NotEmptyText();

            CardConsumeRecordEntity queryConsumeRecord = CardConsumeRecordRepository.GetByMerchantIDAndRequestID(model.MerchantID, model.MerchantRequestID);

            if (queryConsumeRecord == null)
                throw new MSCException(ErrorCode.NotFoundCardConsume);

            var queryQueue = CardConsumeMerchantNotifyQueueRepository.GetByFilter(new QueryCardConsumeMerchantNotifyQueueRequestViewModel
                {
                    MerchantIDs = queryConsumeRecord.MerchantID.ToListCollection(),
                    MerchantRequestID = queryConsumeRecord.MerchantRequestID
                })
                .FirstOrDefault();

            if (queryQueue == null)
                throw new MSCException(ErrorCode.NotFoundNotifyQueue);

            if (model.IsSkipStatusCheck == false && queryQueue.Status == (int)NotifyQueueStatusEnum.Success)
                throw new MSCException(ErrorCode.InvalidStatus);

            Task.Factory.StartNew(() =>
            {
                var newModel = new NotifyCardDepositRequestModel
                {
                    CardConsumeRecordID = queryConsumeRecord.ID,
                    NotifyApiUrl = queryQueue.NotifyUrl
                };

                TriggerExternalAPIAndUpdateStatus(newModel, queryQueue.Id, queryConsumeRecord);
            });
        }


        private void TriggerExternalAPIAndUpdateStatusAsync(NotifyCardDepositRequestModel model, long notify_queue_id, CardConsumeRecordEntity queryCardConsume)
        {
            Task.Factory.StartNew(() =>
            {
                TriggerExternalAPIAndUpdateStatus(model, notify_queue_id, queryCardConsume);
            });
        }


        private void TriggerExternalAPIAndUpdateStatus(NotifyCardDepositRequestModel model, long notify_queue_id, CardConsumeRecordEntity queryCardConsume)
        {
            //透過uuid可以匹配 ExternalApiRequestLog 跟 LastTimeRequestUUID
            string uuid = Guid.NewGuid().ToString();

            try
            {
                var paras = new RequestParameters
                {
                    BaseUrl = model.NotifyApiUrl,
                    ApiUrl = APIRequestHelper.No_Url_Keyword,
                    LogAction = WriteRequestLog,
                    Method = HttpMethodType.HttpPost,
                    AdvancedSetting = new AdvancedSetting
                    {
                        Log_Group = queryCardConsume.OrderNo,
                        Log_ApiCategory = "CardSuccessNotify",
                        Log_RequestUUID = ""
                    }
                };

                var queryMerchant = MerchantRepository.GetById(queryCardConsume.MerchantID);

                var requestBody = new CardSuccessNotifyRequestModel
                {
                    RequestInfo = new CardSuccessNotifyRequestInfo_V1_Model
                    {
                        TS = DateTimeHelper.GetCurrentTimeTimestamp(),
                        CardNo = queryCardConsume.CardNo,
                        IPAddress = queryCardConsume.IPAddress,
                        CreateTime = DateTimeHelper.ToISO_8601(queryCardConsume.CreateTime),
                        MemberCode = queryCardConsume.MemberCode,
                        MerchantID = queryCardConsume.MerchantID,
                        MSCOrderID = queryCardConsume.OrderNo,
                        MerchantRequestID = queryCardConsume.MerchantRequestID,
                        RandomID = Guid.NewGuid().ToString(),
                        Point = queryCardConsume.Point
                    },
                };

                var signKeyResult = SignKeyRuleHelper.GenerateVendorApiIntegrationSignKey(requestBody.RequestInfo, queryMerchant.ApiSecretKey);
                requestBody.SignKey = signKeyResult.AfterMD5;

                //立即執行一次api
                var result = APIRequestHelper.RequestApiResult<string>(paras, requestBody);

                bool isContainsOK = (result.ResponseText ?? "").Contains("OK");

                //若沒有拿到OK回應或拿到exception
                if (!isContainsOK || result.Exception != null)
                {
                    string msg = result.Exception?.Message + "response: " + result.ResponseText;

                    //更新狀態 以及 更新某些response資訊
                    CardConsumeMerchantNotifyQueueRepository.UpdateStatusAndLastExecutedInfo(notify_queue_id, msg, uuid);
                }
                else
                {
                    //更新狀態 以及 更新某些response資訊
                    CardConsumeMerchantNotifyQueueRepository.UpdateStatus(notify_queue_id, NotifyQueueStatusEnum.Success);
                }
            }
            catch (Exception ex)
            {
                //照理說不會進到這邊 更新狀態 以及 更新某些response資訊
                CardConsumeMerchantNotifyQueueRepository.UpdateStatusAndLastExecutedInfo(notify_queue_id, ex.GetMessageAndErrorLine(), uuid);
            }
        }

        private void WriteRequestLog(ApiRequestLogTemporaryInfoModel model)
        {
            ExternalApiRequestLogRefService.Create(model);
        }

    }
}
