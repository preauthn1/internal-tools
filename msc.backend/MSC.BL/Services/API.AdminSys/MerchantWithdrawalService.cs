using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.BL.Services
{
    public class MerchantWithdrawalService : BaseAdminService, IMerchantWithdrawalService
    {

        public MerchantWithdrawalRepository MerchantWithdrawalRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public MerchantWalletRepository MerchantWalletRepository { get; set; }


        public PagedList<MerchantWithdrawalViewModel> GetByPagingFilter(int userID, QueryMerchantWithdrawalApiRequestParameterModel filter)
        {
            filter.MerchantIDs = FilterAvailableMerchants(userID, filter.MerchantIDs);

            var query = MerchantWithdrawalRepository.GetPagingByFilter(filter);

            var paging = ConvertToNewTypePagingObject<MerchantWithdrawalEntity, MerchantWithdrawalViewModel>(query);

            paging.Items = query.Items
                .Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return paging;
        }


        public void Create(int userID, MerchantWithdrawalCreateRequestModel model)
        {
            model.RuleFor(z => z.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.Amount).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.AccountMemo).Rule_NotEmptyText();

            var currentUser = GetUserNameByUserID(userID);


            string orderNo = RandomHelper.GetOrderNo("MW");

            // Step 1. 扣除金幣、寫入交易紀錄
            var result = MerchantWalletRepository.UpdateWalletAndWriteLog(new UpdateMerchantWalletSPRequestParameter
            {
                SourceAmount = model.Amount * -1,
                FeeRate = 0,
                MerchantID = model.MerchantID,
                RefNo = orderNo,
                TransactionType = TransactionTypeEnum.MerchantWithdrawal
            });

            if(result.Code != ErrorCode.Success)
            {
                throw new MSCException(result.Code);
            }

            // Step 2. 寫入下發申請紀錄
            var resultID = MerchantWithdrawalRepository.Create(new MerchantWithdrawalEntity
            {
                MerchantID = model.MerchantID,
                Type = 1,
                Status = MerchantWithdrawalStatusEnum.Complete.ToInteger(),
                AccountMemo = model.AccountMemo,
                Memo = model.Memo,
                OrderNo = orderNo,
                CreateBy = currentUser,
                UpdateTime = DateTimeHelper.GetCurrentTime(),
                CreateTime = DateTimeHelper.GetCurrentTime(),
                Amount = model.Amount,
                UpdateBy = currentUser,
            });

            return;
        }


        private MerchantWithdrawalViewModel ConvertToViewModels(MerchantWithdrawalEntity queryEntity)
        {
            var convert = queryEntity.ReturnInjecteObject<MerchantWithdrawalViewModel>();
            convert.MerchantName = MerchantRepository.GetSingleByCache(queryEntity.MerchantID)?.DisplayName;

            return convert;
        }


    }
}
