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
    public class ServerRobotConfigController : BaseController
    {
        public IServerRobotConfigService ServerRobotConfigService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServerRobotConfig_Search)]
        [Route("server/robot-config")]
        [HttpGet]
        public PagedList<ServerRobotConfigViewModel> GetServerRobotConfigPagingByFilter([FromQuery]QueryServerRobotConfigApiRequestParameterModel filter)
        {
            var query = ServerRobotConfigService.GetByPagingFilter(filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServerRobotConfig_Search)]
        [Route("server/robot-config/{serverID}")]
        [HttpGet]
        public ServerRobotConfigViewModel GetByID(int serverID)
        {
            var query = ServerRobotConfigService.GetByID(serverID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServerRobotConfig_Create)]
        [Route("server/robot-config")]
        [HttpPost]
        public void Create(ServerRobotConfigCreateRequestModel model)
        {
            ServerRobotConfigService.Create(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServerRobotConfig_Update)]
        [Route("server/robot-config")]
        [HttpPut]
        public void Update(ServerRobotConfigUpdateRequestModel model)
        {
            ServerRobotConfigService.Update(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServerRobotConfig_Delete)]
        [Route("server/robot-config/{targetID}")]
        [HttpDelete]
        public void Delete(int targetID)
        {
            ServerRobotConfigService.Delete(targetID);
        }

    }
}
