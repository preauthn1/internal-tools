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
    public class RoleController : BaseController
    {
        public IRoleService RoleService { get; set; }



        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingRoles_Search)]
        [Route("roles")]
        [HttpGet]
        public PagedList<RoleViewModel> GetRolePagingByFilter([FromQuery]QueryRoleApiRequestParameterModel filter)
        {
            var query = RoleService.GetByPagingFilter(filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingRoles_Search)]
        [Route("roles/{roleID}")]
        [HttpGet]
        public RoleViewModel GetByID(int roleID)
        {
            var query = RoleService.GetByID(roleID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingRoles_Create)]
        [Route("roles")]
        [HttpPost]
        public void Create(RoleCreateRequestModel model)
        {
            RoleService.Create(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingRoles_Update)]
        [Route("roles")]
        [HttpPut]
        public void Update(RoleUpdateRequestModel model)
        {
            RoleService.Update(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingRoles_Delete)]
        [Route("roles/{roleID}")]
        [HttpDelete]
        public void Delete(int roleID)
        {
            RoleService.Delete(roleID);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingRoles_Update , AdminSysPermissionKey.SystemSettingRoles_Rename)]
        [Route("roles/{roleID}/permissions")]
        [HttpPut]
        public void UpdateRolePermissions(int roleID , [FromBody]RolePermissionUpdateModel model)
        {
            model = model ?? new RolePermissionUpdateModel();
            model.RoleID = roleID;

            RoleService.UpdateRolePermissions(model);
        }


    }
}
