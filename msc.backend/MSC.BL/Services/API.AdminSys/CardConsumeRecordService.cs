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

    public class CardConsumeRecordService : BaseAdminService, ICardConsumeRecordService
    {
        public CardConsumeRecordRepository CardConsumeRecordRepository { get; set; }

        public CardConsumeMerchantNotifyQueueRepository CardConsumeMerchantNotifyQueueRepository { get; set; }
        

        public MerchantRepository MerchantRepository { get; set; }


        public PagedList<CardConsumeRecordViewModel> GetByPagingFilter(int userID , QueryCardConsumeRecordApiRequestParameterModel filter)
        {
            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailableMerchants(userID, filter.MerchantIDs);

            var query = CardConsumeRecordRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<CardConsumeRecordEntity, CardConsumeRecordViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            //額外欄位設置
            SetExtraAttributeValue(viewModels.Items);

            return viewModels;
        }


        public CardConsumeRecordViewModel GetByID(int userID, int targetID)
        {
            var queryCardConsumeRecord = CardConsumeRecordRepository.GetById(targetID);

            if (queryCardConsumeRecord == null)
                throw new MSCException(ErrorCode.NotFoundData);

            CheckHasMerchantViewPermission(userID, queryCardConsumeRecord.MerchantID);

            var convert = ConvertToViewModels(queryCardConsumeRecord);
            return convert;
        }

             
        private CardConsumeRecordViewModel ConvertToViewModels(CardConsumeRecordEntity queryCardConsumeRecord)
        {
            var convert = queryCardConsumeRecord.ReturnInjecteObject<CardConsumeRecordViewModel>();
            convert.MerchantName = MerchantRepository.GetSingleByCache(convert.MerchantID)?.DisplayName;

            return convert;
        }


        private void SetExtraAttributeValue(List<CardConsumeRecordViewModel> viewModels)
        {
            List<string> requestIDs = viewModels.Select(z => z.MerchantRequestID).ToList();

            var queryNotify = CardConsumeMerchantNotifyQueueRepository.GetByFilter(new QueryCardConsumeMerchantNotifyQueueRequestViewModel
            {
                MerchantRequestIDs = requestIDs
            });

            foreach (var item in viewModels)
            {
                var matchNotify = queryNotify.FirstOrDefault(z => z.MerchantRequestID == item.MerchantRequestID && z.MerchantID == item.MerchantID);

                if (matchNotify != null)
                {
                    item.IsNotifyCompleted = matchNotify.Status == NotifyQueueStatusEnum.Success.ToInteger();
                }
            }

            return;
        }



    }
}
