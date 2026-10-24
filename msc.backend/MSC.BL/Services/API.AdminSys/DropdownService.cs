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

    public class DropdownService : BaseAdminService, IDropdownService
    {
        public MerchantRepository MerchantRepository { get; set; }

        public UserMerchantMappingRepository UserMerchantMappingRepository { get; set; }

        public RoleRepository RoleRepository { get; set; }


        public List<DropdownSelectValueViewModel> GetRoles()
        {
            var query = RoleRepository.GetAll();
            return query
                    .Select(x => new DropdownSelectValueViewModel
                    {
                        Text = x.DisplayName,
                        Value = x.RoleID.ToString()
                    })
                    .ToList();
        }


        public List<DropdownSelectValueViewModel> GetTransactionTypes()
        {
            return ExtensionOfEnum.GetEnumAllValue<TransactionTypeEnum>()
                .Select(x => new DropdownSelectValueViewModel
                {
                    Text = MultiLanguageTextHelper.GetEnums_i18nText(x),
                    Value = x.ToIntegerString()
                })
                .ToList();
        }


        public List<DropdownSelectValueViewModel> GetUserOperateTypes()
        {
            return ExtensionOfEnum.GetEnumAllValue<UserOperateTypeEnum>()
                .Select(x => new DropdownSelectValueViewModel
                {
                    Text = MultiLanguageTextHelper.GetEnums_i18nText(x),
                    Value = x.ToIntegerString()
                })
                .ToList();
        }

        public List<DropdownSelectValueViewModel> GetCurrency()
        {
            return ExtensionOfEnum.GetEnumAllValue<CurrencyEnum>()
                .Select(x => new DropdownSelectValueViewModel
                {
                    Text = x.ToString(),
                    Value = x.ToString()
                })
                .ToList();
        }

        

        public List<DropdownSelectValueViewModel> GetNormalStatus()
        {
            var status = new List<NormalActiveStatus>
            { 
                NormalActiveStatus.Enable,
                NormalActiveStatus.Disable
            };

            return status.Select(x => new DropdownSelectValueViewModel
                {
                    Text = MultiLanguageTextHelper.GetEnums_i18nText(x),
                    Value = x.ToIntegerString()
                })
                .ToList();
        }


        public List<DropdownSelectValueViewModel> GetMerchants(int userID)
        {
            List<int> merchantIDs = GetUserMerchantIDsByCache(userID);

            var query = MerchantRepository.GetByIds(merchantIDs);

            return query.Select(x=> new DropdownSelectValueViewModel
                {
                    Text = x.DisplayName,
                    Value = x.MerchantID.ToString()
                })
                .ToList();
        }


        public List<DropdownSelectValueViewModel> GetMerchantsWithCurrency(int userID)
        {
            List<int> merchantIDs = GetUserMerchantIDsByCache(userID);

            var query = MerchantRepository.GetByIds(merchantIDs);

            return query.Select(x => new DropdownSelectValueViewModel
                {
                    Text = $"{x.Currency}｜{x.DisplayName}",
                    Value = x.MerchantID.ToString()
                })
                .ToList();
        }



    }
}
