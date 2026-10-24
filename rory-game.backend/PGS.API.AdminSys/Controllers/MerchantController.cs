using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.Paging;
using target.Web.Modules.Filter;
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class MerchantController : BaseController
    {
        public IMerchantService MerchantService { get; set; }


        public IMerchantCustomizedSettingService MerchantCustomizedSettingService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingMerchant_Search)]
        [Route("merchants")]
        [HttpGet]
        public PagedList<MerchantViewModel> GetMerchantPagingByFilter([FromQuery] QueryMerchantApiRequestParameterModel filter)
        {
            var query = MerchantService.GetByPagingFilter(filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingMerchant_Search)]
        [Route("merchants/{merchantID}")]
        [HttpGet]
        public MerchantViewModel GetByID(int merchantID)
        {
            var query = MerchantService.GetByID(merchantID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingMerchant_Search)]
        [Route("merchants/copy-info/{merchantID}")]
        [HttpGet]
        public MerchantCopyInfoPreviewViewModel GetCopyPreviewInfo(int merchantID)
        {
            var query = MerchantService.GetCopyPreviewInfo(merchantID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingMerchant_Create)]
        [Route("merchants")]
        [HttpPost]
        public void Create(MerchantCreateRequestModel model)
        {

            MerchantService.Create(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingMerchant_Update)]
        [Route("merchants")]
        [HttpPut]
        public void Update(MerchantUpdateRequestModel model)
        {
            MerchantService.Update(model);
        }

        

        [Route("v2/merchants/{merchantID}/advance-setting")]
        [HttpGet]
        public SystemConfigSettingPageViewModel GetMerchantAdvanceSettings_v2(int merchantID)
        {
            int userID = GetCurrentUserID();
            var result = MerchantCustomizedSettingService.GetMerchantAdvanceSettings(userID, merchantID);
            return result;
        }


        [Route("v2/merchants/advance-setting")]
        [HttpPost]
        public void UpdateMerchantAdvanceSettings_v2(SystemConfigSettingUpdateRequestModel model)
        {
            int userID = GetCurrentUserID();
            MerchantCustomizedSettingService.UpdateSettings(userID, model);
            return;
        }


    }
}
