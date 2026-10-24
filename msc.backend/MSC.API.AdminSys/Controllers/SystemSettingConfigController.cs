using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.Enums;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Web.Modules.Filter;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class SystemSettingConfigController : BaseController
    {       
        public ISystemSettingConfigService SystemSettingConfigService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_SystemSettingConfig_List)]
        [Route("system-settings")]
        [HttpGet]
        public SystemConfigSettingPageViewModel GetGlobalSystemSettings()
        {
            int userID = GetCurrentUserID();
            var query = SystemSettingConfigService.GetGlobalSystemSettings();
            return query;
        }

        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_SystemSettingConfig_List)]
        [Route("system-settings")]
        [HttpPost]
        public void UpdateGlobalSystemSettings([FromBody] List<SystemConfigSettingUpdateRequestItemModel> models)
        {
            int userID = GetCurrentUserID();
            SystemSettingConfigService.UpdateGlobalSystemSettings(userID , models);
            return;
        }


    }
}
