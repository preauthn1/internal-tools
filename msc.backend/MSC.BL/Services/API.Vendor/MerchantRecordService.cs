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

namespace MSC.BL.Services
{

    public class MerchantRecordService : BaseVendorService, IMerchantRecordService
    {
        public CardLinkRequestRecordRepository CardLinkRequestRecordRepository { get; set; }
        public CardConsumeRecordRepository CardConsumeRecordRepository { get; set; }
        public CardConsumeMerchantNotifyQueueRepository CardConsumeMerchantNotifyQueueRepository { get; set; }


        public MerchantRequestRecordViewModel GetMerchantRequestRecord(VendorApiQueryMerchantRequestRecordRequestModel model)
        {
            model.RuleFor(x => x.SignKey).Rule_NotEmptyText();
            model.RuleFor(x => x.RequestInfo);

            var request = model.RequestInfo;

            request.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            request.RuleFor(x => x.TS).Rule_NotZeroNotMinus().Rule_IsValidTimestamp();
            request.RuleFor(x => x.RandomID).Rule_NotEmptyText();
            request.RuleFor(x => x.MerchantRequestID).Rule_NotEmptyText();

            var queryMerchant = GetAndCheckIsValidMerchant(model.RequestInfo);

            CheckSignKeyIsValid(model.SignKey, model.RequestInfo, queryMerchant.ApiSecretKey);

            var queryRequest = CardLinkRequestRecordRepository.GetByMerchantIDAndRequestID(request.MerchantID, request.MerchantRequestID);

            if (queryRequest == null)
            {
                return null;
            }

            var result = new MerchantRequestRecordViewModel
            {
                MerchantID = queryRequest.MerchantID,
                CreateTime = queryRequest.CreateTime,
                NotifyUrl = queryRequest.NotifyUrl,
                ReturnUrl = queryRequest.ReturnUrl,
                MemberCode = queryRequest.MemberCode,
                IPAddress = queryRequest.IPAddress,
                MerchantRequestID = queryRequest.MerchantRequestID,
                MerchantOrderAmount = queryRequest.Amount,
                CardInputInfo = new CardInputInfoViewModel
                {
                    IsSuccess = false,
                },
                NotifyInfo = new SuccessNotifyInfoViewModel
                {
                    IsSuccess = false,
                }
            };


            if (queryRequest.IsDisabled)
            {
                var queryConsume = CardConsumeRecordRepository.GetByMerchantIDAndRequestID(request.MerchantID, request.MerchantRequestID);
                var queryNotify = CardConsumeMerchantNotifyQueueRepository.GetByMerchantIDAndRequestID(request.MerchantID, request.MerchantRequestID);


                if (queryConsume != null)
                {
                    result.CardInputInfo = new CardInputInfoViewModel
                    {
                        Amount = queryConsume.Point,
                        CardNo = queryConsume.CardNo,
                        CreateTime = queryConsume.CreateTime,
                        IPAddress = queryConsume.IPAddress,
                        IsSuccess = true
                    };


                }

                if (queryNotify != null)
                {
                    result.NotifyInfo = new SuccessNotifyInfoViewModel
                    {
                        LastNotifyTime = queryNotify.UpdateTime,
                        LastTimeResponseText = queryNotify.LastTimeResponse,
                        TotalRequestTimes = queryNotify.RetryTimes,
                        IsSuccess = queryNotify.Status == (int)NotifyQueueStatusEnum.Success
                    };
                }
            }


            return result;
        }


    }
}
