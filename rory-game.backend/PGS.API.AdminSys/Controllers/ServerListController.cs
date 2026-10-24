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
    public class ServerListController : BaseController
    {
        public IServerListService ServerListService { get; set; }


        public IServerMerchantMappingService ServerMerchantMappingService { get; set; }


        public IGameHealthSummaryService GameHealthSummaryService { get; set; }


        public IViewPermissionCheckService ViewPermissionCheckService { get; set; }



        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Search)]
        [Route("servers")]
        [HttpGet]
        public PagedList<ServerListViewModel> GetServerListPagingByFilter([FromQuery]QueryServerListApiRequestParameterModel filter)
        {
            var query = ServerListService.GetByPagingFilter(filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagementServersHealthSummary)]
        [Route("servers/health/summary")]
        [HttpGet]
        public GameHealthSummaryViewModel GetEachServerHealthySummary([FromQuery]GameHealthSummaryRequestModel model)
        {
            int userID = GetCurrentUserID();
            ViewPermissionCheckService.CheckViewMerchantPermission(userID , model?.MerchantID ?? 0);
            var query = GameHealthSummaryService.GetEachServerHealthySummary(userID , model);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Search)]
        [Route("servers/{serverID}")]
        [HttpGet]
        public ServerListViewModel GetByID(int serverID)
        {
            var query = ServerListService.GetByID(serverID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Create)]
        [Route("servers")]
        [HttpPost]
        public void Create(ServerListCreateRequestModel model)
        {
            ServerListService.Create(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Update)]
        [Route("servers")]
        [HttpPut]
        public void Update(ServerListUpdateRequestModel model)
        {
            ServerListService.Update(model);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Delete)]
        [Route("servers/{serverID}")]
        [HttpDelete]
        public void Delete(int serverID)
        {
            ServerListService.Delete(serverID);
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Dispatch)]
        [Route("servers/all-merchant-mapping")]
        [HttpGet]
        public AllServerMerchantMappingViewModel GetAllServerMerchantMappingsByCurrency()
        {
            var query = ServerMerchantMappingService.GetAllServerMerchantMappings();
            return query;
        }

        
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Dispatch)]
        [Route("servers/merchant-mapping/revers")]
        [HttpPut]
        public ReversServerMerchantResultViewModel ReversMerchantMappings(ReversServerMerchantRequestViewModel model)
        {
            var query = ServerMerchantMappingService.ReversMerchantMappings(model);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Dispatch)]
        [Route("servers/{serverID}/availableMerchants")]
        [HttpGet]
        public ServerMerchantMappingViewModel GetAvaialbleMerchants(int serverID)
        {
            var query = ServerMerchantMappingService.GetServerMerchantMappings(serverID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.SystemSettingServers_Dispatch)]
        [Route("servers/{serverID}/availableMerchants")]
        [HttpPut]
        public void UpdateServerMerchantMappings(int serverID ,[FromBody]ServerAvailableMerchantUpdateRequestModel model)
        {
            ServerMerchantMappingService.UpdateServerMerchantMappings(serverID , model);
            return;
        }



    }
}
