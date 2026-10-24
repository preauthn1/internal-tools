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
    public class VendorMaintenanceController : BaseController
    {
        public IVendorMaintenanceService VendorMaintenanceService { get; set; }


        [Route("vendor-maintenance")]
        [HttpGet]
        public PagedList<VendorMaintenanceViewModel> GetVendorMaintenancePagingByFilter([FromQuery] QueryVendorMaintenanceApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = VendorMaintenanceService.GetByPagingFilter(userID, filter);
            return query;
        }


        [Route("vendor-maintenance")]
        [HttpPost]
        public void Create(VendorMaintenanceCreateViewModel model)
        {
            int userID = GetCurrentUserID();
            VendorMaintenanceService.Create(userID, model);
        }



        [Route("vendor-maintenance/{targetID}")]
        [HttpDelete]
        public void Delete(int targetID)
        {
            int userID = GetCurrentUserID();
            VendorMaintenanceService.Delete(userID, targetID);
        }

    }
}
