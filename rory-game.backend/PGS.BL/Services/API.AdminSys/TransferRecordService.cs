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
using example.invalid;

namespace target.BL.Services
{

    public class TransferRecordService : BaseAdminSysService, ITransferRecordService
    {
        public TransferRecordRepository TransferRecordRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }


        public PagedList<TransferRecordViewModel> GetByPagingFilter(int userID , QueryTransferRecordApiRequestParameterModel filter)
        {
            var currentUser = GetCurrentOwnPermissionByCache(userID);

            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailablePermissionMerchantIDs(userID, filter.MerchantIDs);
            filter.PlayerIDs = filter.PlayerID.ToIntegerListCollection();

            var query = TransferRecordRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<TransferRecordEntity, TransferRecordViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        public Stream ExportExcelByPagingFilter(int userID, QueryTransferRecordApiRequestParameterModel filter)
        {
            filter.PageSize = 10000;
            var query = GetByPagingFilter(userID , filter);

            var excelStream = ExcelDataConvertHelper.ConvertToGameUserGameInfoExcelStream(query.Items);
            return excelStream;
        }
        

        private TransferRecordViewModel ConvertToViewModels(TransferRecordEntity queryTransferRecord)
        {
            var convert = queryTransferRecord.ReturnInjecteObject<TransferRecordViewModel>();
            convert.MerchantName = MerchantRepository.GetSingleByCache(queryTransferRecord.MerchantID)?.Name;
            convert.TransferTypeText = MultiLanguageTextHelper.GetEnumsCultureText((WalletTransferTypeEnum)queryTransferRecord.TransferType);

            return convert;
        }

    }
}
