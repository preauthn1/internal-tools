using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IMerchantCustomizedSettingService
    {
        SystemConfigSettingPageViewModel GetMerchantAdvanceSettings(int userID, int merchantID);

        void UpdateSettings(int userID, SystemConfigSettingUpdateRequestModel model);
    }
}
