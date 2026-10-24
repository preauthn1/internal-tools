using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.Utility;
using MSC.DA.Repository;
using MSC.Models.Table;
using System.Transactions;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Models.ViewModels.API.Point;

namespace MSC.BL.Services
{

    public class CardLinkRequestRecordService : BaseAdminService, ICardLinkRequestRecordService
    {
        public CardLinkRequestRecordRepository CardLinkRequestRecordRepository { get; set; }

        public CardConsumeRecordRepository CardConsumeRecordRepository { get; set; }

        public CardConsumeMerchantNotifyQueueRepository CardConsumeMerchantNotifyQueueRepository { get; set; }


        public MerchantRepository MerchantRepository { get; set; }


        public PagedList<CardLinkRequestRecordViewModel> GetByPagingFilter(int userID, QueryCardLinkRequestRecordApiRequestParameterModel filter)
        {
            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailableMerchants(userID, filter.MerchantIDs);

            var query = CardLinkRequestRecordRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<CardLinkRequestRecordEntity, CardLinkRequestRecordViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            //設置額外數值
            SetExtraAttributeValue(viewModels.Items);

            return viewModels;
        }


        public MerchantRequestProcessDetailInfoViewModel GetProcessDetail(int userID , MerchantRequestProcessDetailInfoRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.MerchantRequestID).Rule_NotEmptyText();

            var queryStep1 = CardLinkRequestRecordRepository.GetByMerchantIDAndRequestID(model.MerchantID , model.MerchantRequestID);

            if (queryStep1 == null)
                throw new MSCException(ErrorCode.NotFoundData);

            CheckHasMerchantViewPermission(userID, queryStep1.MerchantID);

            var queryStep2 = CardConsumeRecordRepository.GetByMerchantIDAndRequestID(model.MerchantID, model.MerchantRequestID);
            var queryStep3 = CardConsumeMerchantNotifyQueueRepository.GetByMerchantIDAndRequestID(model.MerchantID, model.MerchantRequestID);

            var result = new MerchantRequestProcessDetailInfoViewModel
            {
                StepRequestLink = ConvertTo(queryStep1 != null, true, queryStep1.CreateTime),
                StepCardConsume = ConvertTo(queryStep2 != null, queryStep2 != null, queryStep2?.CreateTime),
                StepNotify = ConvertTo(queryStep3 != null, queryStep3?.Status == (int)NotifyQueueStatusEnum.Success, queryStep3?.UpdateTime),       
                ReturnUrl = queryStep1.ReturnUrl,
                NotifyUrl = queryStep1.NotifyUrl,
                MerchantRequestID = model.MerchantRequestID,
                CreateTime = queryStep1.CreateTime,
                MerchantID = queryStep1.MerchantID,
                MerchantName = MerchantRepository.GetSingleByCache(queryStep1.MerchantID)?.DisplayName,
                MemberCode = queryStep1.MemberCode,
            };


            if(queryStep2 != null)
            {
                result.CardNo = queryStep2.CardNo;
                result.PointAmount = queryStep2.Point;
            }

            return result;
        }


        private RequestStepInfoViewModel ConvertTo(bool isExist , bool? isSuccess , DateTime? lastTime)
        {
            if (!isExist)
                return null;

            return new RequestStepInfoViewModel
            {
                IsSuccess = isSuccess.Value,
                LastTime = lastTime.Value
            };
        }


        public CardLinkRequestRecordViewModel GetByID(int targetID)
        {
            var queryCardLinkRequestRecord = CardLinkRequestRecordRepository.GetById(targetID);

            if (queryCardLinkRequestRecord == null)
                throw new MSCException(ErrorCode.NotFoundData);

            var convert = ConvertToViewModels(queryCardLinkRequestRecord);
            return convert;
        }



        private CardLinkRequestRecordViewModel ConvertToViewModels(CardLinkRequestRecordEntity queryCardLinkRequestRecord)
        {
            var convert = queryCardLinkRequestRecord.ReturnInjecteObject<CardLinkRequestRecordViewModel>();
            convert.MerchantName = MerchantRepository.GetSingleByCache(convert.MerchantID)?.DisplayName;
            convert.IsNotifyCompleted = false;
            return convert;
        }

        private void SetExtraAttributeValue(List<CardLinkRequestRecordViewModel> viewModels)
        {
            List<string> requestIDs = viewModels.Where(z=> z.IsDisabled).Select(z => z.MerchantRequestID).ToList();

            var queryNotify = CardConsumeMerchantNotifyQueueRepository.GetByFilter(new QueryCardConsumeMerchantNotifyQueueRequestViewModel
            {
                MerchantRequestIDs = requestIDs
            });

            foreach(var item in viewModels)
            {
                var matchNotify = queryNotify.FirstOrDefault(z => z.MerchantRequestID == item.MerchantRequestID && z.MerchantID == item.MerchantID);

                if(matchNotify != null)
                {
                    item.IsNotifyExist = true;
                    item.IsNotifyCompleted = matchNotify.Status == NotifyQueueStatusEnum.Success.ToInteger();
                }
            }

            return;
        }


    }
}
