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
using target.Models.ViewModels.API.AdminSys;
using Microsoft.Extensions.Caching.Memory;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class DropdownController : BaseController
    {
        public IAdminSysDropdownService AdminSysDropdownService { get; set; }

        /// <summary>
        /// 组织类型
        /// </summary>
        /// <returns></returns>
        [Route("dropdown/organizations")]
        [HttpGet]
        public List<DropdownValueViewModel> GetAvailableOrgs()
        {
            int userID = GetCurrentUserID();
            var query = AdminSysDropdownService.GetAvailableOrgs(userID);
            return query;
        }

        /// <summary>
        /// 商户
        /// </summary>
        [Route("dropdown/merchants")]
        [HttpGet]
        public List<DropdownValueViewModel> GetAvailableMerchants([FromQuery]QueryDropdownMerchantRequestViewModel filter)
        {
            int userID = GetCurrentUserID();
            var query = AdminSysDropdownService.GetAvailableMerchants(userID , filter);
            return query;
        }

        /// <summary>
        /// 商户分類
        /// </summary>
        [Route("dropdown/merchant-types")]
        [HttpGet]
        public List<DropdownValueViewModel> GetMerchantTypes()
        {
            int userID = GetCurrentUserID();
            var query = AdminSysDropdownService.GetMerchantTypes();
            return query;
        }

        /// <summary>
        /// 內部測試用商户
        /// </summary>
        [Route("dropdown/internal-merchants")]
        [HttpGet]
        public List<DropdownValueViewModel> GetInternalMerchants()
        {
            int userID = GetCurrentUserID();
            var query = AdminSysDropdownService.GetInternalMerchants(userID);
            return query;
        }


        /// <summary>
        /// 簡易查詢的項目
        /// </summary>
        [Route("dropdown/dynamic-query-items")]
        [HttpGet]
        public List<DropdownValueViewModel> GetDynamicQueryItems()
        {
            int userID = GetCurrentUserID();
            var query = AdminSysDropdownService.GetDynamicQueryItems();
            return query;
        }


        /// <summary>
        /// 遊戲大廳房間類型
        /// </summary>
        [Route("dropdown/room-levels")]
        [HttpGet]
        public List<DropdownValueViewModel> GetRoomLevels()
        {
            var query = AdminSysDropdownService.GetRoomLevels();
            return query;
        }

        
        /// <summary>
        /// 遊戲分類
        /// </summary>
        [Route("dropdown/game-types")]
        [HttpGet]
        public List<DropdownValueViewModel> GetGameTypes()
        {
            var query = AdminSysDropdownService.GetGameTypes();
            return query;
        }


        /// <summary>
        /// 投注單狀態
        /// </summary>
        [Route("dropdown/bet-record-status")]
        [HttpGet]
        public List<DropdownValueViewModel> GetBetRecordStatus()
        {
            var query = AdminSysDropdownService.GetBetRecordStatus();
            return query;
        }


        /// <summary>
        /// 遊戲三方廠商
        /// </summary>
        [Route("dropdown/game-vendors")]
        [HttpGet]
        public List<DropdownValueViewModel> GetGameVendorTypes()
        {
            var query = AdminSysDropdownService.GetGameVendorTypes();
            return query;
        }


        /// <summary>
        /// 角色类型
        /// </summary>
        /// <returns></returns>
        [Route("dropdown/roles")]
        [HttpGet]
        public List<DropdownValueViewModel> GetAvailableRoles()
        {
            int userID = GetCurrentUserID();
            var query = AdminSysDropdownService.GetAvailableRoles(userID);
            return query;
        }

        /// <summary>
        /// 游戏类型
        /// </summary>
        /// <returns></returns>
        [Route("dropdown/games")]
        [HttpGet]
        public List<DropdownValueViewModel> GetGames()
        {
            int userID = GetCurrentUserID();
            var query = AdminSysDropdownService.GetGames(userID);
            return query;
        }

        /// <summary>
        /// 游戏伺服器类型
        /// </summary>
        [Route("dropdown/servers")]
        [HttpGet]
        public List<DropdownValueViewModel> GetServers([FromQuery]QueryServerListApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = AdminSysDropdownService.GetServers(userID , filter);
            return query;
        }

        /// <summary>
        /// 时区类型
        /// </summary>
        [Route("dropdown/timezone-areas")]
        [HttpGet]
        public List<DropdownValueViewModel> GetTimeZoneAreas()
        {
            var query = AdminSysDropdownService.GetTimeZoneAreas();
            return query;
        }


        /// <summary>
        /// 交易類型
        /// </summary>
        [Route("dropdown/transactionTypes")]
        [HttpGet]
        public List<DropdownValueViewModel> GetTransactionTypes()
        {
            var query = AdminSysDropdownService.GetTransactionTypes();
            return query;
        }


        /// <summary>
        /// 三方遊戲商
        /// </summary>
        [Route("dropdown/third-party-vendors")]
        [HttpGet]
        public List<DropdownValueViewModel> GetThirdPartyVendors()
        {
            var query = AdminSysDropdownService.GetThirdPartyVendors();
            return query;
        }


        /// <summary>
        /// 取得內部測試帳號
        /// </summary>
        [Route("dropdown/internal-test-accounts")]
        [HttpGet]
        public List<DropdownValueViewModel> GetInternalTestPlayerAccounts()
        {
            var query = AdminSysDropdownService.GetInternalTestPlayerAccounts();
            return query;
        }

    }
}
