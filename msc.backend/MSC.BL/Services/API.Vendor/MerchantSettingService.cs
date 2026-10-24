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

    public class MerchantSettingService : BaseVendorService, IMerchantSettingService
    {


        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }


        public MerchantAvailableAmountViewModel GetAvialbleAmountSetting(VendorApiMerchantAmountSettingRequestModel model)
        {
            var requestInfo = model?.RequestInfo;

            model.RuleFor(x => x.SignKey).Rule_NotEmptyText();
            requestInfo.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            requestInfo.RuleFor(x => x.RandomID).Rule_NotEmptyText();
            requestInfo.RuleFor(x => x.TS).Rule_NotZeroNotMinus().Rule_IsValidTimestamp();

            var queryMerchant = GetAndCheckIsValidMerchant(model.RequestInfo);
            
            CheckSignKeyIsValid(model.SignKey , model.RequestInfo, queryMerchant.ApiSecretKey);

            int merchantID = queryMerchant.MerchantID;
            var querySetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.Basic_AvailableAmounts);


            if (querySetting == null || querySetting.SettingValue.IsNullOrWhiteSpace())
            {
                return new MerchantAvailableAmountViewModel
                {
                    MerchantID = merchantID,
                    Amounts = new List<decimal>()
                };
            }


            return new MerchantAvailableAmountViewModel
            {
                MerchantID = merchantID,
                Amounts = querySetting.SettingValueListDecimal,
            };
        }


    }
}
