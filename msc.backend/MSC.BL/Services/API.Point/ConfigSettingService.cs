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
using MSC.Utility;
using Lck.Utility.Extensions.Helper;

namespace MSC.BL.Services
{

    public class ConfigSettingService : BaseService, IConfigSettingService
    {

        public SystemSettingConfigRepository SystemSettingConfigRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }


        public List<ShoppingStoreViewModel> GetCardShoppingStores()
        {
            //此API已廢棄
            return new List<ShoppingStoreViewModel>();
        }

        public string GetShoppingsIframeUrl()
        {
            var query = SystemSettingConfigRepository.GetSingleByCache(SystemSettingConfigEnum.ShoppingsIframeUrl);

            if (query == null)
                return string.Empty;

            return query.SettingValue;
        }


        public CardInputPageSettingViewModel GetCardInputPageSetting(int merchantID)
        {
            var querySetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.IsBackButtonVisible);
            var querySettingStyle = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.Style_BrandingTypeLayout);

            var result = new CardInputPageSettingViewModel
            {
                WebSiteBrandingType = WebSiteBrandingTypeEnum.OGC.ToInteger(),
                IsBackButtonVisible = true
            };

            if (querySetting != null && !string.IsNullOrWhiteSpace(querySetting.SettingValue))
            {
                result.IsBackButtonVisible = querySetting.SettingValueOfBoolean;
            }

            if (querySetting != null && !string.IsNullOrWhiteSpace(querySetting.SettingValue))
            {
                result.WebSiteBrandingType = querySettingStyle.SettingValueOfInt;
            }

            return result;
        }



    }
}
