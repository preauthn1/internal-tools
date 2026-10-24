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

namespace target.BL.Services
{

    public class TransferRecordVendorService : BaseVendorService, ITransferRecordVendorService
    {
        public TransferRecordRepository TransferRecordRepository { get; set; }


        public TransferRecordVendorViewModel GetTransferByMerchant(QueryTransferRecordRequestModel filter)
        {
            this.GetMerchantAndCheckIsValid(filter.MerchantID, filter.MerchantSecretKey);
            
            //不允許兩個都為空
            if(string.IsNullOrWhiteSpace(filter.MerchantTransactionNo) &&
               string.IsNullOrWhiteSpace(filter.PGSTransactionNo))
            {
                throw new TerproException(ApiErrorCode.InvalidParameter, "MerchantTransactionNo must have value.");
            }
            

            var queryTransfer = TransferRecordRepository.GetByFilter(new QueryTransferRecordRequestViewModel
            {
                MerchantIDs = filter.MerchantID.ToListCollection(),
                OrderNo = filter.PGSTransactionNo,
                MerchantOrderNo = filter.MerchantTransactionNo
            });

            var convert = ConvertToViewModel(queryTransfer.FirstOrDefault());
            return convert;
        }
        

        private TransferRecordVendorViewModel ConvertToViewModel(TransferRecordEntity query)
        {
            if (query == null)
                return null;

            var convert = query.ReturnInjecteObject<TransferRecordVendorViewModel>();

            convert.PGSTransactionNo = query.OrderNo;
            convert.MerchantTransactionNo = query.MerchantOrderNo;
            return convert;
        }




    }
}
