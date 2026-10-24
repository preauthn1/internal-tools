using game.Entity.Game;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using target.BL.Interface;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Web.Modules.Filter;
using System;
using System.Collections.Generic;
using example.invalid;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class VendorConfigController : BaseController
    {
        public IVendorConfigService VendorConfigService { get; set; }


        [Route("vendor-config")]
        [HttpGet]
        public PagedList<VendorConfigViewModel> GetVendorConfigPagingByFilter([FromQuery] QueryVendorConfigApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = VendorConfigService.GetByPagingFilter(userID, filter);
            return query;
        }


        [Route("vendor-config/{vendorConfigID}/settings")]
        [HttpGet]
        public SystemConfigSettingPageViewModel GetVendorConfigSettins(int vendorConfigID)
        {
            int userID = GetCurrentUserID();
            var query = VendorConfigService.GetVendorConfigSettins(vendorConfigID);
            return query;
        }


        [Route("vendor-config/{vendorConfigID}/sync-unassign-merchant")]
        [HttpGet]
        public void SyncDefaultVendorConfigMapping(int vendorConfigID)
        {
            int userID = GetCurrentUserID();
            VendorConfigService.SyncDefaultVendorConfigMapping(vendorConfigID);
            return;
        }


        [Route("vendor-config/merchants")]
        [HttpDelete]
        public void DeleteVendorConfigMerchantMapping([FromBody]VendorConfigDeleteRequestModel model)
        {
            int userID = GetCurrentUserID();
            VendorConfigService.DeleteVendorConfigMerchantMapping(model);
            return;
        }


    }
}
