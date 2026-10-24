using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;

namespace target.BL.Interface
{
    public interface IVendorConfigService
    {
        PagedList<VendorConfigViewModel> GetByPagingFilter(int userID, QueryVendorConfigApiRequestParameterModel filter);

        SystemConfigSettingPageViewModel GetVendorConfigSettins(int vendorConfigID);

        void SyncDefaultVendorConfigMapping(int vendorID);

        void DeleteVendorConfigMerchantMapping(VendorConfigDeleteRequestModel model);
    }
}
