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
    public class VendorController : BaseController
    {
        public IVendorService VendorService { get; set; }


        [Route("vendors")]
        [HttpGet]
        public PagedList<VendorViewModel> GetVendorPagingByFilter([FromQuery] QueryVendorApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = VendorService.GetByPagingFilter(userID , filter);
            return query;
        }


        [Route("vendors/{targetID}")]
        [HttpGet]
        public VendorViewModel GetByID(int targetID)
        {
            var query = VendorService.GetByID(targetID);
            return query;
        }


        [Route("vendors")]
        [HttpPost]
        public void Create(VendorCreateRequestModel model)
        {
            VendorService.Create(model);
        }


        [Route("vendors")]
        [HttpPut]
        public void Update(VendorUpdateRequestModel model)
        {
            VendorService.Update(model);
        }


    }
}
