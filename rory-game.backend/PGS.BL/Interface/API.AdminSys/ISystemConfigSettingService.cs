using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface ISystemConfigSettingService
    {

        SystemConfigSettingPageViewModel GetAllSystemConfigSetting(int userID);

        SystemConfigSettingParseModel GetSystemConfigSetting(SystemConfigSettingEnum settingKey);

        void UpdateSystemSettings(int userID, SystemConfigSettingUpdateRequestModel model);

    }
}
