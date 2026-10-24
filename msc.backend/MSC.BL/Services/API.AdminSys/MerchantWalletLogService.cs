using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Transactions;

namespace MSC.BL.Services
{
    public class MerchantWalletLogService : BaseAdminService, IMerchantWalletLogService
    {
        public MerchantWalletLogRepository MerchantWalletLogRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public PagedList<MerchantWalletLogViewModel> GetByPagingFilter(int userID , QueryMerchantWalletLogApiRequestParameterModel filter)
        {
            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailableMerchants(userID, filter.MerchantIDs);

            var query = MerchantWalletLogRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<MerchantWalletLogEntity, MerchantWalletLogViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        public MerchantWalletLogViewModel GetByID(int targetID)
        {
            var queryMerchantWalletLog = MerchantWalletLogRepository.GetById(targetID);

            if (queryMerchantWalletLog == null)
                throw new MSCException(ErrorCode.NotFoundData);

            var convert = ConvertToViewModels(queryMerchantWalletLog);
            return convert;
        }


        private MerchantWalletLogViewModel ConvertToViewModels(MerchantWalletLogEntity query)
        {
            var convert = query.ReturnInjecteObject<MerchantWalletLogViewModel>();

            convert.MerchantName = MerchantRepository.GetSingleByCache(query.MerchantID)?.DisplayName;
            convert.TransactionTypeText = MultiLanguageTextHelper.GetEnums_i18nText((TransactionTypeEnum)query.TransactionType);

            return convert;
        }

    }
}
