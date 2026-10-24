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
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{

    [AdminTokenrAuthorizeFilter]
    public class OrganizationController : BaseController
    {
        public IOrganizationService OrganizationService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingOrganizations_Search)]
        [Route("organizations")]
        [HttpGet]
        public PagedList<OrganizationViewModel> GetPagingByFilter([FromQuery]QueryOrganizationApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = OrganizationService.GetPagingByFilter(filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingOrganizations_Search)]
        [Route("organizations/{orgID}")]
        [HttpGet]
        public OrganizationViewModel GetUserByID(int orgID)
        {
            var query = OrganizationService.GetByID(orgID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingOrganizations_Create)]
        [Route("organizations")]
        [HttpPost]
        public void Create(OrganizationCreateRequestModel model)
        {
            OrganizationService.Create(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingOrganizations_Update)]
        [Route("organizations")]
        [HttpPut]
        public void Update(OrganizationUpdateRequestModel model)
        {
            OrganizationService.Update(model);
        }


    }
}
