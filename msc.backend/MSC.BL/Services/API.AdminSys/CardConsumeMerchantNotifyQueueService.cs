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

    public class CardConsumeMerchantNotifyQueueService : BaseAdminService, ICardConsumeMerchantNotifyQueueService
    {
        public CardConsumeMerchantNotifyQueueRepository CardConsumeMerchantNotifyQueueRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public PagedList<CardConsumeMerchantNotifyQueueViewModel> GetByPagingFilter(int userID , QueryCardConsumeMerchantNotifyQueueApiRequestParameterModel filter)
        {
            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailableMerchants(userID, filter.MerchantIDs);

            var query = CardConsumeMerchantNotifyQueueRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<CardConsumeMerchantNotifyQueueEntity, CardConsumeMerchantNotifyQueueViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        public CardConsumeMerchantNotifyQueueViewModel GetByID(int userID , int notifyQueueID)
        {
            var queryCardConsumeMerchantNotifyQueue = CardConsumeMerchantNotifyQueueRepository.GetById(notifyQueueID);

            if (queryCardConsumeMerchantNotifyQueue == null)
                throw new MSCException(ErrorCode.NotFoundNotifyQueue);

            CheckHasMerchantViewPermission(userID , queryCardConsumeMerchantNotifyQueue.MerchantID);

            var convert = ConvertToViewModels(queryCardConsumeMerchantNotifyQueue);
            return convert;
        }
        

        private CardConsumeMerchantNotifyQueueViewModel ConvertToViewModels(CardConsumeMerchantNotifyQueueEntity model)
        {
            var convert = model.ReturnInjecteObject<CardConsumeMerchantNotifyQueueViewModel>(); 
            convert.MerchantName = MerchantRepository.GetSingleByCache(convert.MerchantID)?.DisplayName;
            convert.StatusText = MultiLanguageTextHelper.GetEnums_i18nText((NotifyQueueStatusEnum)model.Status);
            return convert;
        }

    }
}
