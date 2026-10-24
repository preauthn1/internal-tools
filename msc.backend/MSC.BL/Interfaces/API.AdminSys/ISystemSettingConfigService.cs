using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.BL.Interface
{
    public interface ISystemSettingConfigService
    {
        SystemConfigSettingPageViewModel GetGlobalSystemSettings();

        void UpdateGlobalSystemSettings(int userID, List<SystemConfigSettingUpdateRequestItemModel> models);
    }
}
