using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Web.Modules.Filter;
using target.Models.Paging;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class SystemConfigSettingController : BaseController
    {

        public ISystemConfigSettingService SystemConfigSettingService { get; set; }


        [Route("system-config-setting")]
        [HttpGet]
        public SystemConfigSettingPageViewModel GetAllSystemConfigSetting()
        {
            int userID = GetCurrentUserID();
            var query = SystemConfigSettingService.GetAllSystemConfigSetting(userID);
            return query;
        }


        [Route("system-config-setting")]
        [HttpPost]
        public void UpdateUserLanguage([FromBody]SystemConfigSettingUpdateRequestModel model)
        {
            int userID = GetCurrentUserID();
            SystemConfigSettingService.UpdateSystemSettings(userID , model);
        }

    }
}
